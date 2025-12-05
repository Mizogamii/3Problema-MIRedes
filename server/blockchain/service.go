package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"pbl/style"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RPC_URL          = "http://127.0.0.1:7545"
	//REGISTRO_ADDRESS = "ENDEREÇO_CONTRATO_REGISTRO" //TODO: mover essa contante para variavel do ambiente
    CHAIN_ID         = 1337
)

type EthereumService struct {
	client    *ethclient.Client
	auth      *bind.TransactOpts
	registro  *Registro  
	historico *Historico 
	
    mu        sync.Mutex
}

// Inicia a conexão
func NewEthereumService(privateKeyHex, addrRegistroHex, addrHistoricoHex string) (*EthereumService, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar RPC: %v", err)
	}

	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("chave privada inválida: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(CHAIN_ID)))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar autenticador: %v", err)
	}

	if !common.IsHexAddress(addrRegistroHex) {
		return nil, fmt.Errorf("endereço registro inválido")
	}
	registroInst, err := NewRegistro(common.HexToAddress(addrRegistroHex), client)
	if err != nil {
		return nil, err
	}

	if !common.IsHexAddress(addrHistoricoHex) {
		return nil, fmt.Errorf("endereço historico inválido")
	}
	historicoInst, err := NewHistorico(common.HexToAddress(addrHistoricoHex), client)
	if err != nil {
		return nil, err
	}

	return &EthereumService{
		client:    client,
		registro:  registroInst, 
		historico: historicoInst,
		auth:      auth,
	}, nil
}

func (s *EthereumService) CriarCarta(playerAddress string, id, elemento, tipo string) (string, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    nonce, _ := s.client.PendingNonceAt(context.Background(), s.auth.From)
    s.auth.Nonce = big.NewInt(int64(nonce))

    destinatario := common.HexToAddress(playerAddress)

    tx, err := s.registro.RegistrarCardParaJogador(
        s.auth,       
        destinatario, 
        id,           
        elemento,     
        tipo,         
    )

    if err != nil {
        return "", err
    }

    return tx.Hash().Hex(), nil
}

//grava o resultado do jogo na blockchain
//resultado: 0 (Empate), 1 (Vitoria Jogador 1), 2 (Vitoria Jogador 2)
func (s *EthereumService) RegistrarPartida(jogador1, jogador2 string, resultado int) (string, error) {
	style.PrintVerd("vamos resgistrar a partida :)")
	s.mu.Lock()
	defer s.mu.Unlock()

	if resultado < 0 || resultado > 2 {
		return "", fmt.Errorf("resultado inválido: deve ser 0 (empate), 1 (j1) ou 2 (j2)")
	}

	nonce, _ := s.client.PendingNonceAt(context.Background(), s.auth.From)
	s.auth.Nonce = big.NewInt(int64(nonce))

	tx, err := s.historico.RegistrarPartida(s.auth, jogador1, jogador2, uint8(resultado))
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}