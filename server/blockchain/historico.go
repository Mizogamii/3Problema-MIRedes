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

// HistoricoMetaData contains all meta data concerning the Historico contract.
var HistoricoMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"jogador1\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"jogador2\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"enumHistorico.Resultado\",\"name\":\"resultado\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"NovaPartida\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_j1\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_j2\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_resultado\",\"type\":\"uint8\"}],\"name\":\"registrarPartida\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506105328061001f6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806333896cea14610030575b600080fd5b61004a6004803603810190610045919061027d565b61004c565b005b60028160ff161115610093576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161008a9061038b565b60405180910390fd5b7f56cc181ac353627e55fa5eccf2e64c91f74aae5a157dc1cf0acadc1a97237b5783838360ff1660028111156100cc576100cb6103ab565b5b426040516100dd94939291906104a9565b60405180910390a1505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b61015182610108565b810181811067ffffffffffffffff821117156101705761016f610119565b5b80604052505050565b60006101836100ea565b905061018f8282610148565b919050565b600067ffffffffffffffff8211156101af576101ae610119565b5b6101b882610108565b9050602081019050919050565b82818337600083830152505050565b60006101e76101e284610194565b610179565b90508281526020810184848401111561020357610202610103565b5b61020e8482856101c5565b509392505050565b600082601f83011261022b5761022a6100fe565b5b813561023b8482602086016101d4565b91505092915050565b600060ff82169050919050565b61025a81610244565b811461026557600080fd5b50565b60008135905061027781610251565b92915050565b600080600060608486031215610296576102956100f4565b5b600084013567ffffffffffffffff8111156102b4576102b36100f9565b5b6102c086828701610216565b935050602084013567ffffffffffffffff8111156102e1576102e06100f9565b5b6102ed86828701610216565b92505060406102fe86828701610268565b9150509250925092565b600082825260208201905092915050565b7f526573756c7461646f20696e76616c69646f202875736520302c2031206f752060008201527f3229000000000000000000000000000000000000000000000000000000000000602082015250565b6000610375602283610308565b915061038082610319565b604082019050919050565b600060208201905081810360008301526103a481610368565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600081519050919050565b60005b838110156104035780820151818401526020810190506103e8565b60008484015250505050565b600061041a826103da565b6104248185610308565b93506104348185602086016103e5565b61043d81610108565b840191505092915050565b60038110610459576104586103ab565b5b50565b600081905061046a82610448565b919050565b600061047a8261045c565b9050919050565b61048a8161046f565b82525050565b6000819050919050565b6104a381610490565b82525050565b600060808201905081810360008301526104c3818761040f565b905081810360208301526104d7818661040f565b90506104e66040830185610481565b6104f3606083018461049a565b9594505050505056fea2646970667358221220708a27bd6836dfd614a08193d9da252b90e462b051c19452e8e6af833ba511cd64736f6c634300081e0033",
}

// HistoricoABI is the input ABI used to generate the binding from.
// Deprecated: Use HistoricoMetaData.ABI instead.
var HistoricoABI = HistoricoMetaData.ABI

// HistoricoBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use HistoricoMetaData.Bin instead.
var HistoricoBin = HistoricoMetaData.Bin

