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
 * @brief use cobra provide cmd function
 * @file worker_test.go
 * @author: shinyxhh
 * @date 2021-11-30
 */
package worker

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestLocalWorker(t *testing.T) {
	localWorker, err := NewLocalWorker(LocalWorkerConfig{0, 5, 20, 5, time.Millisecond * 5, time.Second * 5})
	assert.NoError(t, err)
	assert.NotNil(t, localWorker)

	var bs []byte
	err = localWorker.SetContext(bs)
	assert.NoError(t, err)

	err = localWorker.BeforeRun()
	assert.NoError(t, err)

	err = localWorker.Do()
	assert.NoError(t, err)

	time.Sleep(time.Second * 5)

	err = localWorker.AfterRun()
	assert.NoError(t, err)
	sent, missed := localWorker.Statistics()
	assert.NotNil(t, sent)
	assert.NotNil(t, missed)

	col, b, _ := localWorker.CheckoutCollector()
	assert.NotNil(t, col)
	assert.NotNil(t, b)

	localWorker.Done()
	localWorker.Teardown()

	l, _ := NewLocalWorker(LocalWorkerConfig{0, 5, 20, 5, time.Millisecond * 5, time.Second * 3})
	l.Do()
	l.cancel()
	time.Sleep(time.Second * 4)

}

func TestLocalNewWorkers(t *testing.T) {
	defer os.RemoveAll("./benchmark")

	localconfig := `
	[engine]
	rate = 1
	duration = "5s"
	cap = 1
	`

	os.MkdirAll("./benchmark/testLocal", 0755)
	ioutil.WriteFile("./benchmark/testLocal/config.toml", []byte(localconfig), 0644)

	viper.AddConfigPath("benchmark/testLocal")
	viper.ReadInConfig()
	worker, err := NewWorkers()
	assert.NotNil(t, worker)
	assert.NoError(t, err)
}

func TestRemoteNewWorkers(t *testing.T) {
	t.Skip()
	defer os.RemoveAll("./.tar.gz")
	defer os.RemoveAll("./benchmark")
	defer os.RemoveAll("./benchmark")

	localconfig := `
	[engine]
	rate = 1
	duration = "5s"
	cap = 1
	`
	remoteconfig := `
	[engine]
	rate = 1
	duration = "5s"
	cap = 1
	urls = ["localhost:8200"]
	`
	os.MkdirAll("./benchmark/testLocal", 0755)
	os.MkdirAll("./benchmark/testRemote", 0755)

	ioutil.WriteFile("./benchmark/testLocal/config.toml", []byte(localconfig), 0644)
	ioutil.WriteFile("./benchmark/testRemote/config.toml", []byte(remoteconfig), 0644)

	config, _ := os.Open("benchmark/testRemote/config.toml")
	viper.ReadConfig(config)
	workers, err := NewWorkers()
	assert.Nil(t, workers)
	assert.Error(t, err)

	viper.Set("__BenchmarkDirPath__", "benchmark/testLocal")
	workers, err = NewWorkers()
	assert.Nil(t, workers)
	assert.Error(t, err)

	viper.Set("engine.urls", `localhost:8100`)
	viper.Set("__BenchmarkDirPath__", "benchmark/testRemote")
	workers, err = NewWorkers()
	assert.NotNil(t, workers)
	assert.NoError(t, err)

}
