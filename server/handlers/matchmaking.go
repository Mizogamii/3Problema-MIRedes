package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"pbl/server/game"
	"pbl/server/models"
	sharedRaft "pbl/server/shared"
	"pbl/server/utils"
	"pbl/shared"
	"pbl/style"

	"github.com/hashicorp/raft"
	"github.com/nats-io/nats.go"
)

// Colocar cliente na fila --> fila local
func HandleJoinQueue(server *models.Server, request shared.Request, nc *nats.Conn, msg *nats.Msg) {
	var entry shared.QueueEntry
	if err := json.Unmarshal(request.Payload, &entry); err != nil {
		log.Printf("Erro ao desserializar payload do JOIN_QUEUE: %v", err)
		resp := shared.Response{
			Status: "error",
			Error:  "Payload inválido",
		}
		data, _ := json.Marshal(resp)
		nc.Publish(msg.Reply, data)
		return
	}

	entry.JoinTime = time.Now() //sobreescreve o horário que o do cliente pode estar com erro

	//Adiciona na fila local
	server.Matchmaking.Mutex.Lock()
	server.Matchmaking.LocalQueue = append(server.Matchmaking.LocalQueue, entry)
	server.Matchmaking.Mutex.Unlock()

	log.Printf("Cliente %s entrou na fila local do servidor %d", entry.Player.UserName, server.ID)

	log.Printf("[DEBUG] Fila local atual: %v", ListLocalQueue(server))

	//Responde para o cliente
	msgData, _ := json.Marshal("Cliente adicionado à fila com sucesso")
	resp := shared.Response{
		Status: "success",
		Data:   msgData,
		Server: server.ID,
	}
	data, _ := json.Marshal(resp)
	nc.Publish(msg.Reply, data)
}

// Monitora a fila local e move jogadores para a fila global se passarem de 5s
func MonitorLocalQueue(server *models.Server, nc *nats.Conn) {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		now := time.Now()

		server.Matchmaking.Mutex.Lock()

		//Move jogadores da fila local para global após 10 segundos de espera
		for i := 0; i < len(server.Matchmaking.LocalQueue); {
			entry := server.Matchmaking.LocalQueue[i]
			waitTime := now.Sub(entry.JoinTime)

			if waitTime > 5*time.Second {
				SendToGlobalQueue(entry, server)
				server.Matchmaking.LocalQueue = append(
					server.Matchmaking.LocalQueue[:i],
					server.Matchmaking.LocalQueue[i+1:]...,
				)
			} else {
				i++
			}
		}

		server.Matchmaking.Mutex.Unlock()

		//Faz matchmaking local
		MatchLocalQueue(server, nc)
	}
}

func MatchLocalQueue(server *models.Server, nc *nats.Conn) {
	server.Matchmaking.Mutex.Lock()
	defer server.Matchmaking.Mutex.Unlock()

	for len(server.Matchmaking.LocalQueue) >= 2 {
		player1 := server.Matchmaking.LocalQueue[0].Player
		player2 := server.Matchmaking.LocalQueue[1].Player

		if player1.Status != "available" || player2.Status != "available" {
			break
		}

		player1.Status = "playing"
		player2.Status = "playing"

		//cria sala
		room := game.CreateRoom(&player1, &player2, nc, server.ID)
		log.Printf("\n[Server %d] - Nova partida criada: %s (%s vs %s)",
			server.ID, room.ID, player1.UserName, player2.UserName)

		//remove da fila
		server.Matchmaking.LocalQueue = server.Matchmaking.LocalQueue[2:]

		//envia mensagem MATCH para os dois clientes
		sendMatchNotification(server, room)
	}
}

func sendMatchNotification(server *models.Server, room *shared.GameRoom) {
	data, err := json.Marshal(room)
	if err != nil {
		log.Printf("Erro ao serializar room: %v", err)
		return
	}

	resp := shared.Response{
		Action: "MATCH",
		Status: "success",
		Data:   data,
		Server: server.ID,
	}

	msgData, _ := json.Marshal(resp)
	nc := server.Matchmaking.Nc

	// Notifica ambos os jogadores sobre o match
	topic1 := fmt.Sprintf("client.%s.inbox", room.Player1.UserId)
	topic2 := fmt.Sprintf("client.%s.inbox", room.Player2.UserId)

	nc.Publish(topic1, msgData)
	nc.Publish(topic2, msgData)

	log.Printf("[Server %d] - Match enviado para %s e %s",
		server.ID, room.Player1.UserName, room.Player2.UserName)
	log.Printf("[Server %d] - Turno inicial: %s", server.ID, room.Turn)
}

// listar as fila local
func ListLocalQueue(server *models.Server) []string {
	server.Matchmaking.Mutex.Lock()
	defer server.Matchmaking.Mutex.Unlock()

	users := make([]string, len(server.Matchmaking.LocalQueue))
	for i, entry := range server.Matchmaking.LocalQueue {
		users[i] = entry.Player.UserName
	}
	return users
}

