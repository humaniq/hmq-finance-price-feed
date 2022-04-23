// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
)

// PriceDataMetaData contains all meta data concerning the PriceData contract.
var PriceDataMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"messageTimestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"currentValueTimestamp\",\"type\":\"uint64\"}],\"name\":\"NotWritten\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"currency\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"}],\"name\":\"Written\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"name\":\"getEth\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"currency\",\"type\":\"string\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"putEth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"currency\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"putPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610a6b806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80633d0f34da1461005157806347d4cc43146101a157806375615db51461024557806391b379b714610302575b600080fd5b61017a6004803603604081101561006757600080fd5b810190602081018135600160201b81111561008157600080fd5b82018360208201111561009357600080fd5b803590602001918460018302840111600160201b831117156100b457600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561010657600080fd5b82018360208201111561011857600080fd5b803590602001918460018302840111600160201b8311171561013957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610442945050505050565b6040805167ffffffffffffffff938416815291909216602082015281519081900390910190f35b61017a600480360360208110156101b757600080fd5b810190602081018135600160201b8111156101d157600080fd5b8201836020820111156101e357600080fd5b803590602001918460018302840111600160201b8311171561020457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610521945050505050565b6103006004803603606081101561025b57600080fd5b810190602081018135600160201b81111561027557600080fd5b82018360208201111561028757600080fd5b803590602001918460018302840111600160201b831117156102a857600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505067ffffffffffffffff83358116945060209093013590921691506105bf9050565b005b6103006004803603608081101561031857600080fd5b810190602081018135600160201b81111561033257600080fd5b82018360208201111561034457600080fd5b803590602001918460018302840111600160201b8311171561036557600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b8111156103b757600080fd5b8201836020820111156103c957600080fd5b803590602001918460018302840111600160201b831117156103ea57600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295505067ffffffffffffffff83358116945060209093013590921691506106bb9050565b600080600080856040518082805190602001908083835b602083106104785780518252601f199092019160209182019101610459565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101842088519094899450925082918401908083835b602083106104d45780518252601f1990920191602091820191016104b5565b51815160209384036101000a6000190180199092169116179052920194855250604051938490030190922054600160401b810467ffffffffffffffff908116999116975095505050505050565b600080600080846040518082805190602001908083835b602083106105575780518252601f199092019160209182019101610538565b51815160001960209485036101000a01908116901991909116179052920194855250604080519485900390910184206208aa8960eb1b8552600385015251928390036023019092205467ffffffffffffffff600160401b820481169650169350505050915091565b6040516391b379b760e01b815267ffffffffffffffff80841660448301528216606482015260806004820190815284516084830152845130926391b379b792879287928792918291602481019160a490910190602088019080838360005b8381101561063557818101518382015260200161061d565b50505050905090810190601f1680156106625780820380516001836020036101000a031916815260200191505b50838103825260038152602001806208aa8960eb1b81525060200195505050505050600060405180830381600087803b15801561069e57600080fd5b505af11580156106b2573d6000803e3d6000fd5b50505050505050565b600080856040518082805190602001908083835b602083106106ee5780518252601f1990920191602091820191016106cf565b51815160209384036101000a6000190180199092169116179052920194855250604051938490038101842088519094899450925082918401908083835b6020831061074a5780518252601f19909201916020918201910161072b565b51815160209384036101000a60001901801990921691161790529201948552506040519384900301909220805490935067ffffffffffffffff9081169085161191505080156107a25750428267ffffffffffffffff16105b156109e15760405180604001604052808367ffffffffffffffff1681526020018467ffffffffffffffff168152506000866040518082805190602001908083835b602083106108025780518252601f1990920191602091820191016107e3565b51815160209384036101000a60001901801990921691161790529201948552506040519384900381018420895190948a9450925082918401908083835b6020831061085e5780518252601f19909201916020918201910161083f565b51815160209384036101000a6000190180199092169116179052920194855250604080519485900382018520865181549784015167ffffffffffffffff1990981667ffffffffffffffff918216176fffffffffffffffff00000000000000001916600160401b98821698909802979097179055948816948401949094525050606080825287519082015286517f8ebbee21ac37f8d356bad48541fda45f095a69124012792b20e903b0d8a669dd9288928892889282918282019160808401919088019080838360005b8381101561093f578181015183820152602001610927565b50505050905090810190601f16801561096c5780820380516001836020036101000a031916815260200191505b50838103825285518152855160209182019187019080838360005b8381101561099f578181015183820152602001610987565b50505050905090810190601f1680156109cc5780820380516001836020036101000a031916815260200191505b509550505050505060405180910390a1610a2e565b80546040805142815267ffffffffffffffff808616602083015290921682820152517ff628b78179ba89052d9db1e70e518a5c9d1b5b1b8c8518171f66004e32edd9ce9181900360600190a15b505050505056fea26469706673582212208331addd276e68df5cd5896b7977afab41a4152a4ba0e4bfff6b35f1852ac14f64736f6c634300060a0033",
}

