/*
@Time : 2019-04-12 14:13
@Author : lmm
*/
package fabric

import (
	"testing"

	fcom "github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"

	"github.com/stretchr/testify/assert"
)

func TestFabric(t *testing.T) {
	t.Skip()
	//new
	op := make(map[string]interface{})
	op["channel"] = "mychannel"
	op["MSP"] = false
	op["instant"] = 2
	op["ccID"] = "1"
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "fabric",
		ConfigPath:   "./../../../benchmark/fabricExample/fabric",
		ContractPath: "github.com/meshplus/hyperbench/benchmark/fabricExample/contract",
		Args:         []interface{}{"init", "A", "123", "B", "234"},
		Options:      op,
	})
	client, err := New(b)
	assert.NotNil(t, client)
	assert.NoError(t, err)

	//deploy
	err = client.DeployContract()
	assert.NoError(t, err)

	//getContext
	context, err := client.GetContext()
	assert.NoError(t, err)
	assert.NotNil(t, context)

	//setContext
	err = client.SetContext(context)
	assert.NoError(t, err)

	//invoke
	txResult := client.Invoke(bcom.Invoke{Func: "query", Args: []interface{}{"A"}})
	assert.Equal(t, txResult.Status, fcom.Success)

	txResult = client.Invoke(bcom.Invoke{Func: "query", Args: []interface{}{"A", "B"}})
	assert.Equal(t, txResult.Status, fcom.Failure)

	client.invoke = false
	txResult = client.Invoke(bcom.Invoke{Func: "query", Args: []interface{}{"A"}})
	assert.Equal(t, txResult.Status, fcom.Success)

	//reset
	err = client.ResetContext()
	assert.NoError(t, err)

	//statistic
	res, err := client.Statistic(bcom.Statistic{From: int64(0), To: int64(1)})
	assert.NotNil(t, res)
	assert.NoError(t, err)

	//string
	s := client.String()
	assert.NotNil(t, s)

	//option
	err = client.Option(bcom.Option{"mode": "query"})
	assert.NoError(t, err)

	err = client.Option(bcom.Option{"mode": "invoke"})
	assert.NoError(t, err)

}
