package recorder

/**
 *  Copyright (C) 2021 HyperBench.
 *  SPDX-License-Identifier: Apache-2.0
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 * @brief Recorder define the service a recorder need provide
 * @file recorder.go
 * @author: Mingmei Liu
 * @date 2021-01-12
 */

import (
	"encoding/csv"
	"os"
	"path"
	"time"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/hyperbench/hyperbench/core/utils"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// Recorder define the service a recorder need provide.
type Recorder interface {
	// Process process input report.
	Process(input fcom.Report)
	ProcessStatistic(rs *fcom.RemoteStatistic)
	// Release source.
	Release()
	processor
}

type processor interface {
	process(report fcom.Report)
	processStatistic(rs *fcom.RemoteStatistic)
	release()
}

// NewRecorder create recoder with config in config.toml.
func NewRecorder() Recorder {
	var ps []processor

	logger := fcom.GetLogger("recd")
	ps = append(ps, newLogProcessor(logger))

	// csv
	if viper.IsSet(fcom.RecorderCsvPath) {
		dirPath := viper.GetString(fcom.RecorderCsvDirPath)
		if dirPath == "" {
			dirPath = "./csv"
		}
		_, err := os.Stat(dirPath)
		if err != nil && !os.IsExist(err) {
			err := os.MkdirAll(dirPath, 0777)
			if err != nil {
				logger.Errorf("make csv dirpath error: %v", err)
			}
		}
		fileName := path.Join(dirPath, time.Now().Format("2006-01-02-15:04:05")+".csv")
		csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			logger.Errorf("open file error: %v", err)
		}
		csvPath = fileName
		if err == nil {
			ps = append(ps, newCSVProcessor(csvFile))
		}
	}

	// influxdb
	if viper.IsSet(fcom.RecorderInflucDBPath) {
		benchmark := viper.GetString(fcom.BenchmarkDirPath)
		url := viper.GetString(fcom.InfluxDBUrlPath)
		db := viper.GetString(fcom.InfluxDBDatabasePath)
		uname := viper.GetString(fcom.InfluxDBUsernamePath)
		pwd := viper.GetString(fcom.InfluxDBPasswordPath)
		influxDB, err := newInfluxdb(benchmark, url, db, uname, pwd)
		if err == nil {
			ps = append(ps, influxDB)
		} else {
			logger.Errorf("int influxdb client error: %v", err)
		}
		if err == nil {
			ps = append(ps, influxDB)
		} else {
			logger.Errorf("int influxdb client error: %v", err)
		}
	}

	return &baseRecorder{
		ps: ps,
	}
}

// Release source.
func (b *baseRecorder) Release() {
	b.release()
}

func (b *baseRecorder) release() {
	for _, r := range b.ps {
		r.release()
	}
}

var (
	csvPath = ""
)

// GetCSVPath return csv path.
func GetCSVPath() string {
	return csvPath
}

type baseRecorder struct {
	ps []processor
}

// Process process input report.
func (b *baseRecorder) Process(input fcom.Report) {
	b.process(input)
}

func (b *baseRecorder) process(report fcom.Report) {
	for _, p := range b.ps {
		p.process(report)
	}
}

// ProcessStatistic process statistic.
func (b *baseRecorder) ProcessStatistic(rs *fcom.RemoteStatistic) {
	b.processStatistic(rs)
}

func (b *baseRecorder) processStatistic(rs *fcom.RemoteStatistic) {
	for _, p := range b.ps {
		p.processStatistic(rs)
	}
}

type logProcessor struct {
	logger *logging.Logger
}

func newLogProcessor(logger *logging.Logger) *logProcessor {
	return &logProcessor{
		logger: logger,
	}
}

func (p *logProcessor) process(report fcom.Report) {
	p.logger.Notice("")
	p.logTitle()
	p.logData("Cur  ", report.Cur)
	p.logData("Sum  ", report.Sum)
	p.logger.Notice("")
}

func (p *logProcessor) processStatistic(sd *fcom.RemoteStatistic) {
	p.logger.Notice("")
	p.logger.Notice("\t\tSent\t\tMissed\t\tTotal\t\tTps")
	p.logger.Noticef("\t\t%v\t\t%v\t\t%v\t\t%.1f", sd.SentTx, sd.MissedTx, sd.SentTx+sd.MissedTx, sd.Tps)
	p.logger.Notice("")
	p.logger.Notice("       From        \t         To           \tBlk\tTx\tCTps\tBps")
	p.logger.Noticef("%s\t%s\t%v\t%v\t%.1f\t%.1f",
		time.Unix(0, sd.Start).Format("2006-01-02 15:04:05"),
		time.Unix(0, sd.End).Format("2006-01-02 15:04:05"),
		sd.BlockNum,
		sd.TxNum,
		sd.CTps,
		sd.Bps,
	)
	p.logger.Notice("")

	//p.logger.Noticef("viper all settings:%v", settings)
	//p.logger.Noticef("viper:%v",viper.GetViper())
}

func (p *logProcessor) logTitle() {

	p.logger.Notice("     \tview\t    \t|\t    \t    \trate\t(/s)\t    \t|\t\tlatency\t(ms)")
	p.logger.Notice("state\tnum \tdu(s)\t|\tsend\tsucc\tfail\tconf\tunkn\t|\tsend\tconf\twrit")
}

func (p *logProcessor) logData(t string, data *fcom.Data) {
	for _, d := range data.Results {
		du := float64(d.Duration) / float64(time.Second)
		p.logger.Noticef("%s\t%d\t%v\t|\t%.1f\t%.1f\t%.1f\t%.1f\t%.1f\t|\t%.1f\t%.1f\t%.1f",
			t,
			d.Num,
			int(du),
			float64(d.Num)/du,
			float64(d.Statuses[fcom.Success])/du,
			float64(d.Statuses[fcom.Failure])/du,
			float64(d.Statuses[fcom.Confirm])/du,
			float64(d.Statuses[fcom.Unknown])/du,
			float64(d.Send.Avg)/float64(time.Millisecond),
			float64(d.Confirm.Avg)/float64(time.Millisecond),
			float64(d.Write.Avg)/float64(time.Millisecond))
	}
}

func (p *logProcessor) release() {
}

type csvProcessor struct {
	writer *csv.Writer
	f      *os.File
}

func newCSVProcessor(f *os.File) *csvProcessor {
	return &csvProcessor{
		writer: csv.NewWriter(f),
		f:      f,
	}
}

func (p *csvProcessor) process(report fcom.Report) {
	p.logData(report.Cur)
	p.logData(report.Sum)
}

func (p *csvProcessor) processStatistic(rs *fcom.RemoteStatistic) {
	_ = p.writer.Write(utils.RemoteStatistic2CSV(nil, rs))
}

func (p *csvProcessor) release() {
	_ = p.f.Close()
}

func (p *csvProcessor) logData(data *fcom.Data) {
	for _, d := range data.Results {
		_ = p.writer.Write(utils.AggData2CSV(nil, data.Type, d))
	}
}
