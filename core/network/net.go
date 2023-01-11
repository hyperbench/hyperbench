// Package network is used to distribute controlling
package network

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
 * @brief Define constants and provide utility functions
 * @file net.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"encoding/hex"
)

const (
	// SetNoncePath set nonce path.
	SetNoncePath = "/set-nonce"
	// UploadPath upload path.
	UploadPath = "/upload"
	// InitPath init path.
	InitPath = "/init"
	// SetContextPath set context path.
	SetContextPath = "/set-context"
	// BeforeRunPath before run path.
	BeforeRunPath = "/before-run"
	// DoPath do path.
	DoPath = "/do"
	// StatisticsPath Statistics path.
	StatisticsPath = "/statistics"
	// AfterRunPath after run path.
	AfterRunPath = "/after-run"
	// TeardownPath teardown path.
	TeardownPath = "/teardown"
	// CheckoutCollectorPath checkout collector path.
	CheckoutCollectorPath = "/checkout-collector"
	// ConfigPath key of configPath
	ConfigPath = "configDir"
	// FileName key of file
	FileName = "file"
	// FilePath key of file path
	FilePath = "filepath"
)

// Bytes2Hex convert bytes to hex.
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

// Hex2Bytes convert hex to bytes.
func Hex2Bytes(h string) []byte {
	b, _ := hex.DecodeString(h)
	return b
}
