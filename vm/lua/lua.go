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
 * @brief Define vm interface
 * @file vm.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

package lua

import (
	"errors"
	base2 "github.com/hyperbench/hyperbench-common/base"
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/hyperbench/hyperbench/plugins/blockchain"
	idex "github.com/hyperbench/hyperbench/plugins/index"
	"github.com/hyperbench/hyperbench/plugins/toolkit"
	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/hyperbench/hyperbench/vm/lua/glua"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/viper"
	lua "github.com/yuin/gopher-lua"
)

// VM the implementation of BaseVM for lua.
type VM struct {
	*base.VMBase

	vm       *lua.LState
	instance *lua.LTable
	meta     *lua.LTable
	client   fcom.Blockchain

	index *idex.Index

	runBatch bool
}

// Type return the type of vm
func (v *VM) Type() string {
	return "lua"
}

// NewVM use given base to create VM.
func NewVM(base *base.VMBase) (vm *VM, err error) {

	vm = &VM{
		VMBase: base,
		index: &idex.Index{
			Worker: base.Ctx.WorkerIdx,
			VM:     base.Ctx.VMIdx,
		},
	}

	vm.vm = lua.NewState()
	defer func() {
		if err != nil {
			vm = nil
		}
	}()

	// inject test case metatable
	vm.injectTestcaseBase()

	// load script
	err = vm.vm.DoFile(base.Path)
	if err != nil {
		return nil, errors2.Wrap(err, "load script fail")
	}
	// get test
	var ok bool
	vm.instance, ok = vm.vm.Get(-1).(*lua.LTable)
	if !ok {
		return nil, errors.New("script's return value is not table")
	}
	vm.vm.Pop(1)

	vm.runBatch = vm.instance.RawGetString(runBatch) != lua.LNil

	vm.Logger.Debugf("runBatch:%v", vm.runBatch)
	return vm, nil
}

// hooks
const (
	beforeDeploy = "BeforeDeploy"
	// nolint
	deployContract = "DeployContract"
	beforeGet      = "BeforeGet"
	beforeSet      = "BeforeSet"
	// nolint
	setContext = "SetContext"
	beforeRun  = "BeforeRun"
	run        = "Run"
	runBatch   = "RunBatch"
	afterRun   = "AfterRun"
	closeBlock = "Close"
)

// builtin
const (
	lNew   = "new"
	lIndex = "__index"
)

// plugins
const (
	testcase = "testcase"
	client   = "blockchain"
	tool     = "toolkit"
	index    = "index"
)

func (v *VM) injectTestcaseBase() {
	mt := v.vm.NewTypeMetatable(testcase)
	v.vm.SetGlobal(testcase, mt)

	var empty lua.LGFunction = func(state *lua.LState) int {
		return 0
	}
	var result lua.LGFunction = func(state *lua.LState) int {
		state.Push(glua.Go2Lua(state, &fcom.Result{}))
		return 1
	}

	v.vm.SetField(mt, lNew, v.vm.NewFunction(func(state *lua.LState) int {
		table := v.vm.NewTable()
		v.vm.SetMetatable(table, v.vm.GetMetatable(lua.LString(testcase)))
		err := v.setPlugins(table)
		if err != nil {
			v.Logger.Errorf("setPlugins fail:%v", err)
		}
		v.vm.Push(table)
		return 1
	}))

	v.vm.SetField(mt, lIndex, v.vm.SetFuncs(v.vm.NewTable(), map[string]lua.LGFunction{
		beforeDeploy: empty,
		beforeGet:    empty,
		beforeSet:    empty,
		beforeRun:    empty,
		run:          result,
		afterRun:     empty,
	}))
	v.meta = mt
}

// BeforeDeploy will call before deploy contract.
func (v *VM) BeforeDeploy() error {
	return v.callInVm(beforeDeploy)
}

// DeployContract deploy contract.
func (v *VM) DeployContract() error {
	return v.client.DeployContract()
}

// BeforeGet will call before get context.
func (v *VM) BeforeGet() error {
	return v.callInVm(beforeGet)
}

// GetContext generate context for execute tx in vm.
func (v *VM) GetContext() ([]byte, error) {
	s, err := v.client.GetContext()
	return []byte(s), err
}

// Statistic statistic remote execute info.
func (v *VM) Statistic(from, to *fcom.ChainInfo) (*fcom.RemoteStatistic, error) {

	return v.client.Statistic(fcom.Statistic{
		From: from,
		To:   to,
	})
}

