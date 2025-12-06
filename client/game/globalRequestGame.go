package game

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"pbl/shared"
	"pbl/style"

	"github.com/nats-io/nats.go"
)

func HandleStartGlobalMatchListener(serverID int, nc *nats.Conn, clientID string, matchChan chan<- MatchInfo) *nats.Subscription {
	// MUDAR PARA O MESMO TÓPICO do match local
	clientTopic := fmt.Sprintf("client.%s.inbox", clientID)

	sub, err := nc.Subscribe(clientTopic, func(msg *nats.Msg) {
		// PRIMEIRO tenta como Response (usado por MATCH)
		var resp shared.Response
		if err := json.Unmarshal(msg.Data, &resp); err == nil {
			// Se for um MATCH (local ou global)
			if resp.Action == "MATCH" {
				var room shared.GameRoom
				if err := json.Unmarshal(resp.Data, &room); err != nil {
					log.Println("Erro ao decodificar sala:", err)
					return
				}

				var opponent shared.User
				if room.Player1.UserId == clientID {
					opponent = *room.Player2
				} else {
					opponent = *room.Player1
				}

				// Determina se é global baseado no ID da sala
				isGlobal := strings.HasPrefix(room.ID, "global-")

				matchChan <- MatchInfo{
					Opponent: opponent,
					Room:     room,
					IsGlobal: isGlobal,
				}

				log.Printf("[Cliente] Nova partida recebida! Sala: %s, Global: %v", room.ID, isGlobal)
				return
			}
		}

		var gameMsg shared.GameMessage
		if err := json.Unmarshal(msg.Data, &gameMsg); err != nil {
			log.Println("Erro ao decodificar mensagem:", err)
			return
		}

		if gameMsg.Type == "GLOBAL_MATCH_CREATED" {
			var room shared.GameRoom
			if err := json.Unmarshal(gameMsg.Data, &room); err != nil {
				log.Println("Erro ao decodificar dados da sala:", err)
				return
			}

			var opponent shared.User
			if room.Player1.UserId == clientID {
				opponent = *room.Player2
			} else {
				opponent = *room.Player1
			}

			matchChan <- MatchInfo{
				Opponent: opponent,
				Room:     room,
				IsGlobal: true,
			}
		}
	})

	if err != nil {
		log.Printf("Erro ao se inscrever no tópico %s: %v", clientTopic, err)
		return nil
	}

	log.Printf("[Cliente] Inscrito no tópico: %s", clientTopic)
	return sub
}

func PlayGlobalGame(nc *nats.Conn, room *shared.GameRoom, currentUser shared.User, opponent shared.User) {

	fmt.Println("\nPARTIDA GLOBAL INICIADA")
	fmt.Print("\n----------------------------------")
	fmt.Printf("\nSala ID: %s", room.ID)
	fmt.Print("\n----------------------------------")
	fmt.Printf("\nPlayer1: %s", room.Player1.UserName)
	fmt.Printf("\nPlayer2: %s", room.Player2.UserName)
	fmt.Print("\n----------------------------------\n")

	isMyTurn := room.Turn == currentUser.UserId
	if isMyTurn {
		fmt.Println("✓ Você começa!")
	} else {
		fmt.Printf("%s começa. Aguarde...\n", opponent.UserName)
	}

	gameOver := false
	alreadyPlayed := false
	gameMsgChan := make(chan shared.GameMessage, 10)

	clientTopic := fmt.Sprintf("client.%s.inbox", currentUser.UserId)
	//log.Printf("[Cliente] Inscrito no tópico: %s", clientTopic)
	
	sub, err := nc.Subscribe(clientTopic, func(msg *nats.Msg) {
		var gameMsg shared.GameMessage
		if err := json.Unmarshal(msg.Data, &gameMsg); err != nil {
			log.Println("Erro ao decodificar mensagem:", err)
			return
		}
		//log.Printf("[Cliente] Mensagem recebida: Type=%s, From=%s", gameMsg.Type, gameMsg.From)
		gameMsgChan <- gameMsg
	})
	if err != nil {
		log.Println("Erro ao criar subscription:", err)
		return
	}
	defer sub.Unsubscribe()


	if isMyTurn {
		card, ok := ChooseCard(currentUser)
		if !ok {
			fmt.Println("Você desistiu da partida.")
			return
		}

		SendCardPlayGlobal(nc, room, currentUser, card)
		fmt.Printf("\nVocê jogou: %s (%s)\n", card.Element, card.Type)
		fmt.Println("Aguardando adversário...")
		alreadyPlayed = true
	}

	for !gameOver {
		select {
		case gameMsg := <-gameMsgChan:
			switch gameMsg.Type {
			case "PLAY_CARD":
				if gameMsg.From == opponent.UserId {
					var card shared.Card
					if len(gameMsg.Data) > 0 {
						if err := json.Unmarshal(gameMsg.Data, &card); err != nil {
							log.Println("Erro ao decodificar carta:", err)
							continue
						}
					}

					//fmt.Printf("\n%s jogou: %s (%s)\n", opponent.UserName, card.Element, card.Type)

					if gameMsg.Turn == currentUser.UserId && !alreadyPlayed {
						fmt.Println("\n✓ Sua vez de jogar!")
						chosenCard, ok := ChooseCard(currentUser)
						if ok {
							SendCardPlayGlobal(nc, room, currentUser, chosenCard)
							fmt.Printf("Você jogou: %s (%s)\n", chosenCard.Element, chosenCard.Type)
							fmt.Println("Aguardando resultado...")
							alreadyPlayed = true
						} else {
							fmt.Println("Você desistiu da partida.")
							return
						}
					} else if alreadyPlayed {
						fmt.Println("Aguardando resultado da rodada...")
					}
				}

			case "ROUND_RESULT":
				fmt.Println("\n--------------------------------")
				fmt.Println("            Resultado           ")
				fmt.Println("--------------------------------")

				if gameMsg.Winner == nil {
					fmt.Println("Empate dos jogadores!")
				} else {
					fmt.Println("Vencedor: ", gameMsg.Winner.UserName)
				}

				alreadyPlayed = false
				gameOver = true

			}	

		case <-time.After(30 * time.Second):
			fmt.Println("\nTimeout: servidor não respondeu.")
			fmt.Println("A partida foi cancelada.")
			gameOver = true
			break
		}
	}

	fmt.Print("Pressione ENTER para voltar ao menu principal...")
	fmt.Scanln()
	style.Clear()
}

func SendCardPlayGlobal(nc *nats.Conn, room *shared.GameRoom, client shared.User, card shared.Card) {
	dataBytes, _ := json.Marshal(card)

	gameMsg := shared.GameMessage{
		Type:   "PLAY_CARD_GLOBAL",
		From:   client.UserId,
		RoomID: room.ID,
		Data:   dataBytes,
	}

	payload, _ := json.Marshal(gameMsg)

	req := shared.Request{
		ClientID: client.UserId,
		Action:   "GAME_MESSAGE_GLOBAL",
		Payload:  payload,
	}

	reqBytes, _ := json.Marshal(req)
	topic := fmt.Sprintf("server.%d.requests", client.ServerID)

	//log.Printf("[DEBUG] Enviando jogada para o servidor %d (sala %s): %+v\n", client.ServerID, room.ID, card)

	nc.Publish(topic, reqBytes)
}
