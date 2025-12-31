package cmd

import (
	"fmt"
	"os"

	"github.com/kreigan/adm-composectl/internal/loader"
)

// ActionRunner handles the common execution flow for stack actions.
type ActionRunner struct {
	action string
}

// NewActionRunner creates a new action runner.
func NewActionRunner(action string) *ActionRunner {
	return &ActionRunner{action: action}
}

// Run executes the action with the given target stack.
func (r *ActionRunner) Run(targetStack string) error {
	logger, err := loader.NewLogger(GetLogFile(), IsVerbose())
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		if err := logger.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()

	logger.Info("Docker Loader started - action: %s", r.action)
	logger.Info("Base directory: %s", GetBaseDir())

	config, err := loader.LoadConfig(GetBaseDir())
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	manager := loader.NewStackManager(GetBaseDir(), config, logger, IsDryRun())

	if err := manager.ExecuteAction(r.action, targetStack); err != nil {
		return fmt.Errorf("%s action failed: %w", r.action, err)
	}

	logger.Info("Docker Loader finished successfully")
	return nil
}

// RunAction is a convenience function for executing an action.
func RunAction(action string, args []string) error {
	targetStack := ""
	if len(args) > 0 {
		targetStack = args[0]
	}

	return NewActionRunner(action).Run(targetStack)
}
