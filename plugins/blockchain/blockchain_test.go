package blockchain

import (
	"fmt"
	"testing"

	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func TestNewBlockchain(t *testing.T) {
	bk, err := NewBlockchain(base.ClientConfig{})
	assert.NotNil(t, bk)
	assert.NoError(t, err)

	bk, err = NewBlockchain(base.ClientConfig{
		ClientType: "eth",
	})
	assert.Nil(t, bk)
	assert.Error(t, err)

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	bk, err = NewBlockchain(base.ClientConfig{
		ClientType: "fabric",
	})
	assert.Nil(t, bk)
	assert.Error(t, err)
}

func TestNewHyperchain(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	bk, err := NewBlockchain(base.ClientConfig{
		ClientType: "flato",
	})
	assert.Nil(t, bk)
	assert.Error(t, err)
}

func TestLua(t *testing.T) {
	L := lua.NewState()
	if err := L.DoString(`
person = {
  name = "Michel",
  age  = "31", -- weakly input
  work_place = "San Jose",
  role = {
    {
      name = "Administrator"
    },
    {
      name = "Operator"
    },
  },
  idx = {1,2}
}
`); err != nil {
		panic(err)
	}
	var m map[string]interface{}
	if err := gluamapper.Map(L.GetGlobal("person").(*lua.LTable), &m); err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", m)
	fmt.Printf("%s\n", m["Idx"])
}

type Person struct {
	Name string `mapstructure:"name"`
}

func (p *Person) Hi(p2 *Person) string {
	return fmt.Sprintf("%s say hi to %s", p.Name, p2.Name)
}

func TestLua2(t *testing.T) {
	L := lua.NewState()
	p := &Person{
		Name: "i",
	}

	L.SetGlobal("p", luar.New(L, p))
	if err := L.DoString(`
print(p:Hi({name="123"}))
`); err != nil {
		panic(err)
	}
	var m map[string]interface{}

	fmt.Printf("%s\n", m)
	fmt.Printf("%s\n", m["Idx"])
}

func TestLua3(t *testing.T) {
	L := lua.NewState()
	p := &Person{
		Name: "i",
	}

	L.SetGlobal("p", luar.New(L, p))
	if err := L.DoString(`

local p1 = 1
local t1 = {}

function t1:process()
    print(p1)
    p1=p1+1
end

return t1
`); err != nil {
		panic(err)
	}
	ret := L.Get(-1).(*lua.LTable)
	L.Pop(-1)

	for i := 0; i < 10; i++ {
		_ = L.CallByParam(lua.P{
			Fn: ret.RawGet(lua.LString("process")),
		})
	}

}
