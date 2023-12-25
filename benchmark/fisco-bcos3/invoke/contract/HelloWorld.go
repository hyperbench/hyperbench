// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package hello

import (
	"math/big"
	"strings"

	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
)

// HelloWorldABI is the input ABI used to generate the binding from.
const HelloWorldABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"initValue\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"v\",\"type\":\"string\"}],\"name\":\"set\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// HelloWorldBin is the compiled bytecode used for deploying new contracts.
var HelloWorldBin = "0x60806040523480156200001157600080fd5b5060405162000924380380620009248339818101604052810190620000379190620002ac565b80600090805190602001906200004f9291906200005f565b5060006001819055505062000362565b8280546200006d906200032c565b90600052602060002090601f016020900481019282620000915760008555620000dd565b82601f10620000ac57805160ff1916838001178555620000dd565b82800160010185558215620000dd579182015b82811115620000dc578251825591602001919060010190620000bf565b5b509050620000ec9190620000f0565b5090565b5b808211156200010b576000816000905550600101620000f1565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b62000178826200012d565b810181811067ffffffffffffffff821117156200019a57620001996200013e565b5b80604052505050565b6000620001af6200010f565b9050620001bd82826200016d565b919050565b600067ffffffffffffffff821115620001e057620001df6200013e565b5b620001eb826200012d565b9050602081019050919050565b60005b8381101562000218578082015181840152602081019050620001fb565b8381111562000228576000848401525b50505050565b6000620002456200023f84620001c2565b620001a3565b90508281526020810184848401111562000264576200026362000128565b5b62000271848285620001f8565b509392505050565b600082601f83011262000291576200029062000123565b5b8151620002a38482602086016200022e565b91505092915050565b600060208284031215620002c557620002c462000119565b5b600082015167ffffffffffffffff811115620002e657620002e56200011e565b5b620002f48482850162000279565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200034557607f821691505b602082108114156200035c576200035b620002fd565b5b50919050565b6105b280620003726000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634ed3885e1461004657806354fd4d50146100765780636d4ce63c14610094575b600080fd5b610060600480360381019061005b919061031c565b6100b2565b60405161006d9190610402565b60405180910390f35b61007e610172565b60405161008b919061043d565b60405180910390f35b61009c610178565b6040516100a99190610402565b60405180910390f35b606060008080546100c290610487565b80601f01602080910402602001604051908101604052809291908181526020018280546100ee90610487565b801561013b5780601f106101105761010080835404028352916020019161013b565b820191906000526020600020905b81548152906001019060200180831161011e57829003601f168201915b5050505050905083836000919061015392919061020a565b506001805461016291906104e8565b6001819055508091505092915050565b60015481565b60606000805461018790610487565b80601f01602080910402602001604051908101604052809291908181526020018280546101b390610487565b80156102005780601f106101d557610100808354040283529160200191610200565b820191906000526020600020905b8154815290600101906020018083116101e357829003601f168201915b5050505050905090565b82805461021690610487565b90600052602060002090601f016020900481019282610238576000855561027f565b82601f1061025157803560ff191683800117855561027f565b8280016001018555821561027f579182015b8281111561027e578235825591602001919060010190610263565b5b50905061028c9190610290565b5090565b5b808211156102a9576000816000905550600101610291565b5090565b600080fd5b600080fd5b600080fd5b600080fd5b600080fd5b60008083601f8401126102dc576102db6102b7565b5b8235905067ffffffffffffffff8111156102f9576102f86102bc565b5b602083019150836001820283011115610315576103146102c1565b5b9250929050565b60008060208385031215610333576103326102ad565b5b600083013567ffffffffffffffff811115610351576103506102b2565b5b61035d858286016102c6565b92509250509250929050565b600081519050919050565b600082825260208201905092915050565b60005b838110156103a3578082015181840152602081019050610388565b838111156103b2576000848401525b50505050565b6000601f19601f8301169050919050565b60006103d482610369565b6103de8185610374565b93506103ee818560208601610385565b6103f7816103b8565b840191505092915050565b6000602082019050818103600083015261041c81846103c9565b905092915050565b6000819050919050565b61043781610424565b82525050565b6000602082019050610452600083018461042e565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b6000600282049050600182168061049f57607f821691505b602082108114156104b3576104b2610458565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006104f382610424565b91506104fe83610424565b9250817f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03831360008312151615610539576105386104b9565b5b817f8000000000000000000000000000000000000000000000000000000000000000038312600083121615610571576105706104b9565b5b82820190509291505056fea2646970667358221220afbf5e9d50254675961e56ef44b66b430784caf7a83d3ed163d037d4cf90eaa364736f6c634300080b0033"

