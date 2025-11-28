package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"pbl/server/blockchain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RPC_URL     = "http://127.0.0.1:7545"
	CHAIN_ID    = 1337
	
	// Substitir pela chave de uma conta ganache
	ADMIN_KEY   = "776b0127939ca656b6b9bb32e038225ee09808fe6b232ec6fae9be3fa80b145a"
)

func main() {
	log.Println("🔌 Conectando ao Ganache...")
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		log.Fatal(err)
	}

	// Preparar a conta do Admin
	keyStr := strings.TrimPrefix(ADMIN_KEY, "0x")
	privateKey, err := crypto.HexToECDSA(keyStr)
	if err != nil {
		log.Fatal(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(CHAIN_ID))
	if err != nil {
		log.Fatal(err)
	}

	// Deploy do registro (cartas)
	log.Println("📜 Iniciando deploy do contrato REGISTRO...")
	addrRegistro, tx, _, err := blockchain.DeployRegistro(auth, client)
	if err != nil {
		log.Fatalf("Erro ao deployar Registro: %v", err)
	}
	
	log.Printf("   Tx enviada: %s", tx.Hash().Hex())
	log.Println("   ⏳ Aguardando mineração...")
	bind.WaitMined(context.Background(), client, tx)
	log.Printf("✅ REGISTRO DEPLOYADO EM: %s", addrRegistro.Hex())


	// Deploy do historico de Partidas
	log.Println("\n📜 Iniciando deploy do contrato HISTORICO...")
	addrHistorico, tx2, _, err := blockchain.DeployHistorico(auth, client)
	if err != nil {
		log.Fatalf("Erro ao deployar Historico: %v", err)
	}

	log.Printf("   Tx enviada: %s", tx2.Hash().Hex())
	log.Println("   ⏳ Aguardando mineração...")
	bind.WaitMined(context.Background(), client, tx2)
	log.Printf("✅ HISTORICO DEPLOYADO EM: %s", addrHistorico.Hex())

	// Resumo para o Makefile
	fmt.Println("\n========================================================")
	fmt.Println(" COPIE ESTAS LINHAS PARA O SEU MAKEFILE (no topo):")
	fmt.Println("========================================================")
	fmt.Printf("CONTRACT_ADDR_REGISTRO := \"%s\"\n", addrRegistro.Hex())
	fmt.Printf("CONTRACT_ADDR_HISTORICO := \"%s\"\n", addrHistorico.Hex())
	fmt.Println("========================================================")
}