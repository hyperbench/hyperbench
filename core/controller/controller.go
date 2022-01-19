package controller

import (
	"context"

	fcom "github.com/meshplus/hyperbench-common/common"
	"github.com/spf13/viper"

	"sync"
	"time"

	"github.com/meshplus/hyperbench/core/collector"
	"github.com/meshplus/hyperbench/core/controller/master"
	"github.com/meshplus/hyperbench/core/controller/worker"
	"github.com/meshplus/hyperbench/core/recorder"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

// Controller is the controller of job
type Controller interface {
	// Prepare prepare
	Prepare() error
	Run() error
}

type WorkerClient struct {
	worker   worker.Worker
	finished bool // check worker finied work
}

// ControllerImpl is the implement of Controller
type ControllerImpl struct {
	master        master.Master
	workerClients []*WorkerClient
	recorder      recorder.Recorder
	reportChan    chan fcom.Report
	curCollector  collector.Collector
	sumCollector  collector.Collector
	logger        *logging.Logger
	start         int64
	end           int64
}

// NewController create Controller.
func NewController() (Controller, error) {

	m, err := master.NewLocalMaster()
	if err != nil {
		return nil, errors.Wrap(err, "can not create master")
	}

	ws, err := worker.NewWorkers()
	if err != nil {
		return nil, errors.Wrap(err, "can not create workers")
	}
	var workerClients []*WorkerClient
	for i := 0; i < len(ws); i++ {
		workerClients = append(workerClients, &WorkerClient{
			ws[i],
			false,
		})
	}

	r := recorder.NewRecorder()

	return &ControllerImpl{
		master:        m,
		workerClients: workerClients,
		//finishedWorker:syncMap,
		recorder:     r,
		logger:       fcom.GetLogger("ctrl"),
		curCollector: collector.NewTDigestSummaryCollector(),
		sumCollector: collector.NewTDigestSummaryCollector(),
		reportChan:   make(chan fcom.Report),
	}, nil
}

// Prepare prepare for job
func (l *ControllerImpl) Prepare() (err error) {

	defer func() {
		// if preparation is failed, then just teardown all workers
		// to avoid that worker is occupied
		if err != nil {
			l.teardownWorkers()
		}
	}()

	l.logger.Notice("ready to prepare")
	err = l.master.Prepare()
	if err != nil {
		return errors.Wrap(err, "master is not ready")
	}

	l.logger.Notice("ready to get context")
	bsCtx, err := l.master.GetContext()
	if err != nil {
		return errors.Wrap(err, "can not get context")
	}

	l.logger.Noticef("ctx: %s", string(bsCtx))
	l.logger.Notice("ready to set context")
	// must ensure all workers ready
	for _, w := range l.workerClients {
		err = w.worker.SetContext(bsCtx)
		if err != nil {
			return errors.Wrap(err, "can not set context")
		}
	}

	return nil
}

// Run start the job
func (l *ControllerImpl) Run() (err error) {
	defer l.teardownWorkers()

	// run all workers
	duration := viper.GetDuration(fcom.EngineDurationPath)
	l.start = time.Now().UnixNano()
	tick := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-tick.C:
				l.end, err = l.master.LogStatus()
				if err != nil {
					l.logger.Error(err)
				}
				return
			}
		}
	}()
	for _, w := range l.workerClients {
		// nolint
		go w.worker.Do()
	}

	// get response
	go l.asyncGetAllResponse()

	for report := range l.reportChan {
		l.recorder.Process(report)
	}

	l.recorder.Release()
	sd, err := l.master.Statistic(l.start, l.end)
	if err != nil {
		l.logger.Notice(err)
	}
	if err == nil {
		l.logStatisticData(sd)
	}

	l.logger.Notice("finish")
	return nil
}

func (l *ControllerImpl) logStatisticData(sd *fcom.RemoteStatistic) {

	l.logger.Notice("")
	l.logger.Notice("       From        \t         To           \tBlk\tTx\tTps\tBps")
	l.logger.Noticef("%s\t%s\t%v\t%v\t%.1f\t%.1f",
		time.Unix(0, sd.Start).Format("2006-01-02 15:04:05"),
		time.Unix(0, sd.End).Format("2006-01-02 15:04:05"),
		sd.BlockNum,
		sd.TxNum,
		sd.Tps,
		sd.Bps,
	)
	l.logger.Notice("")
}

func (l *ControllerImpl) asyncGetAllResponse() {

	workerNum := len(l.workerClients)

	output := make(chan collector.Collector, workerNum)
	close(output)

	time.Sleep(200 * time.Millisecond)

	l.curCollector.Reset()
	l.sumCollector.Reset()
	tick := time.NewTicker(10 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	var finishWg sync.WaitGroup
	finishWg.Add(workerNum)

	go func() {
		finishWg.Wait()
		l.logger.Notice("cancel")
		cancel()
	}()

	for {
		select {
		case <-tick.C:
			// get process all value
			var wg sync.WaitGroup
			wg.Add(workerNum)
			output = make(chan collector.Collector, workerNum)
			for _, w := range l.workerClients {
				go l.getWorkerResponse(w, &wg, &finishWg, output)
			}
			wg.Wait()
			//l.logger.Notice("====got")
			close(output)
			for col := range output {
				_ = l.curCollector.MergeC(col)
				_ = l.sumCollector.MergeC(col)
			}
			l.report()

		case <-ctx.Done():
			//l.logger.Notice("====ctx.done")
			close(l.reportChan)
			return
		}
	}
}

func (l *ControllerImpl) report() {
	report := fcom.Report{
		Cur: l.curCollector.Get(),
		Sum: l.sumCollector.Get(),
	}
	l.reportChan <- report
	l.curCollector.Reset()
}

func (l *ControllerImpl) getWorkerResponse(w *WorkerClient, batchWg *sync.WaitGroup, finishWg *sync.WaitGroup, output chan collector.Collector) {
	if w.finished {
		batchWg.Done()
		return
	}

	col, valid, err := w.worker.CheckoutCollector()
	if err != nil {
		l.logger.Error(err)
		batchWg.Done()
		return
	}
	if !valid {
		w.finished = true
		l.logger.Notice("finishWg done")
		finishWg.Done()
		batchWg.Done()
		return
	}
	output <- col
	batchWg.Done()
}

func (l *ControllerImpl) teardownWorkers() {
	for _, w := range l.workerClients {
		w.worker.Teardown()
	}
}
