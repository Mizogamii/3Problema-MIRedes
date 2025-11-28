package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Iniciando servidor...")

	err := StartServer(
		os.Getenv("ID"),
		os.Getenv("PORT"),
		os.Getenv("PEERS"),
		os.Getenv("NATS_URL"),
		os.Getenv("BLOCKCHAIN_PRIVATE_KEY"),
		os.Getenv("BLOCKCHAIN_CONTRACT_ADDR_REGISTRO"),
		os.Getenv("BLOCKCHAIN_CONTRACT_ADDR_HISTORICO"),
	)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}

}