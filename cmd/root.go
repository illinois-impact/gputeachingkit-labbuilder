package cmd

import (
	"fmt"
	"os"

	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
var (
	verbose string
	RootCmd = &cobra.Command{
		Use:   "gputeachingkit-labbuilder",
		Short: "Helper tools to build lab documentations and interact with pandoc",
	}
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	buildCmd.PersistentFlags().StringVarP(&verbose, "verbose", "v", "fatal", "Choose verbosity level [debug,info,warn,error].")
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	lvl := log.FatalLevel
	switch strings.ToLower(verbose) {
	case "debug":
		lvl = log.DebugLevel
	case "info":
		lvl = log.InfoLevel
	case "warn":
		lvl = log.WarnLevel
	case "error":
		lvl = log.ErrorLevel
	case "fatal":
		lvl = log.FatalLevel
	default:
		lvl = log.FatalLevel
	}
	log.SetLevel(lvl)
}
