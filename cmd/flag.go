package cmd

import (
	"github.com/spf13/cobra"
)

var (
	start = &cobra.Command{
		Use:     "start",
		Short:   "start a benchmark",
		Args:    cobra.ExactArgs(1),
		Example: "hyperbench start benchmark/transfer",
	}

	create = &cobra.Command{
		Use:     "new [testplan Name]",
		Short:   "initialize a test plan",
		Args:    cobra.ExactArgs(1),
		Example: "hyperbench new myTest",
	}

	rootCmd = &cobra.Command{
		Use:     "hyperbench",
		Example: "hyperbench --doc ./doc (generate document to specify path)",
	}

	version = &cobra.Command{
		Use:     "version",
		Short:   "get code version",
		Example: "hyperbench version",
	}

	initDir = &cobra.Command{
		Use:     "init",
		Short:   "init a stress test dir",
		Example: "hyperbench init",
	}

	worker = &cobra.Command{
		Use:     "worker",
		Short:   "start as a worker server ",
		Example: "hyperbench worker",
	}

	enableDebug = rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	document    = rootCmd.PersistentFlags().String("doc", "", "use to create doc and specify the doc path")
	port        = worker.PersistentFlags().IntP("port", "p", 8080, "port of worker")
)
