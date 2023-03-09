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
 * @brief the unit test for base.go
 * @file base_test.go
 * @author: shinyxhh
 * @date 2021-11-30
 */
package base

import (
	"testing"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/stretchr/testify/assert"
)

func TestBaseVm(t *testing.T) {
	t.Skip()
	base := NewVMBase(ConfigBase{
		Path: "",
	})
	Type := base.Type()
	assert.Equal(t, Type, "base")

	err := base.BeforeDeploy()
	assert.NoError(t, err)

	err = base.DeployContract()
	assert.NoError(t, err)

	err = base.BeforeGet()
	assert.NoError(t, err)

	bs, err := base.GetContext()
	assert.NoError(t, err)
	assert.NotNil(t, bs)

	res, err := base.Statistic(nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	err = base.BeforeSet()
	assert.NoError(t, err)

	err = base.SetContext(nil)
	assert.NoError(t, err)

	err = base.BeforeRun()
	assert.NoError(t, err)

	result, err := base.Run(fcom.TxContext{})
	assert.NoError(t, err)
	assert.NotNil(t, result)

	err = base.AfterRun()
	assert.NoError(t, err)

	_, err = base.LogStatus()
	assert.NoError(t, err)

	base.Close()

}
