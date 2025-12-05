package handlers

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"pbl/server/models"
	"pbl/server/utils"
	"pbl/shared"

	"github.com/nats-io/nats.go"
)

func HandleChooseServer(server *models.Server, request shared.Request, nc *nats.Conn, message *nats.Msg) {
	// Pega o server_id do payload do cliente (mesmo que seja esse servidor)
	var payloadData map[string]int
	if err := json.Unmarshal(request.Payload, &payloadData); err != nil {
		log.Printf("[%d] - Erro ao decodificar payload: %v", server.ID, err)
		return
	}

	chosenServerID := payloadData["server_id"]
	log.Printf("[%d] - Cliente %s escolheu este servidor (ID=%d)", server.ID, request.ClientID, chosenServerID)

	// Resposta para o cliente confirmando que ele escolheu o servidor
	response := shared.Response{
		Status: "success",
		Action: "CHOOSE_SERVER",
		Server: server.ID,
	}
	data, _ := json.Marshal(response)

	if message.Reply != "" {
		nc.Publish(message.Reply, data)
	}
}

func HandleLogin(server *models.Server, request shared.Request, nc *nats.Conn, msg *nats.Msg) {
    var user shared.User
	
    if err := json.Unmarshal(request.Payload, &user); err != nil {
        log.Printf("[%d] - Erro ao desserializar login: %v", server.ID, err)
        resp := shared.Response{
            Status: "error",
            Action: "LOGIN_FAIL",
            Error:  "payload inválido",
            Server: server.ID,
        }
        data, _ := json.Marshal(resp)
        nc.Publish(msg.Reply, data)
        return
    }

    server.Mu.Lock()
    defer server.Mu.Unlock()

    //Verifica se já existe um usuário com o mesmo nome online
    for _, existingUser := range server.Users {
        if existingUser.UserName == user.UserName {
            log.Printf("[%d] - Tentativa de login duplicado para '%s'", server.ID, user.UserName)
            resp := shared.Response{
                Status: "error",
                Action: "LOGIN_FAIL",
                Error:  "Usuário já está logado em outro cliente.",
                Server: server.ID,
            }
            data, _ := json.Marshal(resp)
            nc.Publish(msg.Reply, data)
            return
        }
    }

    //Insere cartas padrão 
    user.Cards = []shared.Card{
        {Id: "1", Element: "AGUA", Type: "NORMAL"},
        {Id: "2", Element: "TERRA", Type: "NORMAL"},
        {Id: "3", Element: "FOGO", Type: "NORMAL"},
        {Id: "4", Element: "AR", Type: "NORMAL"},
        {Id: "5", Element: "MATO", Type: "NORMAL"},
    }

    user.Deck = []shared.Card{
        {Id: "1", Element: "AGUA", Type: "NORMAL"},
        {Id: "2", Element: "TERRA", Type: "NORMAL"},
        {Id: "3", Element: "FOGO", Type: "NORMAL"},
        {Id: "4", Element: "AR", Type: "NORMAL"},
    }

	user.ServerID = server.ID
	address,key := utils.GenerateNewWallet()

	user.Address = address
	user.PrivateKey = key

	if server.Blockchain != nil {
        // Fazemos em goroutine para o login ser instantâneo
        go func(destAddr string) {
            log.Printf("[Faucet] Enviando ETH para %s...", destAddr)
            tx, err := server.Blockchain.EnviarEther(destAddr)
            if err != nil {
                log.Printf("❌ [Faucet Erro] %v", err)
            } else {
                log.Printf("✅ [Faucet Sucesso] ETH enviado! Tx: %s", tx)
            }
        }(address)
    }

    //Armazena o usuário logado
    server.Users[request.ClientID] = user
    log.Printf("[%d] - Usuário '%s' conectado com ClientID '%s'", server.ID, user.UserName, request.ClientID)

    resp := shared.Response{
        Status: "success",
        Action: "LOGIN_SUCCESS",
        Data:   utils.MustMarshal(user),
        Server: server.ID,
    }
    data, _ := json.Marshal(resp)
    nc.Publish(msg.Reply, data)
}

func HandleLogout(server *models.Server, request shared.Request, nc *nats.Conn, msg *nats.Msg) {
	server.Mu.Lock()
	user, exists := server.Users[request.ClientID]
	if exists {
		log.Printf("[%d] - Cliente '%s' desconectado (ClientID: %s)", server.ID, user.UserName, request.ClientID)
		delete(server.Users, request.ClientID)
	} else {
		log.Printf("[%d] - Cliente com ClientID '%s' desconectado (usuário não encontrado)", server.ID, request.ClientID)
	}
	server.Mu.Unlock()

	mu.Lock()
	delete(activeClients, request.ClientID)
	mu.Unlock()

	DisconnectClient(server, request.ClientID)

	// Resposta para o cliente
	resp := shared.Response{
		Status: "success",
		Action: "LOGOUT_SUCCESS",
		Server: server.ID,
	}
	data, _ := json.Marshal(resp)
	nc.Publish(msg.Reply, data)
}

const (
	Active ClientState = iota
	WaitingReconnection
)

type ClientState int
type ClientInfo struct {
	ClientID string
	LastSeen time.Time
	State ClientState
}

var (
	activeClients = make(map[string]*ClientInfo)
	mu            = sync.Mutex{}
)

// Função que trata a desconexão de forma genérica
func DisconnectClient(server *models.Server, clientID string) {
	server.Mu.Lock()
	delete(server.Users, clientID)
	server.Mu.Unlock()

	mu.Lock()
	delete(activeClients, clientID)
	mu.Unlock()

	log.Printf("Cliente '%s' caiu ou ficou inativo. Removido do servidor.", clientID)
}

