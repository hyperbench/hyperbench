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
func NewBlockchain(clientConfig base.ClientConfig) (client fcom.Blockchain, err error) {
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
	client, ok := Client.(fcom.Blockchain)
	if !ok {
		return nil, errors.New(fmt.Sprint(reflect.TypeOf(client)) + " is not blockchain.Blockchain")
	}
	return
}