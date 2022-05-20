package glua

import (
	"encoding/json"
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/mitchellh/mapstructure"
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
	return go2Lua(L, value)
}

// TableLua2GoStruct maps the lua table to the given struct pointer.
func TableLua2GoStruct(tbl *lua.LTable, st interface{}) error {
	value, err := Lua2Go(tbl)
	mp, ok := value.(map[string]interface{})
	if !ok {
		return errors.New("arguments #1 must be a table, but got an array")
	}
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           st,
		TagName:          "lua",
		ErrorUnused:      false,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(mp)
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

func go2Lua(L *lua.LState, value interface{}) lua.LValue {
	// check value is struct for Implementation lua.table
	luaValue, ok := go2luaStruct(L, value)
	if ok {
		return luaValue
	}
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case uint8:
		return lua.LNumber(converted)
	case []uint8:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(go2Lua(L, item))
		}
		return arr
	case [3]uint8:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(go2Lua(L, item))
		}
		return arr
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(go2Lua(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), go2Lua(L, item))
		}
		return tbl
	case nil:
		return lua.LNil
	}
	panic("unreachable")
}

// go2luaStruct convert struct for Implementation lua.table  to lua.Table
func go2luaStruct(L *lua.LState, value interface{}) (lua.LValue, bool) {
	switch value.(type) {
	case *fcom.Result:
		return newCommonResult(L, value.(*fcom.Result)), true
	default:
		return nil, false
	}
}

//	 function run()
//	    local i = 0
//	    print("----coro-----")
//	    return i
//    end
func runLuaRunFunc(state *lua.LState, script string) (lua.LValue, error) {
	//exec lua run func
	err := state.DoString(script)
	if err != nil {
		return nil, err
	}
	fn := state.GetGlobal("run").(*lua.LFunction)
	err = state.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: false,
	})
	if err != nil {
		return nil, err
	}
	ret := state.Get(-1)
	state.Pop(1)
	return ret, nil
}
