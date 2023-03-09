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
 * @brief the unit test for vm.go
 * @file vm_test.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */
package vm

import (
	"testing"

	"github.com/hyperbench/hyperbench/vm/base"
	"github.com/stretchr/testify/assert"
)

func TestNewVM(t *testing.T) {
	t.Skip()
	vm, err := NewVM("lua", base.ConfigBase{})
	assert.Nil(t, vm)
	assert.Error(t, err)

	vm, err = NewVM("", base.ConfigBase{})
	assert.NotNil(t, vm)
	assert.NoError(t, err)
}
