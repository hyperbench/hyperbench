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
 * @brief the unit test for master.go
 * @file master_test.go
 * @author: shinyxhh
 * @date 2021-11-30
 */
package master

import (
	fcom "github.com/hyperbench/hyperbench-common/common"
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
	_, err = localMaster.Statistic(&fcom.ChainInfo{TimeStamp: 1}, &fcom.ChainInfo{TimeStamp: 1})
	assert.NoError(t, err)
	_, err = localMaster.LogStatus()
	assert.NoError(t, err)

}
