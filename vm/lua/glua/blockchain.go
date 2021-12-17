package glua

import (
	"fmt"
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	lua "github.com/yuin/gopher-lua"
)

const (
	BlockChain = "blockchain"

	Invoke = "Invoke"
	invoke = "invoke"

	Transfer = "Transfer"
	transfer = "transfer"

	Confirm = "Confirm"
	confirm = "confirm"

	Query = "Query"
	query = "query"

	BCOption = "Option"
	BCoption = "option"

	GetContext = "GetContext"
	getContext = "getContext"

	SetContext = "SetContext"
	setContext = "setContext"

	ResetContext = "ResetContext"
	resetContext = "resetContext"

	Statistic = "Statistic"
	statistic = "statistic"
)

var blockChainMethod = map[string]lua.LGFunction{
	Invoke:       InvokeLua,
	invoke:       InvokeLua,
	Transfer:     TransferLua,
	Confirm:      ConfirmLua,
	confirm:      ConfirmLua,
	Query:        QueryLua,
	query:        QueryLua,
	BCOption:     OptionLua,
	BCoption:     OptionLua,
	GetContext:   GetContextLua,
	getContext:   GetContextLua,
	SetContext:   SetContextLua,
	setContext:   SetContextLua,
	ResetContext: ResetContextLua,
	resetContext: ResetContextLua,
	//Statistic:StatisticLua,
}

func registerBlockchain(L *lua.LState) lua.LValue {
	CommonResultTable := L.NewTypeMetatable(BlockChain)
	L.SetGlobal(BlockChain, CommonResultTable)
	L.SetField(CommonResultTable, index, L.SetFuncs(L.NewTable(), blockChainMethod))
	return CommonResultTable
}

func newBlockchain(L *lua.LState, client blockchain.Blockchain) lua.LValue {
	metatable := L.GetTypeMetatable(BlockChain)
	if metatable == nil || metatable == lua.LNil {
		metatable = registerBlockchain(L)
	}
	ud := L.NewUserData()
	ud.Value = client
	ud.Metatable = metatable
	return ud
}

func checkBlockchain(L *lua.LState) blockchain.Blockchain {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(blockchain.Blockchain); ok {
		return v
	}
	L.ArgError(1, "blockchain Blockchain expected")
	return nil

}

func InvokeLua(L *lua.LState) int {
	var invoke bcom.Invoke
	return invokeLua(L, invoke, func(b blockchain.Blockchain, b2 interface{}, option ...bcom.Option) interface{} {
		return b.Invoke(b2.(bcom.Invoke), option...)
	})
}
func TransferLua(L *lua.LState) int {
	var transfer bcom.Transfer
	return invokeLua(L, transfer, func(b blockchain.Blockchain, b2 interface{}, option ...bcom.Option) interface{} {
		return b.Transfer(b2.(bcom.Transfer), option...)
	})
}

func ConfirmLua(L *lua.LState) int {
	cli := checkBlockchain(L)
	invokeTable := L.CheckUserData(2)
	result, ok := invokeTable.Value.(*common.Result)
	if !ok {
		L.ArgError(1, "*common.Result expected")
	}
	if L.GetTop() == 2 {
		ret := cli.Confirm(result)
		L.Push(newCommonResult(L, ret))
		return 1
	}
	var opts []bcom.Option
	for i := 3; i <= L.GetTop(); i++ {
		table := L.CheckTable(i)
		var map1 bcom.Option
		err := Map(table, &map1)
		if err != nil {
			L.ArgError(1, "common.Option expected")
		}
		opts = append(opts, map1)
	}
	ret := cli.Confirm(result, opts...)
	L.Push(newCommonResult(L, ret))
	return 1
}

func QueryLua(L *lua.LState) int {
	var query bcom.Query
	return invokeLua(L, query, func(b blockchain.Blockchain, i interface{}, option ...bcom.Option) interface{} {
		return b.Query(query, option...)
	})
}

func OptionLua(L *lua.LState) int {
	cli := checkBlockchain(L)
	invokeTable := L.CheckTable(2)
	var opt bcom.Option
	err := Map(invokeTable, &opt)
	if err != nil {
		L.ArgError(1, "Option. expected")
	}
	err = cli.Option(opt)
	if err != nil {
		L.ArgError(1, fmt.Sprintf("Option failed %s", err.Error()))
	}
	return 0
}

func GetContextLua(L *lua.LState) int {
	cli := checkBlockchain(L)
	str, err := cli.GetContext()
	if err != nil {
		L.ArgError(1, fmt.Sprintf("Option failed %s", err.Error()))
	}
	L.Push(lua.LString(str))
	return 1
}

func SetContextLua(L *lua.LState) int {
	cli := checkBlockchain(L)
	contextLua := L.CheckString(2)
	err := cli.SetContext(contextLua)
	if err != nil {
		L.ArgError(1, fmt.Sprintf("SetContext failed %s", err.Error()))
	}
	return 0
}

func ResetContextLua(L *lua.LState) int {
	cli := checkBlockchain(L)
	err := cli.ResetContext()
	if err != nil {
		L.ArgError(1, fmt.Sprintf("ResetContext failed %s", err.Error()))
	}
	return 0
}

func invokeLua(L *lua.LState, arg1Type interface{}, fn func(blockchain.Blockchain, interface{}, ...bcom.Option) interface{}) int {
	cli := checkBlockchain(L)
	invokeTable := L.CheckTable(2)
	err := Map(invokeTable, &arg1Type)
	if err != nil {
		L.ArgError(1, "interface. expected")
	}
	if L.GetTop() == 2 {
		ret := fn(cli, arg1Type)
		L.Push(decode(L, ret))
		return 1
	}
	var opts []bcom.Option
	for i := 3; i <= L.GetTop(); i++ {
		table := L.CheckTable(i)
		var map1 bcom.Option
		err := Map(table, &map1)
		if err != nil {
			L.ArgError(1, "common.Option expected")
		}
		opts = append(opts, map1)
	}
	ret := fn(cli, arg1Type, opts...)
	L.Push(decode(L, ret))
	return 1
}
