package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"pbl/client/utils"
	"pbl/server/models"
	sharedRaft "pbl/server/shared"
	"pbl/shared"
	"pbl/style"

	"github.com/hashicorp/raft"
	"github.com/nats-io/nats.go"
)

func JoinExchangeQueue(server *models.Server, request shared.Request, nc *nats.Conn, msg *nats.Msg) {
	var entry shared.ExchangeRequest
	if err := json.Unmarshal(request.Payload, &entry); err != nil {
		log.Printf("Erro ao desserializar payload do JOIN_EXCHANGE_QUEUE: %v", err)
		resp := shared.Response{
			Status: "error",
			Error:  "Payload inválido",
		}
		data, _ := json.Marshal(resp)
		nc.Publish(msg.Reply, data)
		return
	}

	// Valida se o usuário existe
	server.Mu.Lock()
	user, exists := server.Users[entry.Player.UserId]
	server.Mu.Unlock()

	if !exists {
		resp := shared.Response{
			Status: "error",
			Error:  "Usuário não encontrado no servidor",
		}
		data, _ := json.Marshal(resp)
		nc.Publish(msg.Reply, data)
		return
	}

	// Verifica se já está em alguma fila
	if user.ExchangeStatus != "" {
		resp := shared.Response{
			Status: "error",
			Error:  fmt.Sprintf("Você já está na fila de troca (%s)", user.ExchangeStatus),
		}
		data, _ := json.Marshal(resp)
		nc.Publish(msg.Reply, data)
		return
	}

	// Verifica se o usuário realmente tem a carta oferecida
	hasCard := false
	for _, c := range user.Cards {
		if c.Id == entry.CardOffered.Id {
			hasCard = true
			break
		}
	}

	if !hasCard {
		resp := shared.Response{
			Status: "error",
			Error:  "Usuário não possui a carta oferecida",
		}
		data, _ := json.Marshal(resp)
		nc.Publish(msg.Reply, data)
		return
	}

	// mudar status para na fila local
	server.Mu.Lock()
	user.ExchangeStatus = "LOCALQUEUE"
	server.Users[entry.Player.UserId] = user
	server.Mu.Unlock()

	server.Exchange.Mutex.Lock()
	
	// Verificação extra na fila 
	for _, e := range server.Exchange.LocalQueue {
		if e.Player.UserId == entry.Player.UserId {
			server.Exchange.Mutex.Unlock()
			resp := shared.Response{
				Status: "error",
				Error:  "Você já está na fila de troca",
			}
			data, _ := json.Marshal(resp)
			nc.Publish(msg.Reply, data)
			return
		}
	}
	
	entry.Timestamp = time.Now()
	entry.Player.ExchangeStatus = "LOCALQUEUE"
	server.Exchange.LocalQueue = append(server.Exchange.LocalQueue, entry)
	
	log.Printf("Jogador %s entrou na fila LOCAL (servidor %d)", entry.Player.UserName, server.ID)
	log.Printf("Fila local: %d jogadores", len(server.Exchange.LocalQueue))
	
	server.Exchange.Mutex.Unlock()

	// Tenta match local imediatamente
	MatchLocalExchangeQueue(server, nc)

	resp := shared.Response{
		Status: "success",
		Data:   []byte(`"Cliente adicionado à fila com sucesso"`),
		Server: server.ID,
	}
	data, _ := json.Marshal(resp)
	nc.Publish(msg.Reply, data)
}

