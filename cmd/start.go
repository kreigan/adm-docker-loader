package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [stack-name]",
	Short: "Start Docker Compose stacks",
	Long:  `Start all Docker Compose stacks or a specific stack if provided.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return RunAction("start", args)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
