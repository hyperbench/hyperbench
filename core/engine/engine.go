package engine

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
 * @brief Engine is used to control the rate for send tx
 * @file engine.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"context"
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/hyperbench/hyperbench/core/utils"
	"github.com/op/go-logging"
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
	// Instant the number of a batch
	Instant int64 `mapstructure:"instant"`
	// Duration engine run duration.
	Duration time.Duration `mapstructure:"duration"`
	// Wg Semaphore of localWorker
	Wg *sync.WaitGroup
}

type baseEngine struct {
	BaseEngineConfig

	log      *logging.Logger
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
		log:              fcom.GetLogger("engine"),
	}).adjust()
}

func (b *baseEngine) adjust() *baseEngine {
	b.interval = time.Duration(float64(b.Instant) / float64(b.Rate) * float64(time.Second))
	return b
}

// Run start the engine.
func (b *baseEngine) Run(callback Callback) {
	b.schedule(callback)
}

func (b *baseEngine) schedule(callback Callback) {
	totalBatch, batchCount := utils.DivideAndCeil(int(b.Duration), int(b.interval)), 0
	tick := time.NewTicker(b.interval)
	defer func() {
		tick.Stop()
	}()
	for ; batchCount < totalBatch; batchCount++ {
		for i := int64(0); i < b.Instant; i++ {
			callback()
		}
		<-tick.C
	}
}

// Close close the engine.
func (b *baseEngine) Close() {
	b.cancelFunc()
}