func MatchLocalExchangeQueue(server *models.Server, nc *nats.Conn) {
	server.Exchange.Mutex.Lock()
	defer server.Exchange.Mutex.Unlock()

	for len(server.Exchange.LocalQueue) >= 2 {
		reqPlayer1 := server.Exchange.LocalQueue[0]
		reqPlayer2 := server.Exchange.LocalQueue[1]
		
		//Verifica se os jogadores ainda estão válidos
		server.Mu.Lock()
		user1, exists1 := server.Users[reqPlayer1.Player.UserId]
		user2, exists2 := server.Users[reqPlayer2.Player.UserId]
		server.Mu.Unlock()

		// Se algum jogador não existe mais ou já não está em fila, remove e continua
		if !exists1 || user1.ExchangeStatus == "" {
			log.Printf("[MatchLocal] Player1 %s inválido, removendo da fila", reqPlayer1.Player.UserName)
			server.Exchange.LocalQueue = server.Exchange.LocalQueue[1:]
			continue
		}
		if !exists2 || user2.ExchangeStatus == "" {
			log.Printf("[MatchLocal] Player2 %s inválido, removendo da fila", reqPlayer2.Player.UserName)
			server.Exchange.LocalQueue = append(server.Exchange.LocalQueue[:1], server.Exchange.LocalQueue[2:]...)
			continue
		}

		// "PROCESSING" para evitar duplicação
		server.Mu.Lock()
		user1.ExchangeStatus = "PROCESSING"
		user2.ExchangeStatus = "PROCESSING"
		server.Users[reqPlayer1.Player.UserId] = user1
		server.Users[reqPlayer2.Player.UserId] = user2
		server.Mu.Unlock()
		
		session := CreateExchangeSession(reqPlayer1, reqPlayer2, nc, server.ID)

		log.Printf("Match LOCAL: %s (%s) <-> %s (%s)", 
			reqPlayer1.Player.UserName, reqPlayer1.CardOffered.Element,
			reqPlayer2.Player.UserName, reqPlayer2.CardOffered.Element)

		// Troca das cartas
		SwapCards(server, session)

		// Notifica os jogadores
		notifyPlayers(session, nc)

		// Remove os dois jogadores da fila
		server.Exchange.LocalQueue = server.Exchange.LocalQueue[2:]
		
		log.Printf("Fila local após match: %d jogadores", len(server.Exchange.LocalQueue))
	}
}

func CreateExchangeSession(reqPlayer1, reqPlayer2 shared.ExchangeRequest, nc *nats.Conn, serverID int) *shared.ExchangeSession {
	sessionID := utils.GenerateRoomID(serverID)
	
	player1 := reqPlayer1.Player
	player2 := reqPlayer2.Player
	
	session := &shared.ExchangeSession{
		ID:      sessionID,
		Player1: &player1,
		Player2: &player2,
		Card1:   reqPlayer1.CardOffered,
		Card2:   reqPlayer2.CardOffered,
	}
	return session
}

func notifyPlayers(session *shared.ExchangeSession, nc *nats.Conn) {
	notif1 := shared.ExchangeNotification{
		SessionID: session.ID,
		YouSend:   session.Card1,
		YouGet:    session.Card2,
		Partner:   session.Player2.UserName,
	}
	data1, _ := json.Marshal(notif1)

	notif2 := shared.ExchangeNotification{
		SessionID: session.ID,
		YouSend:   session.Card2,
		YouGet:    session.Card1,
		Partner:   session.Player1.UserName,
	}
	data2, _ := json.Marshal(notif2)

	nc.Publish("exchange.notify." + session.Player1.UserId, data1)
	nc.Publish("exchange.notify." + session.Player2.UserId, data2)

	log.Printf("Notificações enviadas para %s e %s", session.Player1.UserName, session.Player2.UserName)
}

