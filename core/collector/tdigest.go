package collector

import (
	"errors"
	"time"

	json "github.com/json-iterator/go"

	"github.com/influxdata/tdigest"
	fcom "github.com/meshplus/hyperbench-common/common"
)

const (
	w float64 = 1
)

// TDigest defined for use tdigest.TDigest.
type TDigest struct {
	*tdigest.TDigest
}

// NewTDigest create TDigest and return.
func NewTDigest() *TDigest {
	td := TDigest{}
	td.TDigest = tdigest.New()
	return &td
}

func (t *TDigest) getLatency() fcom.Latency {
	return fcom.Latency{
		Avg:  t.avg(),
		P0:   int64(t.Quantile(0)),
		P50:  int64(t.Quantile(0.5)),
		P90:  int64(t.Quantile(0.9)),
		P95:  int64(t.Quantile(0.95)),
		P99:  int64(t.Quantile(0.99)),
		P100: int64(t.Quantile(1)),
	}
}

func (t *TDigest) avg() int64 {
	cs := t.TDigest.Centroids()
	sum := 0.0
	count := 0.0
	for _, c := range cs {
		sum += c.Mean * c.Weight
		count += c.Weight
	}
	if count == 0 {
		return 0
	}
	return int64(sum / count)
}

func (t *TDigest) merge(src *TDigest) {
	t.AddCentroidList(src.Centroids())
}

// UnmarshalJSON implement Unmarshaler.
func (t *TDigest) UnmarshalJSON(bytes []byte) error {
	var list tdigest.CentroidList
	_ = json.Unmarshal(bytes, &list)
	t.TDigest = tdigest.New()
	t.AddCentroidList(list)
	return nil
}

// MarshalJSON implement Marshaler.
func (t *TDigest) MarshalJSON() ([]byte, error) {
	list := t.Centroids()
	return json.Marshal(list)
}

// TDigestDetailsCollector is the implement of Collector processing result in detail
type TDigestDetailsCollector struct {
	Data map[string]*Details
	Time int64
}

// Details detail info for statistic.
type Details struct {
	Label          string
	Num            int
	SendLatency    *TDigest
	ConfirmLatency *TDigest
	WriteLatency   *TDigest
	Status         map[fcom.Status]int
}

// NewDetails create a Details and return.
func NewDetails(label string) *Details {
	return &Details{
		Label:          label,
		Num:            0,
		SendLatency:    NewTDigest(),
		ConfirmLatency: NewTDigest(),
		WriteLatency:   NewTDigest(),
		Status:         make(map[fcom.Status]int),
	}
}

func (d *Details) merge(src *Details) {
	if src.Label != d.Label {
		return
	}
	d.Num += src.Num
	d.SendLatency.merge(src.SendLatency)
	d.ConfirmLatency.merge(src.ConfirmLatency)
	d.WriteLatency.merge(src.WriteLatency)

	for status, num := range src.Status {
		d.Status[status] += num
	}
}

func (d *Details) add(result *fcom.Result) {
	if result == nil || result.BuildTime == 0 {
		return
	}
	d.Num++
	d.Status[result.Status]++
	if result.SendTime != 0 {
		d.SendLatency.Add(float64(result.SendTime-result.BuildTime), w)
	}
	if result.ConfirmTime != 0 {
		d.ConfirmLatency.Add(float64(result.ConfirmTime-result.BuildTime), w)
	}
	if result.WriteTime != 0 {
		d.ConfirmLatency.Add(float64(result.WriteTime-result.BuildTime), w)
	}
}

// NewTDigestDetailsCollector new TDigestDetailsCollector and return.
func NewTDigestDetailsCollector() Collector {
	return newTDigestDetailsCollector()
}

func newTDigestDetailsCollector() *TDigestDetailsCollector {
	return &TDigestDetailsCollector{
		Data: make(map[string]*Details),
	}
}

// Add append result to statistic
func (t *TDigestDetailsCollector) Add(result *fcom.Result) {
	var (
		cur   *Details
		exist bool
	)
	if cur, exist = t.Data[result.Label]; !exist {
		cur = NewDetails(result.Label)
		t.Data[result.Label] = cur
	}
	cur.add(result)
}

// Type return the types of collector
func (t *TDigestDetailsCollector) Type() string {
	return "details"
}

// Serialize generate serialized data to pass through network in remote mode
func (t *TDigestDetailsCollector) Serialize() []byte {
	bs, _ := json.Marshal(t)
	return bs
}

