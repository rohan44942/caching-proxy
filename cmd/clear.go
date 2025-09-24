package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clear the cache",
	Run: func(cmd *cobra.Command, args []string) {
		// Phase 4: hook into cache.Clear()
		fmt.Println("Cache cleared! (stub, will implement in Phase 4)")
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
