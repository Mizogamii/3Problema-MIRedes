// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// RegistroMetaData contains all meta data concerning the Registro contract.
var RegistroMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"cards\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"element\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"cardType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"dono\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_jogadorDestino\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_id\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_element\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_cardType\",\"type\":\"string\"}],\"name\":\"registrarCardParaJogador\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_id\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_novoDono\",\"type\":\"address\"}],\"name\":\"transferirCard\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50610da28061001f6000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806333697b66146100465780634dac95a814610062578063e92908d014610095575b600080fd5b610060600480360381019061005b9190610711565b6100b1565b005b61007c6004803603810190610077919061076d565b6101c4565b60405161008c9493929190610844565b60405180910390f35b6100af60048036038101906100aa919061089e565b6103c2565b005b3373ffffffffffffffffffffffffffffffffffffffff166000836040516100d89190610995565b908152602001604051809103902060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610160576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610157906109f8565b60405180910390fd5b806000836040516101719190610995565b908152602001604051809103902060030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b6000818051602081018201805184825260208301602085012081835280955050505050506000915090508060000180546101fd90610a47565b80601f016020809104026020016040519081016040528092919081815260200182805461022990610a47565b80156102765780601f1061024b57610100808354040283529160200191610276565b820191906000526020600020905b81548152906001019060200180831161025957829003601f168201915b50505050509080600101805461028b90610a47565b80601f01602080910402602001604051908101604052809291908181526020018280546102b790610a47565b80156103045780601f106102d957610100808354040283529160200191610304565b820191906000526020600020905b8154815290600101906020018083116102e757829003601f168201915b50505050509080600201805461031990610a47565b80601f016020809104026020016040519081016040528092919081815260200182805461034590610a47565b80156103925780601f1061036757610100808354040283529160200191610392565b820191906000526020600020905b81548152906001019060200180831161037557829003601f168201915b5050505050908060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905084565b600073ffffffffffffffffffffffffffffffffffffffff166000846040516103ea9190610995565b908152602001604051809103902060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610472576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161046990610ac4565b60405180910390fd5b60405180608001604052808381526020018281526020018481526020018573ffffffffffffffffffffffffffffffffffffffff168152506000846040516104b99190610995565b908152602001604051809103902060008201518160000190816104dc9190610c9a565b5060208201518160010190816104f29190610c9a565b5060408201518160020190816105089190610c9a565b5060608201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505050505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6105c082610577565b810181811067ffffffffffffffff821117156105df576105de610588565b5b80604052505050565b60006105f2610559565b90506105fe82826105b7565b919050565b600067ffffffffffffffff82111561061e5761061d610588565b5b61062782610577565b9050602081019050919050565b82818337600083830152505050565b600061065661065184610603565b6105e8565b90508281526020810184848401111561067257610671610572565b5b61067d848285610634565b509392505050565b600082601f83011261069a5761069961056d565b5b81356106aa848260208601610643565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106de826106b3565b9050919050565b6106ee816106d3565b81146106f957600080fd5b50565b60008135905061070b816106e5565b92915050565b6000806040838503121561072857610727610563565b5b600083013567ffffffffffffffff81111561074657610745610568565b5b61075285828601610685565b9250506020610763858286016106fc565b9150509250929050565b60006020828403121561078357610782610563565b5b600082013567ffffffffffffffff8111156107a1576107a0610568565b5b6107ad84828501610685565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60005b838110156107f05780820151818401526020810190506107d5565b60008484015250505050565b6000610807826107b6565b61081181856107c1565b93506108218185602086016107d2565b61082a81610577565b840191505092915050565b61083e816106d3565b82525050565b6000608082019050818103600083015261085e81876107fc565b9050818103602083015261087281866107fc565b9050818103604083015261088681856107fc565b90506108956060830184610835565b95945050505050565b600080600080608085870312156108b8576108b7610563565b5b60006108c6878288016106fc565b945050602085013567ffffffffffffffff8111156108e7576108e6610568565b5b6108f387828801610685565b935050604085013567ffffffffffffffff81111561091457610913610568565b5b61092087828801610685565b925050606085013567ffffffffffffffff81111561094157610940610568565b5b61094d87828801610685565b91505092959194509250565b600081905092915050565b600061096f826107b6565b6109798185610959565b93506109898185602086016107d2565b80840191505092915050565b60006109a18284610964565b915081905092915050565b7f566f6365206e616f2065206f20646f6e6f000000000000000000000000000000600082015250565b60006109e26011836107c1565b91506109ed826109ac565b602082019050919050565b60006020820190508181036000830152610a11816109d5565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610a5f57607f821691505b602082108103610a7257610a71610a18565b5b50919050565b7f4944206a6120656d2075736f0000000000000000000000000000000000000000600082015250565b6000610aae600c836107c1565b9150610ab982610a78565b602082019050919050565b60006020820190508181036000830152610add81610aa1565b9050919050565b60008190508160005260206000209050919050565b60006020601f8301049050919050565b600082821b905092915050565b600060088302610b467fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610b09565b610b508683610b09565b95508019841693508086168417925050509392505050565b6000819050919050565b6000819050919050565b6000610b97610b92610b8d84610b68565b610b72565b610b68565b9050919050565b6000819050919050565b610bb183610b7c565b610bc5610bbd82610b9e565b848454610b16565b825550505050565b600090565b610bda610bcd565b610be5818484610ba8565b505050565b5b81811015610c0957610bfe600082610bd2565b600181019050610beb565b5050565b601f821115610c4e57610c1f81610ae4565b610c2884610af9565b81016020851015610c37578190505b610c4b610c4385610af9565b830182610bea565b50505b505050565b600082821c905092915050565b6000610c7160001984600802610c53565b1980831691505092915050565b6000610c8a8383610c60565b9150826002028217905092915050565b610ca3826107b6565b67ffffffffffffffff811115610cbc57610cbb610588565b5b610cc68254610a47565b610cd1828285610c0d565b600060209050601f831160018114610d045760008415610cf2578287015190505b610cfc8582610c7e565b865550610d64565b601f198416610d1286610ae4565b60005b82811015610d3a57848901518255600182019150602085019450602081019050610d15565b86831015610d575784890151610d53601f891682610c60565b8355505b6001600288020188555050505b50505050505056fea2646970667358221220cb4a797ce932fb327ba6b886c1501cd1369b9efb5db3c8e16542e5613482369a64736f6c634300081e0033",
}

