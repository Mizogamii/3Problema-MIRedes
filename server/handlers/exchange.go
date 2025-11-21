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
	
	entry.Timestamp = time.Now()
	server.Exchange.LocalQueue = append(server.Exchange.LocalQueue, entry)
	log.Printf("Jogador %s entrou na fila de troca do servidor %d", entry.Player.UserName, server.ID)
	log.Printf("Fila local atual tem %d jogadores", len(server.Exchange.LocalQueue))
	
	server.Exchange.Mutex.Unlock()

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
		
		session := CreateExchangeSession(reqPlayer1, reqPlayer2, nc, server.ID)

		log.Printf("Sessão de troca criada: %s", session.ID)
		log.Printf("Match LOCAL: %s (%s) <-> %s (%s)", 
			reqPlayer1.Player.UserName, reqPlayer1.CardOffered.Element,
			reqPlayer2.Player.UserName, reqPlayer2.CardOffered.Element)

		// Troca das cartas
		SwapCards(server, session)

		// Notifica os jogadores
		notifyPlayers(session, nc)

		// Remove os dois jogadores da fila
		server.Exchange.LocalQueue = server.Exchange.LocalQueue[2:]
	}
}

func CreateExchangeSession(reqPlayer1, reqPlayer2 shared.ExchangeRequest, nc *nats.Conn, serverID int) *shared.ExchangeSession {
	sessionID := utils.GenerateRoomID(serverID)
	
	// Cria ponteiros para os players
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
	// Notificação para o Player1
	notif1 := shared.ExchangeNotification{
		SessionID: session.ID,
		YouSend:   session.Card1,
		YouGet:    session.Card2,
		Partner:   session.Player2.UserName,
	}
	data1, _ := json.Marshal(notif1)

	// Notificação para o Player2
	notif2 := shared.ExchangeNotification{
		SessionID: session.ID,
		YouSend:   session.Card2,
		YouGet:    session.Card1,
		Partner:   session.Player1.UserName,
	}
	data2, _ := json.Marshal(notif2)

	nc.Publish("exchange.notify." + session.Player1.UserId, data1)
	nc.Publish("exchange.notify." + session.Player2.UserId, data2)

	log.Printf("Notificação enviada para %s e %s", session.Player1.UserName, session.Player2.UserName)
}

func SwapCards(server *models.Server, s *shared.ExchangeSession) {
	server.Mu.Lock()
	defer server.Mu.Unlock()

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

	server.Users[s.Player1.UserId] = player1
	server.Users[s.Player2.UserId] = player2

	log.Printf("[SwapCards] Troca realizada: %s recebeu %s, %s recebeu %s",
		player1.UserName, s.Card2.Element,
		player2.UserName, s.Card1.Element)
}

