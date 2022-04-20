//Copyright 2021 Mingmei Liu
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

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
			Instant:  viper.GetInt64("engine.instant"),
			Wait:     viper.GetDuration("engine.wait"),
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
