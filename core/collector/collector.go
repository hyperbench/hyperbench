package collector

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
 * @brief Collector is used to collect result and generate statistic data group by label
 * @file collector.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	fcom "github.com/hyperbench/hyperbench-common/common"
)

// Collector is used to collect result and generate statistic data group by label
// Collector may not be implement concurrently safe, so you should receive data in a goroutine
type Collector interface {
	// Type return the types of collector
	Type() string

	// Add append result to statistic
	Add([]*fcom.Result)

	// Serialize generate serialized data to pass through network in remote mode
	Serialize() []byte

	// Merge merge serialized data
	Merge([]byte) error

	// MergeC try to merge a Collector, if it can not do this, just raise a error
	MergeC(Collector) error

	// Get get current statistic data group by label
	Get() *fcom.Data

	// Reset reset data should reset the time window and clean data
	Reset()
}
