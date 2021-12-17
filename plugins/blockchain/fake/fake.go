package fake

import (
	"fmt"
	"github.com/meshplus/hyperbench/common"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
)

type FakeChain struct {
	Name string
}

func New() (client *FakeChain, err error) {
	return &FakeChain{"data"}, nil
}

func (chain *FakeChain) DeployContract() error {
	return nil
}

func (chain *FakeChain) Invoke(invoke bcom.Invoke, ops ...bcom.Option) *common.Result {
	fmt.Printf("invoke:%v\n",invoke)
	fmt.Printf("ops:%v\n",ops)

	return &common.Result{Status: common.Success}
}

// Transfer a amount of money from a account to the other one
func (chain *FakeChain) Transfer(t bcom.Transfer, ops ...bcom.Option) *common.Result {
	return &common.Result{}
}

// Confirm check the result of `Invoke` or `Transfer`
func (chain *FakeChain) Confirm(rt *common.Result, ops ...bcom.Option) *common.Result {
	return &common.Result{}
}

// Query do some query
func (chain *FakeChain) Query(bq bcom.Query, ops ...bcom.Option) interface{} {
	return nil
}

// Option pass the options to affect the action of client
func (chain *FakeChain) Option(bcom.Option) error {
	return nil
}

// GetContext Generate TxContext based on New/Init/DeployContract
// GetContext will only be run in master
// return the information how to invoke the contract, maybe include
// contract address, abi or so.
// the return value will be send to worker to tell them how to invoke the contract
func (chain *FakeChain) GetContext() (string, error) {
	return "", nil
}

// SetContext set test context into go client
// SetContext will be run once per worker
func (chain *FakeChain) SetContext(ctx string) error {
	return nil
}

// ResetContext reset test group context in go client
func (chain *FakeChain) ResetContext() error {
	return nil
}

// Statistic query the statistics information in the time interval defined by
// nanosecond-level timestamps `from` and `to`
func (chain *FakeChain) Statistic(statistic bcom.Statistic) (*common.RemoteStatistic, error) {
	return &common.RemoteStatistic{}, nil
}
