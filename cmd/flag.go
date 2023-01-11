package cmd

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
 * @brief Define the provided cmd commands and parameters
 * @file flag.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"github.com/spf13/cobra"
)

var (
	start = &cobra.Command{
		Use:     "start",
		Short:   "start a benchmark",
		Args:    cobra.ExactArgs(1),
		Example: "hyperbench start benchmark/transfer",
	}

	create = &cobra.Command{
		Use:     "new [testplan Name]",
		Short:   "initialize a test plan",
		Args:    cobra.ExactArgs(1),
		Example: "hyperbench new myTest",
	}

	rootCmd = &cobra.Command{
		Use:     "hyperbench",
		Example: "hyperbench --doc ./doc (generate document to specify path)",
	}

	version = &cobra.Command{
		Use:     "version",
		Short:   "get code version",
		Example: "hyperbench version",
	}

	initDir = &cobra.Command{
		Use:     "init",
		Short:   "init a stress test dir",
		Example: "hyperbench init",
	}

	worker = &cobra.Command{
		Use:     "worker",
		Short:   "start as a worker server ",
		Example: "hyperbench worker",
	}

	enableDebug = rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	document    = rootCmd.PersistentFlags().String("doc", "", "use to create doc and specify the doc path")
	port        = worker.PersistentFlags().IntP("port", "p", 8080, "port of worker")
)
