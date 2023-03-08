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
 * @brief use cobra provide cmd function
 * @file recorder_test.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */
package recorder

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRecorder(t *testing.T) {
	viper.Set("recorder.csv", "")
	recorder := NewRecorder()
	assert.NotNil(t, recorder)

	recorder.Process(fcom.Report{
		Cur: &fcom.Data{
			Results: []fcom.AggData{
				{
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
