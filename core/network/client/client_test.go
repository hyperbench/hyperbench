package client_test

import (
	"testing"
	"time"

	"github.com/meshplus/hyperbench/core/network/client"
	"github.com/meshplus/hyperbench/core/network/server"
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
