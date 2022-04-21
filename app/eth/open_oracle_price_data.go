// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

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

// EthMetaData contains all meta data concerning the Eth contract.
var EthMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"priorTimestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"messageTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockTimestamp\",\"type\":\"uint256\"}],\"name\":\"NotWritten\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bts\",\"type\":\"bytes\"}],\"name\":\"RequestGot\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"}],\"name\":\"SourceGot\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"}],\"name\":\"Write\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"value\",\"type\":\"uint64\"}],\"name\":\"WriteGot\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"get\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"source\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"getPrice\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"put\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"source\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610cec806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806338636e9a14610051578063482a61931461018457806376977a3a146102c9578063fc2525ab14610363575b600080fd5b61010f6004803603604081101561006757600080fd5b810190602081018135600160201b81111561008157600080fd5b82018360208201111561009357600080fd5b803590602001918460018302840111600160201b831117156100b457600080fd5b919390929091602081019035600160201b8111156100d157600080fd5b8201836020820111156100e357600080fd5b803590602001918460018302840111600160201b8311171561010457600080fd5b509092509050610407565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610149578181015183820152602001610131565b50505050905090810190601f1680156101765780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102ad6004803603604081101561019a57600080fd5b810190602081018135600160201b8111156101b457600080fd5b8201836020820111156101c657600080fd5b803590602001918460018302840111600160201b831117156101e757600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295949360208101935035915050600160201b81111561023957600080fd5b82018360208201111561024b57600080fd5b803590602001918460018302840111600160201b8311171561026c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610574945050505050565b604080516001600160a01b039092168252519081900360200190f35b610347600480360360408110156102df57600080fd5b6001600160a01b038235169190810190604081016020820135600160201b81111561030957600080fd5b82018360208201111561031b57600080fd5b803590602001918460018302840111600160201b8311171561033c57600080fd5b509092509050610661565b604080516001600160401b039092168252519081900360200190f35b6103e16004803603604081101561037957600080fd5b6001600160a01b038235169190810190604081016020820135600160201b8111156103a357600080fd5b8201836020820111156103b557600080fd5b803590602001918460018302840111600160201b831117156103d657600080fd5b5090925090506106ba565b604080516001600160401b03938416815291909216602082015281519081900390910190f35b60607fb23429a2c11fa29a7445d8cb8bf9d83255e53490725b5741efcdb93235dbbf0d858560405180806020018281038252848482818152602001925080828437600083820152604051601f909101601f19169092018290039550909350505050a16000806060600061047c89898989610729565b9350935093509350836001600160a01b03167f9f4e154620987b6ba6059ab67ef7491a565c9dd7c5ae25d66f7495bd7e78ab0b8385846040518080602001846001600160401b03166001600160401b03168152602001836001600160401b03166001600160401b03168152602001828103825285818151815260200191508051906020019080838360005b8381101561051f578181015183820152602001610507565b50505050905090810190601f16801561054c5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a2610567848484846109f6565b9998505050505050505050565b60008060008084806020019051606081101561058f57600080fd5b5080516020808301516040938401518a518b84012085517f19457468657265756d205369676e6564204d6573736167653a0a33320000000081860152603c8082019290925286518082039092018252605c81018088528251928601929092206000909252607c810180885282905260ff8316609c82015260bc810186905260dc8101849052955194985091965094509260019260fc8083019392601f198301929081900390910190855afa15801561064b573d6000803e3d6000fd5b5050604051601f19015198975050505050505050565b6001600160a01b0383166000908152602081905260408082209051849084908083838082843791909101948552505060405192839003602001909220546001600160401b03600160401b90910416925050509392505050565b6000806000806000876001600160a01b03166001600160a01b0316815260200190815260200160002085856040518083838082843791909101948552505060405192839003602001909220546001600160401b038082169650600160401b909104169350505050935093915050565b60008060606000806107a489898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080516020601f8d018190048102820181019092528b815292508b91508a908190840183828082843760009201919091525061057492505050565b905060606000606060008c8c60808110156107be57600080fd5b810190602081018135600160201b8111156107d857600080fd5b8201836020820111156107ea57600080fd5b803590602001918460018302840111600160201b8311171561080b57600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092956001600160401b03853516959094909350604081019250602001359050600160201b81111561086e57600080fd5b82018360208201111561088057600080fd5b803590602001918460018302840111600160201b831117156108a157600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080516570726963657360d01b60208083019190915282516006818403018152602683019093528251928101929092208a519a9e50989c50939a5093356001600160401b0316985095968b96506046909201945084935050908401908083835b602083106109555780518252601f199092019160209182019101610936565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405160208183030381529060405280519060200120146109e2576040805162461bcd60e51b815260206004820152601d60248201527f4b696e64206f662064617461206d757374206265202770726963657327000000604482015290519081900360640190fd5b939c919b5099509197509095505050505050565b6001600160a01b038416600090815260208181526040808320905185516060949387929182918401908083835b60208310610a425780518252601f199092019160209182019101610a23565b51815160209384036101000a6000190180199092169116179052920194855250604051938490030190922080549093506001600160401b03908116908816119150508015610a9c575042610e1001856001600160401b0316105b8015610ab057506001600160a01b03861615155b15610c60576040518060400160405280866001600160401b03168152602001846001600160401b0316815250600080886001600160a01b03166001600160a01b03168152602001908152602001600020856040518082805190602001908083835b60208310610b305780518252601f199092019160209182019101610b11565b51815160209384036101000a6000190180199092169116179052920194855250604080519485900382018520865181549784015167ffffffffffffffff199098166001600160401b03918216176fffffffffffffffff00000000000000001916600160401b988216989098029790971790558a861685830152948816948401949094525050606080825286519082015285516001600160a01b038916927f4d3f5aa96531b83f5389343ecd20cd8ac1fba33b64634c1b547a4d85d31540d39288928a92899291829160808301919087019080838360005b83811015610c1f578181015183820152602001610c07565b50505050905090810190601f168015610c4c5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a2610cac565b8054604080516001600160401b03928316815291871660208301524282820152517f7d218dba44a461fb2d7b5fe792128439313d3c48c86d4c3e4981a8eaca831a769181900360600190a15b509194935050505056fea26469706673582212204236d8544b8d67125a061355d95c22901cc38ebb23b5ee8ccb17b0011e9bbd7f64736f6c634300060a0033",
}

