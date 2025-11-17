package exchange

import (
	"log"
	"time"
	"encoding/json"

	"pbl/client/utils"
	"pbl/server/models"
	"pbl/shared"

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
	
	// Adicionar na fila ANTES de tentar fazer match
	entry.Timestamp = time.Now()
	server.Exchange.LocalQueue = append(server.Exchange.LocalQueue, entry)
	log.Printf("Jogador %s entrou na fila de troca do servidor %d", entry.Player.UserName, server.ID)
	log.Printf("Fila atual tem %d jogadores", len(server.Exchange.LocalQueue))
	
	// Agora libera o lock e tenta fazer match
	server.Exchange.Mutex.Unlock()

	// Tenta fazer match (não precisa de goroutine, já que é rápido)
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

	// Processa todos os pares disponíveis
	for len(server.Exchange.LocalQueue) >= 2 {
		reqPlayer1 := server.Exchange.LocalQueue[0]
		reqPlayer2 := server.Exchange.LocalQueue[1]
		
		session := CreateExchangeSession(reqPlayer1, reqPlayer2, nc, server.ID)

		log.Printf("Sessão de troca criada: %s", session.ID)
		log.Printf("Match: %s (%s) <-> %s (%s)", 
			reqPlayer1.Player.UserName, reqPlayer1.CardOffered.Element,
			reqPlayer2.Player.UserName, reqPlayer2.CardOffered.Element)

		SwapCards(server, session)

		notifyPlayers(session, nc)

		// Remove os dois jogadores da fila
		server.Exchange.LocalQueue = server.Exchange.LocalQueue[2:]
	}
}

func CreateExchangeSession(reqPlayer1, reqPlayer2 shared.ExchangeRequest, nc *nats.Conn, serverID int) *shared.ExchangeSession {
	sessionID := utils.GenerateRoomID(serverID)
	session := &shared.ExchangeSession{
		ID:      sessionID,
		Player1: reqPlayer1.Player,
		Player2: reqPlayer2.Player,
		Card1:   reqPlayer1.CardOffered,
		Card2:   reqPlayer2.CardOffered,
	}
	return session
}

func notifyPlayers(session *shared.ExchangeSession, nc *nats.Conn) {
	// Cria notificação personalizada para o Player1
	notif1 := shared.ExchangeNotification{
		SessionID: session.ID,
		YouSend:   session.Card1,
		YouGet:    session.Card2,
		Partner:   session.Player2.UserName,
	}
	data1, _ := json.Marshal(notif1)

	// Cria notificação personalizada para o Player2
	notif2 := shared.ExchangeNotification{
		SessionID: session.ID,
		YouSend:   session.Card2,
		YouGet:    session.Card1,
		Partner:   session.Player1.UserName,
	}
	data2, _ := json.Marshal(notif2)

	// Envia para cada jogador individualmente
	nc.Publish("exchange.notify."+session.Player1.UserId, data1)
	nc.Publish("exchange.notify."+session.Player2.UserId, data2)

	log.Printf("Notificação enviada para %s e %s", session.Player1.UserName, session.Player2.UserName)
}

func SwapCards(server *models.Server, s *shared.ExchangeSession) {
    server.Mu.Lock()
    defer server.Mu.Unlock()

    // Pega as refs REAIS dos jogadores
    player1 := server.Users[s.Player1.UserId]
    player2 := server.Users[s.Player2.UserId]

    // Remove carta de player1
    for i, card := range player1.Cards {
        if card.Id == s.Card1.Id {
            player1.Cards = append(player1.Cards[:i], player1.Cards[i+1:]...)
            break
        }
    }

    // Remove carta de player2
    for i, c := range player2.Cards {
        if c.Id == s.Card2.Id {
            player2.Cards = append(player2.Cards[:i], player2.Cards[i+1:]...)
            break
        }
    }

    // Adiciona a carta trocada
    player1.Cards = append(player1.Cards, s.Card2)
    player2.Cards = append(player2.Cards, s.Card1)

    // Salva de volta
    server.Users[s.Player1.UserId] = player1
    server.Users[s.Player2.UserId] = player2
}
