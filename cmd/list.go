// Package cmd provides CLI commands for composectl.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all discovered stacks",
	Long:  `List all Docker Compose stacks discovered in the stacks directory.`,
	Args:  cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		if IsDryRun() {
			return fmt.Errorf("--dry-run flag is not applicable for list command")
		}
		return RunAction("list", nil)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
