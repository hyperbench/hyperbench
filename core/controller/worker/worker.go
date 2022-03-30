package worker

import (
	fcom "github.com/hyperbench/hyperbench-common/common"

	"os"

	"github.com/hyperbench/hyperbench/core/collector"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Worker is the interface of worker node
type Worker interface {
	// SetContext set the context of worker passed from Master.
	SetContext([]byte) error

	// BeforeRun call user hook
	BeforeRun() error

	// Do call the workers to running.
	Do() error

	// AfterRun call user hook
	AfterRun() error

	// CheckoutCollector checkout collector.
	CheckoutCollector() (collector.Collector, bool, error)

	// Statistic get the number of sent and missed transactions
	Statistics() (int64, int64)

	// Teardown close the worker manually.
	Teardown()
}

var (
	// ErrConfig config error.
	ErrConfig = errors.New("config error")
)

// NewWorkers generate workers according to config
func NewWorkers() (workers []Worker, err error) {
	defer func() {
		if err != nil {
			for _, w := range workers {
				w.Teardown()
			}
			workers = nil
		}
	}()

	urls := viper.GetStringSlice(fcom.EngineURLsPath)
	if len(urls) == 0 {
		var localWorkerConfig LocalWorkerConfig
		localWorkerConfig.Cap = viper.GetInt64(fcom.EngineCapPath)
		localWorker, err := NewLocalWorker(LocalWorkerConfig{
			Index:    0,
			Cap:      viper.GetInt64(fcom.EngineCapPath),
			Rate:     viper.GetInt64(fcom.EngineRatePath),
			Duration: viper.GetDuration(fcom.EngineDurationPath),
		})
		if err != nil {
			return nil, ErrConfig
		}
		workers = []Worker{
			localWorker,
		}
	} else {
		// create archive for sync benchmark context
		p := viper.GetString(fcom.BenchmarkDirPath)

		target := p + ".tar.gz"
		os.RemoveAll(target)
		err := archiver.Archive([]string{p}, target)
		if err != nil {
			return nil, errors.Wrapf(err, "can not archive: %v", target)
		}
		viper.Set(fcom.BenchmarkArchivePath, target)
		// remove archiver
		// nolint
		defer os.RemoveAll(target)

		workers = make([]Worker, 0, len(urls))
		for idx, urls := range urls {
			worker, err := NewRemoteWorker(idx, urls)
			if err != nil {
				return workers, err
			}
			workers = append(workers, worker)
		}
	}
	return
}
