package glua

import (
	"fmt"
	"github.com/hyperbench/hyperbench-common/base"
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/hyperbench/hyperbench/plugins/toolkit"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"reflect"
	"testing"
)

const (
	Invoke   = "Invoke"
	Transfer = "Transfer"
	Option   = "Option"
	Result   = "Result"
	Query    = "Query"
)

type FakeChain struct {
	Name     string
	base     *base.BlockchainBase
	tempData map[string]interface{}
}

func NewMock() (client *FakeChain, err error) {
	tempMap := make(map[string]interface{})
	return &FakeChain{"fake",
		base.NewBlockchainBase(base.ClientConfig{}),
		tempMap,
	}, nil
}

func (chain *FakeChain) Invoke(invoke fcom.Invoke, ops ...fcom.Option) *fcom.Result {
	chain.tempData[Invoke] = invoke
	chain.tempData[Option] = ops
	return &fcom.Result{
		Label:  "label",
		UID:    "UUID",
		Status: fcom.Success,
		Ret:    []interface{}{"Invoke", "Invoke"},
	}
}

func (chain *FakeChain) Transfer(transfer fcom.Transfer, ops ...fcom.Option) *fcom.Result {
	chain.tempData[Transfer] = transfer
	chain.tempData[Option] = ops

	return &fcom.Result{
		Label:  "label",
		UID:    "UUID",
		Status: fcom.Success,
		Ret:    []interface{}{"Transfer", "Transfer"},
	}
}

func (chain *FakeChain) Confirm(rt *fcom.Result, ops ...fcom.Option) *fcom.Result {
	chain.tempData[Result] = rt
	chain.tempData[Option] = ops
	return &fcom.Result{
		Label:  "Confirm",
		UID:    "UUID",
		Status: fcom.Confirm,
		Ret:    []interface{}{"Confirm", "Confirm"},
	}
}

func (chain *FakeChain) Query(bq fcom.Query, ops ...fcom.Option) interface{} {
	chain.tempData[Query] = bq
	chain.tempData[Option] = ops
	return "Query"
}

func (chain *FakeChain) Option(op fcom.Option) error {
	chain.tempData[Option] = op
	return nil
}

func TestGo2Lua(t *testing.T) {
	tempMap := map[uint]bool{'a': true, 'b': true}
	tempBytes := []byte{'a', 'b', 'c'}
	tempArray := []float64{9.321, 49.321, 0.432}
	tempFunc := func() {
		fmt.Println("gLua test")
	}
	var tempMap2 map[int]int
	var tempFunc2 func()
	tempLuaValue := lua.LBool(true)
	tempArr := [10]string{"array"}
	L := lua.NewState()
	mt := &lua.LTable{}
	mt.RawSetString("__index", L.NewFunction(indexFunc4Slice))
	mt.RawSetString("__tostring", L.NewFunction(lua2string))
	temp := lua.LUserData{Value: [][]string{{"11"}, {"22", "33"}, {"44", "55", "66"}}, Metatable: mt}
	tempUserData := lua.LUserData{Value: temp, Metatable: mt}
	tbl := L.NewTable()
	L.SetGlobal("test", tbl)
	L.SetField(tbl, "map", Go2Lua(L, tempMap))
	L.SetField(tbl, "map2", Go2Lua(L, tempMap2))
	L.SetField(tbl, "byte", Go2Lua(L, tempBytes))
	L.SetField(tbl, "float", Go2Lua(L, tempArray))
	L.SetField(tbl, "func", Go2Lua(L, tempFunc))
	L.SetField(tbl, "func2", Go2Lua(L, tempFunc2))
	L.SetField(tbl, "lua", Go2Lua(L, tempLuaValue))
	L.SetField(tbl, "luaUserData", Go2Lua(L, &tempUserData))
	L.SetField(tbl, "arr", Go2Lua(L, tempArr))
	defer L.Close()
	err := L.DoString(`
		print(test.map[97])
		print(test.byte[1])
		print(test.float[1])
		print(test.func())
		print(test.luaUserData)
		print(test.arr[1])
		print(test.byte["1"])
		print(test.arr["1"])
		print(test.float["1"])
		print(test.map["98"])

    `)
	assert.NoError(t, err)
	err = L.DoString(`
		print(test.byte[0])
    `)
	assert.Error(t, err)

	err = L.DoString(`
		print(test.arr[0])
    `)
	assert.Error(t, err)

	err = L.DoString(`
		print(test.float[0])
    `)
	assert.Error(t, err)

	err = L.DoString(`
		print(test.func("111"))
    `)
	assert.Error(t, err)
}

