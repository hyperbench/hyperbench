package glua

import (
	"fmt"
	idex "github.com/meshplus/hyperbench/plugins/index"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

const (
	Index = "index"

	//field

	Worker = "worker"
	VM     = "vm"
	Engine = "engine"
	Tx     = "tx"
)

func registerIndex(L *lua.LState) lua.LValue {
	IndexTable := L.NewTypeMetatable(Index)
	L.SetGlobal(Index, IndexTable)
	L.SetField(IndexTable, index, L.NewFunction(func(L *lua.LState) int {
		index := checkIndex(L, 1)
		invokeName := L.CheckString(2)
		switch strings.ToLower(invokeName) {
		case Worker:
			L.Push(lua.LNumber(index.Worker))
			return 1
		case VM:
			L.Push(lua.LNumber(index.VM))
			return 1
		case Engine:
			L.Push(lua.LNumber(index.Engine))
			return 1
		case Tx:
			L.Push(lua.LNumber(index.Tx))
			return 1
		}
		L.Push(lua.LNil)
		return 1
	}))
	L.SetField(IndexTable, newIndex, L.NewFunction(func(L *lua.LState) int {
		index := checkIndex(L, 1)
		switch strings.ToLower(L.CheckString(2)) {
		case Worker:
			index.Worker = L.CheckInt64(3)
			return 0
		case VM:
			index.VM = L.CheckInt64(3)
			return 0
		case Engine:
			index.Engine = L.CheckInt64(3)
			return 0
		case Tx:
			index.Tx = L.CheckInt64(3)
			return 0
		}
		return 0
	}))
	L.SetField(IndexTable, toString, L.NewFunction(func(L *lua.LState) int {
		result := checkIndex(L, 1)
		L.Push(lua.LString(fmt.Sprintf("%v", result)))
		return 1
	}))
	L.SetField(IndexTable, eq, L.NewFunction(func(L *lua.LState) int {
		r1 := checkIndex(L, 1)
		r2 := checkIndex(L, 2)
		L.Push(lua.LBool(fmt.Sprintf("%v", r1) == fmt.Sprintf("%v", r2)))
		return 1
	}))
	return IndexTable
}

func newIdexIndex(L *lua.LState, idx *idex.Index) lua.LValue {

	metatable := L.GetTypeMetatable(Index)
	if metatable == nil || metatable == lua.LNil {
		metatable = registerIndex(L)
	}
	ud := L.NewUserData()
	ud.Value = idx
	ud.Metatable = metatable
	return ud

}

func checkIndex(L *lua.LState, n int) *idex.Index {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*idex.Index); ok {
		return v
	}
	L.ArgError(1, "Index expected")
	return nil

}
