package glua

import (
	"github.com/hyperbench/hyperbench/plugins/toolkit"
	lua "github.com/yuin/gopher-lua"
)

func hexLuaFunction(L *lua.LState, kit *toolkit.ToolKit) *lua.LFunction {
	return L.NewFunction(func(state *lua.LState) int {
		input := state.CheckString(1)
		ret := kit.Hex(input)
		state.Push(lua.LString(ret))
		return 1
	})
}

func randStrLuaFunction(L *lua.LState, kit *toolkit.ToolKit) *lua.LFunction {
	return L.NewFunction(func(state *lua.LState) int {
		size := state.CheckInt(1)
		ret := kit.RandStr(uint(size))
		L.Push(lua.LString(ret))
		return 1
	})
}

func randIntLuaFunction(L *lua.LState, kit *toolkit.ToolKit) *lua.LFunction {
	return L.NewFunction(func(state *lua.LState) int {
		min := state.CheckInt(1)
		max := state.CheckInt(2)
		ret := kit.RandInt(min, max)
		L.Push(lua.LNumber(ret))
		return 1
	})
}

func stringLuaFunction(L *lua.LState, kit *toolkit.ToolKit) *lua.LFunction {
	return L.NewFunction(func(state *lua.LState) int {
		argLength := state.GetTop()
		if argLength < 1 {
			panic("args are less than 1")
		}
		input := state.CheckAny(1)
		if argLength == 1 {
			ret := kit.String(input)
			L.Push(lua.LString(ret))
			return 1
		}
		var offsets []int
		for i := 2; i < argLength; i++ {
			offset := state.CheckInt(i)
			offsets = append(offsets, offset)
		}
		ret := kit.String(input, offsets...)
		L.Push(lua.LString(ret))
		return 1
	})
}

func newToolKit(L *lua.LState, kit *toolkit.ToolKit) lua.LValue {
	toolkitTable := L.NewTable()
	toolkitTable.RawSetString("Hex", hexLuaFunction(L, kit))
	toolkitTable.RawSetString("RandStr", randStrLuaFunction(L, kit))
	toolkitTable.RawSetString("RandInt", randIntLuaFunction(L, kit))
	toolkitTable.RawSetString("String", stringLuaFunction(L, kit))
	return toolkitTable
}
