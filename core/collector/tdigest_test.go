package collector

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/influxdata/tdigest"
	"github.com/meshplus/hyperbench/common"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTDigest_Add(b *testing.B) {

	//
	td := tdigest.New()
	rand.Seed(time.Now().Unix())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//n := rand.Int63n(1000)
		td.Add(float64(i), 1)
	}
	b.StopTimer()

}

func BenchmarkTDigestCollector_Add(b *testing.B) {

	// 145ns
	td := NewTDigestDetailsCollector()
	rand.Seed(time.Now().Unix())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//n := rand.Int63n(1000)
		td.Add(&common.Result{
			Label:     "test",
			BuildTime: 1,
			SendTime:  int64(i),
			//Latency: n,
			Status: common.Success,
			Ret:    nil,
		})
	}
	b.StopTimer()
}

func BenchmarkTDigestCollector_Serialize(b *testing.B) {

	// 1587ns
	td := NewTDigestDetailsCollector()
	rand.Seed(time.Now().Unix())
	n := 100000
	for i := 0; i < n; i++ {
		r := rand.Int63n(1000)
		td.Add(&common.Result{
			Label:     "test",
			BuildTime: 0,
			SendTime:  r,
			Status:    common.Success,
			Ret:       nil,
		})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		td.Serialize()
	}
	b.StopTimer()
}

func BenchmarkTDigestCollector_Merge(b *testing.B) {

	// 89us
	td := NewTDigestDetailsCollector()
	rand.Seed(time.Now().Unix())
	n := 100000
	for i := 0; i < n; i++ {
		r := rand.Int63n(1000)
		td.Add(&common.Result{
			Label:     "test",
			BuildTime: 0,
			SendTime:  r,
			Status:    common.Success,
			Ret:       nil,
		})
	}
	bs := td.Serialize()

	ntd := NewTDigestDetailsCollector()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ntd.Merge(bs)
	}
	b.StopTimer()
}

func TestTDigest_MarshalJSON(t *testing.T) {
	ast := assert.New(t)
	var td = NewTDigest()
	n := 10000
	sum := 0
	for i := 1; i <= n; i++ {
		r := rand.Intn(i)
		td.Add(float64(r), 1)
		sum += r
	}
	bytes, err := json.Marshal(td)
	ast.NoError(err)
	ast.NotEqual(0, len(bytes))
	ast.Equal(sum/n, int(td.avg()))
	ast.Equal(n, int(td.Count()))

	var newTd = NewTDigest()
	err = json.Unmarshal(bytes, &newTd)
	ast.NoError(err)

	ast.Equal(td.Centroids(), newTd.Centroids())
}

func TestTDigestDetailsCollector(t *testing.T) {
	col := NewTDigestCollectorBuilder("details")
	col1 := col()
	col2 := col()
	assert.Equal(t, col1.Type(), "details")
	assert.Equal(t, col2.Type(), "details")

	res1 := &common.Result{}
	res2 := &common.Result{
		BuildTime:   time.Now().UnixNano(),
		SendTime:    time.Now().UnixNano(),
		ConfirmTime: time.Now().UnixNano(),
		WriteTime:   time.Now().UnixNano(),
	}

	col1.Add(res1)
	assert.Equal(t, len(col1.Get().Results), 1)
	col2.Add(res2)
	assert.Equal(t, len(col1.Get().Results), 1)

	var bs []byte
	col2.Merge(bs)
	assert.Equal(t, len(col1.Get().Results), 1)

	col1.Add(res2)
	bs = col1.Serialize()
	assert.NotNil(t, bs)
	col2.Merge(bs)
	assert.Equal(t, len(col1.Get().Results), 1)

	col1.MergeC(col2)
	assert.Equal(t, len(col1.Get().Results), 1)

	col1.Reset()
	assert.Equal(t, len(col1.Get().Results), 0)
	col3 := col()
	col3.MergeC(col2)
	assert.Equal(t, len(col3.Get().Results), 1)

}

func TestTDigestSummaryCollector(t *testing.T) {
	co1 := NewTDigestCollectorBuilder("summary")
	co2 := NewTDigestCollectorBuilder("")

	col1 := co1()
	col2 := co2()
	assert.Equal(t, col1.Type(), "summary")
	assert.Equal(t, col2.Type(), "summary")

	res1 := &common.Result{}
	res2 := &common.Result{
		BuildTime: time.Now().UnixNano(),
	}

	col1.Add(res1)
	assert.Equal(t, col1.(*TDigestSummaryCollector).Data.Num, 0)
	col2.Add(res2)
	assert.Equal(t, col2.(*TDigestSummaryCollector).Data.Num, 1)

	var bs []byte
	col2.Merge(bs)
	assert.Equal(t, col2.(*TDigestSummaryCollector).Data.Num, 1)

	col1.Add(res2)
	bs = col1.Serialize()
	assert.NotNil(t, bs)
	col2.Merge(bs)
	assert.Equal(t, col2.(*TDigestSummaryCollector).Data.Num, 2)

	col1.MergeC(col2)
	assert.Equal(t, col1.(*TDigestSummaryCollector).Data.Num, 3)

	col1.Reset()
	assert.Equal(t, col1.(*TDigestSummaryCollector).Data.Num, 0)
	assert.Equal(t, len(col1.Get().Results), 1)

}
