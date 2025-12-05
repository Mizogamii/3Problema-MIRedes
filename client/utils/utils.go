package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"pbl/client/models"
	"pbl/shared"
	"pbl/style"

	"github.com/nats-io/nats.go"
)

func ReadLineSafe() string {
    reader := bufio.NewReader(os.Stdin)
    input, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Erro ao ler input:", err)
        return ""
    }
    return strings.TrimSpace(input)
}

func Contains(slice []int, value int) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
    }
    return false
}

func Clear() {
	nameOS := runtime.GOOS
	fmt.Println("Sistema operacional:", nameOS)

	var cmd *exec.Cmd
	if nameOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func GerarIdAleatorio() int {
	var b [4]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	id := int(binary.LittleEndian.Uint32(b[:]))
	fmt.Println("ID:", id)
	return int(math.Abs(float64(id)))
}

func GenerateRoomID(serverID int) string {
	var b [4]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	id := int(binary.LittleEndian.Uint32(b[:]))
	return fmt.Sprintf("%d-%d", serverID, int(math.Abs(float64(id))))
}

// Helper para converter qualquer struct em json.RawMessage
func MustMarshal(v interface{}) json.RawMessage {
	b, _ := json.Marshal(v)
	return json.RawMessage(b)
}

func HandleClientSeeCards(nc *nats.Conn, server models.ServerInfo, clientID string)[]shared.Card{
	//fmt.Println("Buscando cartas...")
	req := shared.Request{
		ClientID: clientID,
		Action: "SEE_CARDS",
		Payload: nil,
	}
	reqData,_ := json.Marshal(req)

	topic := fmt.Sprintf("server.%d.requests", server.ID)
	msg, err := nc.Request(topic, reqData, 5*time.Second)

	if err != nil {
		log.Printf("Erro na requisição para pegar carta: %v", err)
		return nil// Sai da função imediatamente para evitar o crash
	}

	if msg == nil || msg.Data == nil{
		log.Printf("O servidor retornou uma resposta vazia.")
		return nil
	}

	var response shared.Response
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Printf("Erro ao decodificar resposta do inventário: %v", err)
		return nil
	}

	if response.Status == "success"{
		var inventario shared.Cards
		if err := json.Unmarshal(response.Data, &inventario); err != nil {
			log.Printf("Erro ao decodificar os dados do inventario: %v", err)
			return nil
		}
		MostrarInventario(inventario.Cards)
		return inventario.Cards
	} else {
		msg := fmt.Sprintf("\n[FALHA] Não foi possível ver inventário: %s\n", response.Error)
		style.PrintVerm(msg)
	}

	return nil
}

func HandleChangeDeckClient(nc *nats.Conn, server models.ServerInfo, clientID string, user *shared.User){
	cards := HandleClientSeeCards(nc, server, clientID)
	deck := ChoseDeck(cards)
	MostrarInventario(deck)


	deckCodf, err := json.Marshal(deck)
	if err!=nil{
		fmt.Printf("\nErro ao converter para JSON: %v", err)
		return
	}
	req := shared.Request{
		ClientID: clientID,
		Action: "CHANGE_DECK",
		Payload: deckCodf,
	}
	reqData, _ := json.Marshal(req)
	topic := fmt.Sprintf("server.%d.requests", server.ID)
	msg, err := nc.Request(topic, reqData, 5*time.Second)
	if err != nil{
		fmt.Printf("\nErro ao salvar deck: %v", err)
		return

	}
	var response shared.Response
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		fmt.Printf("\nErro ao decodificar resposta do servidor: %v", err)
		return
	}
	if response.Status == "success"{
		user.Deck = deck 
		style.PrintVerd("Deck salvo!")
	} else {
		style.PrintVerm("Erro ao salvar deck, tente novamente")
	}
}

func ChoseDeck(cards []shared.Card) []shared.Card{
	var selectedCards []int
	var deck []shared.Card
	if cards != nil{
		for i := range(4){
			valida := false
			for ! valida{
				fmt.Printf("Digite o número da %d° carta para o baralho: ", i+1)
				in := ReadLineSafe()
				if inInt, err := strconv.Atoi(in); err == nil {
					if inInt>=0 && inInt<len(cards) && !Contains(selectedCards, inInt){
						deck = append(deck, cards[inInt])
						selectedCards = append(selectedCards, inInt)
						valida = true
					}else{
						style.PrintMag("Valor inválido\n")
					}
				}else{
					style.PrintMag("Valor inválido, digite o número da carta!\n")
				}
			}
		}
	}
	
	return deck
}

func HandleClientSeePackCards(nc *nats.Conn, server models.ServerInfo, clientID string)[]shared.Card{
	//fmt.Println("Buscando cartas...")
	req := shared.Request{
		ClientID: clientID,
		Action: "SEE_CARDS",
		Payload: nil,
	}
	reqData,_ := json.Marshal(req)

	topic := fmt.Sprintf("server.%d.requests", server.ID)
	msg, err := nc.Request(topic, reqData, 5*time.Second)

	if err != nil {
		log.Printf("Erro na requisição para pegar carta: %v", err)
		return nil// Sai da função imediatamente para evitar o crash
	}

	if msg == nil || msg.Data == nil{
		log.Printf("O servidor retornou uma resposta vazia.")
		return nil
	}

	var response shared.Response
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Printf("Erro ao decodificar resposta do inventário: %v", err)
		return nil
	}

	if response.Status == "success"{
		var inventario shared.Cards
		if err := json.Unmarshal(response.Data, &inventario); err != nil {
			log.Printf("Erro ao decodificar os dados do inventario: %v", err)
			return nil
		}
		packCards := inventario.Cards[5:]
		MostrarInventario(packCards)
		return packCards
	} else {
		msg := fmt.Sprintf("\n[FALHA] Não foi possível ver inventário: %s\n", response.Error)
		style.PrintVerm(msg)
	}

	return nil
}
