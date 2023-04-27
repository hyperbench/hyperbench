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
 * @brief Go2Lua convert go values to lua.LValue
 * @file glua_convert.go
 * @author: wangxiaohui
 * @date 2022-03-30
 */

import (
	"fmt"
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"reflect"
)

const (
	indexKey   = "__index"
	fieldsKey  = "fields"
	methodsKey = "methods"
)

var (
	logger         = fcom.GetLogger("glua")
	invalidTypeErr = func(v lua.LValue, target reflect.Type) error {
		return errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
	}
)

//Go2Lua convert go values to lua.LValue
func Go2Lua(L *lua.LState, value interface{}) lua.LValue {
	lValue, err := go2Lua(L, value)
	if err != nil {
		logger.Errorf("convert go to lua err:%v", err)
		return lua.LNil
	}
	return lValue
}

func go2Lua(L *lua.LState, value interface{}) (lua.LValue, error) {
	if value == nil {
		return lua.LNil, nil
	}
	if luaValue, ok := value.(lua.LValue); ok {
		return luaValue, nil
	}
	switch val := reflect.ValueOf(value); val.Kind() {
	case reflect.Bool:
		return lua.LBool(val.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(float64(val.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(val.Float()), nil
	case reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		if val.IsNil() {
			return lua.LNil, nil
		}
		return getUserData(L, val)
	case reflect.Array, reflect.Struct:
		return getUserData(L, val)
	case reflect.Func:
		if val.IsNil() {
			return lua.LNil, nil
		}
		return go2LuaFunction(L, val, false), nil
	case reflect.String:
		return lua.LString(val.String()), nil
	default:
		ud := L.NewUserData()
		ud.Value = val.Interface()
		return ud, nil
	}
}

// convert2MetaTbl generates metatable of specific type
func convert2MetaTbl(L *lua.LState, valueType reflect.Type) (*lua.LTable, error) {
	mt := &lua.LTable{}
	methods := L.CreateTable(0, valueType.NumMethod())
	switch valueType.Kind() {
	case reflect.Array:
		mt = createTblAndSetRaw(L, 7, indexKey, indexFunc4Array)
	case reflect.Slice:
		mt = createTblAndSetRaw(L, 8, indexKey, indexFunc4Slice)
	case reflect.Struct:
		mt = L.CreateTable(0, 6)
		fields := L.CreateTable(0, valueType.NumField())
		injectFields(L, valueType, fields)
		mt.RawSetString(fieldsKey, fields)
	case reflect.Map:
		mt = createTblAndSetRaw(L, 7, indexKey, indexFunc4Map)
	case reflect.Ptr:
		switch valueType.Elem().Kind() {
		case reflect.Struct:
			mt = createTblAndSetRaw(L, 8, indexKey, indexFunc4Struct)
		}
		injectMethods(L, valueType, methods, true)
	default:
		return nil, fmt.Errorf("unexpected kind %v", valueType.Kind().String())
	}
	mt.RawSetString(methodsKey, methods)
	mt.RawSetString("__tostring", L.NewFunction(lua2string))
	return mt, nil
}

func createTblAndSetRaw(L *lua.LState, hcap int, key string, fn func(L *lua.LState) int) *lua.LTable {
	mt := L.CreateTable(0, hcap)
	mt.RawSetString(key, L.NewFunction(fn))
	return mt
}

// lua2string is the LuaFunction converts luaUserData to string
func lua2string(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if stringer, ok := ud.Value.(fmt.Stringer); ok {
		L.Push(lua.LString(stringer.String()))
		return 1
	}
	L.Push(lua.LString(ud.String()))
	return 1
}

// go2LuaFunction converts go function to LuaFunction
func go2LuaFunction(L *lua.LState, fn reflect.Value, isPtrReceiverMethod bool) *lua.LFunction {
	up := L.NewUserData()
	up.Value = fn
	return L.NewClosure(getLuaFunction, up, lua.LBool(isPtrReceiverMethod))
}

// getLuaFunction is the regular LuaFunction for converting go function to LuaFunction
func getLuaFunction(L *lua.LState) int {
	ref := L.Get(lua.UpvalueIndex(1)).(*lua.LUserData).Value.(reflect.Value)
	refType := ref.Type()
	// check arguments of function
	top := L.GetTop()
	expected := refType.NumIn()
	variadic := refType.IsVariadic()
	if !variadic && top != expected {
		L.RaiseError("invalid number of function arguments (%d expected, got %d)", expected, top)
	}
	if variadic && top < expected-1 {
		L.RaiseError("invalid number of function arguments (%d or more expected, got %d)", expected-1, top)
	}

	convertedPtr := false
	var receiver reflect.Value
	var ud lua.LValue

	args := make([]reflect.Value, top)
	for i := 0; i < L.GetTop(); i++ {
		var target reflect.Type
		if variadic && i >= expected-1 {
			target = refType.In(expected - 1).Elem()
		} else {
			target = refType.In(i)
		}
		var arg reflect.Value
		var err error
		if i == 0 && bool(L.Get(lua.UpvalueIndex(2)).(lua.LBool)) {
			ud = L.Get(1)
			v := ud
			// record converted luaValue to avoid repeated work
			convertedRecord := make(map[*lua.LTable]reflect.Value)
			arg, err = getGoValueReflect(L, v, target, convertedRecord, &convertedPtr)
			if err != nil {
				L.ArgError(1, err.Error())
			}
			receiver = arg
		} else {
			v := L.Get(i + 1)
			// record converted luaValue to avoid repeated work
			convertedRecord := make(map[*lua.LTable]reflect.Value)
			arg, err = getGoValueReflect(L, v, target, convertedRecord, nil)
			if err != nil {
				L.ArgError(i+1, err.Error())
			}
		}
		args[i] = arg
	}
	ret := ref.Call(args)

	if convertedPtr {
		ud.(*lua.LUserData).Value = receiver.Elem().Interface()
	}

	for _, val := range ret {
		L.Push(Go2Lua(L, val.Interface()))
	}
	return len(ret)
}

// getGoValueReflect returns go reflect of luaValue
func getGoValueReflect(L *lua.LState, v lua.LValue, target reflect.Type, convertedRecord map[*lua.LTable]reflect.Value, tryConvertPtr *bool) (reflect.Value, error) {
	// check if target type is luaValue
	if target.Implements(reflect.TypeOf((*lua.LValue)(nil)).Elem()) {
		return reflect.ValueOf(v), nil
	}

	switch converted := v.(type) {
	case lua.LString:
		val := reflect.ValueOf(string(converted))
		if !val.Type().ConvertibleTo(target) {
			return reflect.Value{}, invalidTypeErr(v, target)
		}
		return val.Convert(target), nil
	case lua.LBool, lua.LNumber, *lua.LState:
		return convertToTargetType(converted, target, v)
	case *lua.LNilType:
		switch target.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer, reflect.Uintptr:
			return reflect.Zero(target), nil
		}
		return reflect.Value{}, invalidTypeErr(v, target)
	case *lua.LTable:
		// if same value has been converted return the record
		if record := convertedRecord[converted]; record.IsValid() {
			return record, nil
		}
		// check if target is map[interface{}]interface{}
		if emptyInterfaceType := reflect.TypeOf((*interface{})(nil)).Elem(); target == emptyInterfaceType {
			target = reflect.MapOf(emptyInterfaceType, emptyInterfaceType)
		}

		switch {
		case target.Kind() == reflect.Array:
			return convertTblToArrayOrSlice(L, target, converted, convertedRecord, func(target reflect.Type, length int) (reflect.Value, error) {
				if length != target.Len() {
					return reflect.Value{}, invalidTypeErr(v, target)
				}
				return reflect.New(target).Elem(), nil
			})
		case target.Kind() == reflect.Slice:
			return convertTblToArrayOrSlice(L, target, converted, convertedRecord, func(target reflect.Type, length int) (reflect.Value, error) {
				return reflect.MakeSlice(target, length, length), nil
			})
		case target.Kind() == reflect.Map:
			return convertTblToMap(L, target, converted, convertedRecord)
		case target.Kind() == reflect.Ptr && target.Elem().Kind() == reflect.Struct:
			target = target.Elem()
			_, p, err := convertTblToStruct(L, target, converted, convertedRecord)
			return p, err
		case target.Kind() == reflect.Struct:
			s, _, err := convertTblToStruct(L, target, converted, convertedRecord)
			return s, err
		default:
			return reflect.Value{}, invalidTypeErr(v, target)
		}
	case *lua.LUserData:
		val := reflect.ValueOf(converted.Value)
		if tryConvertPtr != nil && val.Kind() != reflect.Ptr && target.Kind() == reflect.Ptr && val.Type() == target.Elem() {
			goValue := reflect.New(target.Elem())
			goValue.Elem().Set(val)
			val = goValue
			*tryConvertPtr = true
		} else {
			if !val.Type().ConvertibleTo(target) {
				return reflect.Value{}, invalidTypeErr(v, target)
			}
			val = val.Convert(target)
			if tryConvertPtr != nil {
				*tryConvertPtr = false
			}
		}
		return val, nil
	default:
		return reflect.Value{}, errors.Errorf("invalid type:%v", v.Type())
	}
}

// convertTblToStruct the first reflect.Value is struct, the second reflect.Value is point.
func convertTblToStruct(L *lua.LState, target reflect.Type, converted *lua.LTable, convertedRecord map[*lua.LTable]reflect.Value) (reflect.Value, reflect.Value, error) {
	s := reflect.New(target)
	convertedRecord[converted] = s
	t := s.Elem()

	for key := lua.LNil; ; {
		var value lua.LValue
		key, value = converted.Next(key)
		if key == lua.LNil {
			break
		}
		if _, ok := key.(lua.LString); !ok {
			continue
		}

		fieldName := key.String()
		field, exist := target.FieldByName(fieldName)
		if !exist {
			return reflect.Value{}, reflect.Value{}, errors.Errorf("type %s has no field %s", target, fieldName)
		}

		lValue, err := getGoValueReflect(L, value, field.Type, convertedRecord, nil)
		if err != nil {
			return reflect.Value{}, reflect.Value{}, err
		}
		t.FieldByIndex(field.Index).Set(lValue)
	}
	return t, s, nil
}

func convertTblToMap(L *lua.LState, target reflect.Type, converted *lua.LTable, convertedRecord map[*lua.LTable]reflect.Value) (reflect.Value, error) {
	keyType := target.Key()
	elemType := target.Elem()
	s := reflect.MakeMap(target)
	convertedRecord[converted] = s

	for key := lua.LNil; ; {
		var value lua.LValue
		key, value = converted.Next(key)
		if key == lua.LNil {
			break
		}

		goKey, err := getGoValueReflect(L, key, keyType, convertedRecord, nil)
		if err != nil {
			return reflect.Value{}, err
		}
		goValue, err := getGoValueReflect(L, value, elemType, convertedRecord, nil)
		if err != nil {
			return reflect.Value{}, err
		}
		s.SetMapIndex(goKey, goValue)
	}
	return s, nil
}

func convertToTargetType(convertValue interface{}, target reflect.Type, v lua.LValue) (reflect.Value, error) {
	val := reflect.ValueOf(convertValue)
	if !val.Type().ConvertibleTo(target) {
		return reflect.Value{}, invalidTypeErr(v, target)
	}
	return val.Convert(target), nil
}

func getUserData(L *lua.LState, val reflect.Value) (lua.LValue, error) {
	ud := L.NewUserData()
	ud.Value = val.Interface()
	metaTbl, err := convert2MetaTbl(L, val.Type())
	if err != nil {
		return nil, err
	}
	ud.Metatable = metaTbl
	return ud, nil
}

func convertTblToArrayOrSlice(L *lua.LState, target reflect.Type, converted *lua.LTable, convertedRecord map[*lua.LTable]reflect.Value, fn func(reflect.Type, int) (reflect.Value, error)) (reflect.Value, error) {
	elemType := target.Elem()
	length := converted.Len()
	s, err := fn(target, length)
	if err != nil {
		return reflect.Value{}, err
	}
	convertedRecord[converted] = s

	for i := 0; i < length; i++ {
		value := converted.RawGetInt(i + 1)
		elemValue, err := getGoValueReflect(L, value, elemType, convertedRecord, nil)
		if err != nil {
			return reflect.Value{}, err
		}
		s.Index(i).Set(elemValue)
	}

	return s, nil
}
