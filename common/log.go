package common

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	logFormatter = `%{color}[%{module}][%{level:.5s}] %{time:15:04:05.000} %{shortfile} %{message} %{color:reset}`
	loggers      = make(map[string]*logging.Logger)
	loggerLock   sync.RWMutex
	backend      logging.LeveledBackend
)

//InitLog init log
func InitLog() string {

	level, err := logging.LogLevel(strings.ToUpper(viper.GetString(LogLevelPath)))
	if err != nil {
		level = logging.NOTICE
	}
	dir := viper.GetString(LogDirPath)
	if dir == "" {
		dir = "./log"
	}
	dump := viper.GetBool(LogDumpPath)
	if dump {
		_, err = os.Stat(dir)
		if err != nil && !os.IsExist(err) {
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				return ""
			}
		}
	}

	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stderrStringFormatter := logging.MustStringFormatter(logFormatter)
	backendFormatter := logging.NewBackendFormatter(stderrBackend, stderrStringFormatter)

	backendStderr := logging.AddModuleLevel(backendFormatter)
	backendStderr.SetLevel(level, "")

	var fileName string

	if !dump {
		backend = logging.MultiLogger(backendStderr)
	} else {
		fileName = path.Join(dir, time.Now().Format("2006-01-02-15:04:05")+".log")
		logFile, err := os.Create(fileName)
		if err != nil {
			log.Fatalln("open file error !")
		}
		logBackend := logging.NewLogBackend(logFile, "", 0)
		var stringFormatter = logging.MustStringFormatter(logFormatter)
		backendFileFormatter := logging.NewBackendFormatter(logBackend, stringFormatter)
		backendFile := logging.AddModuleLevel(backendFileFormatter)
		backendFile.SetLevel(level, "")
		backend = logging.MultiLogger(backendStderr, backendFile)
	}

	loggerLock.Lock()
	for _, logger := range loggers {
		logger.SetBackend(backend)
	}
	loggerLock.Unlock()

	return fileName
}

//GetLogger get logging.Logger with module
func GetLogger(module string) *logging.Logger {
	var logger *logging.Logger

	loggerLock.Lock()
	if loggers[module] != nil {
		logger = loggers[module]
	} else {
		logger = logging.MustGetLogger(module)
		if backend != nil {
			logger.SetBackend(backend)
		}
		loggers[module] = logger
	}
	loggerLock.Unlock()
	return logger
}