func MonitorExchangeLocalQueue(server *models.Server, nc *nats.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	
	log.Printf("[Monitor Exchange] Iniciado para servidor %d", server.ID)

	for range ticker.C {
		now := time.Now()

		server.Exchange.Mutex.Lock()

		// Debug: mostra o estado da fila a cada tick
		if len(server.Exchange.LocalQueue) > 0 {
			log.Printf("[Monitor] Fila LOCAL tem %d jogadores", len(server.Exchange.LocalQueue))
			for idx, e := range server.Exchange.LocalQueue {
				waitTime := now.Sub(e.Timestamp)
				log.Printf("[Monitor]   [%d] %s (carta: %s) - esperando há %.1fs", 
					idx, e.Player.UserName, e.CardOffered.Element, waitTime.Seconds())
			}
		}

		i := 0
		for i < len(server.Exchange.LocalQueue) {
			entry := server.Exchange.LocalQueue[i]
			waitTime := now.Sub(entry.Timestamp)

			log.Printf("[Monitor] Verificando jogador %s - tempo de espera: %.1fs", 
				entry.Player.UserName, waitTime.Seconds())

			// Match local antes de enviar pra global
			if len(server.Exchange.LocalQueue) >= 2 && i < len(server.Exchange.LocalQueue)-1 {
				log.Printf("[Monitor] Tentando match LOCAL (fila tem %d jogadores)", 
					len(server.Exchange.LocalQueue))
				server.Exchange.Mutex.Unlock()
				MatchLocalExchangeQueue(server, nc)
				server.Exchange.Mutex.Lock()
				continue
			}

			// Passou 5 segundos --> envia para fila global
			if waitTime >= 5*time.Second {
				style.PrintVerm("ENVIANDO PARA FILA GLOBAL")
				log.Printf("Jogador: %s", entry.Player.UserName)
				log.Printf("UserID: %s", entry.Player.UserId)
				log.Printf("Carta: %s ║", entry.CardOffered.Element)
				log.Printf("Tempo de espera: %.1f segundos", waitTime.Seconds())
	
				exchangeQueueEntry := shared.ExchangeQueueEntry{
					Player:    entry.Player,
					Card:      entry.CardOffered,
					Timestamp: entry.Timestamp,
				}

				server.Exchange.Mutex.Unlock()
				
				// Envia para global
				SendToExchangeGlobalQueue(exchangeQueueEntry, server)
				
				// Verifica fila global após envio
				globalUsers := ListGlobalExchangeQueue(server)
				log.Printf("[Monitor] Fila GLOBAL após envio: %v (total: %d)", 
					globalUsers, len(globalUsers))
				
				server.Exchange.Mutex.Lock()

				// Remove da fila local
				log.Printf("[Monitor] Removendo %s da fila LOCAL", entry.Player.UserName)
				server.Exchange.LocalQueue = append(
					server.Exchange.LocalQueue[:i],
					server.Exchange.LocalQueue[i+1:]...,
				)
				log.Printf("[Monitor] Fila LOCAL agora tem %d jogadores", 
					len(server.Exchange.LocalQueue))

				// Não incrementa i pois removemos o elemento atual
				continue
			}

			// Se não passou 5 segundos e não há match, passa pro próximo
			i++
		}

		server.Exchange.Mutex.Unlock()
	}
}

func SendToExchangeGlobalQueue(entry shared.ExchangeQueueEntry, server *models.Server) {
	log.Printf("[SendToGlobal] Iniciando envio para fila global")
	log.Printf("[SendToGlobal] Jogador: %s, Carta: %s", entry.Player.UserName, entry.Card.Element)
	log.Printf("[SendToGlobal] É líder? %v", server.Exchange.IsLeader)
	isLeader := server.Raft != nil && server.Raft.State() == raft.Leader
	// Se este servidor é o líder
	if isLeader {
		log.Printf("[SendToGlobal] Este servidor É O LÍDER - adicionando via Raft")
		
		cmdData, _ := json.Marshal(sharedRaft.Command{
			Type: sharedRaft.CommandQueueJoinGlobalExchange,
			Data: utils.MustMarshal(entry),
		})

		log.Printf("[SendToGlobal] Aplicando comando no Raft...")
		future := server.Raft.Apply(cmdData, 5*time.Second)
		if err := future.Error(); err != nil {
			log.Printf("[Exchange Líder] - ERRO ao replicar entrada da fila global via Raft: %v", err)
			return
		}

		log.Printf("[Exchange Líder] - Jogador %s adicionado à fila GLOBAL de trocas", entry.Player.UserName)
		
		// Mostra fila global atual
		globalUsers := ListGlobalExchangeQueue(server)
		log.Printf("[Exchange Líder] - Fila GLOBAL atual: %v (total: %d)", 
			globalUsers, len(globalUsers))

		// Tenta fazer matches após um delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			
			log.Printf("[Exchange Líder] Tentando fazer matches globais...")
			createdSessions := server.FSM.TryMatchExchangeGlobal()

			if len(createdSessions) > 0 {
				log.Printf("[Exchange Líder] - %d sessões de troca criadas!", len(createdSessions))
				// Notifica os servidores sobre as trocas criadas
				for _, session := range createdSessions {
					NotifyServersAboutExchangeMatch(session, server)
				}
			} else {
				log.Printf("[Exchange Líder] - Nenhum match encontrado ainda (precisa de 2+ jogadores)")
			}
		}()

		return
	}

	// Se não é líder, envia para o líder
	log.Printf("[SendToGlobal] Este servidor NÃO é líder - enviando para o líder via HTTP")
	
	leaderAddr := string(server.Raft.Leader())
	if leaderAddr == "" {
		log.Printf("[Exchange Follower] - Nenhum líder disponível, não foi possível enviar %s", entry.Player.UserName)
		return
	}

	log.Printf("[SendToGlobal] Líder encontrado: %s", leaderAddr)
	url := fmt.Sprintf("http://%s/exchange/join-global", leaderAddr)
	log.Printf("[SendToGlobal] URL: %s", url)

	payload := utils.MustMarshal(entry)
	log.Printf("[SendToGlobal] Enviando requisição HTTP POST...")
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("[Exchange Follower] - Erro ao enviar %s para líder: %v", entry.Player.UserName, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Exchange Follower] - Resposta inválida do líder: %d", resp.StatusCode)
	} else {
		log.Printf("[Exchange Follower] - Jogador %s enviado ao líder (%s)", entry.Player.UserName, leaderAddr)
	}
}

