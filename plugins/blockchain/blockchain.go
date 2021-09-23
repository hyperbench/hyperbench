package blockchain

import (
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/base"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"github.com/meshplus/hyperbench/plugins/blockchain/fabric"
	"github.com/meshplus/hyperbench/plugins/blockchain/hyperchain"
)

const (
	clientHpc    = "hyperchain"
	clientFlato  = "flato"
	clientFabric = "fabric"
)

// NewBlockchain create blockchain with different client type.
func NewBlockchain(clientConfig base.ClientConfig) (client Blockchain, err error) {
	clientBase := base.NewBlockchainBase(clientConfig)
	switch clientConfig.ClientType {
	case clientHpc, clientFlato:
		client, err = hyperchain.NewClient(clientBase)
	case clientFabric:
		client, err = fabric.New(clientBase)
	default:
		client = clientBase
	}
	return
}

// Blockchain define the service need provided in blockchain.
type Blockchain interface {

	// DeployContract should deploy contract with config file
	DeployContract() error

	// Invoke just invoke the contract
	Invoke(bcom.Invoke, ...bcom.Option) *common.Result

	// Transfer a amount of money from a account to the other one
	Transfer(bcom.Transfer, ...bcom.Option) *common.Result

	// Confirm check the result of `Invoke` or `Transfer`
	Confirm(*common.Result, ...bcom.Option) *common.Result

	// Query do some query
	Query(bcom.Query, ...bcom.Option) interface{}

	// Option pass the options to affect the action of client
	Option(bcom.Option) error

	// GetContext Generate TxContext based on New/Init/DeployContract
	// GetContext will only be run in master
	// return the information how to invoke the contract, maybe include
	// contract address, abi or so.
	// the return value will be send to worker to tell them how to invoke the contract
	GetContext() (string, error)

	// SetContext set test context into go client
	// SetContext will be run once per worker
	SetContext(ctx string) error

	// ResetContext reset test group context in go client
	ResetContext() error

	// Statistic query the statistics information in the time interval defined by
	// nanosecond-level timestamps `from` and `to`
	Statistic(statistic bcom.Statistic) (*common.RemoteStatistic, error)
}
