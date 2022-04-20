package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/viper"
)

func TestController(t *testing.T) {
	m := make(map[string]interface{})
	m["engine.cap"] = 1
	m["engine.rate"] = 1
	m["engine.instant"] = 1
	m["engine.wait"] = 0
	m["engine.duration"] = 1
	viper.MergeConfigMap(m)
	ctl, err := NewController()
	assert.NoError(t, err)
	assert.NotNil(t, ctl)

	err = ctl.Prepare()
	assert.NoError(t, err)

	err = ctl.Run()
	assert.NoError(t, err)
}
