package glua

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func Test_Go2Lua(t *testing.T) {

	t.Run("demo", func(t *testing.T) {
		str :="demo"
		L := lua.NewState()
		strLua := Go2Lua(L,str)
		assert.Equal(t, strLua,lua.LString("demo"))
	})
}