// Merge merge serialized data
func (t *TDigestDetailsCollector) Merge(bs []byte) (err error) {

	var tmp = newTDigestDetailsCollector()
	err = json.Unmarshal(bs, tmp)
	if err != nil {
		return
	}
	t.merge(tmp)
	return nil
}

// MergeC try to merge a Collector, if it can not do this, just raise a error
func (t *TDigestDetailsCollector) MergeC(collector Collector) (err error) {
	if tmp, ok := collector.(*TDigestDetailsCollector); ok {
		t.merge(tmp)
	}
	return errors.New("error type of collector")
}

// Get get current statistic data group by label
func (t *TDigestDetailsCollector) Get() *fcom.Data {
	data := &fcom.Data{
		Results: make([]fcom.AggData, 0, len(t.Data)),
	}
	now := time.Now().UnixNano()
	duration := now - t.Time
	for _, v := range t.Data {
		r := fcom.AggData{
			Label:    v.Label,
			Time:     t.Time,
			Duration: duration,
			Num:      v.Num,
			Statuses: v.Status,
			Send:     v.SendLatency.getLatency(),
			Confirm:  v.ConfirmLatency.getLatency(),
			Write:    v.ConfirmLatency.getLatency(),
		}
		data.Results = append(data.Results, r)
	}

	return data
}

// Reset reset data should reset the time window and clean data.
func (t *TDigestDetailsCollector) Reset() {
	t.Data = make(map[string]*Details)
	t.Time = time.Now().UnixNano()
}

func (t *TDigestDetailsCollector) merge(src *TDigestDetailsCollector) {
	for k, srcDetails := range src.Data {
		details, exist := t.Data[k]
		if !exist {
			details = NewDetails(k)
			t.Data[k] = details
		}
		details.merge(srcDetails)
	}
}

// TDigestSummaryCollector is the implement of Collector processing result in summary
type TDigestSummaryCollector struct {
	Data *Details
	Time int64
}

// NewTDigestSummaryCollector create TDigestSummaryCollector.
func NewTDigestSummaryCollector() Collector {
	return &TDigestSummaryCollector{
		Data: NewDetails(""),
		Time: time.Now().UnixNano(),
	}
}

// Add append result to statistic
func (t *TDigestSummaryCollector) Add(result *fcom.Result) {
	t.Data.add(result)
}

// Serialize generate serialized data to pass through network in remote mode
func (t *TDigestSummaryCollector) Serialize() []byte {
	bytes, _ := json.Marshal(t)
	return bytes
}

// Merge merge serialized data
func (t *TDigestSummaryCollector) Merge(bytes []byte) (err error) {
	var tmp TDigestSummaryCollector
	err = json.Unmarshal(bytes, &tmp)
	if err != nil {
		return
	}
	t.Data.merge(tmp.Data)
	return
}

// MergeC try to merge a Collector, if it can not do this, just raise a error
func (t *TDigestSummaryCollector) MergeC(collector Collector) (err error) {
	if tmp, ok := collector.(*TDigestSummaryCollector); ok {
		t.Data.merge(tmp.Data)
	}
	return errors.New("error type of collector")
}

// Get get current statistic data group by label
func (t *TDigestSummaryCollector) Get() *fcom.Data {
	data := &fcom.Data{
		Results: make([]fcom.AggData, 0, 1),
	}
	now := time.Now().UnixNano()
	duration := now - t.Time
	v := t.Data
	r := fcom.AggData{
		Label:    v.Label,
		Time:     t.Time,
		Duration: duration,
		Num:      v.Num,
		Statuses: v.Status,
		Send:     v.SendLatency.getLatency(),
		Confirm:  v.ConfirmLatency.getLatency(),
		Write:    v.ConfirmLatency.getLatency(),
	}
	data.Results = append(data.Results, r)
	return data
}

// Reset reset data should reset the time window and clean data.
func (t *TDigestSummaryCollector) Reset() {
	t.Data = NewDetails("")
	t.Time = time.Now().UnixNano()
}

// Type return the types of collector
func (t *TDigestSummaryCollector) Type() string {
	return "summary"
}

// NewTDigestCollectorBuilder create Collector by mode.
func NewTDigestCollectorBuilder(mode string) func() Collector {
	switch mode {
	case "details":
		return NewTDigestDetailsCollector
	case "summary":
		return NewTDigestSummaryCollector
	default:
		return NewTDigestSummaryCollector
	}
}
