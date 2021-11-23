package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	log     *logrus.Logger
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "speedtest-prometheus-exporter",
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		if verbose {
			log.SetLevel(logrus.DebugLevel)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", verbose, "Verbose output")
}

func initConfig() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
}
