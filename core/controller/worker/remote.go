package worker

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
 * @brief Provide RemoteWorker, the agent of remote worker.
 * @file remote.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"github.com/hyperbench/hyperbench/core/network/client"
)

// RemoteWorker is the agent of remote worker.
type RemoteWorker struct {
	*client.Client
}

// NewRemoteWorker create RemoteWorker.
func NewRemoteWorker(index int, url string) (*RemoteWorker, error) {
	c := client.NewClient(index, url)
	err := c.Init()
	if err != nil {
		return nil, err
	}
	return &RemoteWorker{
		Client: c,
	}, nil
}
