package client

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
 * @brief Client is used to communicate with worker by master
 * @file client.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	fcom "github.com/hyperbench/hyperbench-common/common"
	json "github.com/json-iterator/go"

	"github.com/hyperbench/hyperbench/core/collector"
	"github.com/hyperbench/hyperbench/core/network"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	httpProtocol = "http://"
	// network interface
	nonce      = "nonce"
	index      = "index"
	setContext = "set context"
	beforeRun  = "before run"
	do         = "do"
	afterRun   = "after run"
	setNonce   = "set nonce"
	initPath   = "init"
	tearDown   = "teardown"
)

// Client is used to communicate with worker by master.
type Client struct {
	url    string
	nonce  string
	index  int
	logger *logging.Logger
	path   string
	err    error
}

// NewClient create Client.
func NewClient(index int, url string) *Client {
	if !strings.HasPrefix(url, httpProtocol) {
		url = httpProtocol + url
	}
	return &Client{
		url:    url,
		index:  index,
		nonce:  strconv.Itoa(int(time.Now().UnixNano())),
		logger: fcom.GetLogger("client"),
		path:   viper.GetString(fcom.BenchmarkArchivePath),
	}
}

// Init tell worker to init for execute tx.
func (c *Client) Init() error {
	c.err = c.setNonce()
	if c.err != nil {
		return c.err
	}

	defer c.teardownWhileErr()

	c.err = c.upload()
	if c.err != nil {
		return c.err
	}

	c.err = c.init()
	if c.err != nil {
		return c.err
	}
	return nil
}

// SetContext set the context of worker passed from Master.
func (c *Client) SetContext(d []byte) error {
	defer c.teardownWhileErr()
	values := url.Values{
		nonce:     {c.nonce},
		"context": {network.Bytes2Hex(d)},
	}
	c.err = c.callWithValues(setContext, network.SetContextPath, values)
	if c.err != nil {
		return c.err
	}
	return nil
}

// CheckoutCollector checkout collector.
func (c *Client) CheckoutCollector() (collector.Collector, bool, error) {
	var err error

	defer func() {
		if err != nil {
			c.logger.Errorf("check out collector fail: %v", err)
		}
	}()

	var resp *http.Response
	values := url.Values{nonce: {c.nonce}}
	resp, err = http.PostForm(c.url+network.CheckoutCollectorPath, values)
	if err != nil {
		c.logger.Error(err)
		return nil, false, err
	}
	// nolint
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		c.logger.Error(err)
		return nil, false, err
	}

	var bs []byte
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		c.err = err
		c.logger.Error(err)
		return nil, false, err
	}

	var retJSON map[string]interface{}
	err = json.Unmarshal(bs, &retJSON)
	if err != nil {
		c.err = err
		c.logger.Error(err)
		return nil, false, err
	}

	colType, _ := retJSON["type"].(string)
	colData, _ := retJSON["col"].(string)
	colValid, _ := retJSON["valid"].(bool)
	if colValid {
		col := collector.NewTDigestCollectorBuilder(colType)()
		err = col.Merge(network.Hex2Bytes(colData))
		if err != nil {
			c.err = err
			c.logger.Error(err)
			return nil, false, err
		}
		return col, colValid, nil
	}
	return nil, false, nil
}

// Teardown close the worker manually.
func (c *Client) Teardown() {
	err := c.teardown()
	if err != nil {
		c.logger.Errorf("call %s/teardown err:%v", c.url, err)
	}
}

// BeforeRun call user hook
func (c *Client) BeforeRun() error {
	defer c.teardownWhileErr()
	return c.callWithValues(beforeRun, network.BeforeRunPath, url.Values{nonce: {c.nonce}})
}

// Do call the workers to running
func (c *Client) Do() error {
	defer c.teardownWhileErr()
	return c.callWithValues(do, network.DoPath, url.Values{nonce: {c.nonce}})
}

// Statistics get the number of sent and missed transactions
func (c *Client) Statistics() (int64, int64) {
	var resp *http.Response
	values := url.Values{nonce: {c.nonce}}
	resp, err := http.PostForm(c.url+network.StatisticsPath, values)
	if err != nil {
		c.logger.Error(err)
		return 0, 0
	}
	// nolint
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		c.logger.Error(err)
		return 0, 0
	}

	var bs []byte
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		c.err = err
		c.logger.Error(err)
		return 0, 0
	}

	var retJSON map[string]interface{}
	err = json.Unmarshal(bs, &retJSON)
	if err != nil {
		c.err = err
		c.logger.Error(err)
		return 0, 0
	}

	sent, _ := retJSON["sent"].(string)
	missed, _ := retJSON["missed"].(string)
	Sent, _ := strconv.ParseInt(sent, 10, 64)
	Missed, _ := strconv.ParseInt(missed, 10, 64)
	return Sent, Missed
}

// AfterRun call user hook
func (c *Client) AfterRun() error {
	defer c.teardownWhileErr()
	return c.callWithValues(afterRun, network.AfterRunPath, url.Values{nonce: {c.nonce}})
}

func (c *Client) teardownWhileErr() {
	if c.err != nil {
		_ = c.teardown()
	}
}

func (c *Client) setNonce() error {
	return c.callWithValues(setNonce, network.SetNoncePath, url.Values{nonce: {c.nonce}})
}

func (c *Client) upload() error {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	nonceWriter, _ := bodyWriter.CreateFormField(nonce)
	_, _ = nonceWriter.Write([]byte(c.nonce))
	dirWriter, _ := bodyWriter.CreateFormField(network.ConfigPath)
	_, _ = dirWriter.Write([]byte(viper.GetString(fcom.BenchmarkConfigPath)))
	filePathWriter, _ := bodyWriter.CreateFormField(network.FilePath)
	_, _ = filePathWriter.Write([]byte(c.path))
	fileWriter, _ := bodyWriter.CreateFormFile(network.FileName, c.path)

	file, _ := os.Open(c.path)
	//nolint
	defer file.Close()

	_, _ = io.Copy(fileWriter, file)
	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	resp, err := http.Post(c.url+network.UploadPath, contentType, bodyBuffer)
	// nolint
	defer resp.Body.Close()

	if err != nil {
		return errors.Wrapf(err, "%v can not receive file", c.url)
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return errors.Wrapf(err, "%v can not receive file", c.url)
	}

	return nil
}

func (c *Client) init() error {
	return c.callWithValues(initPath, network.InitPath, url.Values{nonce: {c.nonce}, index: {strconv.Itoa(c.index)}})
}

// Testinit used for unit test
func (c *Client) Testinit() error {
	return c.init()
}

// TestsetNonce used for unit test
func (c *Client) TestsetNonce() error {
	return c.setNonce()
}

func (c *Client) callWithValues(method string, path string, values url.Values) (err error) {

	defer func() {
		if err != nil {
			c.err = err
		}
	}()

	var resp *http.Response
	// set nonce
	resp, err = http.PostForm(c.url+path, values)
	if err != nil {
		return errors.Wrapf(err, "%v can not %v", c.url, method)
	}
	// nolint
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return errors.Wrapf(err, "%v can not %v", c.url, method)
	}

	return nil
}

func (c *Client) teardown() error {
	return c.callWithValues(tearDown, network.TeardownPath, url.Values{nonce: {c.nonce}})
}