// ListGlobalExchangeQueue retorna a lista de usuários na fila GLOBAL
func ListGlobalExchangeQueue(server *models.Server) []string {
	server.FSM.GlobalExchangeQueueMu.Lock()
	defer server.FSM.GlobalExchangeQueueMu.Unlock()

	users := make([]string, len(server.FSM.GlobalExchangeQueue))
	for i, entry := range server.FSM.GlobalExchangeQueue {
		users[i] = fmt.Sprintf("%s[%s]", entry.Player.UserName, entry.Card.Element)
	}
	
	if len(users) > 0 {
		log.Printf("[ListGlobalQueue] - Fila global tem %d jogadores: %v", len(users), users)
	} else {
		log.Printf("[ListGlobalQueue] - Fila global está VAZIA")
	}
	
	return users
}


func NotifyServersAboutExchangeMatch(session *shared.ExchangeSession, server *models.Server) {
	log.Printf("[Exchange] Notificando servidores sobre sessão %s", session.ID)

	// Notifica o servidor do Player1
	if session.Server1ID != server.ID {
		go notifyServerAboutExchange(session, session.Server1ID, server)
	}

	// Notifica o servidor do Player2
	if session.Server2ID != server.ID {
		go notifyServerAboutExchange(session, session.Server2ID, server)
	}

	log.Printf("[Exchange] Sessão global criada: %s (%sserver%d <-> %sserver%d)",
		session.ID,
		session.Player1.UserName, session.Server1ID,
		session.Player2.UserName, session.Server2ID)
}

// notifyServerAboutExchange envia a notificação para um servidor específico
func notifyServerAboutExchange(session *shared.ExchangeSession, targetServerID int, server *models.Server) {
	// Busca o endereço do servidor na configuração do Raft
	serverAddr := getServerAddress(targetServerID, server)
	if serverAddr == "" {
		log.Printf("[Exchange] Não foi possível encontrar endereço do servidor %d", targetServerID)
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("[Exchange] Servidor %d retornou status %d", targetServerID, resp.StatusCode)
	} else {
		log.Printf("[Exchange] Servidor %d notificado com sucesso sobre %s", targetServerID, session.ID)
	}
}

// Busca o endereço de um servidor pelo ID
func getServerAddress(serverID int, server *models.Server) string {
	// Obtém a configuração atual do Raft
	future := server.Raft.GetConfiguration()
	if err := future.Error(); err != nil {
		log.Printf("[Exchange] Erro ao obter configuração do Raft: %v", err)
		return ""
	}

	config := future.Configuration()
	for _, srv := range config.Servers {
		expectedID := fmt.Sprintf("%d", serverID)
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

        // Executa a troca localmente
        SwapCards(server, &session)

        // Notifica os usuários envolvidos
        notifyPlayers(&session, nc)

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Swap executado"))
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

        // Apenas o líder pode adicionar globalmente via HTTP
        if server.Raft.State() != raft.Leader {
            http.Error(w, "Apenas o líder aceita join-global", http.StatusForbidden)
            return
        }

        SendToExchangeGlobalQueue(entry, server)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    }
}
