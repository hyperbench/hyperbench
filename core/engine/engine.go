package engine

import (
	"context"
	"sync"
	"time"
)

// Callback will be call in engine run.
type Callback func()

// Engine is used to control the rate for send tx.
type Engine interface {
	// Run start the engine.
	Run(callback Callback)
	// Close close the engine.
	Close()
}

// NewEngine use given baseEngineConf create Engine.
func NewEngine(baseEngineConf BaseEngineConfig) (e Engine) {
	baseEngine := newBaseEngine(baseEngineConf)
	switch baseEngine.Type {
	default:
		e = newConstantEngine(baseEngine)
	}
	return
}

// BaseEngineConfig base engine config.
type BaseEngineConfig struct {
	// Type engine type.
	Type string `mapstructure:"type"`
	// Rate engine call Callback rate.
	Rate int64 `mapstructure:"rate"`
	// Duration engine run duration.
	Duration time.Duration `mapstructure:"duration"`
	// Wg Semaphore of localWorker
	Wg *sync.WaitGroup
}

type baseEngine struct {
	BaseEngineConfig

	batch    int64
	interval time.Duration
	//wg         sync.WaitGroup
	timeoutCtx context.Context
	cancelFunc context.CancelFunc
}

func newBaseEngine(config BaseEngineConfig) *baseEngine {
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), config.Duration)

	return (&baseEngine{
		BaseEngineConfig: config,
		timeoutCtx:       timeoutCtx,
		cancelFunc:       cancelFunc,
	}).adjust()
}

func (b *baseEngine) adjust() *baseEngine {

	// calculate batch and interval
	if b.Rate <= 100 {
		b.batch = 1
		b.interval = time.Second / time.Duration(b.Rate)
	} else {
		b.batch = b.Rate / 10
		b.interval = time.Second / 10
	}

	return b
}

// Run start the engine.
func (b *baseEngine) Run(callback Callback) {
	b.schedule(callback)
}

func (b *baseEngine) schedule(callback Callback) {
	totalBatch, batchCount := int(b.Duration/b.interval), 0
	tick := time.NewTicker(b.interval)
	defer func() {
		tick.Stop()
	}()
	for ; batchCount < totalBatch; batchCount++ {
		<-tick.C
		for i := int64(0); i < b.batch; i++ {
			b.Wg.Add(1)
			go callback()
		}
	}
}

// Close close the engine.
func (b *baseEngine) Close() {
	b.cancelFunc()
}
