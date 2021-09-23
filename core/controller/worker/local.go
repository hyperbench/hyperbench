package worker

import (
	"context"
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/core/collector"
	"github.com/meshplus/hyperbench/core/engine"
	"github.com/meshplus/hyperbench/core/vmpool"
	"github.com/meshplus/hyperbench/vm"
	"sync"
	"sync/atomic"
	"time"
)

// LocalWorker is the local Worker implement
type LocalWorker struct {
	conf      LocalWorkerConfig
	eg        engine.Engine
	pool      vmpool.Pool
	collector collector.Collector
	idx       common.TxIndex
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	resultCh  chan *common.Result
	done      chan struct{}
	colRet    chan collector.Collector
	colReq    chan struct{}
}

// LocalWorkerConfig define the local worker need config.
type LocalWorkerConfig struct {
	Index    int64
	Cap      int64
	Rate     int64
	Duration time.Duration
}

// NewLocalWorker create LocalWorker.
func NewLocalWorker(config LocalWorkerConfig) (*LocalWorker, error) {

	// init engine
	eg := engine.NewEngine(engine.BaseEngineConfig{
		Rate:     config.Rate,
		Duration: config.Duration,
	})

	// init vm pool
	pool, err := vmpool.NewPoolImpl(config.Index, config.Cap)
	if err != nil {
		return nil, err
	}

	// init index
	idx := common.TxIndex{
		EngineIdx: config.Index,
		TxIdx:     -1,
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &LocalWorker{
		conf:      config,
		eg:        eg,
		pool:      pool,
		collector: collector.NewTDigestSummaryCollector(),
		idx:       idx,
		ctx:       ctx,
		cancel:    cancel,
		resultCh:  make(chan *common.Result, 1024),
		done:      make(chan struct{}),
		colReq:    make(chan struct{}),
		colRet:    make(chan collector.Collector),
	}, nil
}

//func engineCreator(t engine.Type) engine

// SetContext set the context of worker passed from Master
func (l *LocalWorker) SetContext(bs []byte) error {
	var err error
	l.pool.Walk(func(v vm.VM) bool {
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

// Do call the workers to running
func (l *LocalWorker) Do() error {

	go l.runEngine()

	go l.runCollector()

	return nil
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

	l.eg.Run(l.asyncJob)

	// close all engines while Do end to ensure all func has been done
	l.wg.Wait()
	close(l.resultCh)
}

func (l *LocalWorker) asyncJob() {
	v := l.pool.Pop()
	if v == nil {
		// if worker can not get vm from pool, just shortcut
		return
	}
	l.wg.Add(1)

	defer func() {
		l.pool.Push(v)
		l.wg.Done()
	}()

	res, err := v.Run(common.TxContext{
		Context: l.ctx,
		TxIndex: l.atomicAddIndex(),
	})
	if err != nil {
		return
	}
	l.resultCh <- res
	return
}

func (l *LocalWorker) atomicAddIndex() (idx common.TxIndex) {
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
func (l *LocalWorker) CheckoutCollector() (collector.Collector, bool) {
	c, b := <-l.colRet
	return c, b
}

// Done close the worker.
func (l *LocalWorker) Done() chan struct{} {
	return l.done
}