func SwapCards(server *models.Server, s *shared.ExchangeSession) {
	server.Mu.Lock()
	defer server.Mu.Unlock()

	if s.Player1 == nil || s.Player2 == nil {
		log.Printf("[SwapCards] ERRO: Player1 ou Player2 está nil!")
		return
	}

	log.Printf("[SwapCards] Iniciando troca: %s <-> %s", s.Player1.UserName, s.Player2.UserName)

	// Player1
	player1, hasPlayer1 := server.Users[s.Player1.UserId]
	if hasPlayer1 {
		// faz troca se ainda está PROCESSING
		if player1.ExchangeStatus != "PROCESSING" && player1.ExchangeStatus != "GLOBALQUEUE" {
			log.Printf("[SwapCards] AVISO: Player1 %s não está em estado válido (status: %s)", 
				player1.UserName, player1.ExchangeStatus)
			return
		}

		cardRemoved := false
		for i, card := range player1.Cards {
			if card.Id == s.Card1.Id {
				player1.Cards = append(player1.Cards[:i], player1.Cards[i+1:]...)
				cardRemoved = true
				log.Printf("[SwapCards] Carta %s removida de %s", s.Card1.Element, player1.UserName)
				break
			}
		}
		
		if !cardRemoved {
			log.Printf("[SwapCards] AVISO: Carta %s não encontrada em %s", s.Card1.Element, player1.UserName)
		}else{
			if player1.PrivateKey != "" && s.Player2.Address != "" {
				triggerBlockchainTransfer(server, player1.PrivateKey, s.Card1.Id, s.Player2.Address)
			} else {
				log.Printf("⚠️ [Blockchain] Não foi possível transferir P1->P2 (falta chave ou endereço).")
			}
		}

		player1.Cards = append(player1.Cards, s.Card2)
		player1.ExchangeStatus = "" // reseta status após troca
		server.Users[s.Player1.UserId] = player1
		
		log.Printf("[SwapCards] %s recebeu %s (%d cartas)", 
			player1.UserName, s.Card2.Element, len(player1.Cards))
	} else {
		log.Printf("[SwapCards] Player1 (%s) não está neste servidor", s.Player1.UserName)
	}

	// Player2
	player2, hasPlayer2 := server.Users[s.Player2.UserId]
	if hasPlayer2 {
		// Só faz troca se ainda está PROCESSING
		if player2.ExchangeStatus != "PROCESSING" && player2.ExchangeStatus != "GLOBALQUEUE" {
			log.Printf("[SwapCards] AVISO: Player2 %s não está em estado válido (status: %s)", 
				player2.UserName, player2.ExchangeStatus)
			return
		}

		cardRemoved := false
		for i, card := range player2.Cards {
			if card.Id == s.Card2.Id {
				player2.Cards = append(player2.Cards[:i], player2.Cards[i+1:]...)
				cardRemoved = true
				log.Printf("[SwapCards] Carta %s removida de %s", s.Card2.Element, player2.UserName)
				break
			}
		}
		
		if !cardRemoved {
			log.Printf("[SwapCards] AVISO: Carta %s não encontrada em %s", s.Card2.Element, player2.UserName)
		}else{
			if player2.PrivateKey != "" && s.Player1.Address != "" {
				triggerBlockchainTransfer(server, player2.PrivateKey, s.Card2.Id, s.Player1.Address)
			} else {
				log.Printf("⚠️ [Blockchain] Não foi possível transferir P2->P1 (falta chave ou endereço).")
			}
		}

		player2.Cards = append(player2.Cards, s.Card1)
		player2.ExchangeStatus = "" // reseta status após troca
		server.Users[s.Player2.UserId] = player2
		
		log.Printf("[SwapCards] %s recebeu %s (%d cartas)", 
			player2.UserName, s.Card1.Element, len(player2.Cards))
	} else {
		log.Printf("[SwapCards] Player2 (%s) não está neste servidor", s.Player2.UserName)
	}

	if hasPlayer1 && hasPlayer2 {
		log.Printf("[SwapCards] ✓ Troca LOCAL completa")
	}
}

