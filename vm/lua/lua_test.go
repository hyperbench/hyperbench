
package lua

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/vm/base"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var tmpdir = "./tmp"

func TestLua(t *testing.T) {
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

	viper.Set(common.ClientScriptPath, scriptPath)

	configBase := base.ConfigBase{
		Ctx: common.VMContext{
			WorkerIdx: 0,
			VMIdx:     1,
		},
		Path: scriptPath,
	}

	vm, err := NewVM(base.NewVMBase(configBase))
	assert.NoError(t, err)

	viper.Set(common.ClientScriptPath, scriptPath2)
	configBase.Path = scriptPath2
	vm2, err := NewVM(base.NewVMBase(configBase))
	assert.NoError(t, err)

	viper.Set(common.ClientTypePath, "eth")
	ethvm, err := NewVM(base.NewVMBase(configBase))
	assert.Error(t, err)
	assert.Nil(t, ethvm)

	res, err := vm.Run(common.TxContext{
		Context: context.Background(),
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	res, err = vm2.Run(common.TxContext{
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

	rs, err := vm.Statistic(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	rs, err = vm2.Statistic(1, 1)
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

	vm.Close()
}

func BenchmarkLua(b *testing.B) {
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
	v.Set(common.ClientScriptPath, scriptPath)

	configBase := base.ConfigBase{
		Ctx: common.VMContext{
			WorkerIdx: 0,
			VMIdx:     1,
		},
	}

	vm, _ := NewVM(base.NewVMBase(configBase))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = vm.Run(common.TxContext{
			Context: context.Background(),
		})
	}
}
