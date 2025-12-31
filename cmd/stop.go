package cmd

import (
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [stack-name]",
	Short: "Stop Docker Compose stacks",
	Long:  `Stop all Docker Compose stacks or a specific stack if provided.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return RunAction("stop", args)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
