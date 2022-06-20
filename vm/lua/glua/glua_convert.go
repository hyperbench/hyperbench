package glua

import (
	"fmt"
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"reflect"
)

//Go2Lua convert go values to lua.LValue
func Go2Lua(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}
	if luaValue, ok := value.(lua.LValue); ok {
		return luaValue
	}
	switch val := reflect.ValueOf(value); val.Kind() {
	case reflect.Bool:
		return lua.LBool(val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(float64(val.Uint()))
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(val.Float())
	case reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		if val.IsNil() {
			return lua.LNil
		}
		fallthrough
	case reflect.Array, reflect.Struct:
		ud := L.NewUserData()
		ud.Value = val.Interface()
		ud.Metatable = convert2MetaTbl(L, val.Type())
		return ud
	case reflect.Func:
		if val.IsNil() {
			return lua.LNil
		}
		return go2LuaFunction(L, val, false)
	case reflect.String:
		return lua.LString(val.String())
	default:
		ud := L.NewUserData()
		ud.Value = val.Interface()
		return ud
	}
}

// convert2MetaTbl generates metatable of specific type
func convert2MetaTbl(L *lua.LState, valueType reflect.Type) *lua.LTable {
	mt := &lua.LTable{}
	methods := L.CreateTable(0, valueType.NumMethod())
	switch valueType.Kind() {
	case reflect.Array:
		mt = L.CreateTable(0, 7)
		mt.RawSetString("__index", L.NewFunction(indexFunc4Array))
	case reflect.Slice:
		mt = L.CreateTable(0, 8)
		mt.RawSetString("__index", L.NewFunction(indexFunc4Slice))
	case reflect.Struct:
		mt = L.CreateTable(0, 6)
		fields := L.CreateTable(0, valueType.NumField())
		injectFields(L, valueType, fields)
		mt.RawSetString("fields", fields)
	case reflect.Map:
		mt = L.CreateTable(0, 7)
		mt.RawSetString("__index", L.NewFunction(indexFunc4Map))
	case reflect.Ptr:
		switch valueType.Elem().Kind() {
		case reflect.Struct:
			mt = L.CreateTable(0, 8)
			mt.RawSetString("__index", L.NewFunction(indexFunc4Struct))
		}
		injectMethods(L, valueType, methods, true)
	default:
		panic("unexpected kind " + valueType.Kind().String())
	}
	mt.RawSetString("methods", methods)
	mt.RawSetString("__tostring", L.NewFunction(lua2string))
	return mt
}

// lua2string is the LuaFunction converts luaUserData to string
func lua2string(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if stringer, ok := ud.Value.(fmt.Stringer); ok {
		L.Push(lua.LString(stringer.String()))
	} else {
		L.Push(lua.LString(ud.String()))
	}
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

	isPtr := false

	switch converted := v.(type) {
	case lua.LBool:
		val := reflect.ValueOf(bool(converted))
		if !val.Type().ConvertibleTo(target) {
			return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
		}
		return val.Convert(target), nil
	case lua.LNumber:
		val := reflect.ValueOf(float64(converted))
		if !val.Type().ConvertibleTo(target) {
			return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
		}
		return val.Convert(target), nil
	case *lua.LNilType:
		switch target.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer, reflect.Uintptr:
			return reflect.Zero(target), nil
		}
		return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
	case *lua.LState:
		val := reflect.ValueOf(converted)
		if !val.Type().ConvertibleTo(target) {
			return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
		}
		return val.Convert(target), nil
	case lua.LString:
		val := reflect.ValueOf(string(converted))
		if !val.Type().ConvertibleTo(target) {
			return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
		}
		return val.Convert(target), nil
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
			elemType := target.Elem()
			length := converted.Len()
			if length != target.Len() {
				return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
			}
			s := reflect.New(target).Elem()
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
		case target.Kind() == reflect.Slice:
			elemType := target.Elem()
			length := converted.Len()
			s := reflect.MakeSlice(target, length, length)
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

		case target.Kind() == reflect.Map:
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

		case target.Kind() == reflect.Ptr && target.Elem().Kind() == reflect.Struct:
			target = target.Elem()
			isPtr = true
			fallthrough
		case target.Kind() == reflect.Struct:
			s := reflect.New(target)
			convertedRecord[converted] = s
			t := s.Elem()

			mt := convert2MetaTbl(L, target)

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
				index := getFieldIndex(mt, fieldName)
				if index == nil {
					return reflect.Value{}, errors.Errorf("type %s has no field %s", target, fieldName)
				}
				field := target.FieldByIndex(index)

				lValue, err := getGoValueReflect(L, value, field.Type, convertedRecord, nil)
				if err != nil {
					return reflect.Value{}, nil
				}
				t.FieldByIndex(field.Index).Set(lValue)
			}
			if isPtr {
				return s, nil
			}
			return t, nil
		}
		return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
	case *lua.LUserData:
		val := reflect.ValueOf(converted.Value)
		if tryConvertPtr != nil && val.Kind() != reflect.Ptr && target.Kind() == reflect.Ptr && val.Type() == target.Elem() {
			goValue := reflect.New(target.Elem())
			goValue.Elem().Set(val)
			val = goValue
			*tryConvertPtr = true
		} else {
			if !val.Type().ConvertibleTo(target) {
				return reflect.Value{}, errors.Errorf("cannot use %v (type %T) as type %s", v, v.Type(), target)
			}
			val = val.Convert(target)
			if tryConvertPtr != nil {
				*tryConvertPtr = false
			}
		}
		return val, nil
	}

	panic("never reaches")
}
