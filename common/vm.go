package common

import (
	"context"
)

// TxContext is the context of transaction (`vm.Run`)
type TxContext struct {
	context.Context
	TxIndex
}

// TxIndex is the unique index of a transaction (`vm.Run`)
type TxIndex struct {
	EngineIdx int64 // EngineIdx is the index of engine (maybe support multiple test stage in future)
	TxIdx     int64 // TxIdx is the index of transaction
}

// VMContext is the context of vm
type VMContext struct {
	WorkerIdx int64 `mapstructure:"worker"` // WorkerIdx is the index of worker
	VMIdx     int64 `mapstructure:"vm"`     // VMIdx is the index of vm
}
