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
 * @file blockchain_test.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */
package blockchain

import (
	"testing"

	"github.com/hyperbench/hyperbench-common/base"
	fcom "github.com/hyperbench/hyperbench-common/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitPlugin(t *testing.T) {
	t.Skip()
	InitPlugin()

	viper.Set(fcom.ClientTypePath, "hyperchain")
	InitPlugin()

	viper.Set(fcom.ClientTypePath, "fabric")
	InitPlugin()

	viper.Set(fcom.ClientTypePath, "eth")
	InitPlugin()

	viper.Set(fcom.ClientTypePath, "xuperchain")
	InitPlugin()

}
func TestNewBlockchain(t *testing.T) {
	t.Skip()
	bk, err := NewBlockchain(base.ClientConfig{})
	assert.NotNil(t, bk)
	assert.NoError(t, err)

	bk, err = NewBlockchain(base.ClientConfig{
		ClientType: "eth",
	})
	assert.Nil(t, bk)
	assert.Error(t, err)

	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	bk, err = NewBlockchain(base.ClientConfig{
		ClientType: "fabric",
	})
	assert.Nil(t, bk)
	assert.Error(t, err)

	bk, err = NewBlockchain(base.ClientConfig{
		ClientType: "xuperchain",
	})
	assert.Nil(t, bk)
	assert.Error(t, err)
}

func TestNewHyperchain(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	bk, err := NewBlockchain(base.ClientConfig{
		ClientType: "flato",
	})
	assert.Nil(t, bk)
	assert.Error(t, err)
}
