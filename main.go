package main

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
 * @brief The entrance of HyperBench
 * @file main.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */
import (
	"fmt"
	"github.com/hyperbench/hyperbench/cmd"
	"os"
	"runtime/pprof"
	"time"
)

func main() {

	err := cmd.InitCmd(debug)
	if err != nil {
		fmt.Println("cmd init fail: ", err)
		return
	}

	err = cmd.GetRootCmd().Execute()
	if err != nil {
		fmt.Println("cmd execute fail: ", err)
	}
}

// Debug related config keys
const (
	PprofRecordDuration = "5s"
	PprofTimeFmt        = "2006-01-02-15-04-05"
)

func debug() {

	duration, err := time.ParseDuration(PprofRecordDuration)
	if err != nil {
		return
	}

	go recordPProf(duration)
}

func recordPProf(duration time.Duration) {
	var (
		cpuFilePath string
		memFilePath string
		cpuFile     *os.File
		memFile     *os.File
	)

	dir := "./debug"
	err := ensurePathExists(dir)
	if err != nil {
		return
	}

	// reset inline function
	reset := func() {
		// use the expected file closed time as the file name's suffix
		timeSuffix := time.Now().Add(duration).Format(PprofTimeFmt)

		cpuFilePath = fmt.Sprint(dir, "/cpu_", timeSuffix)
		memFilePath = fmt.Sprint(dir, "/mem_", timeSuffix)

		cpuFile, _ = os.Create(cpuFilePath)
		memFile, _ = os.Create(memFilePath)

		// start pprof
		err := pprof.StartCPUProfile(cpuFile)
		if err != nil {
			return
		}
	}

	tick := time.NewTicker(duration)

	reset()

	//nolint
	for {
		select {
		case <-tick.C:

			pprof.StopCPUProfile()
			err := cpuFile.Close()
			if err != nil {
				continue
			}
			err = pprof.WriteHeapProfile(memFile)
			if err != nil {
				continue
			}
			err = memFile.Close()
			if err != nil {
				continue
			}

			reset()
		}
	}

}

func ensurePathExists(path string) error {

	_, err := os.Stat(path)
	// already exist
	if err == nil {
		return nil
	}

	// not exist
	if os.IsNotExist(err) {
		// make full path
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}

	return err
}
