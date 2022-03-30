package glua

import (
	fcom "github.com/hyperbench/hyperbench-common/common"
	idex "github.com/hyperbench/hyperbench/plugins/index"
	"github.com/hyperbench/hyperbench/plugins/toolkit"
	lua "github.com/yuin/gopher-lua"
)

func NewClientLValue(L *lua.LState, client fcom.Blockchain) lua.LValue {
	return newBlockchain(L, client)
}

func NewToolKitLValue(L *lua.LState, kit *toolkit.ToolKit) lua.LValue {
	return newToolKit(L, kit)
}

func NewLIndexLValue(L *lua.LState, idx *idex.Index) lua.LValue {
	return newIdexIndex(L, idx)
}

func NewResultLValue(L *lua.LState, r *fcom.Result) lua.LValue {
	return newCommonResult(L, r)
}
