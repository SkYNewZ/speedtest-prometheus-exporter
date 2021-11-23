package cmd

import (
	"fmt"
	"os/exec"

	"github.com/SkYNewZ/speedtest-prometheus-exporter/internal/speedtest"
	"github.com/spf13/cobra"
)

var (
	speedtestPath string
	serverURL     = "http://localhost:3100"
)

// speedtestCmd represents the speedtest command
var speedtestCmd = &cobra.Command{
	Use:   "speedtest",
	Short: "Run speedtest and report to given server",
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		if speedtestPath == "" {
			return fmt.Errorf("missing --speedtest-path")
		}

		if serverURL == "" {
			return fmt.Errorf("missing --server")
		}

		// Check if speedtest path is correct
		if _, err := exec.Command(speedtestPath, "--version").Output(); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Run the speedtest
		results, err := speedtest.Run(cmd.Context(), speedtestPath, log)
		if err != nil {
			return err
		}

		// Push results to server
		return speedtest.Push(cmd.Context(), log, &speedtest.PushConfig{
			URL:     serverURL,
			Results: results,
		})
	},
}

func init() {
	rootCmd.AddCommand(speedtestCmd)
	speedtestCmd.Flags().StringVar(&speedtestPath, "speedtest-path", "/usr/bin/speedtest", "Speedtest cli path")
	speedtestCmd.Flags().StringVar(&serverURL, "server", serverURL, "Server to push results to")
}