func TestLua2Go(t *testing.T) {
	L := lua.NewState()
	m := map[*lua.LTable]reflect.Value{}
	luaValue := lua.LBool(true)
	value, err := getGoValueReflect(L, luaValue, reflect.TypeOf(luaValue), nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, value)

	value, err = getGoValueReflect(L, luaValue, reflect.TypeOf(1), nil, nil)
	assert.Error(t, err)
	assert.NotNil(t, value)

	value, err = getGoValueReflect(L, luaValue, reflect.TypeOf(true), nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, value)

	luaNumber := lua.LNumber(22)
	value, err = getGoValueReflect(L, luaNumber, reflect.TypeOf(true), nil, nil)
	assert.Error(t, err)
	assert.NotNil(t, value)

	value, err = getGoValueReflect(L, lua.LNil, reflect.TypeOf(func() {}), nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, value)
	value, err = getGoValueReflect(L, lua.LNil, reflect.TypeOf(1), nil, nil)
	assert.Error(t, err)
	assert.NotNil(t, value)

	value, err = getGoValueReflect(L, L, reflect.TypeOf(func() {}), nil, nil)
	assert.Error(t, err)
	assert.NotNil(t, value)

	value, err = getGoValueReflect(L, lua.LString("qq"), reflect.TypeOf(func() {}), nil, nil)
	assert.Error(t, err)
	assert.NotNil(t, value)

	mt := &lua.LTable{}
	value, err = getGoValueReflect(L, mt, reflect.TypeOf([1]string{}), m, nil)
	assert.Error(t, err)
	assert.NotNil(t, value)
	mt.RawSetInt(1, lua.LString("1"))
	value, err = getGoValueReflect(L, mt, reflect.TypeOf([1]string{}), m, nil)
	assert.NoError(t, err)
	assert.NotNil(t, value)

	c, _ := NewMock()
	value, err = getGoValueReflect(L, mt, reflect.TypeOf(c), m, nil)
	assert.NoError(t, err)
	assert.NotNil(t, value)
}

func TestBlockchain(t *testing.T) {
	client, _ := NewMock()
	L := lua.NewState()
	tbl := L.NewTable()
	L.SetGlobal("test", tbl)
	L.SetField(tbl, "blockchain", Go2Lua(L, client))
	L.SetField(tbl, "toolkit", Go2Lua(L, toolkit.NewToolKit()))
	defer L.Close()
	err := L.DoString(`
		local ret=test.blockchain:Transfer({
			From = "0",
			To 	 = "1",
			Amount=1,
		})
		print(ret.Status)
		print(ret.Ret[1])
		local p1 = "1"
        local p2 = { "2", "2" }
        local p3 = { "3", "3", "3" }
        ret = test.blockchain:Invoke({
            Func = "typeUint128",
            Args = { p1, p2, p3 },
        })
		print(ret.Ret[1])
		ret=test.blockchain:Confirm(ret)
		print(ret.Status)
		print(ret.Ret[1])
		print(test.blockchain:Query({
            Func = "typeUint128",
            Args = { p1, p2, p3 },
        }))
		test.blockchain:Option({account=1})
		print(test.toolkit:Hex("111"))
		print(test.toolkit:RandStr(10))
		print(test.toolkit:RandInt(0,100))
    `)
	assert.NoError(t, err)
}
