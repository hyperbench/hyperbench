package recorder

import (
	"encoding/csv"
	"fmt"
	"github.com/meshplus/hyperbench/common"
	"github.com/meshplus/hyperbench/core/utils"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

// Recorder define the service a recorder need provide.
type Recorder interface {
	// Process process input report.
	Process(input common.Report)
	// Release source.
	Release()
	processor
}

type processor interface {
	process(report common.Report)
	release()
}

// NewRecorder create recoder with config in config.toml.
func NewRecorder() Recorder {
	var ps []processor

	logger := common.GetLogger("recd")
	ps = append(ps, newLogProcessor(logger))

	// csv
	if viper.IsSet(common.RecorderCsvPath) {
		dirPath := viper.GetString(common.RecorderCsvDirPath)
		if dirPath == "" {
			dirPath = "./csv"
		}

		t := time.Now()
		fileName := fmt.Sprintf("%s.csv", t.Format("2006-12-12-23:59:59"))
		filePath := filepath.Join(dirPath, fileName)
		csvFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
		csvPath = filePath
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
func (b *baseRecorder) Process(input common.Report) {
	b.process(input)
}

func (b *baseRecorder) process(report common.Report) {
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

func (p *logProcessor) process(report common.Report) {
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

func (p *logProcessor) logData(t string, data *common.Data) {
	for _, d := range data.Results {
		du := float64(d.Duration) / float64(time.Second)
		p.logger.Noticef("%s\t%d\t%v\t|\t%.1f\t%.1f\t%.1f\t%.1f\t%.1f\t|\t%.1f\t%.1f\t%.1f",
			t,
			d.Num,
			int(du),
			float64(d.Num)/du,
			float64(d.Statuses[common.Success])/du,
			float64(d.Statuses[common.Failure])/du,
			float64(d.Statuses[common.Confirm])/du,
			float64(d.Statuses[common.Unknown])/du,
			float64(d.Send.Avg)/float64(time.Millisecond),
			float64(d.Confirm.Avg)/float64(time.Millisecond),
			float64(d.Write.Avg)/float64(time.Millisecond))
	}
}

func (p *logProcessor) release() {
	return
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

func (p *csvProcessor) process(report common.Report) {
	p.logData(report.Cur)
	p.logData(report.Sum)
}

func (p *csvProcessor) release() {
	_ = p.f.Close()
}

func (p *csvProcessor) logData(data *common.Data) {
	for _, d := range data.Results {
		_ = p.writer.Write(utils.AggData2CSV(nil, data.Type, d))
	}
}
