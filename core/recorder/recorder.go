package recorder

import (
	"encoding/csv"
	"os"
	"path"
	"time"

	fcom "github.com/meshplus/hyperbench-common/common"

	"github.com/meshplus/hyperbench/core/utils"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

// Recorder define the service a recorder need provide.
type Recorder interface {
	// Process process input report.
	Process(input fcom.Report)
	// Release source.
	Release()
	processor
}

type processor interface {
	process(report fcom.Report)
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

func (p *csvProcessor) release() {
	_ = p.f.Close()
}

func (p *csvProcessor) logData(data *fcom.Data) {
	for _, d := range data.Results {
		_ = p.writer.Write(utils.AggData2CSV(nil, data.Type, d))
	}
}
