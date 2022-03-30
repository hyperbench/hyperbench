package utils

import (
	fcom "github.com/hyperbench/hyperbench-common/common"

	"strconv"
)

// AggData2CSV append data's fields to base,
func AggData2CSV(base []string, t fcom.DataType, data fcom.AggData) []string {
	if base == nil {
		base = make([]string, 0, 30)
	}
	base = append(base,
		string(t),
		data.Label,
		i2s(data.Time),
		i2s(data.Duration),
		i2s(data.Num),
		i2s(data.Statuses[fcom.Failure]),
		i2s(data.Statuses[fcom.Success]),
		i2s(data.Statuses[fcom.Confirm]),
		i2s(data.Statuses[fcom.Unknown]))
	base = Latency2CSV(base, data.Send)
	base = Latency2CSV(base, data.Confirm)
	base = Latency2CSV(base, data.Write)
	return base
}

// Latency2CSV append latency's fields to base
func Latency2CSV(base []string, latency fcom.Latency) []string {
	if base == nil {
		base = make([]string, 0, 7)
	}
	base = append(base,
		i2s(latency.Avg),
		i2s(latency.P0),
		i2s(latency.P50),
		i2s(latency.P90),
		i2s(latency.P95),
		i2s(latency.P99),
		i2s(latency.P100))
	return base
}

func i2s(i interface{}) (s string) {
	switch v := i.(type) {
	case int64:
		return strconv.Itoa(int(v))
	case int32:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	}
	return ""
}
