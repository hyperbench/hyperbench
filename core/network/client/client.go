// net is the
package client

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

	json "github.com/json-iterator/go"
	fcom "github.com/meshplus/hyperbench-common/common"

	"github.com/meshplus/hyperbench/core/collector"
	"github.com/meshplus/hyperbench/core/network"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	httpProtocol = "http://"
)

// Client is used to communicate with worker by master.
type Client struct {
	url      string
	nonce    string
	index    int
	logger   *logging.Logger
	path     string
	err      error
	finished bool
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
		"nonce":   {c.nonce},
		"context": {network.Bytes2Hex(d)},
	}
	c.err = c.callWithValues("set context", network.SetContextPath, values)
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
	values := url.Values{"nonce": {c.nonce}}
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

// Do call the workers to running
func (c *Client) Do() error {
	defer c.teardownWhileErr()
	return c.callWithValues("do", network.DoPath, url.Values{"nonce": {c.nonce}})
}

func (c *Client) teardownWhileErr() {
	if c.err != nil {
		_ = c.teardown()
	}
}

func (c *Client) setNonce() error {
	return c.callWithValues("set nonce", network.SetNoncePath, url.Values{"nonce": {c.nonce}})
}

func (c *Client) upload() error {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	nonceWriter, _ := bodyWriter.CreateFormField("nonce")
	_, _ = nonceWriter.Write([]byte(c.nonce))
	fileWriter, _ := bodyWriter.CreateFormFile("file", c.path)

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
	return c.callWithValues("init", network.InitPath, url.Values{"nonce": {c.nonce}, "index": {strconv.Itoa(c.index)}})
}

func (c *Client) Testinit() error {
	return c.init()
}

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
	return c.callWithValues("teardown", network.TeardownPath, url.Values{"nonce": {c.nonce}})
}