func MonitorExchangeLocalQueue(server *models.Server, nc *nats.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	log.Printf("[Monitor Exchange] Iniciado para servidor %d", server.ID)

	for range ticker.C {
		now := time.Now()
		server.Exchange.Mutex.Lock()

		i := 0
		for i < len(server.Exchange.LocalQueue) {
			entry := server.Exchange.LocalQueue[i]
			waitTime := now.Sub(entry.Timestamp)

			// Verifica se usuário ainda existe e tem status válido
			server.Mu.Lock()
			user, exists := server.Users[entry.Player.UserId]
			server.Mu.Unlock()

			// Se usuário não existe mais ou status foi resetado (troca já feita), remove da fila
			if !exists || user.ExchangeStatus == "" {
				log.Printf("[Monitor] Removendo %s da fila (troca já processada)", 
					entry.Player.UserName)
				server.Exchange.LocalQueue = append(
					server.Exchange.LocalQueue[:i],
					server.Exchange.LocalQueue[i+1:]...,
				)
				continue
			}

			// Se está PROCESSING, aguarda (troca em andamento)
			if user.ExchangeStatus == "PROCESSING" {
				log.Printf("[Monitor] %s está sendo processado, aguardando...", entry.Player.UserName)
				i++
				continue
			}

			// Tenta match local se houver 2+ jogadores
			if len(server.Exchange.LocalQueue) >= 2 {
				log.Printf("[Monitor] Tentando match local (%d jogadores na fila)", 
					len(server.Exchange.LocalQueue))
				server.Exchange.Mutex.Unlock()
				MatchLocalExchangeQueue(server, nc)
				server.Exchange.Mutex.Lock()
				continue
			}

			// Após 5 segundos, envia para fila global
			if waitTime >= 5*time.Second {
				style.PrintVerm("[Monitor] ENVIANDO PARA FILA GLOBAL")
				log.Printf("[Monitor] Jogador: %s (esperou %.1fs)", 
					entry.Player.UserName, waitTime.Seconds())

				exchangeQueueEntry := shared.ExchangeQueueEntry{
					Player:    entry.Player,
					Card:      entry.CardOffered,
					Timestamp: entry.Timestamp,
					ServerID:  server.ID,
				}

				// Atualiza status para GLOBALQUEUE
				server.Mu.Lock()
				user.ExchangeStatus = "GLOBALQUEUE"
				server.Users[entry.Player.UserId] = user
				server.Mu.Unlock()

				// Remove da fila local
				log.Printf("[Monitor] Removendo %s da fila LOCAL", entry.Player.UserName)
				server.Exchange.LocalQueue = append(
					server.Exchange.LocalQueue[:i],
					server.Exchange.LocalQueue[i+1:]...,
				)
				
				server.Exchange.Mutex.Unlock()
				
				// Envia para global
				SendToExchangeGlobalQueue(exchangeQueueEntry, server)
				
				server.Exchange.Mutex.Lock()
				continue
			}
			i++
		}

		server.Exchange.Mutex.Unlock()
	}
}

func SendToExchangeGlobalQueue(entry shared.ExchangeQueueEntry, server *models.Server) {
	log.Printf("[SendToGlobal] Enviando %s (carta: %s, serverID: %d)", 
		entry.Player.UserName, entry.Card.Element, entry.ServerID)
	
	isLeader := server.Raft != nil && server.Raft.State() == raft.Leader
	
	if isLeader {
		log.Printf("[SendToGlobal] Este servidor É O LÍDER")
		
		cmdData, _ := json.Marshal(sharedRaft.Command{
			Type: sharedRaft.CommandQueueJoinGlobalExchange,
			Data: utils.MustMarshal(entry),
		})

		future := server.Raft.Apply(cmdData, 5*time.Second)
		if err := future.Error(); err != nil {
			log.Printf("[SendToGlobal] ERRO ao aplicar Raft: %v", err)
			return
		}

		log.Printf("[SendToGlobal] ✓ Jogador adicionado à fila global")
		
		time.Sleep(150 * time.Millisecond)
		
		createdSessions := server.FSM.TryMatchExchangeGlobal()
		if len(createdSessions) > 0 {
			log.Printf("[SendToGlobal] ✓ %d sessões criadas", len(createdSessions))
			for _, session := range createdSessions {
				NotifyServersAboutExchangeMatch(session, server)
			}
		}
		return
	}

	// Não é líder - envia via HTTP
	log.Printf("[SendToGlobal] NÃO é líder, enviando via HTTP")
	
	leaderAddr := string(server.Raft.Leader())
	if leaderAddr == "" {
		log.Printf("[SendToGlobal] ERRO: Nenhum líder disponível")
		return
	}

	url := fmt.Sprintf("http://%s/exchange/join-global", leaderAddr)
	payload := utils.MustMarshal(entry)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("[SendToGlobal] ERRO ao enviar: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("[SendToGlobal] ✓ Enviado ao líder (%s)", leaderAddr)
	} else {
		log.Printf("[SendToGlobal] Resposta do líder: %d", resp.StatusCode)
	}
}

