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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"vencedor\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"perdedor\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"placar\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"data\",\"type\":\"uint256\"}],\"name\":\"NovaPartida\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_vencedor\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_perdedor\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_placar\",\"type\":\"string\"}],\"name\":\"registrarResultado\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b506103b78061001f6000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063ca8dd47714610030575b600080fd5b61004a600480360381019061004591906101e8565b61004c565b005b7fb9e7345d75df19b69bead6028780b38813ca9fd034ce46ac9bf34185c4b5346f838383426040516100819493929190610327565b60405180910390a1505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6100f5826100ac565b810181811067ffffffffffffffff82111715610114576101136100bd565b5b80604052505050565b600061012761008e565b905061013382826100ec565b919050565b600067ffffffffffffffff821115610153576101526100bd565b5b61015c826100ac565b9050602081019050919050565b82818337600083830152505050565b600061018b61018684610138565b61011d565b9050828152602081018484840111156101a7576101a66100a7565b5b6101b2848285610169565b509392505050565b600082601f8301126101cf576101ce6100a2565b5b81356101df848260208601610178565b91505092915050565b60008060006060848603121561020157610200610098565b5b600084013567ffffffffffffffff81111561021f5761021e61009d565b5b61022b868287016101ba565b935050602084013567ffffffffffffffff81111561024c5761024b61009d565b5b610258868287016101ba565b925050604084013567ffffffffffffffff8111156102795761027861009d565b5b610285868287016101ba565b9150509250925092565b600081519050919050565b600082825260208201905092915050565b60005b838110156102c95780820151818401526020810190506102ae565b60008484015250505050565b60006102e08261028f565b6102ea818561029a565b93506102fa8185602086016102ab565b610303816100ac565b840191505092915050565b6000819050919050565b6103218161030e565b82525050565b6000608082019050818103600083015261034181876102d5565b9050818103602083015261035581866102d5565b9050818103604083015261036981856102d5565b90506103786060830184610318565b9594505050505056fea26469706673582212201061049cfe98983557a4a10d34949ff275a2e78a4a43587b8b45877441fcef5064736f6c634300081e0033",
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

// RegistrarResultado is a paid mutator transaction binding the contract method 0xca8dd477.
//
// Solidity: function registrarResultado(string _vencedor, string _perdedor, string _placar) returns()
func (_Historico *HistoricoTransactor) RegistrarResultado(opts *bind.TransactOpts, _vencedor string, _perdedor string, _placar string) (*types.Transaction, error) {
	return _Historico.contract.Transact(opts, "registrarResultado", _vencedor, _perdedor, _placar)
}

// RegistrarResultado is a paid mutator transaction binding the contract method 0xca8dd477.
//
// Solidity: function registrarResultado(string _vencedor, string _perdedor, string _placar) returns()
func (_Historico *HistoricoSession) RegistrarResultado(_vencedor string, _perdedor string, _placar string) (*types.Transaction, error) {
	return _Historico.Contract.RegistrarResultado(&_Historico.TransactOpts, _vencedor, _perdedor, _placar)
}

// RegistrarResultado is a paid mutator transaction binding the contract method 0xca8dd477.
//
// Solidity: function registrarResultado(string _vencedor, string _perdedor, string _placar) returns()
func (_Historico *HistoricoTransactorSession) RegistrarResultado(_vencedor string, _perdedor string, _placar string) (*types.Transaction, error) {
	return _Historico.Contract.RegistrarResultado(&_Historico.TransactOpts, _vencedor, _perdedor, _placar)
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
	Vencedor string
	Perdedor string
	Placar   string
	Data     *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNovaPartida is a free log retrieval operation binding the contract event 0xb9e7345d75df19b69bead6028780b38813ca9fd034ce46ac9bf34185c4b5346f.
//
// Solidity: event NovaPartida(string vencedor, string perdedor, string placar, uint256 data)
func (_Historico *HistoricoFilterer) FilterNovaPartida(opts *bind.FilterOpts) (*HistoricoNovaPartidaIterator, error) {

	logs, sub, err := _Historico.contract.FilterLogs(opts, "NovaPartida")
	if err != nil {
		return nil, err
	}
	return &HistoricoNovaPartidaIterator{contract: _Historico.contract, event: "NovaPartida", logs: logs, sub: sub}, nil
}

// WatchNovaPartida is a free log subscription operation binding the contract event 0xb9e7345d75df19b69bead6028780b38813ca9fd034ce46ac9bf34185c4b5346f.
//
// Solidity: event NovaPartida(string vencedor, string perdedor, string placar, uint256 data)
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

// ParseNovaPartida is a log parse operation binding the contract event 0xb9e7345d75df19b69bead6028780b38813ca9fd034ce46ac9bf34185c4b5346f.
//
// Solidity: event NovaPartida(string vencedor, string perdedor, string placar, uint256 data)
func (_Historico *HistoricoFilterer) ParseNovaPartida(log types.Log) (*HistoricoNovaPartida, error) {
	event := new(HistoricoNovaPartida)
	if err := _Historico.contract.UnpackLog(event, "NovaPartida", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
