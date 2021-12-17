package glua

import (
	"github.com/meshplus/hyperbench/plugins/toolkit"
	lua "github.com/yuin/gopher-lua"
)

const (
	toolKit  = "toolkit"
	methods  = "methods"
	index    = "__index"
	newIndex = "__newindex"

	// RandStr support method
	RandStr = "RandStr"
	randStr = "randStr"

	RandInt = "RandInt"
	randInt = "randInt"

	_String = "String"
	_string = "string"

	Hex = "Hex"
	hex = "hex"
	// support value
	TestInterface = "TestInterface"
)

var toolKitMethod = map[string]lua.LGFunction{
	RandStr: randStrLua,
	randStr: randStrLua,

	RandInt: RandIntLua,
	randInt: RandIntLua,

	_String: StringLua,
	_string: StringLua,

	Hex: HexLua,
	hex: HexLua,

	TestInterface: TestInterfaceLua,
}

func newTookKitMock(L *lua.LState) lua.LValue {
	return newToolKit(L, toolkit.NewToolKit())
}

func registerTookKit(L *lua.LState) lua.LValue {

	toolKitTable := L.NewTypeMetatable(toolKit)
	L.SetGlobal(toolKit, toolKitTable)
	L.SetField(toolKitTable, methods, L.SetFuncs(L.NewTable(), toolKitMethod))
	L.SetField(toolKitTable, index, L.NewFunction(func(L *lua.LState) int {
		invokeName := L.CheckString(2)
		if _, ok := toolKitMethod[invokeName]; ok {
			L.Push(L.GetField(toolKitTable, methods).(*lua.LTable).RawGetString(invokeName))
			return 1
		}
		kit := checkToolKit(L)
		L.Push(lua.LString(kit.Name))
		return 1
	}))
	L.SetField(toolKitTable, newIndex, L.NewFunction(func(state *lua.LState) int {
		kit := checkToolKit(L)
		name := L.CheckString(3)
		kit.Name = name
		return 0
	}))
	return toolKitTable
}

func newToolKit(L *lua.LState, kit *toolkit.ToolKit) lua.LValue {
	toolkitTable := L.GetTypeMetatable(toolKit)
	if toolkitTable == nil || toolkitTable == lua.LNil {
		toolkitTable = registerTookKit(L)
	}
	ud := L.NewUserData()
	ud.Value = kit
	ud.Metatable = toolkitTable
	return ud
}

func randStrLua(state *lua.LState) int {
	kit := checkToolKit(state)
	size := state.CheckInt(2)
	ret := kit.RandStr(uint(size))
	state.Push(lua.LString(ret))
	return 1
}

func RandIntLua(state *lua.LState) int {
	kit := checkToolKit(state)
	begin := state.CheckInt(2)
	end := state.CheckInt(3)
	ret := kit.RandInt(begin, end)
	state.Push(lua.LNumber(ret))
	return 1
}

func HexLua(state *lua.LState) int {
	kit := checkToolKit(state)
	str := state.CheckString(2)
	ret := kit.Hex(str)
	state.Push(lua.LString(ret))
	return 1
}

//String todo  support {}interface and ...int
func StringLua(state *lua.LState) int {
	kit := checkToolKit(state)
	input := state.CheckAny(2)
	arg1, err := Lua2Go(input)
	if err != nil {

	}
	end := state.CheckUserData(3)
	ret := kit.String(arg1, end.Value.([]int)...)
	state.Push(lua.LString(ret))
	return 1
}

func checkToolKit(L *lua.LState) *toolkit.ToolKit {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*toolkit.ToolKit); ok {
		return v
	}
	L.ArgError(1, "toolkit expected")
	return nil
}

func TestInterfaceLua(state *lua.LState) int {
	kit := checkToolKit(state)
	input := state.CheckUserData(2)
	ret := kit.TestInterface(LuaParams(input))
	state.Push(Go2Lua(state, ret))
	return 1
}
