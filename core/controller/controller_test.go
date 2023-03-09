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
 * @brief the unit test for controller.go
 * @file controller_test.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */
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
	m["verify.enable"] = true
	viper.MergeConfigMap(m)
	ctl, err := NewController()
	assert.NoError(t, err)
	assert.NotNil(t, ctl)

	err = ctl.Prepare()
	assert.NoError(t, err)

	err = ctl.Run()
	assert.NoError(t, err)
}
