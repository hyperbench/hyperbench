package glua

import (
	"fmt"
	"github.com/meshplus/hyperbench/plugins/toolkit"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"testing"
)

var script1 = `
	print("-----1-----")
	print(u.toolkit:RandStr(10))
	print(u.toolkit:randStr(10))
	print("----2-----")
	print(u.toolkit.Name)
	print("----3-----")
	u.toolkit.Name = "hhe"
	print(u.toolkit.Name)
	print(u.toolkit:Hex("aaaaaa"))
	-- print(u.toolkit:String((1,2,3,4),1,2))
	print("TestInterfaceLua",u.toolkit:TestInterface({aa=bb}))

`

func Test_toolkit(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("u")
	L.SetGlobal("u", mt)
	L.SetField(mt, "toolkit", newTookKit(L))
	if err := L.DoString(script1); err != nil {
		panic(err)
	}
}

func Test_toolkit2(t *testing.T) {
	kit := toolkit.NewToolKit()
	fmt.Println("String:", kit.String([]byte{1, 2, 3, 4}, 0, 3))

}

func Test_1234(t *testing.T) {
	var value interface{}
	value = "1"
	//ret := value.(type)
	//fmt.Println(ret)
	switch converted := value.(type) {
	case bool:
		fmt.Println(converted)
	case int:
		fmt.Println(converted)
	case string:
		fmt.Println("str:", converted)

	}
}

func newTookKit(L *lua.LState) lua.LValue {
	u := &toolkit.ToolKit{Name: "aa"}
	return luar.New(L, u)
}
