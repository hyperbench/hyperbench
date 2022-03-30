package glua

import (
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_CommonResult(t *testing.T) {

	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	result := &fcom.Result{
		Label:  "Confirm",
		UID:    "UUID",
		Status: fcom.Success,
		Ret:    []interface{}{"result", "result"},
	}
	cLua := newCommonResult(L, result)
	L.SetField(mt, "result", cLua)
	scripts := []string{`
		function run()
			return case.result
		end
	`}
	for _, script := range scripts {
		lvalue, err := runLuaRunFunc(L, script)
		assert.Nil(t, err)
		idx := &fcom.Result{}
		err = TableLua2GoStruct(lvalue.(*lua.LTable), idx)
		assert.Nil(t, err)
		assert.Equal(t, idx, &fcom.Result{Label: "Confirm", UID: "UUID", BuildTime: 0, SendTime: 0, ConfirmTime: 0, WriteTime: 0, Status: "success", Ret: []interface{}{"result", "result"}})
	}
}
