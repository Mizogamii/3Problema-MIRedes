package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"pbl/server/models"
	sharedRaft "pbl/server/shared"
	"pbl/server/utils"
	"pbl/shared"
	"pbl/style"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"github.com/nats-io/nats.go"
)

func HandleDrawCard(server *models.Server, request shared.Request, nc *nats.Conn, message *nats.Msg) {
	userWallet := ""
	server.Mu.Lock()
	if user, ok := server.Users[request.ClientID]; ok {
		userWallet = user.Address
	}
	server.Mu.Unlock()

	if server.Raft.State() == raft.Leader {
		result, err := processDrawCardRequest(server, request.ClientID)
		if err != nil {
			respondWithError(nc, message, err.Error())
			return
		}
		
		// Salva no inventário local
		saveCardToLocalUser(server, request.ClientID, result)

		style.PrintCian("hora de testar a carteira")
		if userWallet != "" {
			style.PrintMag("tem carteira(lider)")
			triggerBlockchainCriarCarta(server, userWallet, result.Id, result.Element, result.Type)
		} else {
			log.Printf("[Aviso] Líder processou carta mas user %s não tem wallet local.", request.ClientID)
		}

		respondWithSuccess(nc, message, result)
		return
	}

	leaderAddr := string(server.Raft.Leader())
	if leaderAddr == "" {
		respondWithError(nc, message, "Líder não disponível no momento.")
		return
	}

	leaderURL := fmt.Sprintf("http://%s/leader/draw-card", leaderAddr)

	payload := map[string]string{
		"clientID":      request.ClientID,
		"player_wallet": userWallet,
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(leaderURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		respondWithError(nc, message, fmt.Sprintf("Falha ao comunicar com líder: %v", err))
		return
	}
	defer resp.Body.Close()

	// Processa a resposta do líder
	var leaderResponse shared.Response
	if err := json.NewDecoder(resp.Body).Decode(&leaderResponse); err != nil {
		respondWithError(nc, message, "Resposta inválida do líder.")
		return
	}

	if leaderResponse.Status == "success" {
		var drawnData shared.CardDrawnData
		if err := json.Unmarshal(leaderResponse.Data, &drawnData); err != nil {
			respondWithError(nc, message, "Dados inválidos do líder.")
			return
		}
		saveCardToLocalUser(server, request.ClientID, drawnData.Card)
	}

	finalResponseBytes, _ := json.Marshal(leaderResponse)
	nc.Publish(message.Reply, finalResponseBytes)
}


func LeaderDrawCardHandler(server *models.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if server.Raft.State() != raft.Leader {
			http.Error(w, "Eu não sou o líder", http.StatusServiceUnavailable)
			return
		}

		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Payload inválido", http.StatusBadRequest)
			return
		}
		
		clientID := payload["clientID"]
		playerWallet := payload["player_wallet"]

		result, err := processDrawCardRequest(server, clientID)
		
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			response := shared.Response{Status: "error", Error: err.Error(), Server: server.ID}
			json.NewEncoder(w).Encode(response)
			return
		}

		style.PrintCian("hora de testar a carteira")
		if playerWallet != "" {
			style.PrintMag("tem carteira")
			triggerBlockchainCriarCarta(server, playerWallet, result.Id, result.Element, result.Type)
		} else {
			log.Printf("[Aviso Blockchain] Líder recebeu pedido de %s sem wallet.", clientID)
		}

		saveCardToLocalUser(server, clientID, result)

		responseData := shared.CardDrawnData{Card: result, RequestID: "forwarded"}
		responseBytes, _ := json.Marshal(responseData)
		response := shared.Response{Status: "success", Action: "CARD_DRAWN", Data: responseBytes, Server: server.ID}
		json.NewEncoder(w).Encode(response)
	}
}

// saveCardToLocalUser adiciona a carta ao inventário do usuário no servidor local.
func saveCardToLocalUser(server *models.Server, clientID string, card shared.Card) {
	server.Mu.Lock()
	defer server.Mu.Unlock()

	if user, ok := server.Users[clientID]; ok {
		user.Cards = append(user.Cards, card)
		server.Users[clientID] = user
		log.Printf("[%d] Carta '%s' adicionada ao inventário local do cliente %s.", server.ID, card.Type, clientID)
	}
}

func processDrawCardRequest(server *models.Server, clientID string) (shared.Card, error) {
	requestID := uuid.New().String()

	payload := sharedRaft.DrawCardPayload{PlayerID: clientID, RequestID: requestID}
	payloadBytes, _ := json.Marshal(payload)
	cmd := sharedRaft.Command{Type: sharedRaft.CommandOpenPack, Data: payloadBytes}
	cmdBytes, _ := json.Marshal(cmd)

	future := server.Raft.Apply(cmdBytes, 500*time.Millisecond)
	if err := future.Error(); err != nil {
		log.Printf("[%d] Erro ao aplicar comando Raft 'DrawCard': %v", server.ID, err)
		return shared.Card{}, fmt.Errorf("erro interno ao processar a jogada")
	}

	responseValue := future.Response()
	if strValue, ok := responseValue.(string); ok && strValue == "STOCK_EMPTY" {
		return shared.Card{}, fmt.Errorf("o estoque de cartas acabou")
	}

	drawnCard, ok := responseValue.(shared.Card)
	if !ok {
		return shared.Card{}, fmt.Errorf("erro inesperado no tipo de resposta do Raft (esperava shared.Card)")
	}

	log.Printf("[%d] Carta '%s' reservada para o cliente %s (RequestID: %s).", server.ID, drawnCard.Type, clientID, requestID)
	go claimCard(server, requestID)

	return drawnCard, nil
}

