package glua

import (
	idex "github.com/hyperbench/hyperbench/plugins/index"
	lua "github.com/yuin/gopher-lua"
)

func newIdexIndex(L *lua.LState, idx *idex.Index) lua.LValue {
	idxTable := L.NewTable()
	idxTable.RawSetString("Worker", lua.LNumber(idx.Worker))
	idxTable.RawSetString("VM", lua.LNumber(idx.VM))
	idxTable.RawSetString("Engine", lua.LNumber(idx.Engine))
	idxTable.RawSetString("Tx", lua.LNumber(idx.Tx))
	return idxTable
}
