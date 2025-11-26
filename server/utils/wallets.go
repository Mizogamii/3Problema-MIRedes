package utils

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"pbl/style"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// cria um par de chaves novo.
func GenerateNewWallet() (string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString := hexutil.Encode(privateKeyBytes)[2:] // Remove o "0x"

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("erro no cast da chave publica")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()//[2:]

	style.PrintVerm("CRIADA CHAVE E ENDEREÇO")
	fmt.Printf("%s -- %s\n", address, privateKeyString)

	return address, privateKeyString
}