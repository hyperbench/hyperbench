package blockchain

import (
	"errors"
	"fmt"
	"plugin"
	"reflect"

	"github.com/meshplus/hyperbench-common/base"
	fcom "github.com/meshplus/hyperbench-common/common"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var plugins *plugin.Plugin
var log *logging.Logger

// InitPlugin initiate plugin file before init master and worker
func InitPlugin() {
	log = fcom.GetLogger("blockchain")
	filePath := viper.GetString(fcom.ClientPluginPath)

	p, err := plugin.Open(filePath)
	if err != nil {
		log.Errorf("plugin failed: %v", err)
	}
	plugins = p
}

// NewBlockchain create blockchain with different client type.
func NewBlockchain(clientConfig base.ClientConfig) (client Blockchain, err error) {
	clientBase := base.NewBlockchainBase(clientConfig)
	newFunc, err := plugins.Lookup("New")
	if err != nil {
		log.Errorf("plugin failed: %v", err)
	}
	New, _ := newFunc.(func(blockchainBase *base.BlockchainBase) (client interface{}, err error))
	Client, err := New(clientBase)
	if err != nil {
		return nil, err
	}
	client, ok := Client.(Blockchain)
	if !ok {
		return nil, errors.New(fmt.Sprint(reflect.TypeOf(client)) + " is not blockchain.Blockchain")
	}
	return
}

// Blockchain define the service need provided in blockchain.
type Blockchain interface {

	// DeployContract should deploy contract with config file
	DeployContract() error

	// Invoke just invoke the contract
	Invoke(fcom.Invoke, ...fcom.Option) *fcom.Result

	// Transfer a amount of money from a account to the other one
	Transfer(fcom.Transfer, ...fcom.Option) *fcom.Result

	// Confirm check the result of `Invoke` or `Transfer`
	Confirm(*fcom.Result, ...fcom.Option) *fcom.Result

	// Query do some query
	Query(fcom.Query, ...fcom.Option) interface{}

	// Option pass the options to affect the action of client
	Option(fcom.Option) error

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
	Statistic(statistic fcom.Statistic) (*fcom.RemoteStatistic, error)

	// LogStatus records blockheight and time
	LogStatus() (int64, error)
}