// DeployHelloWorld deploys a new contract, binding an instance of HelloWorld to it.
func DeployHelloWorld(auth *bind.TransactOpts, backend bind.ContractBackend, initValue string) (common.Address, *types.Receipt, *HelloWorld, error) {
	parsed, err := abi.JSON(strings.NewReader(HelloWorldABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, receipt, contract, err := bind.DeployContract(auth, parsed, common.FromHex(HelloWorldBin), backend, initValue)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, receipt, &HelloWorld{HelloWorldCaller: HelloWorldCaller{contract: contract}, HelloWorldTransactor: HelloWorldTransactor{contract: contract}, HelloWorldFilterer: HelloWorldFilterer{contract: contract}}, nil
}

func AsyncDeployHelloWorld(auth *bind.TransactOpts, handler func(*types.Receipt, error), backend bind.ContractBackend, initValue string) (*types.Transaction, error) {
	parsed, err := abi.JSON(strings.NewReader(HelloWorldABI))
	if err != nil {
		return nil, err
	}

	tx, err := bind.AsyncDeployContract(auth, handler, parsed, common.FromHex(HelloWorldBin), backend, initValue)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// HelloWorld is an auto generated Go binding around a Solidity contract.
type HelloWorld struct {
	HelloWorldCaller     // Read-only binding to the contract
	HelloWorldTransactor // Write-only binding to the contract
	HelloWorldFilterer   // Log filterer for contract events
}

// HelloWorldCaller is an auto generated read-only Go binding around a Solidity contract.
type HelloWorldCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelloWorldTransactor is an auto generated write-only Go binding around a Solidity contract.
type HelloWorldTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelloWorldFilterer is an auto generated log filtering Go binding around a Solidity contract events.
type HelloWorldFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelloWorldSession is an auto generated Go binding around a Solidity contract,
// with pre-set call and transact options.
type HelloWorldSession struct {
	Contract     *HelloWorld       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HelloWorldCallerSession is an auto generated read-only Go binding around a Solidity contract,
// with pre-set call options.
type HelloWorldCallerSession struct {
	Contract *HelloWorldCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// HelloWorldTransactorSession is an auto generated write-only Go binding around a Solidity contract,
// with pre-set transact options.
type HelloWorldTransactorSession struct {
	Contract     *HelloWorldTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// HelloWorldRaw is an auto generated low-level Go binding around a Solidity contract.
type HelloWorldRaw struct {
	Contract *HelloWorld // Generic contract binding to access the raw methods on
}

// HelloWorldCallerRaw is an auto generated low-level read-only Go binding around a Solidity contract.
type HelloWorldCallerRaw struct {
	Contract *HelloWorldCaller // Generic read-only contract binding to access the raw methods on
}

// HelloWorldTransactorRaw is an auto generated low-level write-only Go binding around a Solidity contract.
type HelloWorldTransactorRaw struct {
	Contract *HelloWorldTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHelloWorld creates a new instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorld(address common.Address, backend bind.ContractBackend) (*HelloWorld, error) {
	contract, err := bindHelloWorld(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HelloWorld{HelloWorldCaller: HelloWorldCaller{contract: contract}, HelloWorldTransactor: HelloWorldTransactor{contract: contract}, HelloWorldFilterer: HelloWorldFilterer{contract: contract}}, nil
}

// NewHelloWorldCaller creates a new read-only instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorldCaller(address common.Address, caller bind.ContractCaller) (*HelloWorldCaller, error) {
	contract, err := bindHelloWorld(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HelloWorldCaller{contract: contract}, nil
}

// NewHelloWorldTransactor creates a new write-only instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorldTransactor(address common.Address, transactor bind.ContractTransactor) (*HelloWorldTransactor, error) {
	contract, err := bindHelloWorld(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HelloWorldTransactor{contract: contract}, nil
}

// NewHelloWorldFilterer creates a new log filterer instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorldFilterer(address common.Address, filterer bind.ContractFilterer) (*HelloWorldFilterer, error) {
	contract, err := bindHelloWorld(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HelloWorldFilterer{contract: contract}, nil
}

// bindHelloWorld binds a generic wrapper to an already deployed contract.
func bindHelloWorld(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(HelloWorldABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HelloWorld *HelloWorldRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HelloWorld.Contract.HelloWorldCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HelloWorld *HelloWorldRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	return _HelloWorld.Contract.HelloWorldTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HelloWorld *HelloWorldRaw) TransactWithResult(opts *bind.TransactOpts, result interface{}, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
	return _HelloWorld.Contract.HelloWorldTransactor.contract.TransactWithResult(opts, result, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HelloWorld *HelloWorldCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _HelloWorld.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HelloWorld *HelloWorldTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
	return _HelloWorld.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HelloWorld *HelloWorldTransactorRaw) TransactWithResult(opts *bind.TransactOpts, result interface{}, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
	return _HelloWorld.Contract.contract.TransactWithResult(opts, result, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(string)
func (_HelloWorld *HelloWorldCaller) Get(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _HelloWorld.contract.Call(opts, out, "get")
	return *ret0, err
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(string)
func (_HelloWorld *HelloWorldSession) Get() (string, error) {
	return _HelloWorld.Contract.Get(&_HelloWorld.CallOpts)
}

// Get is a free data retrieval call binding the contract method 0x6d4ce63c.
//
// Solidity: function get() constant returns(string)
func (_HelloWorld *HelloWorldCallerSession) Get() (string, error) {
	return _HelloWorld.Contract.Get(&_HelloWorld.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(int256)
func (_HelloWorld *HelloWorldCaller) Version(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _HelloWorld.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(int256)
func (_HelloWorld *HelloWorldSession) Version() (*big.Int, error) {
	return _HelloWorld.Contract.Version(&_HelloWorld.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(int256)
func (_HelloWorld *HelloWorldCallerSession) Version() (*big.Int, error) {
	return _HelloWorld.Contract.Version(&_HelloWorld.CallOpts)
}

// Set is a paid mutator transaction binding the contract method 0x4ed3885e.
//
// Solidity: function set(string v) returns(string)
func (_HelloWorld *HelloWorldTransactor) Set(opts *bind.TransactOpts, v string) (string, *types.Transaction, *types.Receipt, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	transaction, receipt, err := _HelloWorld.contract.TransactWithResult(opts, out, "set", v)
	return *ret0, transaction, receipt, err
}

func (_HelloWorld *HelloWorldTransactor) AsyncSet(handler func(*types.Receipt, error), opts *bind.TransactOpts, v string) (*types.Transaction, error) {
	return _HelloWorld.contract.AsyncTransact(opts, handler, "set", v)
}

// Set is a paid mutator transaction binding the contract method 0x4ed3885e.
//
// Solidity: function set(string v) returns(string)
func (_HelloWorld *HelloWorldSession) Set(v string) (string, *types.Transaction, *types.Receipt, error) {
	return _HelloWorld.Contract.Set(&_HelloWorld.TransactOpts, v)
}

func (_HelloWorld *HelloWorldSession) AsyncSet(handler func(*types.Receipt, error), v string) (*types.Transaction, error) {
	return _HelloWorld.Contract.AsyncSet(handler, &_HelloWorld.TransactOpts, v)
}

// Set is a paid mutator transaction binding the contract method 0x4ed3885e.
//
// Solidity: function set(string v) returns(string)
func (_HelloWorld *HelloWorldTransactorSession) Set(v string) (string, *types.Transaction, *types.Receipt, error) {
	return _HelloWorld.Contract.Set(&_HelloWorld.TransactOpts, v)
}

func (_HelloWorld *HelloWorldTransactorSession) AsyncSet(handler func(*types.Receipt, error), v string) (*types.Transaction, error) {
	return _HelloWorld.Contract.AsyncSet(handler, &_HelloWorld.TransactOpts, v)
}