// LogStatus records blockheight and time
func (v *VM) LogStatus() (chainInfo *fcom.ChainInfo, err error) {
	return v.client.LogStatus()
}

// Verify check the relative time of transaction
func (v *VM) Verify(res *fcom.Result, ops ...fcom.Option) *fcom.Result {
	return v.client.Verify(res)
}

// VerifyBatch check the relative time of txs
func (v *VM) VerifyBatch(res ...*fcom.Result) []*fcom.Result {
	return v.client.VerifyBatch(res...)
}

// BeforeSet will call before set context.
func (v *VM) BeforeSet() error {
	return v.callInVm(beforeSet)
}

func (v *VM) callInVm(methodName string) error {
	fn := v.instance.RawGetString(methodName)
	if fn != lua.LNil {
		return v.vm.CallByParam(lua.P{
			Fn: fn,
		}, v.instance)
	}
	return nil
}

// SetContext set context for execute tx in vm, the ctx is generated by GetContext.
func (v *VM) SetContext(ctx []byte) error {
	return v.client.SetContext(string(ctx))
}

// BeforeRun will call once before run.
func (v *VM) BeforeRun() error {
	return v.callInVm(beforeRun)
}

// Run create and send tx to client.
func (v *VM) Run(ctx fcom.TxContext) (*fcom.Result, error) {
	ud, err := v.run(ctx, run)
	if err != nil {
		return nil, err
	}
	res, ok := ud.Value.(*fcom.Result)
	if !ok {
		v.Logger.Debugf("returned user data is not result")
		return nil, errors.New("returned user data is not result")
	}
	return res, nil
}

func (v *VM) IsRunBatch() bool {
	return v.runBatch
}

func (v *VM) RunBatch(ctx fcom.TxContext) ([]*fcom.Result, error) {
	ud, err := v.run(ctx, runBatch)
	if err != nil {
		return nil, err
	}
	res, ok := ud.Value.([]*fcom.Result)
	if !ok {
		v.Logger.Debugf("returned user data is not result")
		return nil, errors.New("returned user data is not []*result")
	}
	return res, nil
}

// AfterRun will call once after run.
func (v *VM) AfterRun() error {
	return v.callInVm(afterRun)
}

// Close close vm.
func (v *VM) Close() {
	err := v.callInVm(closeBlock)
	if err != nil {
		v.Logger.Errorf("close err:%v", err)
	}
	v.vm.Close()
}

func (v *VM) setPlugins(table *lua.LTable) (err error) {

	clientType, clientConfigPath := viper.GetString(fcom.ClientTypePath), viper.GetString(fcom.ClientConfigPath)
	options := viper.GetStringMap(fcom.ClientOptionPath)
	contractPath := viper.GetString(fcom.ClientContractPath)
	args, _ := viper.Get(fcom.ClientContractArgsPath).([]interface{})
	v.client, err = blockchain.NewBlockchain(base2.ClientConfig{
		ClientType:   clientType,
		ConfigPath:   clientConfigPath,
		ContractPath: contractPath,
		Args:         args,
		Options:      options,
		VmID:         int(v.index.VM),
		WorkerID:     int(v.index.Worker),
	})

	if err != nil {
		return err
	}

	lClient := glua.Go2Lua(v.vm, v.client)
	lToolKit := glua.Go2Lua(v.vm, toolkit.NewToolKit())
	lIndex := glua.Go2Lua(v.vm, v.index)
	v.vm.SetField(table, client, lClient)
	v.vm.SetField(table, tool, lToolKit)
	v.vm.SetField(table, index, lIndex)

	return nil
}

func (v *VM) run(ctx fcom.TxContext, funcName string) (*lua.LUserData, error) {
	v.index.Engine = ctx.EngineIdx
	// todo：批量发送时，txIndex有问题
	v.index.Tx = ctx.TxIdx

	err := v.vm.CallByParam(lua.P{
		Fn:      v.instance.RawGetString(funcName),
		NRet:    1,
		Protect: false,
	}, v.instance)

	if err != nil {
		v.Logger.Error(err)
		return nil, err
	}
	val := v.vm.Get(-1)
	v.vm.Pop(1)
	ud, ok := val.(*lua.LUserData)
	if !ok {
		return nil, errors.New("returned val is not user data")
	}
	return ud, nil
}
