package base

import (
	"github.com/meshplus/hyperbench/common"
	bcom "github.com/meshplus/hyperbench/plugins/blockchain/common"
	"github.com/op/go-logging"
)

// BlockchainBase the base implementation of blockChain.
type BlockchainBase struct {
	ClientConfig
	Logger *logging.Logger
}

// ClientConfig define the filed for client config.
type ClientConfig struct {
	// client
	ClientType string `mapstructure:"type"`
	ConfigPath string `mapstructure:"config"`

	// contract
	ContractPath string        `mapstructure:"contract"`
	Args         []interface{} `mapstructure:"args"`

	// options
	Options map[string]interface{} `mapstructure:"options"`
}

// Confirm confirm invoke result for tx.
func (b *BlockchainBase) Confirm(result *common.Result, ops ...bcom.Option) *common.Result {
	return result
}

// DeployContract send tx for deploy contract.
func (b *BlockchainBase) DeployContract() error {
	return nil
}

// Invoke send tx for invoke contract.
func (b *BlockchainBase) Invoke(bcom.Invoke, ...bcom.Option) *common.Result {
	return &common.Result{}
}

// Transfer send tx for transfer.
func (b *BlockchainBase) Transfer(bcom.Transfer, ...bcom.Option) *common.Result {
	return &common.Result{}
}

// Query query info.
func (b *BlockchainBase) Query(bcom.Query, ...bcom.Option) interface{} {
	return nil
}

// Option receive some options.
func (b *BlockchainBase) Option(bcom.Option) error {
	return nil
}

// GetContext get context for execute tx in vm.
func (b *BlockchainBase) GetContext() (string, error) {
	return "", nil
}

// SetContext set context for execute tx in vm.
func (b *BlockchainBase) SetContext(ctx string) error {
	return nil
}

// ResetContext reset context.
func (b *BlockchainBase) ResetContext() error {
	return nil
}

// Statistic statistic remote execute result.
func (b *BlockchainBase) Statistic(statistic bcom.Statistic) (*common.RemoteStatistic, error) {
	return &common.RemoteStatistic{}, nil
}

// NewBlockchainBase new blockchain base.
func NewBlockchainBase(clientConfig ClientConfig) *BlockchainBase {
	return &BlockchainBase{
		ClientConfig: clientConfig,
		Logger:       common.GetLogger("client"),
	}
}
