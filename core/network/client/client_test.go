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
 * @file client_test.go
 * @author: shinyxhh
 * @date 2021-11-30
 */
package client_test

import (
	"testing"
	"time"

	"github.com/hyperbench/hyperbench/core/network/client"
	"github.com/hyperbench/hyperbench/core/network/server"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	svr := server.NewServer(8085)
	assert.NotNil(t, svr)

	go svr.Start()

	cli := client.NewClient(0, "localhost:8085")
	assert.NotNil(t, cli)

	err := cli.Init()
	assert.Error(t, err)

	m := make(map[string]interface{})
	m["engine.urls"] = `"localhost:8085"`
	m["engine.rate"] = 1
	m["engine.duration"] = 5
	m["engine.cap"] = 1
	m["engine.instant"] = 1
	m["engine.wait"] = 1
	m["client.plugin"] = "hyperchain.so"
	viper.MergeConfigMap(m)

	err = cli.TestsetNonce()
	assert.NoError(t, err)

	err = cli.Testinit()
	assert.NoError(t, err)

	err = cli.SetContext(nil)
	assert.NoError(t, err)

	err = cli.BeforeRun()
	assert.NoError(t, err)

	go cli.Do()

	go cli.CheckoutCollector()

	time.Sleep(time.Second * 2)

	err = cli.AfterRun()
	assert.NoError(t, err)

	sent, missed := cli.Statistics()
	assert.NotNil(t, sent)
	assert.NotNil(t, missed)

	cli.Teardown()

	time.Sleep(time.Second * 2)
}