// finaliza a transação, removendo a carta da área de pendentes
func claimCard(server *models.Server, requestID string) {
	log.Printf("[%d] Reivindicando carta para o RequestID: %s", server.ID, requestID)

	payload := sharedRaft.ClaimCardPayload{RequestID: requestID}
	payloadBytes, _ := json.Marshal(payload)

	cmd := sharedRaft.Command{
		Type: sharedRaft.CommandClaimCard,
		Data: payloadBytes,
	}
	cmdBytes, _ := json.Marshal(cmd)

	future := server.Raft.Apply(cmdBytes, 500*time.Millisecond)
	if err := future.Error(); err != nil {
		log.Printf("[%d] ERRO CRÍTICO: Falha ao reivindicar a carta para o RequestID %s: %v", server.ID, requestID, err)
	}
}

func HandleSeeCards(server *models.Server, request shared.Request, nc *nats.Conn, message *nats.Msg){
	server.Mu.Lock()
	cards := shared.Cards {
		Cards : server.Users[request.ClientID].Cards,
	}
	server.Mu.Unlock()
	resp := shared.Response{
		Status: "success",
		Action: "SEE_CARDS",
		Data: utils.MustMarshal(cards),
		Server: server.ID,
	}
	data, _ := json.Marshal(resp)
	nc.Publish(message.Reply, data)
}

func HandleSeeDeck(server *models.Server, request shared.Request, nc *nats.Conn, message *nats.Msg){
	server.Mu.Lock()
	deck := shared.Cards{
		Cards: server.Users[request.ClientID].Deck,
	}
	server.Mu.Unlock()
	resp := shared.Response{
		Status: "success",
		Action: "SEE_DECK",
		Data: utils.MustMarshal(deck),
		Server: server.ID,
	}
	data, _ := json.Marshal(resp)
	nc.Publish(message.Reply, data)
}

func HandleChangeDeck(server *models.Server, request shared.Request, nc *nats.Conn, msg *nats.Msg){
	var deck []shared.Card
	if err := json.Unmarshal(request.Payload, &deck); err != nil{
		log.Printf("[%d] - Erro ao desserializar o deck: %v", server.ID, err)
		resp := shared.Response{
            Status: "error",
            Action: "CHANGE_DECK_FAIL",
            Error:  "payload inválido",
            Server: server.ID,
        }
		data, _ := json.Marshal(resp)
        nc.Publish(msg.Reply, data)
        return
	}

	server.Mu.Lock()
	defer server.Mu.Unlock()

	user := server.Users[request.ClientID]
	user.Deck = deck
	server.Users[request.ClientID] = user

	resp := shared.Response{
            Status: "success",
            Action: "CHANGE_DECK",
            Server: server.ID,
        }
	data, _ := json.Marshal(resp)
	nc.Publish(msg.Reply, data)
	fmt.Printf("%s deck atualizado\n", request.ClientID)
}

// Função auxiliar para enviar respostas de sucesso
func respondWithSuccess(nc *nats.Conn, msg *nats.Msg, card shared.Card) {
	responseData := shared.CardDrawnData{Card: card, RequestID: "client-facing-id"}
	responseBytes, _ := json.Marshal(responseData)
	response := shared.Response{Status: "success", Action: "CARD_DRAWN", Data: responseBytes}
	finalBytes, _ := json.Marshal(response)
	nc.Publish(msg.Reply, finalBytes)
}

// Função auxiliar para enviar respostas de erro
func respondWithError(nc *nats.Conn, msg *nats.Msg, errorMsg string) {
	response := shared.Response{Status: "error", Error: errorMsg}
	data, _ := json.Marshal(response)
	nc.Publish(msg.Reply, data)
}

func triggerBlockchainCriarCarta(server *models.Server, wallet, cardID, element, cardType string) {
    style.PrintAz("entrou na função")
	if server.Blockchain == nil {
		style.PrintVerm("cade a blockchain???")
        return
    }
    
    go func() {
		style.PrintAma("tentando blckchain")
        log.Printf("⛓️ [Blockchain] Iniciando mint para %s...", wallet)
        txHash, err := server.Blockchain.CriarCarta(wallet, cardID, element, cardType)
        if err != nil {
            log.Printf("[Blockchain] Erro: %v", err)
        } else {
            log.Printf("[Blockchain] Sucesso! Tx: %s", txHash)
        }
    }()
}