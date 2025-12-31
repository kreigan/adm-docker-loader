package cmd

import (
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:     "restart [stack-name]",
	Aliases: []string{"reload"},
	Short:   "Restart Docker Compose stacks",
	Long:    `Restart all Docker Compose stacks or a specific stack if provided.`,
	Args:    cobra.MaximumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return RunAction("restart", args)
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
