package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"

	fcom "github.com/meshplus/hyperbench-common/common"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Rate     int           `mapstruct:"rate"`
	Duration time.Duration `mapstruct:"duration"`
}

func TestUtils(t *testing.T) {
	a := AggData2CSV(nil, "", fcom.AggData{
		Label: "11",
	})
	assert.NotNil(t, a)

	b := Latency2CSV(nil, fcom.Latency{})
	assert.NotNil(t, b)

	i := uint(1)
	c := i2s(i)
	assert.NotNil(t, c)

	j := int32(1)
	d := i2s(j)
	assert.NotNil(t, d)
}

func TestReflectParam(t *testing.T) {

	configTxt := `
[[schedules]]
type = "constant"
duration = "10s"
rate = 100

[[schedules]]
type = "monotonic"
duration = "10s"
rate = 100
[schedules.option]
step = -10
interval = "1s"

[[schedules]]
type = "auto"
duration = "10s"
rate = 100
[schedules.option]
failed-threshold = 0.95
adjust-factor = 0.95
`
	viper.SetConfigType("toml")
	_ = viper.ReadConfig(strings.NewReader(configTxt))
	viper.Sub("schedules")
	v := viper.Get("schedules")
	param := v.([]interface{})[0].(map[string]interface{})
	//fmt.Println(v)
	//fmt.Println(rv)
	conf := &TestConfig{}
	//mapstructure.Decode(param, conf)

	dc := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     conf,
	}
	d, err := mapstructure.NewDecoder(&dc)
	fmt.Println("====", err)
	_ = d.Decode(param)
	fmt.Println(dc.Result)

}

func TReflectParam(fn interface{}, param map[string]interface{}) []interface{} {
	t := reflect.TypeOf(fn)
	ret := make([]interface{}, 0, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		valType := t.In(i)
		val := reflect.New(valType)
		fmt.Println(val.Kind(), val.Type())
		fmt.Println("===", val.NumField())
		itf := val.Interface()
		fmt.Println("===", itf)
		_ = mapstructure.Decode(param, &itf)
		//fmt.Println("===", itf)
		//bs, _ := json.Marshal(itf)
		//fmt.Println("---", err.Error())
		//fmt.Println("---", string(bs))
		ret = append(ret, itf)
	}
	return ret
}

func TestRand(t *testing.T) {
	_, _ = rand.Int(rand.Reader, big.NewInt(1000))
}