// RegistroABI is the input ABI used to generate the binding from.
// Deprecated: Use RegistroMetaData.ABI instead.
var RegistroABI = RegistroMetaData.ABI

// RegistroBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RegistroMetaData.Bin instead.
var RegistroBin = RegistroMetaData.Bin

// DeployRegistro deploys a new Ethereum contract, binding an instance of Registro to it.
func DeployRegistro(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Registro, error) {
	parsed, err := RegistroMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RegistroBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Registro{RegistroCaller: RegistroCaller{contract: contract}, RegistroTransactor: RegistroTransactor{contract: contract}, RegistroFilterer: RegistroFilterer{contract: contract}}, nil
}

// Registro is an auto generated Go binding around an Ethereum contract.
type Registro struct {
	RegistroCaller     // Read-only binding to the contract
	RegistroTransactor // Write-only binding to the contract
	RegistroFilterer   // Log filterer for contract events
}

// RegistroCaller is an auto generated read-only Go binding around an Ethereum contract.
type RegistroCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistroTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RegistroTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistroFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RegistroFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RegistroSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RegistroSession struct {
	Contract     *Registro         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RegistroCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RegistroCallerSession struct {
	Contract *RegistroCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// RegistroTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RegistroTransactorSession struct {
	Contract     *RegistroTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RegistroRaw is an auto generated low-level Go binding around an Ethereum contract.
type RegistroRaw struct {
	Contract *Registro // Generic contract binding to access the raw methods on
}

// RegistroCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RegistroCallerRaw struct {
	Contract *RegistroCaller // Generic read-only contract binding to access the raw methods on
}

// RegistroTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RegistroTransactorRaw struct {
	Contract *RegistroTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRegistro creates a new instance of Registro, bound to a specific deployed contract.
func NewRegistro(address common.Address, backend bind.ContractBackend) (*Registro, error) {
	contract, err := bindRegistro(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Registro{RegistroCaller: RegistroCaller{contract: contract}, RegistroTransactor: RegistroTransactor{contract: contract}, RegistroFilterer: RegistroFilterer{contract: contract}}, nil
}

// NewRegistroCaller creates a new read-only instance of Registro, bound to a specific deployed contract.
func NewRegistroCaller(address common.Address, caller bind.ContractCaller) (*RegistroCaller, error) {
	contract, err := bindRegistro(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RegistroCaller{contract: contract}, nil
}

// NewRegistroTransactor creates a new write-only instance of Registro, bound to a specific deployed contract.
func NewRegistroTransactor(address common.Address, transactor bind.ContractTransactor) (*RegistroTransactor, error) {
	contract, err := bindRegistro(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RegistroTransactor{contract: contract}, nil
}

// NewRegistroFilterer creates a new log filterer instance of Registro, bound to a specific deployed contract.
func NewRegistroFilterer(address common.Address, filterer bind.ContractFilterer) (*RegistroFilterer, error) {
	contract, err := bindRegistro(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RegistroFilterer{contract: contract}, nil
}

// bindRegistro binds a generic wrapper to an already deployed contract.
func bindRegistro(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RegistroMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registro *RegistroRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registro.Contract.RegistroCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registro *RegistroRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registro.Contract.RegistroTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registro *RegistroRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registro.Contract.RegistroTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Registro *RegistroCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Registro.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Registro *RegistroTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Registro.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Registro *RegistroTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Registro.Contract.contract.Transact(opts, method, params...)
}

// Cards is a free data retrieval call binding the contract method 0x4dac95a8.
//
// Solidity: function cards(string ) view returns(string element, string cardType, string id, address dono)
func (_Registro *RegistroCaller) Cards(opts *bind.CallOpts, arg0 string) (struct {
	Element  string
	CardType string
	Id       string
	Dono     common.Address
}, error) {
	var out []interface{}
	err := _Registro.contract.Call(opts, &out, "cards", arg0)

	outstruct := new(struct {
		Element  string
		CardType string
		Id       string
		Dono     common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Element = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.CardType = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Id = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Dono = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Cards is a free data retrieval call binding the contract method 0x4dac95a8.
//
// Solidity: function cards(string ) view returns(string element, string cardType, string id, address dono)
func (_Registro *RegistroSession) Cards(arg0 string) (struct {
	Element  string
	CardType string
	Id       string
	Dono     common.Address
}, error) {
	return _Registro.Contract.Cards(&_Registro.CallOpts, arg0)
}

// Cards is a free data retrieval call binding the contract method 0x4dac95a8.
//
// Solidity: function cards(string ) view returns(string element, string cardType, string id, address dono)
func (_Registro *RegistroCallerSession) Cards(arg0 string) (struct {
	Element  string
	CardType string
	Id       string
	Dono     common.Address
}, error) {
	return _Registro.Contract.Cards(&_Registro.CallOpts, arg0)
}

// RegistrarCardParaJogador is a paid mutator transaction binding the contract method 0xe92908d0.
//
// Solidity: function registrarCardParaJogador(address _jogadorDestino, string _id, string _element, string _cardType) returns()
func (_Registro *RegistroTransactor) RegistrarCardParaJogador(opts *bind.TransactOpts, _jogadorDestino common.Address, _id string, _element string, _cardType string) (*types.Transaction, error) {
	return _Registro.contract.Transact(opts, "registrarCardParaJogador", _jogadorDestino, _id, _element, _cardType)
}

// RegistrarCardParaJogador is a paid mutator transaction binding the contract method 0xe92908d0.
//
// Solidity: function registrarCardParaJogador(address _jogadorDestino, string _id, string _element, string _cardType) returns()
func (_Registro *RegistroSession) RegistrarCardParaJogador(_jogadorDestino common.Address, _id string, _element string, _cardType string) (*types.Transaction, error) {
	return _Registro.Contract.RegistrarCardParaJogador(&_Registro.TransactOpts, _jogadorDestino, _id, _element, _cardType)
}

// RegistrarCardParaJogador is a paid mutator transaction binding the contract method 0xe92908d0.
//
// Solidity: function registrarCardParaJogador(address _jogadorDestino, string _id, string _element, string _cardType) returns()
func (_Registro *RegistroTransactorSession) RegistrarCardParaJogador(_jogadorDestino common.Address, _id string, _element string, _cardType string) (*types.Transaction, error) {
	return _Registro.Contract.RegistrarCardParaJogador(&_Registro.TransactOpts, _jogadorDestino, _id, _element, _cardType)
}

// TransferirCard is a paid mutator transaction binding the contract method 0x33697b66.
//
// Solidity: function transferirCard(string _id, address _novoDono) returns()
func (_Registro *RegistroTransactor) TransferirCard(opts *bind.TransactOpts, _id string, _novoDono common.Address) (*types.Transaction, error) {
	return _Registro.contract.Transact(opts, "transferirCard", _id, _novoDono)
}

// TransferirCard is a paid mutator transaction binding the contract method 0x33697b66.
//
// Solidity: function transferirCard(string _id, address _novoDono) returns()
func (_Registro *RegistroSession) TransferirCard(_id string, _novoDono common.Address) (*types.Transaction, error) {
	return _Registro.Contract.TransferirCard(&_Registro.TransactOpts, _id, _novoDono)
}

// TransferirCard is a paid mutator transaction binding the contract method 0x33697b66.
//
// Solidity: function transferirCard(string _id, address _novoDono) returns()
func (_Registro *RegistroTransactorSession) TransferirCard(_id string, _novoDono common.Address) (*types.Transaction, error) {
	return _Registro.Contract.TransferirCard(&_Registro.TransactOpts, _id, _novoDono)
}
