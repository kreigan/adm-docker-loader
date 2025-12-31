package cmd

import (
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down [stack-name]",
	Short: "Take down Docker Compose stacks",
	Long: `Take down all Docker Compose stacks or a specific stack if provided. 
This removes containers, networks, and volumes.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return RunAction("down", args)
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
