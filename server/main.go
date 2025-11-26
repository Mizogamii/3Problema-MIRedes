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
	)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}

}