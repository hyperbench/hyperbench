package glua

/**
 *  Copyright (C) 2021 HyperBench.
 *  SPDX-License-Identifier: Apache-2.0
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 * @brief Access the LuaFunction converted from go
 * @file index.go
 * @author: wangxiaohui
 * @date 2022-03-30
 */

import (
	lua "github.com/yuin/gopher-lua"
	"reflect"
)

// injectMethods injects methods of given struct in lua
func injectMethods(L *lua.LState, valueType reflect.Type, tbl *lua.LTable, ptrReceiver bool) {
	for i := 0; i < valueType.NumMethod(); i++ {
		method := valueType.Method(i)
		fn := go2LuaFunction(L, method.Func, ptrReceiver)
		tbl.RawSetString(method.Name, fn)
	}
}

// injectFields injects fields of given struct in lua
func injectFields(L *lua.LState, valueType reflect.Type, tbl *lua.LTable) {
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		ud := L.NewUserData()
		ud.Value = field.Index
		tbl.RawSetString(field.Name, ud)
	}
}

// indexFunc4Map is the LuaFunction for accessing map converted from go
func indexFunc4Map(L *lua.LState) int {
	ref, mt := getParameter(L, 1)
	key := L.CheckAny(2)
	// record converted luaValue to avoid repeated work
	convertedRecord := make(map[*lua.LTable]reflect.Value)
	convertedKey, err := getGoValueReflect(L, key, ref.Type().Key(), convertedRecord, nil)
	if err == nil {
		item := ref.MapIndex(convertedKey)
		if item.IsValid() {
			L.Push(Go2Lua(L, item.Interface()))
			return 1
		}
	}

	if luaString, ok := key.(lua.LString); ok {
		if fn := getFunc(mt, string(luaString)); fn != nil {
			L.Push(fn)
			return 1
		}
	}

	return 0
}

// indexFunc4Array is the LuaFunction for accessing array converted from go
func indexFunc4Array(L *lua.LState) int {
	ref, mt := getParameter(L, 1)
	ref = reflect.Indirect(ref)
	key := L.CheckAny(2)
	return indexFunc4ArrayOrSlice(L, ref, key, mt)
}

func indexFunc4ArrayOrSlice(L *lua.LState, ref reflect.Value, key lua.LValue, mt *lua.LTable) int {
	switch converted := key.(type) {
	case lua.LNumber:
		index := int(converted)
		if index < 1 || index > ref.Len() {
			logger.Errorf("index[%v] out of range %v", index, ref.Len())
			L.ArgError(2, "index out of range")
		}
		value := ref.Index(index - 1)
		if (value.Kind() == reflect.Struct || value.Kind() == reflect.Array) && value.CanAddr() {
			value = value.Addr()
		}
		L.Push(Go2Lua(L, value.Interface()))
	case lua.LString:
		if fn := getFunc(mt, string(converted)); fn != nil {
			L.Push(fn)
			return 1
		}
		return 0
	default:
		logger.Errorf("invalid type:%v, support type is: a number or string", key.Type())
		L.ArgError(2, "must be a number or string")
	}
	return 1
}

// indexFunc4Slice is the LuaFunction for accessing slice converted from go
func indexFunc4Slice(L *lua.LState) int {
	ref, mt := getParameter(L, 1)
	key := L.CheckAny(2)
	return indexFunc4ArrayOrSlice(L, ref, key, mt)
}

// indexFunc4Struct is the LuaFunction for accessing struct converted from go
func indexFunc4Struct(L *lua.LState) int {
	ref, mt := getParameter(L, 1)
	key := L.CheckString(2)

	if fn := getFunc(mt, key); fn != nil {
		L.Push(fn)
		return 1
	}

	ref = ref.Elem()
	switch typ := reflect.TypeOf(ref.Interface()); typ.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct:
		mt, _ = convert2MetaTbl(L, typ)
	}
	index := getFieldIndex(mt, key)
	if index == nil {
		return 0
	}
	field := ref.FieldByIndex(index)
	if !field.CanInterface() {
		L.RaiseError("cannot interface field " + key)
	}
	L.Push(Go2Lua(L, field.Interface()))
	return 1
}

// getParameter returns value in the stack of luaState by index
func getParameter(L *lua.LState, index int) (ref reflect.Value, mt *lua.LTable) {
	ud := L.CheckUserData(index)
	ref = reflect.ValueOf(ud.Value)
	mt = ud.Metatable.(*lua.LTable)
	return
}

// getFieldIndex returns indexes of fields of go struct
func getFieldIndex(m *lua.LTable, name string) []int {
	fields := m.RawGetString(fieldsKey).(*lua.LTable)
	if index := fields.RawGetString(name); index != lua.LNil {
		return index.(*lua.LUserData).Value.([]int)
	}
	return nil
}

// getFunc returns function of specific key in map、slice、array or struct
func getFunc(m *lua.LTable, name string) lua.LValue {
	methods := m.RawGetString(methodsKey).(*lua.LTable)
	if fn := methods.RawGetString(name); fn != lua.LNil {
		return fn
	}
	return nil
}
