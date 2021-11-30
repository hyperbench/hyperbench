package base

import (
	"testing"

	fcom "github.com/meshplus/hyperbench/common"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	b := NewBlockchainBase(ClientConfig{})
	assert.NotNil(t, b)

	err := b.DeployContract()
	assert.NoError(t, err)

	res := b.Confirm(&fcom.Result{})
	assert.NotNil(t, res)

	res = b.Invoke(bcom.Invoke{})
	assert.NotNil(t, res)

	res = b.Transfer(bcom.Transfer{})
	assert.NotNil(t, res)

	result := b.Query(bcom.Query{})
	assert.Nil(t, result)

	err = b.Option(bcom.Option{})
	assert.NoError(t, err)

	s, err := b.GetContext()
	assert.Equal(t, s, "")
	assert.NoError(t, err)

	err = b.SetContext("")
	assert.NoError(t, err)

	err = b.ResetContext()
	assert.NoError(t, err)

	rs, err := b.Statistic(bcom.Statistic{From: int64(0), To: int64(1)})
	assert.NotNil(t, rs)
	assert.NoError(t, err)

}
