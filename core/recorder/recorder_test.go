package recorder

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRecorder(t *testing.T) {
	viper.Set(fcom.RecorderCsvPath, "")
	recorder := NewRecorder()
	assert.NotNil(t, recorder)

	recorder.Process(fcom.Report{
		Cur: &fcom.Data{
			Results: []fcom.AggData{
				fcom.AggData{
					Label: "11",
				},
			},
		},
		Sum: &fcom.Data{},
	})

	assert.NotNil(t, GetCSVPath())

	recorder.Release()

	os.RemoveAll("./csv")

}

func BenchmarkWrite(b *testing.B) {
	times := 100
	b.Run("bytes", func(b *testing.B) {
		var pool sync.Pool
		pool.New = func() interface{} {
			return bytes.NewBuffer(nil)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := pool.Get().(*bytes.Buffer)
			buf.Reset()
			for j := 0; j < times; j++ {
				buf.WriteString("1234567890")
			}
			_ = buf.String()
			pool.Put(buf)
		}
	})

	b.Run("bytes-no-str", func(b *testing.B) {
		var pool sync.Pool
		pool.New = func() interface{} {
			return bytes.NewBuffer(nil)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := pool.Get().(*bytes.Buffer)
			buf.Reset()
			for j := 0; j < times; j++ {
				buf.WriteString("1234567890")
			}
			pool.Put(buf)
		}
	})

	b.Run("builder", func(b *testing.B) {
		var pool sync.Pool
		for i := 0; i < 4; i++ {
			pool.Put(&strings.Builder{})
		}
		pool.New = func() interface{} {
			return &strings.Builder{}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := pool.Get().(*strings.Builder)
			buf.Reset()
			for j := 0; j < times; j++ {
				buf.WriteString("1234567890")
			}
			_ = buf.String()
			pool.Put(buf)
		}
	})

}

func TestLogTile(t *testing.T) {
	newLogProcessor(fcom.GetLogger("test")).logTitle()
}

func TestInfluxdb(t *testing.T) {
	t.Skip()
	viper.Set(fcom.RecorderInflucDBPath, "")
	viper.Set(fcom.BenchmarkDirPath, "http://172.22.66.159:8086")
	viper.Set(fcom.InfluxDBUrlPath, "hyperbench")
	viper.Set(fcom.InfluxDBDatabasePath, "admin")
	viper.Set(fcom.InfluxDBUsernamePath, "configured")
	viper.Set(fcom.InfluxDBPasswordPath, "user1")
	recorder := NewRecorder()
	assert.NotNil(t, recorder)
	recorder.Process(fcom.Report{
		Cur: &fcom.Data{
			Results: []fcom.AggData{
				{
					Label:    "11",
					Time:     time.Now().UnixNano(),
					Duration: 1e10,
					Num:      100,
					Statuses: map[fcom.Status]int{
						fcom.Success: 98,
						fcom.Failure: 2,
					},
				},
			},
		},
		Sum: &fcom.Data{},
	})
	recorder.Release()
}
