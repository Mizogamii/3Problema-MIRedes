package exchange

import (
	"encoding/json"
	"log"
	"time"

	"pbl/server/models"
	"pbl/shared"

	//"github.com/hashicorp/raft"
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

	// Verifica se o usuário já está na fila da troca para evitar duplicação
	server.Exchange.Mutex.Lock()
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
	server.Exchange.Mutex.Unlock()

	// Adicionar na fila
	entry.Timestamp = time.Now()

	server.Exchange.Mutex.Lock()
	server.Exchange.LocalQueue = append(server.Exchange.LocalQueue, entry)
	log.Println("Fila atual:", server.Exchange.LocalQueue)
	server.Exchange.Mutex.Unlock()

	log.Printf("Jogador %s entrou na fila de troca do servidor %d", entry.Player.UserName, server.ID)

	resp := shared.Response{
		Status: "success",
		Data:   []byte(`"Cliente adicionado à fila com sucesso"`),
		Server: server.ID,
	}
	data, _ := json.Marshal(resp)
	nc.Publish(msg.Reply, data)
}
