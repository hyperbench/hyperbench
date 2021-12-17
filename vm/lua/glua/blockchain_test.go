package glua

import (
	"github.com/meshplus/hyperbench/plugins/blockchain/fake"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
	"testing"
)

var s1 = `
	 ret=case.blockchain:Invoke({
		   func="123",
		   args={"123", "123"}
		},{aa="aa"},{bb="bb"})
`

func Test_client(t *testing.T) {

	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	client, _ := fake.New()
	cLua := newBlockchain(L, client)
	L.SetField(mt, "blockchain", cLua)
	if err := L.DoString(s1); err != nil {
		panic(err)
	}
}

//handle BenchmarkInvokeLua-8   	   44344	     28033 ns/op
//auto   BenchmarkInvokeLua-8   	   41647	     26326 ns/op
func BenchmarkInvokeLua(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	client, _ := fake.New()
	cLua := luar.New(L, client)
	L.SetField(mt, "blockchain", cLua)

	L1 := lua.NewState()
	defer L1.Close()
	mt1 := L1.NewTypeMetatable("case")
	L1.SetGlobal("case", mt1)
	client1, _ := fake.New()
	cLua1 := newBlockchain(L1, client1)
	L1.SetField(mt1, "blockchain", cLua1)

	for i := 0; i < b.N; i++ {
		luar2(L1)
		luar1(L)
	}

}

func luar1(L *lua.LState) {
	if err := L.DoString(s1); err != nil {
		panic(err)
	}
}

func luar2(L *lua.LState) {
	if err := L.DoString(s1); err != nil {
		panic(err)
	}
}
