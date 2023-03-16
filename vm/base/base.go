package base

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
 * @brief VMBase the base vm for support base config
 * @file vm.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/op/go-logging"
)

// VMBase the base vm for support base config.
type VMBase struct {
	ConfigBase
	Logger *logging.Logger
}

// Type return the vm type.
func (v *VMBase) Type() string {
	return "base"
}

// Close close vm.
func (v *VMBase) Close() {
}

// BeforeDeploy will call before deploy contract.
func (v *VMBase) BeforeDeploy() error {
	return nil
}

// DeployContract deploy contract.
func (v *VMBase) DeployContract() error {
	return nil
}

// BeforeGet will call before get context.
func (v *VMBase) BeforeGet() error {
	return nil
}

// GetContext generate context for execute tx in vm.
func (v *VMBase) GetContext() ([]byte, error) {
	return []byte(""), nil
}

// Statistic statistic remote execute info.
func (v *VMBase) Statistic(from, to *fcom.ChainInfo) (*fcom.RemoteStatistic, error) {
	return &fcom.RemoteStatistic{}, nil
}

// LogStatus records blockheight and time
func (v *VMBase) LogStatus() (*fcom.ChainInfo, error) {
	return nil, nil
}

// Verify check the relative time of transaction
func (v *VMBase) Verify(*fcom.Result, ...fcom.Option) *fcom.Result {
	return nil
}

func (v *VMBase) VerifyBatch(res ...*fcom.Result) []*fcom.Result {
	return nil
}

// BeforeSet will call before set context.
func (v *VMBase) BeforeSet() error {
	return nil
}

// SetContext set context for execute tx in vm, the ctx is generated by GetContext.
func (v *VMBase) SetContext(ctx []byte) error {
	return nil
}

// BeforeRun will call once before run.
func (v *VMBase) BeforeRun() error {
	return nil
}

// Run create and send tx to client.
func (v *VMBase) Run(ctx fcom.TxContext) (*fcom.Result, error) {
	return &fcom.Result{}, nil
}

func (v *VMBase) IsRunBatch() bool { return false }

// RunBatch create and send batch tx to client.
func (v *VMBase) RunBatch(ctx fcom.TxContext) ([]*fcom.Result, error) {
	return []*fcom.Result{{}}, nil
}

// AfterRun will call once after run.
func (v *VMBase) AfterRun() error {
	return nil
}

// ConfigBase define base config in vm.
type ConfigBase struct {
	// Path is the path of script file
	Path string
	// Ctx is the context of vm
	Ctx fcom.VMContext
}

// NewVMBase use given config create VMBase.
func NewVMBase(config ConfigBase) *VMBase {
	return &VMBase{
		ConfigBase: config,
		Logger:     fcom.GetLogger("vm"),
	}
}
