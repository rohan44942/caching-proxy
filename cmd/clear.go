package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear-cache",
	Short: "Clear the cache",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cache cleared!")
		// Phase 4: actual cache clearing logic will go here
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