func HandleJoinGlobalExchange(server *models.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var entry shared.ExchangeQueueEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "Payload inválido", http.StatusBadRequest)
			return
		}

		if server.Raft.State() != raft.Leader {
			http.Error(w, "Apenas o líder aceita join-global", http.StatusForbidden)
			return
		}

		SendToExchangeGlobalQueue(entry, server)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func NotifyServersAboutExchangeMatch(session *shared.ExchangeSession, server *models.Server) {
	log.Printf("[Exchange] Match: %s (S%d) <-> %s (S%d)", 
		session.Player1.UserName, session.Server1ID,
		session.Player2.UserName, session.Server2ID)

	if session.Server1ID == server.ID {
		SwapCards(server, session)
		
		notif1 := shared.ExchangeNotification{
			SessionID: session.ID,
			YouSend:   session.Card1,
			YouGet:    session.Card2,
			Partner:   session.Player2.UserName,
		}
		data1, _ := json.Marshal(notif1)
		server.Matchmaking.Nc.Publish("exchange.notify."+session.Player1.UserId, data1)
	}

	if session.Server2ID == server.ID {
		SwapCards(server, session)
		
		notif2 := shared.ExchangeNotification{
			SessionID: session.ID,
			YouSend:   session.Card2,
			YouGet:    session.Card1,
			Partner:   session.Player1.UserName,
		}
		data2, _ := json.Marshal(notif2)
		server.Matchmaking.Nc.Publish("exchange.notify."+session.Player2.UserId, data2)
	}

	if session.Server1ID != server.ID {
		go notifyServerAboutExchange(session, session.Server1ID, server)
	}

	if session.Server2ID != server.ID {
		go notifyServerAboutExchange(session, session.Server2ID, server)
	}
}

func notifyServerAboutExchange(session *shared.ExchangeSession, targetServerID int, server *models.Server) {
	serverAddr := getServerAddress(targetServerID, server)
	if serverAddr == "" {
		return
	}

	url := fmt.Sprintf("http://%s/exchange/execute-swap", serverAddr)
	payload := utils.MustMarshal(session)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("[Exchange] Erro ao notificar servidor %d: %v", targetServerID, err)
		return
	}
	defer resp.Body.Close()
}

func getServerAddress(serverID int, server *models.Server) string {
	future := server.Raft.GetConfiguration()
	if err := future.Error(); err != nil {
		return ""
	}

	config := future.Configuration()
	expectedID := fmt.Sprintf("%d", serverID)
	
	for _, srv := range config.Servers {
		if string(srv.ID) == expectedID {
			return string(srv.Address)
		}
	}
	return ""
}

func HandleExecuteSwap(server *models.Server, nc *nats.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var session shared.ExchangeSession
		if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
			http.Error(w, "Payload inválido", http.StatusBadRequest)
			return
		}

		SwapCards(server, &session)

		server.Mu.Lock()
		_, hasPlayer1 := server.Users[session.Player1.UserId]
		_, hasPlayer2 := server.Users[session.Player2.UserId]
		server.Mu.Unlock()

		if hasPlayer1 {
			notif1 := shared.ExchangeNotification{
				SessionID: session.ID,
				YouSend:   session.Card1,
				YouGet:    session.Card2,
				Partner:   session.Player2.UserName,
			}
			data1, _ := json.Marshal(notif1)
			nc.Publish("exchange.notify."+session.Player1.UserId, data1)
		}

		if hasPlayer2 {
			notif2 := shared.ExchangeNotification{
				SessionID: session.ID,
				YouSend:   session.Card2,
				YouGet:    session.Card1,
				Partner:   session.Player1.UserName,
			}
			data2, _ := json.Marshal(notif2)
			nc.Publish("exchange.notify."+session.Player2.UserId, data2)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Swap executado"))
	}
}

func triggerBlockchainTransfer(server *models.Server, senderPrivateKey, cardID, targetAddress string) {
	if server.Blockchain == nil {
		return
	}
	go func() {
		log.Printf("🔄 [Blockchain] Transferindo carta %s para %s...", cardID, targetAddress)
		
		txHash, err := server.Blockchain.TransferirCarta(senderPrivateKey, cardID, targetAddress)
		
		if err != nil {
			log.Printf("❌ [Blockchain] Erro na troca: %v", err)
		} else {
			log.Printf("✅ [Blockchain] Carta transferida! Tx: %s", txHash)
		}
	}()
}