// PriceDataABI is the input ABI used to generate the binding from.
// Deprecated: Use PriceDataMetaData.ABI instead.
var PriceDataABI = PriceDataMetaData.ABI

// PriceDataBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PriceDataMetaData.Bin instead.
var PriceDataBin = PriceDataMetaData.Bin

// DeployPriceData deploys a new Ethereum contract, binding an instance of PriceData to it.
func DeployPriceData(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PriceData, error) {
	parsed, err := PriceDataMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PriceDataBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PriceData{PriceDataCaller: PriceDataCaller{contract: contract}, PriceDataTransactor: PriceDataTransactor{contract: contract}, PriceDataFilterer: PriceDataFilterer{contract: contract}}, nil
}

// PriceData is an auto generated Go binding around an Ethereum contract.
type PriceData struct {
	PriceDataCaller     // Read-only binding to the contract
	PriceDataTransactor // Write-only binding to the contract
	PriceDataFilterer   // Log filterer for contract events
}

// PriceDataCaller is an auto generated read-only Go binding around an Ethereum contract.
type PriceDataCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceDataTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PriceDataTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceDataFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PriceDataFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PriceDataSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PriceDataSession struct {
	Contract     *PriceData        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PriceDataCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PriceDataCallerSession struct {
	Contract *PriceDataCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// PriceDataTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PriceDataTransactorSession struct {
	Contract     *PriceDataTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PriceDataRaw is an auto generated low-level Go binding around an Ethereum contract.
type PriceDataRaw struct {
	Contract *PriceData // Generic contract binding to access the raw methods on
}

// PriceDataCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PriceDataCallerRaw struct {
	Contract *PriceDataCaller // Generic read-only contract binding to access the raw methods on
}

// PriceDataTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PriceDataTransactorRaw struct {
	Contract *PriceDataTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPriceData creates a new instance of PriceData, bound to a specific deployed contract.
func NewPriceData(address common.Address, backend bind.ContractBackend) (*PriceData, error) {
	contract, err := bindPriceData(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PriceData{PriceDataCaller: PriceDataCaller{contract: contract}, PriceDataTransactor: PriceDataTransactor{contract: contract}, PriceDataFilterer: PriceDataFilterer{contract: contract}}, nil
}

// NewPriceDataCaller creates a new read-only instance of PriceData, bound to a specific deployed contract.
func NewPriceDataCaller(address common.Address, caller bind.ContractCaller) (*PriceDataCaller, error) {
	contract, err := bindPriceData(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PriceDataCaller{contract: contract}, nil
}

// NewPriceDataTransactor creates a new write-only instance of PriceData, bound to a specific deployed contract.
func NewPriceDataTransactor(address common.Address, transactor bind.ContractTransactor) (*PriceDataTransactor, error) {
	contract, err := bindPriceData(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PriceDataTransactor{contract: contract}, nil
}

// NewPriceDataFilterer creates a new log filterer instance of PriceData, bound to a specific deployed contract.
func NewPriceDataFilterer(address common.Address, filterer bind.ContractFilterer) (*PriceDataFilterer, error) {
	contract, err := bindPriceData(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PriceDataFilterer{contract: contract}, nil
}

// bindPriceData binds a generic wrapper to an already deployed contract.
func bindPriceData(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PriceDataABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PriceData *PriceDataRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceData.Contract.PriceDataCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PriceData *PriceDataRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceData.Contract.PriceDataTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PriceData *PriceDataRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceData.Contract.PriceDataTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PriceData *PriceDataCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PriceData.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PriceData *PriceDataTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PriceData.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PriceData *PriceDataTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PriceData.Contract.contract.Transact(opts, method, params...)
}

// GetEth is a free data retrieval call binding the contract method 0x47d4cc43.
//
// Solidity: function getEth(string symbol) view returns(uint64, uint64)
func (_PriceData *PriceDataCaller) GetEth(opts *bind.CallOpts, symbol string) (uint64, uint64, error) {
	var out []interface{}
	err := _PriceData.contract.Call(opts, &out, "getEth", symbol)

	if err != nil {
		return *new(uint64), *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return out0, out1, err

}

// GetEth is a free data retrieval call binding the contract method 0x47d4cc43.
//
// Solidity: function getEth(string symbol) view returns(uint64, uint64)
func (_PriceData *PriceDataSession) GetEth(symbol string) (uint64, uint64, error) {
	return _PriceData.Contract.GetEth(&_PriceData.CallOpts, symbol)
}

// GetEth is a free data retrieval call binding the contract method 0x47d4cc43.
//
// Solidity: function getEth(string symbol) view returns(uint64, uint64)
func (_PriceData *PriceDataCallerSession) GetEth(symbol string) (uint64, uint64, error) {
	return _PriceData.Contract.GetEth(&_PriceData.CallOpts, symbol)
}

// GetPrice is a free data retrieval call binding the contract method 0x3d0f34da.
//
// Solidity: function getPrice(string symbol, string currency) view returns(uint64, uint64)
func (_PriceData *PriceDataCaller) GetPrice(opts *bind.CallOpts, symbol string, currency string) (uint64, uint64, error) {
	var out []interface{}
	err := _PriceData.contract.Call(opts, &out, "getPrice", symbol, currency)

	if err != nil {
		return *new(uint64), *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return out0, out1, err

}

// GetPrice is a free data retrieval call binding the contract method 0x3d0f34da.
//
// Solidity: function getPrice(string symbol, string currency) view returns(uint64, uint64)
func (_PriceData *PriceDataSession) GetPrice(symbol string, currency string) (uint64, uint64, error) {
	return _PriceData.Contract.GetPrice(&_PriceData.CallOpts, symbol, currency)
}

// GetPrice is a free data retrieval call binding the contract method 0x3d0f34da.
//
// Solidity: function getPrice(string symbol, string currency) view returns(uint64, uint64)
func (_PriceData *PriceDataCallerSession) GetPrice(symbol string, currency string) (uint64, uint64, error) {
	return _PriceData.Contract.GetPrice(&_PriceData.CallOpts, symbol, currency)
}

// PutEth is a paid mutator transaction binding the contract method 0x75615db5.
//
// Solidity: function putEth(string symbol, uint64 value, uint64 timestamp) returns()
func (_PriceData *PriceDataTransactor) PutEth(opts *bind.TransactOpts, symbol string, value uint64, timestamp uint64) (*types.Transaction, error) {
	return _PriceData.contract.Transact(opts, "putEth", symbol, value, timestamp)
}

// PutEth is a paid mutator transaction binding the contract method 0x75615db5.
//
// Solidity: function putEth(string symbol, uint64 value, uint64 timestamp) returns()
func (_PriceData *PriceDataSession) PutEth(symbol string, value uint64, timestamp uint64) (*types.Transaction, error) {
	return _PriceData.Contract.PutEth(&_PriceData.TransactOpts, symbol, value, timestamp)
}

// PutEth is a paid mutator transaction binding the contract method 0x75615db5.
//
// Solidity: function putEth(string symbol, uint64 value, uint64 timestamp) returns()
func (_PriceData *PriceDataTransactorSession) PutEth(symbol string, value uint64, timestamp uint64) (*types.Transaction, error) {
	return _PriceData.Contract.PutEth(&_PriceData.TransactOpts, symbol, value, timestamp)
}

// PutPrice is a paid mutator transaction binding the contract method 0x91b379b7.
//
// Solidity: function putPrice(string symbol, string currency, uint64 value, uint64 timestamp) returns()
func (_PriceData *PriceDataTransactor) PutPrice(opts *bind.TransactOpts, symbol string, currency string, value uint64, timestamp uint64) (*types.Transaction, error) {
	return _PriceData.contract.Transact(opts, "putPrice", symbol, currency, value, timestamp)
}

// PutPrice is a paid mutator transaction binding the contract method 0x91b379b7.
//
// Solidity: function putPrice(string symbol, string currency, uint64 value, uint64 timestamp) returns()
func (_PriceData *PriceDataSession) PutPrice(symbol string, currency string, value uint64, timestamp uint64) (*types.Transaction, error) {
	return _PriceData.Contract.PutPrice(&_PriceData.TransactOpts, symbol, currency, value, timestamp)
}

// PutPrice is a paid mutator transaction binding the contract method 0x91b379b7.
//
// Solidity: function putPrice(string symbol, string currency, uint64 value, uint64 timestamp) returns()
func (_PriceData *PriceDataTransactorSession) PutPrice(symbol string, currency string, value uint64, timestamp uint64) (*types.Transaction, error) {
	return _PriceData.Contract.PutPrice(&_PriceData.TransactOpts, symbol, currency, value, timestamp)
}

// PriceDataNotWrittenIterator is returned from FilterNotWritten and is used to iterate over the raw logs and unpacked data for NotWritten events raised by the PriceData contract.
type PriceDataNotWrittenIterator struct {
	Event *PriceDataNotWritten // Event containing the contract specifics and raw log

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
func (it *PriceDataNotWrittenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceDataNotWritten)
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
		it.Event = new(PriceDataNotWritten)
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
func (it *PriceDataNotWrittenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceDataNotWrittenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceDataNotWritten represents a NotWritten event raised by the PriceData contract.
type PriceDataNotWritten struct {
	BlockTimestamp        *big.Int
	MessageTimestamp      uint64
	CurrentValueTimestamp uint64
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterNotWritten is a free log retrieval operation binding the contract event 0xf628b78179ba89052d9db1e70e518a5c9d1b5b1b8c8518171f66004e32edd9ce.
//
// Solidity: event NotWritten(uint256 blockTimestamp, uint64 messageTimestamp, uint64 currentValueTimestamp)
func (_PriceData *PriceDataFilterer) FilterNotWritten(opts *bind.FilterOpts) (*PriceDataNotWrittenIterator, error) {

	logs, sub, err := _PriceData.contract.FilterLogs(opts, "NotWritten")
	if err != nil {
		return nil, err
	}
	return &PriceDataNotWrittenIterator{contract: _PriceData.contract, event: "NotWritten", logs: logs, sub: sub}, nil
}

// WatchNotWritten is a free log subscription operation binding the contract event 0xf628b78179ba89052d9db1e70e518a5c9d1b5b1b8c8518171f66004e32edd9ce.
//
// Solidity: event NotWritten(uint256 blockTimestamp, uint64 messageTimestamp, uint64 currentValueTimestamp)
func (_PriceData *PriceDataFilterer) WatchNotWritten(opts *bind.WatchOpts, sink chan<- *PriceDataNotWritten) (event.Subscription, error) {

	logs, sub, err := _PriceData.contract.WatchLogs(opts, "NotWritten")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceDataNotWritten)
				if err := _PriceData.contract.UnpackLog(event, "NotWritten", log); err != nil {
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

// ParseNotWritten is a log parse operation binding the contract event 0xf628b78179ba89052d9db1e70e518a5c9d1b5b1b8c8518171f66004e32edd9ce.
//
// Solidity: event NotWritten(uint256 blockTimestamp, uint64 messageTimestamp, uint64 currentValueTimestamp)
func (_PriceData *PriceDataFilterer) ParseNotWritten(log types.Log) (*PriceDataNotWritten, error) {
	event := new(PriceDataNotWritten)
	if err := _PriceData.contract.UnpackLog(event, "NotWritten", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PriceDataWrittenIterator is returned from FilterWritten and is used to iterate over the raw logs and unpacked data for Written events raised by the PriceData contract.
type PriceDataWrittenIterator struct {
	Event *PriceDataWritten // Event containing the contract specifics and raw log

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
func (it *PriceDataWrittenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PriceDataWritten)
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
		it.Event = new(PriceDataWritten)
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
func (it *PriceDataWrittenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PriceDataWrittenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PriceDataWritten represents a Written event raised by the PriceData contract.
type PriceDataWritten struct {
	Symbol   string
	Currency string
	Value    uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWritten is a free log retrieval operation binding the contract event 0x8ebbee21ac37f8d356bad48541fda45f095a69124012792b20e903b0d8a669dd.
//
// Solidity: event Written(string symbol, string currency, uint64 value)
func (_PriceData *PriceDataFilterer) FilterWritten(opts *bind.FilterOpts) (*PriceDataWrittenIterator, error) {

	logs, sub, err := _PriceData.contract.FilterLogs(opts, "Written")
	if err != nil {
		return nil, err
	}
	return &PriceDataWrittenIterator{contract: _PriceData.contract, event: "Written", logs: logs, sub: sub}, nil
}

// WatchWritten is a free log subscription operation binding the contract event 0x8ebbee21ac37f8d356bad48541fda45f095a69124012792b20e903b0d8a669dd.
//
// Solidity: event Written(string symbol, string currency, uint64 value)
func (_PriceData *PriceDataFilterer) WatchWritten(opts *bind.WatchOpts, sink chan<- *PriceDataWritten) (event.Subscription, error) {

	logs, sub, err := _PriceData.contract.WatchLogs(opts, "Written")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PriceDataWritten)
				if err := _PriceData.contract.UnpackLog(event, "Written", log); err != nil {
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

// ParseWritten is a log parse operation binding the contract event 0x8ebbee21ac37f8d356bad48541fda45f095a69124012792b20e903b0d8a669dd.
//
// Solidity: event Written(string symbol, string currency, uint64 value)
func (_PriceData *PriceDataFilterer) ParseWritten(log types.Log) (*PriceDataWritten, error) {
	event := new(PriceDataWritten)
	if err := _PriceData.contract.UnpackLog(event, "Written", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
