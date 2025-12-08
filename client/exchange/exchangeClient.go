package exchange

import (
	"fmt"
	"log"
	"time"

	"encoding/json"

	"pbl/client/models"
	"pbl/shared"
	"pbl/style"

	"pbl/client/utils"

	"github.com/nats-io/nats.go"
)

//Local
func HandleExchange(nc *nats.Conn, server models.ServerInfo, user *shared.User){
    cards := utils.HandleClientSeePackCards(nc, server, user.UserId)
    //m := fmt.Sprintf("cartas do cliente: %i", len(cards))
    //style.PrintCian(m)
    if cards == nil || len(cards) == 0 {
        fmt.Println("Você não tem cartas para trocar!\nabra pacotes para conseguir mais cartas")
        return
    }

    exchangeCardIndex := utils.Troca()
    if exchangeCardIndex < 0 || exchangeCardIndex >= len(cards){
        fmt.Println("Carta inválida!")
        return
    }

    notifChan := make(chan shared.ExchangeNotification, 1)
    topic := fmt.Sprintf("exchange.notify.%s", user.UserId)
    
    sub, err := nc.Subscribe(topic, func(msg *nats.Msg) {
        var notif shared.ExchangeNotification
        if err := json.Unmarshal(msg.Data, &notif); err != nil {
            log.Printf("Erro ao deserializar notificação: %v", err)
            return
        }
        select {
        case notifChan <- notif:
        default:
        }
    })
    
    if err != nil {
        fmt.Println("Erro ao se inscrever para notificações:", err)
        return
    }
    defer sub.Unsubscribe()

    // Prepara o request de troca
    requestExchange := shared.ExchangeRequest{
        Player: *user,
        CardOffered: cards[exchangeCardIndex],
        ServerID: server.ID,
        Timestamp: time.Now(),
    }

    data, _ := json.Marshal(requestExchange)

    request := shared.Request{
        ClientID: user.UserId,
        Action: "EXCHANGE_REQUEST",
        Payload: data,
    }

    send, _ := json.Marshal(request)
    serverTopic := fmt.Sprintf("server.%d.requests", server.ID)

    // Envia request e aguarda confirmação de entrada na fila
    msg, err := nc.Request(serverTopic, send, 5*time.Second)
    if err != nil {
        fmt.Println("Erro ao enviar o pedido de troca:", err)
        return
    }

    var response shared.Response
    if err := json.Unmarshal(msg.Data, &response); err != nil {
        fmt.Println("Erro ao decodificar resposta:", err)
        return
    }

    if response.Status != "success" {
        fmt.Println("ERRO:", response.Error)
        return
    }

    // Confirmação de entrada na fila
    fmt.Println("\n✓ Você entrou na fila de troca!")
    fmt.Println("Aguardando outro jogador...")

    // Aguarda a notificação de match (com timeout)
    var notif shared.ExchangeNotification
    select {
    case notif = <-notifChan:
        // Match encontrado!
    case <-time.After(5 * time.Minute):
        fmt.Println("\nTempo limite excedido. Nenhum parceiro encontrado.")
        return
    }
    
	style.Clear()
    fmt.Println("\n------------------------------------------")
    fmt.Println("        Você foi pareado para troca!       ")
    fmt.Println("------------------------------------------")
    fmt.Printf("Sua carta enviada: ")
    utils.PrintCartaCor(notif.YouSend)
    fmt.Print("\n")
    fmt.Printf("Carta recebida: ")
    utils.PrintCartaCor(notif.YouGet)
    fmt.Print("\n")
	fmt.Printf("Recebido de: %v\n", notif.Partner)
    fmt.Println("------------------------------------------")
    
    utils.HandleClientSeeCards(nc, server, user.UserId)
    //para obrigar o usuario a trocar o deck
    style.PrintAma("\n    Troca realizada com sucesso!")
    fmt.Print("\n\n")
    style.PrintMag("    Modifique o seu deck para as\npróximas partidas:")
    utils.HandleChangeDeckClient(nc, server, user.UserId, user)

    fmt.Println("\nPressione ENTER para voltar ao menu...")
	fmt.Scanln()

}
