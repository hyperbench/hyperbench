package common

// Report contains need report data about tx been executed info.
type Report struct {
	Cur *Data
	Sum *Data
}

// DataType data type.
type DataType string

const (
	// Cur current data type.
	Cur DataType = "Cur"
	// Sum sum data type.
	Sum DataType = "Sum"
)

// Data define the data which will be reported.
type Data struct {
	// Type data type.
	Type DataType
	// Results the slice of AggData.
	Results []AggData
}

// AggData aggregation data.
type AggData struct {
	// key
	Label    string
	Time     int64
	Duration int64

	// request
	Num      int
	Statuses map[Status]int

	// latency
	Send    Latency
	Confirm Latency
	Write   Latency
}

// Latency is the percent of latency in increase order
type Latency struct {
	Avg  int64 // Avg is the average latency
	P0   int64 // P0  is the minimal latency
	P50  int64 // P50 is the median of latency
	P90  int64 // P90 is the latency exactly larger than 90%
	P95  int64 // P90 is the latency exactly larger than 95%
	P99  int64 // P99 is the latency exactly larger than 95%
	P100 int64 // P100 is the maximal latency
}

// RemoteStatistic remote statistic data.
type RemoteStatistic struct {
	// time info
	Start int64
	End   int64

	// block info
	BlockNum int

	// transaction info
	TxNum int
}
