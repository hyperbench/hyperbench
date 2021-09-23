package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"math/big"
	"reflect"
	"strings"
	"testing"
	"time"
)

type TestConfig struct {
	Rate     int           `mapstruct:"rate"`
	Duration time.Duration `mapstruct:"duration"`
}

func MyFunc(*TestConfig) {}

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
	//rv := TReflectParam(MyFunc, param)
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
