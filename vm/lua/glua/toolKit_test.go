package glua

import (
	"fmt"
	"github.com/hyperbench/hyperbench/plugins/toolkit"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_toolKit(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	mt := L.NewTypeMetatable("case")
	L.SetGlobal("case", mt)
	L.SetField(mt, "toolkit", newToolKit(L, toolkit.NewToolKit()))

	t.Run("RandStr", func(t *testing.T) {
		scripts := []string{
			`
		function run()
		return case.toolkit.RandStr(10)
		end
		`,
		}
		for _, script := range scripts {
			lvalue, err := runLuaRunFunc(L, script)
			assert.Nil(t, err)
			result, err := Lua2Go(lvalue)
			assert.Nil(t, err)
			assert.Equal(t, len(fmt.Sprint(result)), 10)
		}
	})
	t.Run("RandStr", func(t *testing.T) {
		scripts := []string{
			`
		function run()
		return case.toolkit.RandInt(50,51)
		end
		`,
		}
		for _, script := range scripts {
			lvalue, err := runLuaRunFunc(L, script)
			assert.Nil(t, err)
			result, err := Lua2Go(lvalue)
			assert.Nil(t, err)
			assert.Equal(t, fmt.Sprint(result), "50")
		}
	})
	t.Run("Hex", func(t *testing.T) {
		scripts := []string{
			`
		function run()
		return case.toolkit.Hex("aaaaaa")
		end
		`,
		}
		for _, script := range scripts {
			lvalue, err := runLuaRunFunc(L, script)
			assert.Nil(t, err)
			result, err := Lua2Go(lvalue)
			assert.Nil(t, err)
			assert.Equal(t, fmt.Sprint(result), "616161616161")
		}

	})
	t.Run("String", func(t *testing.T) {
		toolkit.NewToolKit().String("")
		scripts := []string{
			`
		function run()
		return case.toolkit.String("aaaaaa",1,2)
		end
		`,
			`
		function run()
		return case.toolkit.String("aaaaaa")
		end
		`,
		}
		for _, script := range scripts {
			lvalue, err := runLuaRunFunc(L, script)
			assert.Nil(t, err)
			result, err := Lua2Go(lvalue)
			assert.Nil(t, err)
			assert.Equal(t, fmt.Sprint(result), "")
		}

	})

}
