package cmd

import "github.com/spf13/cobra"

var pubMode string
var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start auction WebSocket server",
	Run: func(cmd *cobra.Command, args []string) {
		switch pubMode {
		case "inmemory":
			startInMemoryServer(port)
		case "redis":
			startRedisServer(port)
		default:
			startInMemoryServer(port)
		}
	},
}

func init() {
	serveCmd.Flags().StringVar(&pubMode, "pub", "inmemory", "inmemory | redis")
	serveCmd.Flags().StringVar(&port, "port", "8080", "HTTP port")
	rootCmd.AddCommand(serveCmd)
}
