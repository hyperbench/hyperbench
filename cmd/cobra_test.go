/*
@Time : 2019-03-08 11:20
@Author : lmm
*/
package cmd

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitCmd(t *testing.T) {
	err := InitCmd(func() {})
	assert.NoError(t, err)
	command := GetRootCmd()
	d := []string{"init", "--debug=true", "--doc=1"}
	command.SetArgs(d)
	_, err = command.ExecuteC()
	assert.NoError(t, err)
	assert.Equal(t, *enableDebug, true)
	assert.Equal(t, *document, "1")
	d = []string{"init", "--debug=false", "--doc="}
	command.SetArgs(d)
	_, err = command.ExecuteC()
	assert.NoError(t, err)
}

func TestStart(t *testing.T) {
	defer os.RemoveAll("./benchmark")
	config := `
	[engine]
	rate = 10
	duration = "5s"
	cap = 1
	urls = ["localhost:8080","localhost:8085","localhost:8082"]

	[client]
	script = "testData/ethInvoke/script.lua"
	type = ""
	`
	os.MkdirAll("./benchmark/ethInvoke", 0755)
	ioutil.WriteFile("./benchmark/ethInvoke/config.toml", []byte(config), 0644)

	err := InitCmd(func() {})
	assert.NoError(t, err)
	command := GetRootCmd()
	d := []string{"start", ""}
	command.SetArgs(d)
	_, err = command.ExecuteC()
	assert.NoError(t, err)
	d = []string{"start", "./benchmark/ethInvoke"}
	command.SetArgs(d)
	_, err = command.ExecuteC()
	assert.NoError(t, err)
	viper.Set("engine.urls", "")
	d = []string{"start", "./benchmark/ethInvoke"}
	command.SetArgs(d)
	_, err = command.ExecuteC()
	assert.NoError(t, err)
}

func TestWorker(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	err := InitCmd(func() {})
	assert.NoError(t, err)
	command := GetRootCmd()
	d := []string{"worker", "-p", "8090"}
	command.SetArgs(d)
	go command.ExecuteC()
	time.Sleep(time.Second * 1)
}

func TestVersion(t *testing.T) {
	err := InitCmd(func() {})
	assert.NoError(t, err)
	command := GetRootCmd()
	d := []string{"version"}
	command.SetArgs(d)
	_, err = command.ExecuteC()
	assert.NoError(t, err)
}

func TestInit(t *testing.T) {
	err := InitCmd(func() {})
	assert.NoError(t, err)
	command := GetRootCmd()
	d := []string{"init"}
	command.SetArgs(d)
	_, err = command.ExecuteC()
	assert.NoError(t, err)
}
