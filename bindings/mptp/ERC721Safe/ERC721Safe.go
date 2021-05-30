// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERC721Safe

import (
	"math/big"
	"strings"

	ethereum "github.com/celo-org/celo-blockchain"
	"github.com/celo-org/celo-blockchain/accounts/abi"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ERC721SafeABI is the input ABI used to generate the binding from.
const ERC721SafeABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"fundERC721\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"constant\":false}]"

// ERC721SafeBin is the compiled bytecode used for deploying new contracts.
var ERC721SafeBin = "0x608060405234801561001057600080fd5b5061010c806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80637354298014602d575b600080fd5b606060048036036060811015604157600080fd5b506001600160a01b038135811691602081013590911690604001356062565b005b604080516323b872dd60e01b81526001600160a01b03848116600483015230602483015260448201849052915185928316916323b872dd91606480830192600092919082900301818387803b15801560b957600080fd5b505af115801560cc573d6000803e3d6000fd5b505050505050505056fea26469706673582212205b1be1ab12c25bb7f48d3a5d588765d7eabe5fb1764e6de2268bf3f2981d732f64736f6c63430006040033"

// DeployERC721Safe deploys a new Ethereum contract, binding an instance of ERC721Safe to it.
func DeployERC721Safe(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC721Safe, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721SafeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC721SafeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC721Safe{ERC721SafeCaller: ERC721SafeCaller{contract: contract}, ERC721SafeTransactor: ERC721SafeTransactor{contract: contract}, ERC721SafeFilterer: ERC721SafeFilterer{contract: contract}}, nil
}

// ERC721Safe is an auto generated Go binding around an Ethereum contract.
type ERC721Safe struct {
	ERC721SafeCaller     // Read-only binding to the contract
	ERC721SafeTransactor // Write-only binding to the contract
	ERC721SafeFilterer   // Log filterer for contract events
}

// ERC721SafeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC721SafeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721SafeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC721SafeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721SafeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC721SafeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC721SafeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC721SafeSession struct {
	Contract     *ERC721Safe       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC721SafeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC721SafeCallerSession struct {
	Contract *ERC721SafeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ERC721SafeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC721SafeTransactorSession struct {
	Contract     *ERC721SafeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ERC721SafeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC721SafeRaw struct {
	Contract *ERC721Safe // Generic contract binding to access the raw methods on
}

// ERC721SafeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC721SafeCallerRaw struct {
	Contract *ERC721SafeCaller // Generic read-only contract binding to access the raw methods on
}

// ERC721SafeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC721SafeTransactorRaw struct {
	Contract *ERC721SafeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC721Safe creates a new instance of ERC721Safe, bound to a specific deployed contract.
func NewERC721Safe(address common.Address, backend bind.ContractBackend) (*ERC721Safe, error) {
	contract, err := bindERC721Safe(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC721Safe{ERC721SafeCaller: ERC721SafeCaller{contract: contract}, ERC721SafeTransactor: ERC721SafeTransactor{contract: contract}, ERC721SafeFilterer: ERC721SafeFilterer{contract: contract}}, nil
}

// NewERC721SafeCaller creates a new read-only instance of ERC721Safe, bound to a specific deployed contract.
func NewERC721SafeCaller(address common.Address, caller bind.ContractCaller) (*ERC721SafeCaller, error) {
	contract, err := bindERC721Safe(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721SafeCaller{contract: contract}, nil
}

// NewERC721SafeTransactor creates a new write-only instance of ERC721Safe, bound to a specific deployed contract.
func NewERC721SafeTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC721SafeTransactor, error) {
	contract, err := bindERC721Safe(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC721SafeTransactor{contract: contract}, nil
}

// NewERC721SafeFilterer creates a new log filterer instance of ERC721Safe, bound to a specific deployed contract.
func NewERC721SafeFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC721SafeFilterer, error) {
	contract, err := bindERC721Safe(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC721SafeFilterer{contract: contract}, nil
}

// bindERC721Safe binds a generic wrapper to an already deployed contract.
func bindERC721Safe(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721SafeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// ParseERC721SafeABI parses the ABI
func ParseERC721SafeABI() (*abi.ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC721SafeABI))
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Safe *ERC721SafeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721Safe.Contract.ERC721SafeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Safe *ERC721SafeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Safe.Contract.ERC721SafeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Safe *ERC721SafeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Safe.Contract.ERC721SafeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC721Safe *ERC721SafeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC721Safe.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC721Safe *ERC721SafeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC721Safe.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC721Safe *ERC721SafeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC721Safe.Contract.contract.Transact(opts, method, params...)
}

// FundERC721 is a paid mutator transaction binding the contract method 0x73542980.
//
// Solidity: function fundERC721(address tokenAddress, address owner, uint256 tokenID) returns()
func (_ERC721Safe *ERC721SafeTransactor) FundERC721(opts *bind.TransactOpts, tokenAddress common.Address, owner common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _ERC721Safe.contract.Transact(opts, "fundERC721", tokenAddress, owner, tokenID)
}

// FundERC721 is a paid mutator transaction binding the contract method 0x73542980.
//
// Solidity: function fundERC721(address tokenAddress, address owner, uint256 tokenID) returns()
func (_ERC721Safe *ERC721SafeSession) FundERC721(tokenAddress common.Address, owner common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _ERC721Safe.Contract.FundERC721(&_ERC721Safe.TransactOpts, tokenAddress, owner, tokenID)
}

// FundERC721 is a paid mutator transaction binding the contract method 0x73542980.
//
// Solidity: function fundERC721(address tokenAddress, address owner, uint256 tokenID) returns()
func (_ERC721Safe *ERC721SafeTransactorSession) FundERC721(tokenAddress common.Address, owner common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _ERC721Safe.Contract.FundERC721(&_ERC721Safe.TransactOpts, tokenAddress, owner, tokenID)
}

// TryParseLog attempts to parse a log. Returns the parsed log, evenName and whether it was succesfull
func (_ERC721Safe *ERC721SafeFilterer) TryParseLog(log types.Log) (eventName string, event interface{}, ok bool, err error) {
	eventName, ok, err = _ERC721Safe.contract.LogEventName(log)
	if err != nil || !ok {
		return "", nil, false, err
	}

	switch eventName {
	}
	if err != nil {
		return "", nil, false, err
	}

	return eventName, event, ok, nil
}
