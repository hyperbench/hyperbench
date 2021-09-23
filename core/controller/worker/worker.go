package worker

import (
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/core/collector"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

// Worker is the interface of worker node
type Worker interface {
	// SetContext set the context of worker passed from Master.
	SetContext([]byte) error

	// Do call the workers to running.
	Do() error

	// CheckoutCollector checkout collector.
	CheckoutCollector() (collector.Collector, bool)

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

	urls := viper.GetStringSlice(common.EngineURLsPath)
	if len(urls) == 0 {
		var localWorkerConfig LocalWorkerConfig
		localWorkerConfig.Cap = viper.GetInt64(common.EngineCapPath)
		localWorker, err := NewLocalWorker(LocalWorkerConfig{
			Index:    0,
			Cap:      viper.GetInt64(common.EngineCapPath),
			Rate:     viper.GetInt64(common.EngineRatePath),
			Duration: viper.GetDuration(common.EngineDurationPath),
		})
		if err != nil {
			return nil, ErrConfig
		}
		workers = []Worker{
			localWorker,
		}
	} else {
		// create archive for sync benchmark context
		p := viper.GetString(common.BenchmarkDirPath)

		target := p + ".tar.gz"
		err := archiver.Archive([]string{p}, target)
		if err != nil {
			return nil, errors.Wrapf(err, "can not archive: %v", target)
		}
		viper.Set(common.BenchmarkArchivePath, target)
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
