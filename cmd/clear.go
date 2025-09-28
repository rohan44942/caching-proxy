package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/rohan44942/caching-proxy/internal/cache"
)

var clearCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clear the cache",
	Run: func(cmd *cobra.Command, args []string) {
		if cache.GlobalCache != nil {
			cache.GlobalCache.Clear()
			fmt.Println("Cache cleared successfully!")
		} else {
			fmt.Println("Cache not initialized (start server first).")
		}
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
