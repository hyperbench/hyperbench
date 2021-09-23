package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

// Debug related config keys
const (
	PprofRecordDuration = "5s"
	PprofTimeFmt        = "2006-01-02-15-04-05"
)

func debug() {

	duration, err := time.ParseDuration(PprofRecordDuration)
	if err != nil {
		return
	}

	go recordPProf(duration)
}

func recordPProf(duration time.Duration) {
	var (
		cpuFilePath string
		memFilePath string
		cpuFile     *os.File
		memFile     *os.File
	)

	dir := "./debug"
	err := ensurePathExists(dir)
	if err != nil {
		return
	}

	// reset inline function
	reset := func() {
		// use the expected file closed time as the file name's suffix
		timeSuffix := time.Now().Add(duration).Format(PprofTimeFmt)

		cpuFilePath = fmt.Sprint(dir, "/cpu_", timeSuffix)
		memFilePath = fmt.Sprint(dir, "/mem_", timeSuffix)

		cpuFile, _ = os.Create(cpuFilePath)
		memFile, _ = os.Create(memFilePath)

		// start pprof
		err := pprof.StartCPUProfile(cpuFile)
		if err != nil {
			return
		}
	}

	tick := time.NewTicker(duration)

	reset()

	//nolint
	for {
		select {
		case <-tick.C:

			pprof.StopCPUProfile()
			err := cpuFile.Close()
			if err != nil {
				continue
			}
			err = pprof.WriteHeapProfile(memFile)
			if err != nil {
				continue
			}
			err = memFile.Close()
			if err != nil {
				continue
			}

			reset()
		}
	}

}

func ensurePathExists(path string) error {

	_, err := os.Stat(path)
	// already exist
	if err == nil {
		return nil
	}

	// not exist
	if os.IsNotExist(err) {
		// make full path
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}

	return err
}
