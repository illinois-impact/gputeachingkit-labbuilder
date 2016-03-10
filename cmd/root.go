package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// This represents the base command when called without any subcommands
var (
	debug   bool
	RootCmd = &cobra.Command{
		Use:   "wgx-pandoc",
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
	buildCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output.")
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
