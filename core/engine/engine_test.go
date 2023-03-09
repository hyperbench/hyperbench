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
 * @brief the unit test for engine.go
 * @file engine_test.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */
package engine

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	engine1 := NewEngine(BaseEngineConfig{
		Rate:     1,
		Duration: time.Millisecond * 500,
	})
	assert.NotNil(t, engine1)

	engine2 := NewEngine(BaseEngineConfig{
		Rate:     101,
		Duration: time.Second * 1,
	})
	assert.NotNil(t, engine2)

	engine1.Run(func() {})

	engine1.Close()

	engine2.Run(func() {})
}

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go test1(ctx)
	go test2(ctx)
	cancel()
	fmt.Println("main1")
	cancel()
	fmt.Println("main2")
	time.Sleep(time.Second)
	fmt.Println("main3")
}

func test1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("test1")
			return
		}
	}
}

func test2(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("test2")
			return
		}
	}
}
