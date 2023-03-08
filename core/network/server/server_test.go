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
 * @file server_test.go
 * @author: shinyxhh
 * @date 2021-11-30
 */
package server

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/hyperbench/hyperbench/core/network/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	defer os.RemoveAll("./benchmark")
	os.MkdirAll("./benchmark/111", 0755)
	ioutil.WriteFile("./benchmark/111/a.toml", []byte(""), 0644)
	viper.Set(fcom.BenchmarkArchivePath, "benchmark/111")
	svr := NewServer(0)
	assert.NotNil(t, svr)

	go svr.Start()

	cli := client.NewClient(0, "localhost:8080")
	assert.NotNil(t, cli)

	go os.RemoveAll("./benchmark")

	err := cli.Init()
	assert.Error(t, err)

	m := make(map[string]interface{})
	m["engine.urls"] = `"localhost:8080"`
	m["engine.rate"] = 1
	m["engine.duration"] = 5
	m["engine.cap"] = 1
	m["engine.instant"] = 1
	m["engine.wait"] = 0
	viper.MergeConfigMap(m)

	err = cli.TestsetNonce()
	assert.NoError(t, err)
	err = cli.Testinit()
	assert.NoError(t, err)

	go cli.Do()

	go cli.CheckoutCollector()

	time.Sleep(time.Second * 2)

	cli.Teardown()
}

func TestSetNonce(t *testing.T) {
	svr := NewServer(0)
	go svr.Start()
	cli := client.NewClient(1, "localhost:8080")

	cli.TestsetNonce()
	err := cli.TestsetNonce()
	assert.Error(t, err)
	cli.Teardown()
	err = callWithValues("set nonce", "/set-nonce", url.Values{})
	assert.Error(t, err)
	err = callWithValues("set nonce", "/set-nonce", url.Values{"nonce": {"q"}})
	assert.Error(t, err)
}
func TestInit(t *testing.T) {
	svr := NewServer(0)
	go svr.Start()
	cli := client.NewClient(2, "localhost:8080")

	err := cli.Testinit()
	assert.Error(t, err)

	cli.TestsetNonce()

	err = cli.Testinit()
	assert.Error(t, err)
	cli.Teardown()

}

func TestSetContext(t *testing.T) {
	svr := NewServer(0)
	go svr.Start()
	cli := client.NewClient(3, "localhost:8080")

	cli.TestsetNonce()
	err := cli.SetContext(nil)
	assert.Error(t, err)
	err = cli.SetContext(nil)
	assert.Error(t, err)
}

func TestDo(t *testing.T) {
	svr := NewServer(0)
	go svr.Start()
	cli := client.NewClient(4, "localhost:8080")

	cli.TestsetNonce()
	err := cli.BeforeRun()
	assert.Error(t, err)
	err = cli.Do()
	assert.Error(t, err)
	err = cli.AfterRun()
	assert.Error(t, err)
	sent, missed := cli.Statistics()
	assert.NotNil(t, sent)
	assert.NotNil(t, missed)

}

func TestCheckoutCollector(t *testing.T) {
	svr := NewServer(0)
	go svr.Start()
	cli := client.NewClient(5, "localhost:8080")

	cli.TestsetNonce()
	col, b, err := cli.CheckoutCollector()
	assert.Error(t, err)
	assert.Nil(t, col)
	assert.Equal(t, b, false)

	cli.Teardown()

	col, b, err = cli.CheckoutCollector()
	assert.Error(t, err)
	assert.Nil(t, col)
	assert.Equal(t, b, false)

}

func callWithValues(method string, path string, values url.Values) (err error) {
	var resp *http.Response
	resp, _ = http.PostForm("http://localhost:8080"+path, values)
	if err != nil {
		return errors.Wrapf(err, "%v can not %v", "http://localhost:8080", method)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return errors.Wrapf(err, "%v can not %v", "http://localhost:8080", method)
	}

	return nil
}
