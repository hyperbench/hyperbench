package base

import (
	"testing"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/stretchr/testify/assert"
)

func TestBaseVm(t *testing.T) {
	t.Skip()
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

	res, err := base.Statistic(nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	err = base.BeforeSet()
	assert.NoError(t, err)

	err = base.SetContext(nil)
	assert.NoError(t, err)

	err = base.BeforeRun()
	assert.NoError(t, err)

	result, err := base.Run(fcom.TxContext{})
	assert.NoError(t, err)
	assert.NotNil(t, result)

	err = base.AfterRun()
	assert.NoError(t, err)

	_, err = base.LogStatus()
	assert.NoError(t, err)

	base.Close()

}
