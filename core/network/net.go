// Package network is used to distribute controlling
package network

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
