package toolkit

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
 * @brief ToolKit is the set of tool plugins
 * @file toolkit.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
)

// ToolKit is the set of tool plugins.
// exported methods and fields can be inject into vm.
type ToolKit struct {
	seed *rand.Rand
}

// NewToolKit create a new ToolKit instant as a plugin to vm
func NewToolKit() *ToolKit {
	r := rand.New(rand.NewSource(1))
	return &ToolKit{
		r,
	}
}

// SetRandSeed reset the rand seed of randStr
func (t *ToolKit) SetRandSeed(s int64) {
	t.seed = rand.New(rand.NewSource(s))
}

// RandStr generate a random string with specific length
func (t *ToolKit) RandStr(l uint) string {
	return randomString(l)
}

// RandStr generate a random string with specific length using seed setted
func (t *ToolKit) RandStrSeed(l uint) string {
	return randomStringWithSeed(t.seed, l)
}

// RandInt generate a random int in specific range
func (t *ToolKit) RandInt(min, max int) int {
	return randomInt(min, max)
}

// String convert to string
func (t *ToolKit) String(input interface{}, offsets ...int) string {

	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Ptr:
		return t.String(v.Elem().Interface(), offsets...)
	case reflect.Slice:
		bs, ok := input.([]byte)
		if !ok {
			return ""
		}
		l := len(bs)
		start, end := 0, l
		if len(offsets) > 0 && offsets[0] <= l {
			start = offsets[0]
		}
		if len(offsets) > 1 && offsets[1] <= l+1 {
			start = offsets[1]
		}
		fmt.Println(bs, start, end)
		return string(bs[start:end])
	case reflect.Array:
		l := v.Type().Len()
		start, end := 0, l
		if len(offsets) > 0 && offsets[0] <= l {
			start = offsets[0]
		}
		if len(offsets) > 1 && offsets[1] <= l+1 {
			start = offsets[1]
		}
		ss := make([]byte, 0, end-start)
		for i := start; i < end; i++ {
			ss = append(ss, v.Index(i).Interface().(byte))
		}
		return string(ss)
	}
	fmt.Println("=====")
	return ""
}

// Hex encode string as hex string
func (t *ToolKit) Hex(input string) string {
	return hex.EncodeToString([]byte(input))
}