//Parte global
func LeaderJoinGlobalQueueHandler(server *models.Server) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if server.Raft.State() != raft.Leader {
            http.Error(w, "Não sou o líder", http.StatusBadRequest)
            return
        }

        var entry shared.QueueEntry
        if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
            http.Error(w, "Payload inválido", http.StatusBadRequest)
            return
        }

        cmd := sharedRaft.Command{
            Type: sharedRaft.CommandQueueJoinGlobal,
            Data: utils.MustMarshal(entry),
        }

        cmdBytes := utils.MustMarshal(cmd)

        future := server.Raft.Apply(cmdBytes, 5*time.Second)
        if err := future.Error(); err != nil {
            http.Error(w, "Erro ao aplicar comando via Raft", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Cliente adicionado à fila global"))
        
        // Tentar criar partidas e notificar
        go func() {
            time.Sleep(50 * time.Millisecond)
            createdRooms := server.FSM.TryMatchPlayers()
            
            // Notificar sobre cada sala criada
            for _, room := range createdRooms {
                NotifyServersAboutMatch(room, server)
            }
        }()
    }
}

func SendToGlobalQueue(entry shared.QueueEntry, server *models.Server) {
	server.Matchmaking.IsLeader = server.Raft.State() == raft.Leader
	eLiderCara := fmt.Sprintf("é lider: %v", server.Matchmaking.IsLeader)
	style.PrintAma(eLiderCara)

    if server.Matchmaking.IsLeader {
		style.PrintAz("entrou na condição")
        cmdData, _ := json.Marshal(sharedRaft.Command{
            Type: sharedRaft.CommandQueueJoinGlobal,
            Data: utils.MustMarshal(entry),
        })
        future := server.Raft.Apply(cmdData, 5*time.Second)
        if err := future.Error(); err != nil {
            log.Printf("[Líder] Erro ao replicar fila global via Raft: %v", err)
            return
        }
        
        // Tentar criar partidas e notificar
        /*go func() {
            time.Sleep(50 * time.Millisecond)
            createdRooms := server.FSM.TryMatchPlayers()
            
            for _, room := range createdRooms {
                NotifyServersAboutMatch(room, server)
            }
        }()*/
        
        return
    }

    // Caso contrário, envia para o líder via HTTP
    leaderAddr := string(server.Raft.Leader())
    if leaderAddr == "" {
        log.Printf("[Follower] Nenhum líder disponível no momento")
        return
    }

    url := fmt.Sprintf("http://%s/leader/join-global-queue", leaderAddr)

    payload := utils.MustMarshal(entry)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        log.Printf("[Follower] Erro ao enviar cliente %s para líder: %v", entry.Player.UserName, err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("[Follower] Resposta inválida do líder: %d", resp.StatusCode)
    } else {
        log.Printf("[Follower] Cliente %s enviado para líder (server%s)", entry.Player.UserName, leaderAddr)
    }
}

func NotifyServersAboutMatch(room *shared.GameRoom, server *models.Server) {
    payload, err := json.Marshal(room)
    if err != nil {
        log.Printf("[Notify] Erro ao serializar sala: %v", err)
        return
    }

    // Servidores que precisam receber a notificação
    targetServers := map[int]bool{
        room.Server1ID: true,
        room.Server2ID: true,
    }

    log.Printf("[Notify] Sala %s: Server1ID=%d, Server2ID=%d, LíderAtual=%d", 
        room.ID, room.Server1ID, room.Server2ID, server.ID)

    // 1. NOTIFICA O PRÓPRIO SERVIDOR SE ELE ESTIVER NA PARTIDA
    if targetServers[server.ID] {
        style.PrintAma(fmt.Sprintf("[Notify] ✓ Notificando servidor LOCAL (server%d)", server.ID))
        utils.NotifyClients(*room, server)
    }

    // 2. NOTIFICA APENAS OS PEERS QUE ESTÃO NA PARTIDA
    for _, peer := range server.Peers {
        if !targetServers[peer.ID] {
            style.PrintVerm(fmt.Sprintf("[Notify] ✗ Pulando server%d (não está na partida)", peer.ID))
            continue
        }

        url := peer.URL + "/notify-match"
        peerID := peer.ID
        
        go func(url string, serverID int) {
            resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
            if err != nil {
                style.PrintVerm(fmt.Sprintf("[Notify] ✗ Erro ao enviar para server%d: %v", serverID, err))
                return
            }
            resp.Body.Close()
            style.PrintMag(fmt.Sprintf("NOTIFICAÇÃO ENVIADA para server%d!", serverID))
            log.Printf("[Notify] Notificação enviada para server%d (%s)", serverID, url)
        }(url, peerID)
    }
}