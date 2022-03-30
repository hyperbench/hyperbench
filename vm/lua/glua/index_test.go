package glua

import (
	"github.com/hyperbench/hyperbench/plugins/index"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_index(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	passIdx := &index.Index{1, 1, 1, 1}

	cLua := newIdexIndex(L, passIdx)
	L.SetField(mt, "index", cLua)
	scripts := []string{`
		function run()
            case.index.Worker=2
            case.index.VM=2
            case.index.Engine=2
            case.index.Tx=2
			return case.index
		end
	`}
	for _, script := range scripts {
		lvalue, err := runLuaRunFunc(L, script)
		assert.Nil(t, err)
		idx := &index.Index{}
		err = TableLua2GoStruct(lvalue.(*lua.LTable), idx)
		assert.Nil(t, err)
		assert.Equal(t, idx, &index.Index{Worker: 2, VM: 2, Engine: 2, Tx: 2})
	}

}
