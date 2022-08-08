package glua

import (
	"github.com/hyperbench/hyperbench-common/base"
	fcom "github.com/hyperbench/hyperbench-common/common"
)

const (
	Invoke         = "Invoke"
	Transfer       = "Transfer"
	Option         = "Option"
	Result         = "Result"
	Query          = "Query"
	SetContext     = "SetContext"
	Statistic      = "Statistic"
	DeployContract = "DeployContract"
)

type FakeChain struct {
	Name     string
	base     *base.BlockchainBase
	tempData map[string]interface{}
}

func NewMock() (client *FakeChain, err error) {
	tempMap := make(map[string]interface{})
	return &FakeChain{"fake",
		base.NewBlockchainBase(base.ClientConfig{}),
		tempMap,
	}, nil
}

func (chain *FakeChain) DeployContract() error {
	chain.tempData[DeployContract] = DeployContract
	return nil
}

func (chain *FakeChain) Invoke(invoke fcom.Invoke, ops ...fcom.Option) *fcom.Result {
	chain.tempData[Invoke] = invoke
	chain.tempData[Option] = ops
	return &fcom.Result{
		Label:  "label",
		UID:    "UUID",
		Status: fcom.Success,
		Ret:    []interface{}{"demo", "demo"},
	}
}

func (chain *FakeChain) Transfer(transfer fcom.Transfer, ops ...fcom.Option) *fcom.Result {
	chain.tempData[Transfer] = transfer
	chain.tempData[Option] = ops

	return &fcom.Result{
		Label:  "label",
		UID:    "UUID",
		Status: fcom.Success,
		Ret:    []interface{}{"demo", "demo"},
	}
}

func (chain *FakeChain) Confirm(rt *fcom.Result, ops ...fcom.Option) *fcom.Result {
	chain.tempData[Result] = rt
	chain.tempData[Option] = ops
	return &fcom.Result{
		Label:  "Confirm",
		UID:    "UUID",
		Status: fcom.Confirm,
		Ret:    []interface{}{"Confirm", "Confirm"},
	}
}

func (chain *FakeChain) Verify(rt *fcom.Result, ops ...fcom.Option) *fcom.Result {
	chain.tempData[Result] = rt
	chain.tempData[Option] = ops
	return &fcom.Result{
		Label:  "Confirm",
		UID:    "UUID",
		Status: fcom.Confirm,
		Ret:    []interface{}{"Confirm", "Confirm"},
	}
}

func (chain *FakeChain) Query(bq fcom.Query, ops ...fcom.Option) interface{} {
	chain.tempData[Query] = bq
	chain.tempData[Option] = ops
	return "nil"
}

func (chain *FakeChain) Option(op fcom.Option) error {
	chain.tempData[Option] = op
	return nil
}

func (chain *FakeChain) GetContext() (string, error) {
	return "GetContext", nil
}

func (chain *FakeChain) SetContext(ctx string) error {
	chain.tempData[SetContext] = ctx
	return nil
}

func (chain *FakeChain) ResetContext() error {
	return nil
}

func (chain *FakeChain) Statistic(statistic fcom.Statistic) (*fcom.RemoteStatistic, error) {
	chain.tempData[Statistic] = statistic
	return &fcom.RemoteStatistic{}, nil
}

func (chain *FakeChain) LogStatus() (*fcom.ChainInfo, error) {
	return nil, nil
}
