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

package worker

import (
	"context"
	"github.com/op/go-logging"
	"math/rand"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/hyperbench/hyperbench/core/collector"
	"github.com/hyperbench/hyperbench/core/engine"
	"github.com/hyperbench/hyperbench/core/vmpool"
	"github.com/hyperbench/hyperbench/plugins/blockchain"
	"github.com/hyperbench/hyperbench/vm"
	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/spf13/viper"
)

// LocalWorker is the local Worker implement
type LocalWorker struct {
	conf           LocalWorkerConfig
	eg             engine.Engine
	pool           *vmpool.PoolImp
	collector      collector.Collector
	idx            fcom.TxIndex
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	resultCh       chan *fcom.Result
	done           chan struct{}
	colRet         chan collector.Collector
	colReq         chan struct{}
	txInterval     int64          // interval for sample ,every txInterval chose one tx randomly
	verifyVM       vm.WorkerVM    // vm used for verification to make sure query from same node
	verifyIndexMap map[int64]bool // store the sample result
	log            *logging.Logger
}

// LocalWorkerConfig define the local worker need config.
type LocalWorkerConfig struct {
	Index    int64         // Index the index of localWorker
	Cap      int64         // Cap the number of vm
	Rate     int64         // Rate the number of tx will be sent per second
	Instant  int64         // Instant the number of a batch
	Wait     time.Duration // Wait maximum time to wait before get vm
	Duration time.Duration // Duration time of pressure test
}

// NewLocalWorker create LocalWorker.
func NewLocalWorker(config LocalWorkerConfig) (*LocalWorker, error) {
	blockchain.InitPlugin()

	localWorker := LocalWorker{
		collector:      collector.NewTDigestSummaryCollector(),
		resultCh:       make(chan *fcom.Result, 1024),
		done:           make(chan struct{}),
		colReq:         make(chan struct{}),
		colRet:         make(chan collector.Collector),
		verifyIndexMap: make(map[int64]bool),
		log:            fcom.GetLogger("localWorker"),
	}
	// init engine
	eg := engine.NewEngine(engine.BaseEngineConfig{
		Rate:     config.Rate,
		Duration: config.Duration,
		Instant:  config.Instant,
		Wg:       &localWorker.wg,
	})

	// init vm pool
	pool, err := vmpool.NewPoolImp(config.Index, config.Rate, config.Cap, localWorker.run)
	if err != nil {
		return nil, err
	}

	// init verification
	var txInterval int64
	verifyEnable := viper.GetBool(fcom.VerifyEnablePath)
	if verifyEnable {
		// calculate txInterval
		samplePercentage := viper.GetFloat64(fcom.VerifyPercentagePath)

		if samplePercentage <= 0 {
			samplePercentage = 0.001
		} else if samplePercentage > 1 {
			samplePercentage = 1
		}
		txInterval = int64(1 / samplePercentage)
		// init vm for verification
		scriptPath := viper.GetString(fcom.ClientScriptPath)
		vmType := strings.TrimPrefix(path.Ext(scriptPath), ".")
		configBase := base.ConfigBase{
			Path: scriptPath,
			Ctx: fcom.VMContext{
				WorkerIdx: config.Index,
				VMIdx:     config.Cap,
			},
		}
		verifyVM, err := vm.NewVM(vmType, configBase)
		if err != nil {
			return nil, err
		}
		localWorker.verifyVM = verifyVM
		// sample transactions
		localWorker.sample()
	}
	// init index
	idx := fcom.TxIndex{
		EngineIdx: config.Index,
		TxIdx:     -1,
		MissIdx:   0,
	}
	ctx, cancel := context.WithCancel(context.Background())
	localWorker.conf = config
	localWorker.eg = eg
	localWorker.pool = pool
	localWorker.idx = idx
	localWorker.ctx = ctx
	localWorker.cancel = cancel
	localWorker.txInterval = txInterval

	return &localWorker, nil
}

// SetContext set the context of worker passed from Master
func (l *LocalWorker) SetContext(bs []byte) (err error) {
	l.pool.AsyncWalk(func(v vm.VM) bool {
		if err = v.BeforeSet(); err != nil {
			return true
		}
		if err = v.SetContext(bs); err != nil {
			return true
		}
		return false
	})
	return err
}

// BeforeRun call user hook
func (l *LocalWorker) BeforeRun() (err error) {
	l.pool.AsyncWalk(func(v vm.VM) bool {
		if err = v.BeforeRun(); err != nil {
			return true
		}
		return false
	})
	return err
}

// Do call the workers to running
func (l *LocalWorker) Do() error {
	go l.runEngine()

	go l.runCollector()
	return nil
}

// AfterRun call user hook
func (l *LocalWorker) AfterRun() (err error) {
	l.pool.AsyncWalk(func(v vm.VM) bool {
		if err = v.AfterRun(); err != nil {
			return true
		}
		return false
	})
	return err
}

// Statistics get the number of sent and missed transactions
func (l *LocalWorker) Statistics() (int64, int64) {
	return l.idx.TxIdx + 1, l.idx.MissIdx
}

func (l *LocalWorker) runCollector() {

	defer func() {
		close(l.done)
		close(l.colRet)
	}()
	l.collector.Reset()
	for {
		select {
		case <-l.ctx.Done():
			return
		case result, valid := <-l.resultCh:
			if !valid {
				// engine stop
				l.colRet <- l.collector
				return
			}
			l.collector.Add(result)
		case l.colRet <- l.collector:
			l.collector = collector.NewTDigestSummaryCollector()
			l.collector.Reset()
		}
	}
}

func (l *LocalWorker) runEngine() {
	l.eg.Run(l.job)
	// close all engines while Do end to ensure all func has been done
	l.pool.Close()
	close(l.resultCh)
	l.log.Noticef("close resultCh")
}

func (l *LocalWorker) job() {
	err := l.pool.Push()
	if err != nil {
		atomic.AddInt64(&l.idx.MissIdx, 1)
		// if worker can not get vm from pool, just shortcut
		return
	}
}

func (l *LocalWorker) run(v vm.VM) {
	txContext := fcom.TxContext{
		Context: l.ctx,
		TxIndex: l.atomicAddIndex(),
	}
	res, err := v.Run(txContext)
	if err != nil {
		return
	}
	// if enable verification and this tx is chosen, verify the tx
	if l.txInterval > 0 && l.verifyIndexMap[txContext.TxIdx] {
		l.verifyVM.Verify(res)
	}
	l.resultCh <- res
}

func (l *LocalWorker) atomicAddIndex() (idx fcom.TxIndex) {
	idx.EngineIdx = atomic.LoadInt64(&l.idx.EngineIdx)
	idx.TxIdx = atomic.AddInt64(&l.idx.TxIdx, 1)
	return
}

// Teardown close the worker manually.
func (l *LocalWorker) Teardown() {
	l.eg.Close()
	l.cancel()
}

// CheckoutCollector checkout collector.
func (l *LocalWorker) CheckoutCollector() (collector.Collector, bool, error) {
	c, b := <-l.colRet
	return c, b, nil
}

// Done close the worker.
func (l *LocalWorker) Done() chan struct{} {
	return l.done
}

// sample chose indexes of transactions to be verified
func (l *LocalWorker) sample() {
	txNum, current := l.conf.Rate*int64(l.conf.Duration/time.Second), int64(0)

	for current < txNum {
		tmp := current + l.txInterval
		if tmp > txNum {
			tmp = txNum
		}
		index := current + rand.Int63n(tmp-current)
		l.verifyIndexMap[index] = true
		current = tmp
	}
}
