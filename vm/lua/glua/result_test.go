package glua

import (
	"github.com/meshplus/hyperbench/common"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"testing"
	"time"
)

var scriptCommonResut = `
	p = u.new();
	print(p)
	print("label:",p.label)
	print("uid:",p.UID)
	print("buildtime:",p.BuildTime)
	print("sendtime:",p.SendTime)
	print("confirm time:",p.ConfirmTime)
	print("write Time:",p.WriteTime)
	p.status="failure"
	print("status:",p.status)
	print("ret:",p.Ret)
`

func Test_CommonResult(t *testing.T) {

	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("u")
	L.SetGlobal("u", mt)

	L.SetField(mt, "new", L.NewFunction(func(L *lua.LState) int {
		new := luar.New(L, &common.Result{
			Label:       "heheh",
			UID:         "uid",
			BuildTime:   time.Now().Unix(),
			SendTime:    time.Now().Unix(),
			ConfirmTime: time.Now().Unix(),
			WriteTime:   time.Now().Unix(),
			Status:      common.Success,
			Ret:         []interface{}{[]byte("hehehehe")},
		})
		L.Push(new)
		return 1
	}))
	if err := L.DoString(scriptCommonResut); err != nil {
		panic(err)
	}
}
