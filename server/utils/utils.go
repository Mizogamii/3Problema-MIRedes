package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"pbl/server/models"

	"log"
	"pbl/shared"
)

//Descobrir o IP do pc que tá rodando o servidor
func LocalIP() (string, error) {
ifaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("erro ao listar interfaces de rede: %v", err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("não foi possível detectar um IP local válido")
}

// Helper para converter qualquer struct em json.RawMessage
func MustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return json.RawMessage(b)
}

func NotifyClients(room shared.GameRoom, server *models.Server) {
    _, err := json.Marshal(room)
    if err != nil {
        log.Printf("[NotifyClients] Erro ao serializar sala: %v", err)
        return
    }

    // Cria mensagem no formato esperado pelo listener global
    gameMsg := shared.GameMessage{
        Type: "GLOBAL_MATCH_CREATED",
        Data: MustMarshal(room),
    }
    msgData := MustMarshal(gameMsg)
    
    nc := server.Matchmaking.Nc

    // Notifica Player1 se estiver neste servidor
    if room.Server1ID == server.ID {
        topic1 := fmt.Sprintf("server.%d.client.%s", server.ID, room.Player1.UserId)
        if err := nc.Publish(topic1, msgData); err != nil {
            log.Printf("[NotifyClients] - Erro ao notificar Player1: %v", err)
        } else {
            log.Printf("[Server %d] - Match GLOBAL enviado para %s via %s", 
                server.ID, room.Player1.UserName, topic1)
        }
    }

    // Notifica Player2 se estiver neste servidor
    if room.Server2ID == server.ID {
        topic2 := fmt.Sprintf("server.%d.client.%s", server.ID, room.Player2.UserId)
        if err := nc.Publish(topic2, msgData); err != nil {
            log.Printf("[NotifyClients] Erro ao notificar Player2: %v", err)
        } else {
            log.Printf("[Server %d] - Match GLOBAL enviado para %s via %s", 
                server.ID, room.Player2.UserName, topic2)
        }
    }
    
    log.Printf("[NotifyClients] Notificação concluída para sala %s", room.ID)
}