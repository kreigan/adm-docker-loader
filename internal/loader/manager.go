package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

// StackManager manages Docker Compose stacks.
type StackManager struct {
	repo    *StackRepository
	compose *ComposeClient
	logger  *Logger
}

// NewStackManager creates a new stack manager.
func NewStackManager(baseDir string, config *Config, logger *Logger, dryRun bool) *StackManager {
	executor := NewDockerExecutor(logger, dryRun)
	compose := NewComposeClient(executor, logger, config)
	repo := NewStackRepository(baseDir, logger, compose)

	return &StackManager{
		repo:    repo,
		compose: compose,
		logger:  logger,
	}
}

// ExecuteAction executes the specified action on stacks.
func (m *StackManager) ExecuteAction(action, targetStack string) error {
	act := Action(action)
	if !act.IsValid() {
		return fmt.Errorf("unrecognized action: %s", action)
	}

	stacks, err := m.getStacks(targetStack)
	if err != nil {
		return err
	}

	if len(stacks) == 0 {
		m.logger.Warning("No stacks found")
		return nil
	}

	return m.performAction(act, stacks)
}

func (m *StackManager) getStacks(targetStack string) ([]*Stack, error) {
	if targetStack != "" {
		stack, err := m.repo.FindByName(targetStack)
		if err != nil {
			return nil, err
		}
		return []*Stack{stack}, nil
	}

	stacks, err := m.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("discovering stacks: %w", err)
	}

	return stacks, nil
}

func (m *StackManager) performAction(action Action, stacks []*Stack) error {
	switch action {
	case ActionList:
		return m.listStacks(stacks)
	case ActionStart:
		return m.executeWithDuplicateCheck(stacks, m.startStack)
	case ActionStop:
		return m.executeWithDuplicateCheck(stacks, m.stopStack)
	case ActionDown:
		return m.executeWithDuplicateCheck(stacks, m.downStack)
	case ActionRestart, ActionReload:
		return m.executeWithDuplicateCheck(stacks, m.restartStack)
	default:
		return fmt.Errorf("unrecognized action: %s", action)
	}
}

func (m *StackManager) executeWithDuplicateCheck(stacks []*Stack, fn func(*Stack) error) error {
	if err := CheckDuplicates(stacks); err != nil {
		return err
	}

	for _, stack := range stacks {
		if err := fn(stack); err != nil {
			return err
		}
	}

	return nil
}

func (m *StackManager) listStacks(stacks []*Stack) error {
	WarnDuplicates(stacks)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ORDER\tSTACK\tSTATUS\tPATH")
	fmt.Fprintln(w, "-----\t-----\t------\t----")

	for _, stack := range stacks {
		dirName := filepath.Base(stack.Dir)
		order := strings.SplitN(dirName, "-", 2)[0]
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", order, stack.Name, stack.Status, stack.Dir)
	}

	//nolint:errcheck // Flush error is non-critical for display purposes
	w.Flush()
	return nil
}

func (m *StackManager) startStack(stack *Stack) error {
	m.logger.Console("==> Starting stack: %s", stack.Name)
	m.logger.Info("Starting stack: %s", stack.Name)

	stackConfig, err := LoadStackConfig(stack.Dir)
	if err != nil {
		m.logger.Warning("Failed to load stack config: %v", err)
		stackConfig = &StackConfig{}
	}

	if m.compose.HasContainers(stack) {
		m.logger.Debug("Containers exist for stack %s, using 'start'", stack.Name)
		return m.compose.Start(stack, stackConfig)
	}

	m.logger.Debug("No containers found for stack %s, using 'up'", stack.Name)
	return m.compose.Up(stack, stackConfig)
}

func (m *StackManager) stopStack(stack *Stack) error {
	m.logger.Console("==> Stopping stack: %s", stack.Name)
	m.logger.Info("Stopping stack: %s", stack.Name)

	stackConfig, err := LoadStackConfig(stack.Dir)
	if err != nil {
		m.logger.Warning("Failed to load stack config: %v", err)
		stackConfig = &StackConfig{}
	}

	return m.compose.Stop(stack, stackConfig)
}

func (m *StackManager) downStack(stack *Stack) error {
	m.logger.Console("==> Taking down stack: %s", stack.Name)
	m.logger.Info("Taking down stack: %s", stack.Name)

	stackConfig, err := LoadStackConfig(stack.Dir)
	if err != nil {
		m.logger.Warning("Failed to load stack config: %v", err)
		stackConfig = &StackConfig{}
	}

	return m.compose.Down(stack, stackConfig)
}

func (m *StackManager) restartStack(stack *Stack) error {
	if err := m.stopStack(stack); err != nil {
		return err
	}
	return m.startStack(stack)
}
