package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"pbl/style"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	RPC_URL          = "http://127.0.0.1:7545"
    CHAIN_ID         = 1337
)

type EthereumService struct {
	client    *ethclient.Client
	auth      *bind.TransactOpts
	registro  *Registro  
	historico *Historico 
	privateKey *ecdsa.PrivateKey 

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
		privateKey: privateKey,
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
func (s *EthereumService) RegistrarPartida(jogador1, jogador2 string, resultado string) (string, error) {
	style.PrintVerd("vamos resgistrar a partida :)")
	s.mu.Lock()
	defer s.mu.Unlock()

	/*if resultado < 0 || resultado > 2 {
		return "", fmt.Errorf("resultado inválido: deve ser 0 (empate), 1 (j1) ou 2 (j2)")
	}*/

	nonce, _ := s.client.PendingNonceAt(context.Background(), s.auth.From)
	s.auth.Nonce = big.NewInt(int64(nonce))

	tx, err := s.historico.RegistrarPartida(s.auth, jogador1, jogador2, resultado)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

//move uma carta de um usuário para outro.
func (s *EthereumService) TransferirCarta(senderPrivateKeyHex string, cardID string, newOwnerAddressHex string) (string, error) {
	
	if len(senderPrivateKeyHex) > 2 && senderPrivateKeyHex[:2] == "0x" {
		senderPrivateKeyHex = senderPrivateKeyHex[2:]
	}
	userKey, err := crypto.HexToECDSA(senderPrivateKeyHex)
	if err != nil {
		return "", fmt.Errorf("chave privada do usuário inválida: %v", err)
	}

	userAuth, err := bind.NewKeyedTransactorWithChainID(userKey, big.NewInt(int64(CHAIN_ID)))
	if err != nil {
		return "", fmt.Errorf("erro ao criar auth do usuário: %v", err)
	}

	nonce, err := s.client.PendingNonceAt(context.Background(), userAuth.From)
	if err != nil {
		return "", fmt.Errorf("falha ao obter nonce do usuário: %v", err)
	}
	userAuth.Nonce = big.NewInt(int64(nonce))

	if !common.IsHexAddress(newOwnerAddressHex) {
		return "", fmt.Errorf("endereço de destino inválido")
	}
	newOwner := common.HexToAddress(newOwnerAddressHex)

	tx, err := s.registro.TransferirCard(userAuth, cardID, newOwner)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

//envia 0.5 eth para novos usuarios poderem realizar trocas
func (s *EthereumService) EnviarEther(destinatarioHex string) (string, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if !common.IsHexAddress(destinatarioHex) {
        return "", fmt.Errorf("endereço inválido")
    }
    toAddress := common.HexToAddress(destinatarioHex)

    fromAddress := crypto.PubkeyToAddress(s.privateKey.PublicKey)
    nonce, err := s.client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        return "", err
    }

    value := new(big.Int)
    value.SetString("50000000000000000", 10) // 0.05 ETH

    gasLimit := uint64(21000) 
    gasPrice, err := s.client.SuggestGasPrice(context.Background())
    if err != nil {
        return "", err
    }

    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

    chainID := big.NewInt(int64(CHAIN_ID))
    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), s.privateKey)
    if err != nil {
        return "", err
    }

    err = s.client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        return "", err
    }

    s.auth.Nonce = big.NewInt(int64(nonce + 1))

    return signedTx.Hash().Hex(), nil
}
