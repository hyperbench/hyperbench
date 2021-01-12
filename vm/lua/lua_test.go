package lua

import (
	"context"
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/vm/base"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
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
	_ = os.Mkdir(tmpdir, 0755)
	scriptPath := filepath.Join(tmpdir, "test.lua")
	_ = ioutil.WriteFile(scriptPath, []byte(script), 0644)
	// nolint
	defer os.RemoveAll(tmpdir)

	v := viper.New()
	v.Set(common.ClientScriptPath, scriptPath)

	configBase := base.ConfigBase{
		Ctx: common.VMContext{
			WorkerIdx: 0,
			VMIdx:     1,
		},
		Path: scriptPath,
	}

	ast := assert.New(t)
	vm, err := NewVM(base.NewVMBase(configBase))
	ast.NoError(err)
	res, err := vm.Run(common.TxContext{
		Context: context.Background(),
	})

	ast.NoError(err)
	ast.NotNil(res)
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
	_ = os.Mkdir(tmpdir, 0755)
	scriptPath := filepath.Join(tmpdir, "test.lua")
	_ = ioutil.WriteFile(scriptPath, []byte(script), 0644)
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