// EthABI is the input ABI used to generate the binding from.
// Deprecated: Use EthMetaData.ABI instead.
var EthABI = EthMetaData.ABI

// EthBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EthMetaData.Bin instead.
var EthBin = EthMetaData.Bin

// DeployEth deploys a new Ethereum contract, binding an instance of Eth to it.
func DeployEth(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Eth, error) {
	parsed, err := EthMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EthBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Eth{EthCaller: EthCaller{contract: contract}, EthTransactor: EthTransactor{contract: contract}, EthFilterer: EthFilterer{contract: contract}}, nil
}

// Eth is an auto generated Go binding around an Ethereum contract.
type Eth struct {
	EthCaller     // Read-only binding to the contract
	EthTransactor // Write-only binding to the contract
	EthFilterer   // Log filterer for contract events
}

// EthCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthSession struct {
	Contract     *Eth              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthCallerSession struct {
	Contract *EthCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthTransactorSession struct {
	Contract     *EthTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthRaw struct {
	Contract *Eth // Generic contract binding to access the raw methods on
}

// EthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthCallerRaw struct {
	Contract *EthCaller // Generic read-only contract binding to access the raw methods on
}

// EthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthTransactorRaw struct {
	Contract *EthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEth creates a new instance of Eth, bound to a specific deployed contract.
func NewEth(address common.Address, backend bind.ContractBackend) (*Eth, error) {
	contract, err := bindEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Eth{EthCaller: EthCaller{contract: contract}, EthTransactor: EthTransactor{contract: contract}, EthFilterer: EthFilterer{contract: contract}}, nil
}

// NewEthCaller creates a new read-only instance of Eth, bound to a specific deployed contract.
func NewEthCaller(address common.Address, caller bind.ContractCaller) (*EthCaller, error) {
	contract, err := bindEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthCaller{contract: contract}, nil
}

// NewEthTransactor creates a new write-only instance of Eth, bound to a specific deployed contract.
func NewEthTransactor(address common.Address, transactor bind.ContractTransactor) (*EthTransactor, error) {
	contract, err := bindEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthTransactor{contract: contract}, nil
}

// NewEthFilterer creates a new log filterer instance of Eth, bound to a specific deployed contract.
func NewEthFilterer(address common.Address, filterer bind.ContractFilterer) (*EthFilterer, error) {
	contract, err := bindEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthFilterer{contract: contract}, nil
}

// bindEth binds a generic wrapper to an already deployed contract.
func bindEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth *EthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eth.Contract.EthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth *EthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth.Contract.EthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth *EthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth.Contract.EthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth *EthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth *EthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth *EthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0xfc2525ab.
//
// Solidity: function get(address source, string key) view returns(uint64, uint64)
func (_Eth *EthCaller) Get(opts *bind.CallOpts, source common.Address, key string) (uint64, uint64, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "get", source, key)

	if err != nil {
		return *new(uint64), *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return out0, out1, err

}

// Get is a free data retrieval call binding the contract method 0xfc2525ab.
//
// Solidity: function get(address source, string key) view returns(uint64, uint64)
func (_Eth *EthSession) Get(source common.Address, key string) (uint64, uint64, error) {
	return _Eth.Contract.Get(&_Eth.CallOpts, source, key)
}

// Get is a free data retrieval call binding the contract method 0xfc2525ab.
//
// Solidity: function get(address source, string key) view returns(uint64, uint64)
func (_Eth *EthCallerSession) Get(source common.Address, key string) (uint64, uint64, error) {
	return _Eth.Contract.Get(&_Eth.CallOpts, source, key)
}

// GetPrice is a free data retrieval call binding the contract method 0x76977a3a.
//
// Solidity: function getPrice(address source, string key) view returns(uint64)
func (_Eth *EthCaller) GetPrice(opts *bind.CallOpts, source common.Address, key string) (uint64, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "getPrice", source, key)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetPrice is a free data retrieval call binding the contract method 0x76977a3a.
//
// Solidity: function getPrice(address source, string key) view returns(uint64)
func (_Eth *EthSession) GetPrice(source common.Address, key string) (uint64, error) {
	return _Eth.Contract.GetPrice(&_Eth.CallOpts, source, key)
}

// GetPrice is a free data retrieval call binding the contract method 0x76977a3a.
//
// Solidity: function getPrice(address source, string key) view returns(uint64)
func (_Eth *EthCallerSession) GetPrice(source common.Address, key string) (uint64, error) {
	return _Eth.Contract.GetPrice(&_Eth.CallOpts, source, key)
}

// Source is a free data retrieval call binding the contract method 0x482a6193.
//
// Solidity: function source(bytes message, bytes signature) pure returns(address)
func (_Eth *EthCaller) Source(opts *bind.CallOpts, message []byte, signature []byte) (common.Address, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "source", message, signature)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Source is a free data retrieval call binding the contract method 0x482a6193.
//
// Solidity: function source(bytes message, bytes signature) pure returns(address)
func (_Eth *EthSession) Source(message []byte, signature []byte) (common.Address, error) {
	return _Eth.Contract.Source(&_Eth.CallOpts, message, signature)
}

// Source is a free data retrieval call binding the contract method 0x482a6193.
//
// Solidity: function source(bytes message, bytes signature) pure returns(address)
func (_Eth *EthCallerSession) Source(message []byte, signature []byte) (common.Address, error) {
	return _Eth.Contract.Source(&_Eth.CallOpts, message, signature)
}

// Put is a paid mutator transaction binding the contract method 0x38636e9a.
//
// Solidity: function put(bytes message, bytes signature) returns(string)
func (_Eth *EthTransactor) Put(opts *bind.TransactOpts, message []byte, signature []byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "put", message, signature)
}

