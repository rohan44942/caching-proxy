package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "caching-proxy",
    Short: "A caching reverse proxy server",
    Long:  `caching-proxy forwards requests to an origin server and caches responses.`,
}

// Execute runs the root command
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func init() {
    // child commands are registered here
}
