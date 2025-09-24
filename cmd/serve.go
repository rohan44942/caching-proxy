package cmd

import (
    "log"

    "github.com/spf13/cobra"
    "github.com/rohan44942/caching-proxy/internal/config"
    "github.com/rohan44942/caching-proxy/internal/server"
)

var port int
var origin string

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the caching proxy server",
    Run: func(cmd *cobra.Command, args []string) {
        cfg := config.Config{
            Port:   port,
            Origin: origin,
        }

        if err := server.Start(cfg); err != nil {
            log.Fatalf("server error: %v", err)
        }
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)

    serveCmd.Flags().IntVarP(&port, "port", "p", 3000, "Port to run the proxy server")
    serveCmd.Flags().StringVarP(&origin, "origin", "o", "", "Origin server URL")
    serveCmd.MarkFlagRequired("origin")
}
