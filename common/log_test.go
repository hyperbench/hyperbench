package common

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	log1 := GetLogger("")
	assert.NotNil(t, log1)

	log2 := GetLogger("register")
	assert.Equal(t, log2.Module, "register")

	log3 := GetLogger("register")
	assert.NotNil(t, log3)
}

func TestInitLog(t *testing.T) {
	m := make(map[string]interface{})
	m["recorder.log.dump"] = true
	m["recorder.log.fir"] = ""
	viper.MergeConfigMap(m)
	file := InitLog()
	assert.NotNil(t, file)
	defer os.RemoveAll("./log")

}
