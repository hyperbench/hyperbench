package collector

import "github.com/meshplus/hyperbench/common"

// Collector is used to collect result and generate statistic data group by label
// Collector may not be implement concurrently safe, so you should receive data in a goroutine
type Collector interface {
	// Type return the types of collector
	Type() string

	// Add append result to statistic
	Add(*common.Result)

	// Serialize generate serialized data to pass through network in remote mode
	Serialize() []byte

	// Merge merge serialized data
	Merge([]byte) error

	// MergeC try to merge a Collector, if it can not do this, just raise a error
	MergeC(Collector) error

	// Get get current statistic data group by label
	Get() *common.Data

	// Reset reset data should reset the time window and clean data
	Reset()
}
