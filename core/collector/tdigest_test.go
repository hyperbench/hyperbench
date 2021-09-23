package collector

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/tdigest"
	"github.com/meshplus/hyperbench/common"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
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

func Test(t *testing.T) {
	td := tdigest.NewWithCompression(1000)
	for _, x := range []float64{1, 2, 3, 4, 5, 5, 4, 3, 2, 1} {
		td.Add(x, 1)
	}

	// Compute Quantiles
	fmt.Println("50th", td.Quantile(0.5))
	fmt.Println("75th", td.Quantile(0.75))
	fmt.Println("90th", td.Quantile(0.9))
	fmt.Println("99th", td.Quantile(0.99))

	// Compute CDFs
	fmt.Println("CDF(1) = ", td.CDF(1))
	fmt.Println("CDF(2) = ", td.CDF(2))
	fmt.Println("CDF(3) = ", td.CDF(3))
	fmt.Println("CDF(4) = ", td.CDF(4))
	fmt.Println("CDF(5) = ", td.CDF(5))
}
