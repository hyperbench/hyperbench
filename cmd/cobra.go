//Package cmd use cobra provide cmd function
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	fcom "github.com/hyperbench/hyperbench-common/common"

	"github.com/hyperbench/hyperbench/core/controller"
	"github.com/hyperbench/hyperbench/core/network/server"
	"github.com/hyperbench/hyperbench/filesystem"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var (
	logger   *logging.Logger
	debugF   func()
	branch   string
	commitID string
	date     string
)

//InitCmd init cobra command
func InitCmd(debug func()) error {
	logger = fcom.GetLogger("cmd")
	debugF = debug

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
	}

	addRun()

	rootCmd.AddCommand(start,
		create,
		initDir,
		version,
		worker)
	return nil
}

// GetRootCmd return the root cobra.Command
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func addRun() {

	start.Run = func(cmd *cobra.Command, args []string) {
		if args[0] == "" {
			logger.Info("Specify test plan first")
			return
		}
		runBenchmark(args[0])
	}

	create.Run = func(cmd *cobra.Command, args []string) {
		// todo: initialize a test plan
	}

	initDir.Run = func(cmd *cobra.Command, args []string) {
		if err := filesystem.Unpack(""); err != nil {
			panic("init stress test dir error")
		}
	}

	version.Run = func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hyperbench Version:\n%s-%s-%s\n", branch, date, commitID)
	}

	worker.Run = func(cmd *cobra.Command, args []string) {
		fcom.InitLog()
		svr := server.NewServer(*port)
		err := svr.Start()
		if err != nil {
			fmt.Printf("can not start server: %v", err)
		}
	}

}

func initConfig() {
	if *enableDebug {
		fmt.Println("========DEBUG MODE========")
		debugF()
		time.Sleep(5 * time.Second)
	}

	if *document != "" {
		err := doc.GenMarkdownTree(rootCmd, *document)
		if err != nil {
			logger.Error("create doc error. ", err)
		}
	}
}

func initBenchmark(dir string) {
	// support for specifying configuration files or specifying test case directories
	// if use test case directory, it will use file named "config" for configuration file
	s, err := os.Stat(dir)
	if err != nil {
		logger.Critical("can not find test case: %v", err)
		return
	}
	isFile := !s.IsDir()
	var path string
	if isFile {
		viper.SetConfigFile(dir)
		path = dir[0 : strings.LastIndex(dir, "/")+1]
		viper.Set(fcom.BenchmarkConfigPath, dir)
	} else {
		path = dir
	}
	viper.AddConfigPath(path)

	err = viper.ReadInConfig()
	if err != nil {
		logger.Critical("can not read in config: %v", err)
		return
	}
	viper.Set(fcom.BenchmarkDirPath, path)

	fcom.InitLog()
}

func runBenchmark(dir string) {
	initBenchmark(dir)

	ctrl, err := controller.NewController()
	if err != nil {
		logger.Criticalf("can not create controller: %v", err)
		return
	}

	err = ctrl.Prepare()
	if err != nil {
		logger.Criticalf("prepare error: %v", err)
		return
	}

	err = ctrl.Run()
	if err != nil {
		logger.Criticalf("prepare error: %v", err)
		return
	}
}