// Put is a paid mutator transaction binding the contract method 0x38636e9a.
//
// Solidity: function put(bytes message, bytes signature) returns(string)
func (_Eth *EthSession) Put(message []byte, signature []byte) (*types.Transaction, error) {
	return _Eth.Contract.Put(&_Eth.TransactOpts, message, signature)
}

// Put is a paid mutator transaction binding the contract method 0x38636e9a.
//
// Solidity: function put(bytes message, bytes signature) returns(string)
func (_Eth *EthTransactorSession) Put(message []byte, signature []byte) (*types.Transaction, error) {
	return _Eth.Contract.Put(&_Eth.TransactOpts, message, signature)
}

// EthNotWrittenIterator is returned from FilterNotWritten and is used to iterate over the raw logs and unpacked data for NotWritten events raised by the Eth contract.
type EthNotWrittenIterator struct {
	Event *EthNotWritten // Event containing the contract specifics and raw log

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
func (it *EthNotWrittenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthNotWritten)
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
		it.Event = new(EthNotWritten)
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
func (it *EthNotWrittenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthNotWrittenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthNotWritten represents a NotWritten event raised by the Eth contract.
type EthNotWritten struct {
	PriorTimestamp   uint64
	MessageTimestamp *big.Int
	BlockTimestamp   *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNotWritten is a free log retrieval operation binding the contract event 0x7d218dba44a461fb2d7b5fe792128439313d3c48c86d4c3e4981a8eaca831a76.
//
// Solidity: event NotWritten(uint64 priorTimestamp, uint256 messageTimestamp, uint256 blockTimestamp)
func (_Eth *EthFilterer) FilterNotWritten(opts *bind.FilterOpts) (*EthNotWrittenIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "NotWritten")
	if err != nil {
		return nil, err
	}
	return &EthNotWrittenIterator{contract: _Eth.contract, event: "NotWritten", logs: logs, sub: sub}, nil
}

// WatchNotWritten is a free log subscription operation binding the contract event 0x7d218dba44a461fb2d7b5fe792128439313d3c48c86d4c3e4981a8eaca831a76.
//
// Solidity: event NotWritten(uint64 priorTimestamp, uint256 messageTimestamp, uint256 blockTimestamp)
func (_Eth *EthFilterer) WatchNotWritten(opts *bind.WatchOpts, sink chan<- *EthNotWritten) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "NotWritten")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthNotWritten)
				if err := _Eth.contract.UnpackLog(event, "NotWritten", log); err != nil {
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

// ParseNotWritten is a log parse operation binding the contract event 0x7d218dba44a461fb2d7b5fe792128439313d3c48c86d4c3e4981a8eaca831a76.
//
// Solidity: event NotWritten(uint64 priorTimestamp, uint256 messageTimestamp, uint256 blockTimestamp)
func (_Eth *EthFilterer) ParseNotWritten(log types.Log) (*EthNotWritten, error) {
	event := new(EthNotWritten)
	if err := _Eth.contract.UnpackLog(event, "NotWritten", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthRequestGotIterator is returned from FilterRequestGot and is used to iterate over the raw logs and unpacked data for RequestGot events raised by the Eth contract.
type EthRequestGotIterator struct {
	Event *EthRequestGot // Event containing the contract specifics and raw log

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
func (it *EthRequestGotIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthRequestGot)
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
		it.Event = new(EthRequestGot)
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
func (it *EthRequestGotIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthRequestGotIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthRequestGot represents a RequestGot event raised by the Eth contract.
type EthRequestGot struct {
	Bts []byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRequestGot is a free log retrieval operation binding the contract event 0xb23429a2c11fa29a7445d8cb8bf9d83255e53490725b5741efcdb93235dbbf0d.
//
// Solidity: event RequestGot(bytes bts)
func (_Eth *EthFilterer) FilterRequestGot(opts *bind.FilterOpts) (*EthRequestGotIterator, error) {

	logs, sub, err := _Eth.contract.FilterLogs(opts, "RequestGot")
	if err != nil {
		return nil, err
	}
	return &EthRequestGotIterator{contract: _Eth.contract, event: "RequestGot", logs: logs, sub: sub}, nil
}

// WatchRequestGot is a free log subscription operation binding the contract event 0xb23429a2c11fa29a7445d8cb8bf9d83255e53490725b5741efcdb93235dbbf0d.
//
// Solidity: event RequestGot(bytes bts)
func (_Eth *EthFilterer) WatchRequestGot(opts *bind.WatchOpts, sink chan<- *EthRequestGot) (event.Subscription, error) {

	logs, sub, err := _Eth.contract.WatchLogs(opts, "RequestGot")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthRequestGot)
				if err := _Eth.contract.UnpackLog(event, "RequestGot", log); err != nil {
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

// ParseRequestGot is a log parse operation binding the contract event 0xb23429a2c11fa29a7445d8cb8bf9d83255e53490725b5741efcdb93235dbbf0d.
//
// Solidity: event RequestGot(bytes bts)
func (_Eth *EthFilterer) ParseRequestGot(log types.Log) (*EthRequestGot, error) {
	event := new(EthRequestGot)
	if err := _Eth.contract.UnpackLog(event, "RequestGot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthSourceGotIterator is returned from FilterSourceGot and is used to iterate over the raw logs and unpacked data for SourceGot events raised by the Eth contract.
type EthSourceGotIterator struct {
	Event *EthSourceGot // Event containing the contract specifics and raw log

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
func (it *EthSourceGotIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthSourceGot)
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
		it.Event = new(EthSourceGot)
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
func (it *EthSourceGotIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthSourceGotIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthSourceGot represents a SourceGot event raised by the Eth contract.
type EthSourceGot struct {
	Source common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSourceGot is a free log retrieval operation binding the contract event 0xe6130684ab03fd6501681914b2548545e54d1fa773213d2ab65af85235d1681d.
//
// Solidity: event SourceGot(address indexed source)
func (_Eth *EthFilterer) FilterSourceGot(opts *bind.FilterOpts, source []common.Address) (*EthSourceGotIterator, error) {

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "SourceGot", sourceRule)
	if err != nil {
		return nil, err
	}
	return &EthSourceGotIterator{contract: _Eth.contract, event: "SourceGot", logs: logs, sub: sub}, nil
}

// WatchSourceGot is a free log subscription operation binding the contract event 0xe6130684ab03fd6501681914b2548545e54d1fa773213d2ab65af85235d1681d.
//
// Solidity: event SourceGot(address indexed source)
func (_Eth *EthFilterer) WatchSourceGot(opts *bind.WatchOpts, sink chan<- *EthSourceGot, source []common.Address) (event.Subscription, error) {

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "SourceGot", sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthSourceGot)
				if err := _Eth.contract.UnpackLog(event, "SourceGot", log); err != nil {
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

// ParseSourceGot is a log parse operation binding the contract event 0xe6130684ab03fd6501681914b2548545e54d1fa773213d2ab65af85235d1681d.
//
// Solidity: event SourceGot(address indexed source)
func (_Eth *EthFilterer) ParseSourceGot(log types.Log) (*EthSourceGot, error) {
	event := new(EthSourceGot)
	if err := _Eth.contract.UnpackLog(event, "SourceGot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthWriteIterator is returned from FilterWrite and is used to iterate over the raw logs and unpacked data for Write events raised by the Eth contract.
type EthWriteIterator struct {
	Event *EthWrite // Event containing the contract specifics and raw log

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
func (it *EthWriteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthWrite)
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
		it.Event = new(EthWrite)
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
func (it *EthWriteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthWriteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthWrite represents a Write event raised by the Eth contract.
type EthWrite struct {
	Source    common.Address
	Key       string
	Timestamp uint64
	Value     uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWrite is a free log retrieval operation binding the contract event 0x4d3f5aa96531b83f5389343ecd20cd8ac1fba33b64634c1b547a4d85d31540d3.
//
// Solidity: event Write(address indexed source, string key, uint64 timestamp, uint64 value)
func (_Eth *EthFilterer) FilterWrite(opts *bind.FilterOpts, source []common.Address) (*EthWriteIterator, error) {

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "Write", sourceRule)
	if err != nil {
		return nil, err
	}
	return &EthWriteIterator{contract: _Eth.contract, event: "Write", logs: logs, sub: sub}, nil
}

// WatchWrite is a free log subscription operation binding the contract event 0x4d3f5aa96531b83f5389343ecd20cd8ac1fba33b64634c1b547a4d85d31540d3.
//
// Solidity: event Write(address indexed source, string key, uint64 timestamp, uint64 value)
func (_Eth *EthFilterer) WatchWrite(opts *bind.WatchOpts, sink chan<- *EthWrite, source []common.Address) (event.Subscription, error) {

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "Write", sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthWrite)
				if err := _Eth.contract.UnpackLog(event, "Write", log); err != nil {
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

// ParseWrite is a log parse operation binding the contract event 0x4d3f5aa96531b83f5389343ecd20cd8ac1fba33b64634c1b547a4d85d31540d3.
//
// Solidity: event Write(address indexed source, string key, uint64 timestamp, uint64 value)
func (_Eth *EthFilterer) ParseWrite(log types.Log) (*EthWrite, error) {
	event := new(EthWrite)
	if err := _Eth.contract.UnpackLog(event, "Write", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthWriteGotIterator is returned from FilterWriteGot and is used to iterate over the raw logs and unpacked data for WriteGot events raised by the Eth contract.
type EthWriteGotIterator struct {
	Event *EthWriteGot // Event containing the contract specifics and raw log

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
func (it *EthWriteGotIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthWriteGot)
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
		it.Event = new(EthWriteGot)
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
func (it *EthWriteGotIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthWriteGotIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthWriteGot represents a WriteGot event raised by the Eth contract.
type EthWriteGot struct {
	Source    common.Address
	Key       string
	Timestamp uint64
	Value     uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWriteGot is a free log retrieval operation binding the contract event 0x9f4e154620987b6ba6059ab67ef7491a565c9dd7c5ae25d66f7495bd7e78ab0b.
//
// Solidity: event WriteGot(address indexed source, string key, uint64 timestamp, uint64 value)
func (_Eth *EthFilterer) FilterWriteGot(opts *bind.FilterOpts, source []common.Address) (*EthWriteGotIterator, error) {

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "WriteGot", sourceRule)
	if err != nil {
		return nil, err
	}
	return &EthWriteGotIterator{contract: _Eth.contract, event: "WriteGot", logs: logs, sub: sub}, nil
}

// WatchWriteGot is a free log subscription operation binding the contract event 0x9f4e154620987b6ba6059ab67ef7491a565c9dd7c5ae25d66f7495bd7e78ab0b.
//
// Solidity: event WriteGot(address indexed source, string key, uint64 timestamp, uint64 value)
func (_Eth *EthFilterer) WatchWriteGot(opts *bind.WatchOpts, sink chan<- *EthWriteGot, source []common.Address) (event.Subscription, error) {

	var sourceRule []interface{}
	for _, sourceItem := range source {
		sourceRule = append(sourceRule, sourceItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "WriteGot", sourceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthWriteGot)
				if err := _Eth.contract.UnpackLog(event, "WriteGot", log); err != nil {
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

// ParseWriteGot is a log parse operation binding the contract event 0x9f4e154620987b6ba6059ab67ef7491a565c9dd7c5ae25d66f7495bd7e78ab0b.
//
// Solidity: event WriteGot(address indexed source, string key, uint64 timestamp, uint64 value)
func (_Eth *EthFilterer) ParseWriteGot(log types.Log) (*EthWriteGot, error) {
	event := new(EthWriteGot)
	if err := _Eth.contract.UnpackLog(event, "WriteGot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
