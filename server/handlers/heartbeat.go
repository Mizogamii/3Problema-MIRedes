package handlers

import(
	"log"
	"fmt"
	"time"
	"encoding/json"

	"pbl/shared"
	"pbl/server/models"
	"github.com/nats-io/nats.go"
)


const (
	heartbeatInterval = 10 * time.Second
	disconnectTimeout = 30 * time.Second
)

//Heartbeat
func StartHeartbeatMonitor(server *models.Server, nc *nats.Conn) {
	go func() {
		for {
			time.Sleep(heartbeatInterval)
			now := time.Now()

			mu.Lock()
			for id, c := range activeClients {
				// Se passou muito tempo desde o último pong
				if now.Sub(c.LastSeen) > disconnectTimeout {
					log.Printf("Cliente '%s' inativo por %v. Desconectando...", id, now.Sub(c.LastSeen))
					DisconnectClient(server, id)
					delete(activeClients, id)
					continue
				}

				// Se estiver em estado ativo, envia ping
				if c.State == Active {
					c.State = WaitingReconnection
					nc.Publish(fmt.Sprintf("client.%s.ping", id), []byte("ping"))
				}
			}
			mu.Unlock()
		}
	}()
}


func HandlePing(server *models.Server, req shared.Request, nc *nats.Conn, msg *nats.Msg) {
	clientID := req.ClientID
	// Atualiza último ping
	mu.Lock()
	if c, ok := activeClients[clientID]; ok {
		c.LastSeen = time.Now()
		c.State = Active
	}
	mu.Unlock()

	// Responde com PONG
	response := shared.Response{
		Status: "success",
		Action: "PONG",
	}
	respData, _ := json.Marshal(response)
	clientTopic := fmt.Sprintf("client.%s.inbox", clientID)
	nc.Publish(clientTopic, respData)
}

func HandleHeartbeat(serverID int, request shared.Request, nc *nats.Conn, msg *nats.Msg) {
	clientID := request.ClientID

	mu.Lock()
	if c, ok := activeClients[clientID]; ok {
		c.LastSeen = time.Now()
		c.State = Active
	} else {
		activeClients[clientID] = &ClientInfo{
			ClientID: clientID,
			LastSeen: time.Now(),
			State:    Active,
		}
	}
	mu.Unlock()
}