// DeployHistorico deploys a new Ethereum contract, binding an instance of Historico to it.
func DeployHistorico(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Historico, error) {
	parsed, err := HistoricoMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(HistoricoBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Historico{HistoricoCaller: HistoricoCaller{contract: contract}, HistoricoTransactor: HistoricoTransactor{contract: contract}, HistoricoFilterer: HistoricoFilterer{contract: contract}}, nil
}

// Historico is an auto generated Go binding around an Ethereum contract.
type Historico struct {
	HistoricoCaller     // Read-only binding to the contract
	HistoricoTransactor // Write-only binding to the contract
	HistoricoFilterer   // Log filterer for contract events
}

// HistoricoCaller is an auto generated read-only Go binding around an Ethereum contract.
type HistoricoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoricoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HistoricoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoricoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HistoricoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HistoricoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HistoricoSession struct {
	Contract     *Historico        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HistoricoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HistoricoCallerSession struct {
	Contract *HistoricoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// HistoricoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HistoricoTransactorSession struct {
	Contract     *HistoricoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// HistoricoRaw is an auto generated low-level Go binding around an Ethereum contract.
type HistoricoRaw struct {
	Contract *Historico // Generic contract binding to access the raw methods on
}

// HistoricoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HistoricoCallerRaw struct {
	Contract *HistoricoCaller // Generic read-only contract binding to access the raw methods on
}

// HistoricoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HistoricoTransactorRaw struct {
	Contract *HistoricoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHistorico creates a new instance of Historico, bound to a specific deployed contract.
func NewHistorico(address common.Address, backend bind.ContractBackend) (*Historico, error) {
	contract, err := bindHistorico(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Historico{HistoricoCaller: HistoricoCaller{contract: contract}, HistoricoTransactor: HistoricoTransactor{contract: contract}, HistoricoFilterer: HistoricoFilterer{contract: contract}}, nil
}

// NewHistoricoCaller creates a new read-only instance of Historico, bound to a specific deployed contract.
func NewHistoricoCaller(address common.Address, caller bind.ContractCaller) (*HistoricoCaller, error) {
	contract, err := bindHistorico(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HistoricoCaller{contract: contract}, nil
}

// NewHistoricoTransactor creates a new write-only instance of Historico, bound to a specific deployed contract.
func NewHistoricoTransactor(address common.Address, transactor bind.ContractTransactor) (*HistoricoTransactor, error) {
	contract, err := bindHistorico(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HistoricoTransactor{contract: contract}, nil
}

// NewHistoricoFilterer creates a new log filterer instance of Historico, bound to a specific deployed contract.
func NewHistoricoFilterer(address common.Address, filterer bind.ContractFilterer) (*HistoricoFilterer, error) {
	contract, err := bindHistorico(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HistoricoFilterer{contract: contract}, nil
}

// bindHistorico binds a generic wrapper to an already deployed contract.
func bindHistorico(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HistoricoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Historico *HistoricoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Historico.Contract.HistoricoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Historico *HistoricoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Historico.Contract.HistoricoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Historico *HistoricoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Historico.Contract.HistoricoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Historico *HistoricoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Historico.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Historico *HistoricoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Historico.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Historico *HistoricoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Historico.Contract.contract.Transact(opts, method, params...)
}

// RegistrarPartida is a paid mutator transaction binding the contract method 0x33896cea.
//
// Solidity: function registrarPartida(string _j1, string _j2, uint8 _resultado) returns()
func (_Historico *HistoricoTransactor) RegistrarPartida(opts *bind.TransactOpts, _j1 string, _j2 string, _resultado uint8) (*types.Transaction, error) {
	return _Historico.contract.Transact(opts, "registrarPartida", _j1, _j2, _resultado)
}

// RegistrarPartida is a paid mutator transaction binding the contract method 0x33896cea.
//
// Solidity: function registrarPartida(string _j1, string _j2, uint8 _resultado) returns()
func (_Historico *HistoricoSession) RegistrarPartida(_j1 string, _j2 string, _resultado uint8) (*types.Transaction, error) {
	return _Historico.Contract.RegistrarPartida(&_Historico.TransactOpts, _j1, _j2, _resultado)
}

// RegistrarPartida is a paid mutator transaction binding the contract method 0x33896cea.
//
// Solidity: function registrarPartida(string _j1, string _j2, uint8 _resultado) returns()
func (_Historico *HistoricoTransactorSession) RegistrarPartida(_j1 string, _j2 string, _resultado uint8) (*types.Transaction, error) {
	return _Historico.Contract.RegistrarPartida(&_Historico.TransactOpts, _j1, _j2, _resultado)
}

// HistoricoNovaPartidaIterator is returned from FilterNovaPartida and is used to iterate over the raw logs and unpacked data for NovaPartida events raised by the Historico contract.
type HistoricoNovaPartidaIterator struct {
	Event *HistoricoNovaPartida // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *HistoricoNovaPartidaIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HistoricoNovaPartida)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(HistoricoNovaPartida)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *HistoricoNovaPartidaIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HistoricoNovaPartidaIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HistoricoNovaPartida represents a NovaPartida event raised by the Historico contract.
type HistoricoNovaPartida struct {
	Jogador1  string
	Jogador2  string
	Resultado uint8
	Data      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterNovaPartida is a free log retrieval operation binding the contract event 0x56cc181ac353627e55fa5eccf2e64c91f74aae5a157dc1cf0acadc1a97237b57.
//
// Solidity: event NovaPartida(string jogador1, string jogador2, uint8 resultado, uint256 data)
func (_Historico *HistoricoFilterer) FilterNovaPartida(opts *bind.FilterOpts) (*HistoricoNovaPartidaIterator, error) {

	logs, sub, err := _Historico.contract.FilterLogs(opts, "NovaPartida")
	if err != nil {
		return nil, err
	}
	return &HistoricoNovaPartidaIterator{contract: _Historico.contract, event: "NovaPartida", logs: logs, sub: sub}, nil
}

// WatchNovaPartida is a free log subscription operation binding the contract event 0x56cc181ac353627e55fa5eccf2e64c91f74aae5a157dc1cf0acadc1a97237b57.
//
// Solidity: event NovaPartida(string jogador1, string jogador2, uint8 resultado, uint256 data)
func (_Historico *HistoricoFilterer) WatchNovaPartida(opts *bind.WatchOpts, sink chan<- *HistoricoNovaPartida) (event.Subscription, error) {

	logs, sub, err := _Historico.contract.WatchLogs(opts, "NovaPartida")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HistoricoNovaPartida)
				if err := _Historico.contract.UnpackLog(event, "NovaPartida", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNovaPartida is a log parse operation binding the contract event 0x56cc181ac353627e55fa5eccf2e64c91f74aae5a157dc1cf0acadc1a97237b57.
//
// Solidity: event NovaPartida(string jogador1, string jogador2, uint8 resultado, uint256 data)
func (_Historico *HistoricoFilterer) ParseNovaPartida(log types.Log) (*HistoricoNovaPartida, error) {
	event := new(HistoricoNovaPartida)
	if err := _Historico.contract.UnpackLog(event, "NovaPartida", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
