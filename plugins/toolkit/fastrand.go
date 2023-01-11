package toolkit

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
 * @brief Provide rand tool function
 * @file fastrand.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	randCh     = make(chan *rand.Rand, runtime.NumCPU())
	randChOnce sync.Once
)

const (
	chars    = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsLen = len(chars)
)

func initRandCh() {
	for i := 0; i < runtime.NumCPU(); i++ {
		randCh <- rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}

func randomString(l uint) string {
	randChOnce.Do(initRandCh)

	r := <-randCh
	s := make([]byte, l)
	for i := 0; i < int(l); i++ {
		s[i] = chars[r.Intn(charsLen)]
	}
	randCh <- r
	return string(s)
}

func randomInt(min, max int) int {
	randChOnce.Do(initRandCh)
	r := <-randCh
	i := r.Intn(max-min) + min
	randCh <- r
	return i
}
