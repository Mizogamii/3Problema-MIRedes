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
	client   *ethclient.Client
	registro *Registro
	auth     *bind.TransactOpts
    mu       sync.Mutex 
}

// Inicia a conexão
func NewEthereumService(privateKey string, registro_add string) (*EthereumService, error) {
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return nil, err
	}

	privK, _ := crypto.HexToECDSA(privateKey)
	auth, _ := bind.NewKeyedTransactorWithChainID(privK, big.NewInt(int64(CHAIN_ID)))

	m := fmt.Sprintf("endereço do registro %s", registro_add)
	style.PrintMag(m)
	address := common.HexToAddress(registro_add)
	instance, err := NewRegistro(address, client)
	if err != nil {
		return nil, err
	}

	return &EthereumService{
		client:   client,
		registro: instance,
		auth:     auth,
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