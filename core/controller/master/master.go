package master

/**
 *  Copyright (C) 2021 HyperBench.
 *  SPDX-License-Identifier: Apache-2.0
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 * @brief Provide Master and LocalMaster, the master node.
 * @file master.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"path/filepath"
	"strings"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/hyperbench/hyperbench/plugins/blockchain"
	"github.com/hyperbench/hyperbench/vm"
	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Master is the interface of mater node
type Master interface {
	// Prepare is used to prepare
	Prepare() error

	// GetContext generate the context, which will be passed to Worker
	GetContext() ([]byte, error)

	// Statistic query the remote statistic data from chain
	Statistic(from, to *fcom.ChainInfo) (*fcom.RemoteStatistic, error)

	// LogStatus records blockheight and time
	LogStatus() (*fcom.ChainInfo, error)
}

// LocalMaster is the implement of master in local
type LocalMaster struct {
	masterVM vm.VM
	params   []string
}

// Prepare is used to prepare
func (m *LocalMaster) Prepare() (err error) {
	// call user hook
	err = m.masterVM.BeforeDeploy()
	if err != nil {
		return errors.Wrap(err, "can not call user hook `BeforeDeploy`")
	}

	// prepare contract
	err = m.masterVM.DeployContract()
	if err != nil {
		return errors.Wrap(err, "can not deploy contract")
	}
	return nil
}

// GetContext generate the context, which will be passed to Worker
func (m *LocalMaster) GetContext() ([]byte, error) {
	err := m.masterVM.BeforeGet()
	if err != nil {
		return nil, err
	}
	return m.masterVM.GetContext()
}

// Statistic query the remote statistic data from chain
func (m *LocalMaster) Statistic(from, to *fcom.ChainInfo) (*fcom.RemoteStatistic, error) {
	return m.masterVM.Statistic(from, to)
}

// LogStatus records blockheight and time
func (m *LocalMaster) LogStatus() (chainInfo *fcom.ChainInfo, err error) {
	return m.masterVM.LogStatus()
}

// NewLocalMaster create LocalMaster.
func NewLocalMaster() (*LocalMaster, error) {
	blockchain.InitPlugin()

	params := viper.GetStringSlice(fcom.ClientContractArgsPath)
	scriptPath := viper.GetString(fcom.ClientScriptPath)
	vmType := strings.TrimPrefix(filepath.Ext(scriptPath), ".")
	masterVM, err := vm.NewVM(vmType, base.ConfigBase{
		Path: scriptPath,
		Ctx: fcom.VMContext{
			WorkerIdx: 0,
			VMIdx:     0,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "can not create local master")
	}

	return &LocalMaster{
		masterVM: masterVM,
		params:   params,
	}, nil
}
