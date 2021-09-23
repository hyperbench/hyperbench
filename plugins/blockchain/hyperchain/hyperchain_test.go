package hyperchain

import (
	"fmt"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"testing"
)

func TestTransfer(t *testing.T) {
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType: "hyperchain",
		ConfigPath: "./../../../benchmark/local/hyperchain",
		//ConfigPath: "benchmark/new/hpc",
	})
	hpc, _ := NewClient(b)
	ret := hpc.Transfer(bcom.Transfer{
		From:   "0",
		To:     "1",
		Amount: 0,
		Extra:  "",
	})
	fmt.Println(ret)
}

func TestInvoke(t *testing.T) {
	b := base.NewBlockchainBase(base.ClientConfig{
		ClientType:   "hyperchain",
		ConfigPath:   "./../../../benchmark/evmType/hyperchain",
		ContractPath: "./../../../benchmark/evmType/hyperchain",
		Args:         nil,
		Options:      nil,
	})
	hpc, _ := NewClient(b)
	_ = hpc.DeployContract()
	hpc.Invoke(bcom.Invoke{
		Func: "typeUint8",
		Args: []interface{}{"8", []interface{}{"8", "8", "8"}, []interface{}{"8", "8", "8"}},
	})
}
