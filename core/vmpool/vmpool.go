package vmpool

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
 * @brief Virtual machine pool, managing virtual machines
 * @file vmpool.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/hyperbench/hyperbench/core/utils"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"path"
	"strings"
	"sync"

	"github.com/hyperbench/hyperbench/vm"
	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/spf13/viper"
)

type WrapVM struct {
	vm vm.VM
	ch chan struct{}
}

type PoolImp struct {
	cap          int
	closeCh      chan struct{}
	wg           sync.WaitGroup
	job          func(v vm.VM)
	currentIndex int
	vms          []*WrapVM
	log          *logging.Logger
}

func (i *PoolImp) startListenVM(nvm *WrapVM) {
	defer i.wg.Done()
	for {
		select {
		case <-i.closeCh:
			return
		case _, ok := <-nvm.ch:
			if !ok {
				return
			}
			i.job(nvm.vm)
		}
	}
}

func (i *PoolImp) Push() error {
	// find a not pull vm, push struct{} into vm.ch
	for j := 0; j < len(i.vms); j++ {
		i.currentIndex = (i.currentIndex + j) % len(i.vms)
		wrapVM := i.vms[i.currentIndex]
		if len(wrapVM.ch) < i.cap {
			wrapVM.ch <- struct{}{}
			i.addCurrentIndex()
			return nil
		}
	}
	return errors.New("vm is too busy")
}

func (i *PoolImp) addCurrentIndex() {
	i.currentIndex++
	i.currentIndex = i.currentIndex % len(i.vms)
}

func (i *PoolImp) Walk(wf func(v vm.VM) bool) {
	for _, wvm := range i.vms {
		exist := wf(wvm.vm)
		if exist {
			return
		}
	}
}

func (i *PoolImp) AsyncWalk(wf func(v vm.VM) bool) {
	wg := sync.WaitGroup{}
	for _, wvm := range i.vms {
		wg.Add(1)
		go func(wvm *WrapVM) {
			defer wg.Done()
			exist := wf(wvm.vm)
			if exist {
				return
			}
		}(wvm)
	}
	wg.Wait()
}

func (i *PoolImp) Close() {
	for _, wvm := range i.vms {
		wvm.vm.Close()
		close(wvm.ch)
	}
	i.wg.Wait()
}

func NewPoolImp(workerID int64, tps, cap int64, job func(v vm.VM)) (*PoolImp, error) {
	chCap := utils.DivideAndCeil(int(tps), int(cap))
	p := &PoolImp{
		cap:     chCap,
		job:     job,
		closeCh: make(chan struct{}),
	}
	scriptPath := viper.GetString(fcom.ClientScriptPath)
	t := strings.TrimPrefix(path.Ext(scriptPath), ".")
	configBase := base.ConfigBase{
		Path: scriptPath,
		Ctx: fcom.VMContext{
			WorkerIdx: workerID,
			VMIdx:     0,
		},
	}
	configBase.Ctx.WorkerIdx = workerID
	var i int64
	p.log = fcom.GetLogger("pool")
	p.log.Notice(workerID, cap, scriptPath, t)
	for i = 0; i < cap; i++ {
		nvm, err := vm.NewVM(t, configBase)
		if err != nil {
			return nil, err
		}
		configBase.Ctx.VMIdx++
		// generate each vm with given index
		p.wg.Add(1)
		wvm := &WrapVM{
			vm: nvm,
			ch: make(chan struct{}, chCap),
		}
		p.vms = append(p.vms, wvm)
		go p.startListenVM(wvm)
	}

	return p, nil
}
