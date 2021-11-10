package common

// Status tx execute status.
type Status string

const (
	// Failure means client sent transaction unsuccessfully
	Failure Status = "failure"
	// Success means client sent transaction successfully
	Success Status = "success"
	// Unknown means sent transaction isn't written in block stably for asynchronous processing
	Unknown Status = "unknown"
	// Confirm means sent transaction is written in block stably for asynchronous processing
	Confirm Status = "confirm"
)

// Builtin
const (
	InvalidLabel         = ""
	BuiltinTransferLabel = "__transfer"
	InvalidUID           = "0x0"
)

// Result define the filed for describe tx invoke result.
type Result struct {
	// Label is the group name of sent transaction
	Label string `mapstructure:"label"`
	// UID is the unique id of sent transaction
	UID interface{} `mapstructure:"uid"`
	// BuildTime is the client start time when client constructs transaction
	BuildTime int64 `mapstructure:"build"`
	// SendTime is the client time when client receives transaction response
	SendTime int64 `mapstructure:"send"`
	// ConfirmTime is the client time when client check transaction is written in block stably
	ConfirmTime int64 `mapstructure:"confirm"`
	// WriteTime is the node time when transaction is written in block stably
	WriteTime int64 `mapstructure:"write"`
	// Status marks the status of transaction
	Status Status `mapstructure:"status"`
	// Ret is the return value of transaction
	Ret []interface{} `mapstructure:"ret"`
}
