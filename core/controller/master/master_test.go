package master

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLocalMaster(t *testing.T) {
	config := `
	[client]
	type = ""
	contract = "testData/contract"
	`

	defer os.RemoveAll("./benchmark")

	os.Mkdir("./benchmark", 0755)

	ioutil.WriteFile("./benchmark/config.toml", []byte(config), 0644)

	viper.AddConfigPath("benchmark")
	viper.ReadInConfig()
	localMaster, err := NewLocalMaster()
	assert.NoError(t, err)
	bs, err := localMaster.GetContext()
	assert.NoError(t, err)
	assert.NotNil(t, bs)
	err = localMaster.Prepare()
	assert.NoError(t, err)
	_, err = localMaster.Statistic(1, 1)
	assert.NoError(t, err)
	_, err = localMaster.LogStatus()
	assert.NoError(t, err)

}
