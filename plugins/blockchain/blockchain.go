//Copyright 2021 Xiaohui Wang
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package blockchain

import (
	"errors"
	"fmt"
	"plugin"
	"reflect"

	"github.com/hyperbench/hyperbench-common/base"
	fcom "github.com/hyperbench/hyperbench-common/common"
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
