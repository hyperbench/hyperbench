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
 * @brief the unit test for net.go
 * @file net_test.go
 * @author: shinyxhh
 * @date 2021-11-30
 */
package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNet(t *testing.T) {
	s := Bytes2Hex(nil)
	assert.Equal(t, s, "")
	bs := Hex2Bytes("nil")
	assert.NotNil(t, bs)
}
