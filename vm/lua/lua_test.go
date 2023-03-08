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
 * @brief use cobra provide cmd function
 * @file lua_test.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */
package lua

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var tmpdir = "./tmp"

func TestLua(t *testing.T) {
	t.Skip()
	script := `
case = testcase.new()

function case:Run(ctx) 
	print(case.toolkit:RandStr(10))
    ret = case.blockchain:Invoke({
    	func="123",
		args={"123", "123"},
    })
	case.blockchain:Option({
    	func="123",
		args={"123", "123"},
    })
	return ret
end

return case
`
	script2 := `
case = testcase.new()

function case:BeforeDeploy()
end

function case:BeforeGet()
end

function case:BeforeSet()
end

function case:BeforeRun()
end

function case:AfterRun()
end
function case:Run(ctx)
end
return case
	`
	os.Mkdir(tmpdir, 0755)
	scriptPath := filepath.Join(tmpdir, "test.lua")
	scriptPath2 := filepath.Join(tmpdir, "test2.lua")
	ioutil.WriteFile(scriptPath, []byte(script), 0644)
	ioutil.WriteFile(scriptPath2, []byte(script2), 0644)
	// nolint
	defer os.RemoveAll(tmpdir)

	viper.Set(fcom.ClientScriptPath, scriptPath)

	configBase := base.ConfigBase{
		Ctx: fcom.VMContext{
			WorkerIdx: 0,
			VMIdx:     1,
		},
		Path: scriptPath,
	}

	vm, err := NewVM(base.NewVMBase(configBase))
	assert.NoError(t, err)

	viper.Set(fcom.ClientScriptPath, scriptPath2)
	configBase.Path = scriptPath2
	vm2, err := NewVM(base.NewVMBase(configBase))
	assert.NoError(t, err)

	viper.Set(fcom.ClientTypePath, "eth")
	ethvm, err := NewVM(base.NewVMBase(configBase))
	assert.Error(t, err)
	assert.Nil(t, ethvm)

	res, err := vm.Run(fcom.TxContext{
		Context: context.Background(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	res, err = vm2.Run(fcom.TxContext{
		Context: context.Background(),
	})
	assert.Error(t, err)
	assert.Nil(t, res)

	err = vm.BeforeDeploy()
	assert.NoError(t, err)
	err = vm2.BeforeDeploy()
	assert.NoError(t, err)

	err = vm.DeployContract()
	assert.NoError(t, err)
	err = vm2.DeployContract()
	assert.NoError(t, err)

	err = vm.BeforeGet()
	assert.NoError(t, err)
	err = vm2.BeforeGet()
	assert.NoError(t, err)

	bs, err := vm.GetContext()
	assert.NoError(t, err)
	assert.NotNil(t, bs)
	bs, err = vm2.GetContext()
	assert.NoError(t, err)
	assert.NotNil(t, bs)

	rs, err := vm.Statistic(&fcom.ChainInfo{TimeStamp: 1}, &fcom.ChainInfo{TimeStamp: 1})
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	rs, err = vm2.Statistic(&fcom.ChainInfo{TimeStamp: 1}, &fcom.ChainInfo{TimeStamp: 1})
	assert.NoError(t, err)
	assert.NotNil(t, rs)

	err = vm.BeforeSet()
	assert.NoError(t, err)
	err = vm2.BeforeSet()
	assert.NoError(t, err)

	err = vm.SetContext(nil)
	assert.NoError(t, err)
	err = vm2.SetContext(nil)
	assert.NoError(t, err)

	err = vm.BeforeRun()
	assert.NoError(t, err)
	err = vm2.BeforeRun()
	assert.NoError(t, err)

	err = vm.AfterRun()
	assert.NoError(t, err)
	err = vm2.AfterRun()
	assert.NoError(t, err)

	_, err = vm.LogStatus()
	assert.NoError(t, err)

	vm.Close()
}

func BenchmarkLua(b *testing.B) {
	b.Skip()
	script := `
case = testcase.new()

function case:Run(ctx)
    ret = case.blockchain:Invoke({
       func="123",
	   args={"123", "123"},
    })
    --case.blockchain:Option({
    --    func="123",
	--    args={"123", "123"},
    --})
	-- ret = case.blockchain:Confirm(ret)
	return {}
end

return case
`
	os.Mkdir(tmpdir, 0755)
	scriptPath := filepath.Join(tmpdir, "test.lua")
	ioutil.WriteFile(scriptPath, []byte(script), 0644)
	// nolint
	defer os.RemoveAll(tmpdir)

	v := viper.New()
	v.Set(fcom.ClientScriptPath, scriptPath)

	configBase := base.ConfigBase{
		Ctx: fcom.VMContext{
			WorkerIdx: 0,
			VMIdx:     1,
		},
	}

	vm, _ := NewVM(base.NewVMBase(configBase))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(fcom.TxContext{
			Context: context.Background(),
		})
	}
}
