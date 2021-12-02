package base

import (
	"testing"

	"github.com/meshplus/hyperbench/common"
	"github.com/stretchr/testify/assert"
)

func TestBaseVm(t *testing.T) {
	base := NewVMBase(ConfigBase{
		Path: "",
	})
	Type := base.Type()
	assert.Equal(t, Type, "base")

	err := base.BeforeDeploy()
	assert.NoError(t, err)

	err = base.DeployContract()
	assert.NoError(t, err)

	err = base.BeforeGet()
	assert.NoError(t, err)

	bs, err := base.GetContext()
	assert.NoError(t, err)
	assert.NotNil(t, bs)

	res, err := base.Statistic(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	err = base.BeforeSet()
	assert.NoError(t, err)

	err = base.SetContext(nil)
	assert.NoError(t, err)

	err = base.BeforeRun()
	assert.NoError(t, err)

	result, err := base.Run(common.TxContext{})
	assert.NoError(t, err)
	assert.NotNil(t, result)

	err = base.AfterRun()
	assert.NoError(t, err)

	base.Close()

}
