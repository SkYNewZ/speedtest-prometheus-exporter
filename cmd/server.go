package cmd

import (
	"github.com/SkYNewZ/speedtest-prometheus-exporter/internal/server"

	"github.com/spf13/cobra"
)

var (
	port uint16 = 3100
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Receive speedtest results and expose them as Prometheus metrics",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run(cmd.Context(), port, log)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().Uint16Var(&port, "port", port, "Listen port for the server")
}
