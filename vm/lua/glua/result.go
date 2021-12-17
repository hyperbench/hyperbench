package glua

import (
	"fmt"
	"github.com/meshplus/hyperbench/common"
	lua "github.com/yuin/gopher-lua"
	"strings"
)

const (
	toString     = "__tostring"
	eq           = "__eq"
	CommonResult = "Result"

	//field

	Label       = "label"
	UID         = "uid"
	BuildTime   = "buildtime"
	SendTime    = "sendtime"
	ConfirmTime = "confirmtime"
	WriteTime   = "writetime"
	Status      = "status"
	Ret         = "ret"
)

func registerCommonResult(L *lua.LState) lua.LValue {
	CommonResultTable := L.NewTypeMetatable(CommonResult)
	L.SetGlobal(CommonResult, CommonResultTable)
	L.SetField(CommonResultTable, index, L.NewFunction(func(L *lua.LState) int {
		result := checkCommonResult(L, 1)
		invokeName := L.CheckString(2)
		switch strings.ToLower(invokeName) {
		case Label:
			L.Push(lua.LString(result.Label))
			return 1
		case UID:
			L.Push(lua.LString(result.UID))
			return 1
		case BuildTime:
			L.Push(lua.LNumber(result.BuildTime))
			return 1
		case SendTime:
			L.Push(lua.LNumber(result.SendTime))
			return 1
		case ConfirmTime:
			L.Push(lua.LNumber(result.ConfirmTime))
			return 1
		case WriteTime:
			L.Push(lua.LNumber(result.WriteTime))
			return 1
		case Status:
			L.Push(lua.LString(result.Status))
			return 1
		case Ret:
			//todo support ret
			L.Push(lua.LString("ret"))
			return 1
		}
		L.Push(lua.LString(""))
		return 1
	}))
	L.SetField(CommonResultTable, newIndex, L.NewFunction(func(state *lua.LState) int {
		result := checkCommonResult(L, 1)
		switch strings.ToLower(L.CheckString(2)) {
		case Label:
			result.Label = L.CheckString(3)
			return 0
		case UID:
			result.UID = L.CheckString(3)
			return 0
		case BuildTime:
			result.BuildTime = L.CheckInt64(3)
			return 0
		case SendTime:
			result.SendTime = L.CheckInt64(3)
			return 0
		case ConfirmTime:
			result.ConfirmTime = L.CheckInt64(3)
			return 0
		case WriteTime:
			result.WriteTime = L.CheckInt64(3)
			return 0
		case Status:
			result.Status = common.Status(L.CheckString(3))
			return 0
		case Ret:
			//todo support ret
			return 0
		}
		return 0
	}))
	L.SetField(CommonResultTable, toString, L.NewFunction(func(state *lua.LState) int {
		fmt.Println(toString)
		result := checkCommonResult(L, 1)
		L.Push(lua.LString(fmt.Sprintf("%v", result)))
		return 1
	}))
	L.SetField(CommonResultTable, eq, L.NewFunction(func(state *lua.LState) int {
		r1 := checkCommonResult(L, 1)
		r2 := checkCommonResult(L, 2)
		L.Push(lua.LBool(fmt.Sprintf("%v", r1) == fmt.Sprintf("%v", r2)))
		return 1
	}))
	return CommonResultTable
}

func newCommonResult(L *lua.LState, r *common.Result) lua.LValue {
	metatable := L.GetTypeMetatable(CommonResult)
	if metatable == nil || metatable == lua.LNil {
		metatable = registerCommonResult(L)
	}
	ud := L.NewUserData()
	ud.Value = r
	ud.Metatable = metatable
	return ud
}

func checkCommonResult(L *lua.LState, n int) *common.Result {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*common.Result); ok {
		return v
	}
	L.ArgError(1, "Result expected")
	return nil

}
