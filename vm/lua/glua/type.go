package glua

import (
	"encoding/json"
	"github.com/meshplus/hyperbench/common"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"reflect"
)

//Go2Lua convert go interface val to lua.LValue value and reutrn
func Go2Lua(L *lua.LState, val interface{}) lua.LValue {
	var (
		jsonBytes []byte
		err       error
	)
	if jsonBytes, err = json.Marshal(val); err != nil {
		return lua.LNil
	}
	// parse json to lua type
	var value interface{}
	if err = json.Unmarshal(jsonBytes, &value); err != nil {
		return lua.LNil
	}
	return decode(L, value)
}

// Lua2Go convert lua.LValue to corresponding type
func Lua2Go(value lua.LValue) (interface{}, error) {
	switch val := value.(type) {
	case lua.LString:
		return string(val), nil
	case lua.LNumber:
		return float64(val), nil
	case lua.LBool:
		return bool(val), nil
	case *lua.LTable:
		l := val.Len()
		if l > 0 {
			// process as []interface{}
			s := make([]interface{}, l)
			var err error

			for i := 0; i < l; i++ {
				s[i], err = Lua2Go(val.RawGetInt(i + 1))
				if err != nil {
					return nil, err
				}
			}
			return s, err
		}

		// process as map[string]interface{}
		m := make(map[string]interface{})
		var err error
		val.ForEach(func(k lua.LValue, v lua.LValue) {
			// there is no shortcut for `ForEach`
			// just use err check to reduce useless operation
			if strKey, ok := k.(lua.LString); ok && err == nil {
				m[strKey.String()], err = Lua2Go(v)
			}
		})

		// if err is set in cb of `ForEach`
		// it will return here
		if err != nil {
			return nil, err
		}

		return m, nil

	case *lua.LUserData:
		return val.Value, nil
	default:
		return nil, errors.Errorf("do not support type: %v", reflect.TypeOf(val).String())
	}
}

//LuaParams convert lua.LValue to string or string slice or user data interface
func LuaParams(value lua.LValue) interface{} {
	switch val := value.(type) {
	case lua.LString:
		return val.String()
	case lua.LNumber:
		return int(val)
	case lua.LBool:
		return bool(val)
	case *lua.LTable:
		arr := make([]interface{}, val.Len())
		for j := 1; j <= val.Len(); j++ {
			val := val.RawGetInt(j)
			arr[j-1] = LuaParams(val)
		}
		return arr
	case *lua.LUserData:
		return val.Value
	default:
		return ""
	}
}

func decode(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(decode(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), decode(L, item))
		}
		return tbl
	case *common.Result:
		return newCommonResult(L, value.(*common.Result))
	case nil:
		return lua.LNil
	}
	panic("unreachable")
}
