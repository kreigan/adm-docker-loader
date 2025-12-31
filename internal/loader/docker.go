package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DockerExecutor handles Docker command execution.
type DockerExecutor interface {
	Run(args []string) error
	RunQuiet(args []string) ([]byte, error)
}

// DefaultDockerExecutor executes Docker commands on the host.
type DefaultDockerExecutor struct {
	logger *Logger
	dryRun bool
}

// NewDockerExecutor creates a new Docker executor.
func NewDockerExecutor(logger *Logger, dryRun bool) *DefaultDockerExecutor {
	return &DefaultDockerExecutor{
		logger: logger,
		dryRun: dryRun,
	}
}

// Run executes a Docker command with TTY passthrough.
func (e *DefaultDockerExecutor) Run(args []string) error {
	if e.dryRun {
		e.logger.Info("[DRY-RUN] Would execute: docker %s", strings.Join(args, " "))
		return nil
	}

	e.logger.Debug("Executing: docker %s", strings.Join(args, " "))

	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker command failed: %w", err)
	}

	return nil
}

// RunQuiet executes a Docker command and returns output without TTY.
func (e *DefaultDockerExecutor) RunQuiet(args []string) ([]byte, error) {
	if e.dryRun {
		return nil, nil
	}

	cmd := exec.Command("docker", args...)
	return cmd.Output()
}

// ComposeClient handles Docker Compose operations for a stack.
type ComposeClient struct {
	executor DockerExecutor
	logger   *Logger
	config   *Config
}

// NewComposeClient creates a new Compose client.
func NewComposeClient(executor DockerExecutor, logger *Logger, config *Config) *ComposeClient {
	return &ComposeClient{
		executor: executor,
		logger:   logger,
		config:   config,
	}
}

// Up brings up a stack.
func (c *ComposeClient) Up(stack *Stack, stackConfig *StackConfig) error {
	config := c.config.MergeStackConfig(stackConfig)
	args := c.buildArgs(stack, config, "up")
	args = append(args, config.UpArgs...)
	return c.executor.Run(args)
}

// Start starts a stopped stack.
func (c *ComposeClient) Start(stack *Stack, stackConfig *StackConfig) error {
	config := c.config.MergeStackConfig(stackConfig)
	args := c.buildArgs(stack, config, "start")
	return c.executor.Run(args)
}

// Stop stops a stack.
func (c *ComposeClient) Stop(stack *Stack, stackConfig *StackConfig) error {
	config := c.config.MergeStackConfig(stackConfig)
	args := c.buildArgs(stack, config, "stop")
	args = append(args, "--timeout", fmt.Sprintf("%d", c.config.Timeout))
	return c.executor.Run(args)
}

// Down takes down a stack.
func (c *ComposeClient) Down(stack *Stack, stackConfig *StackConfig) error {
	config := c.config.MergeStackConfig(stackConfig)
	args := c.buildArgs(stack, config, "down")
	args = append(args, config.DownArgs...)
	return c.executor.Run(args)
}

// HasContainers checks if a stack has any containers.
func (c *ComposeClient) HasContainers(stack *Stack) bool {
	args := []string{
		"compose",
		"--project-directory", stack.Dir,
		"--project-name", stack.Name,
		"ps", "-a", "-q",
	}

	output, err := c.executor.RunQuiet(args)
	if err != nil {
		c.logger.Debug("Failed to check containers for stack %s: %v", stack.Name, err)
		return false
	}

	return strings.TrimSpace(string(output)) != ""
}

// GetProjectStatuses returns status for all Docker Compose projects.
func (c *ComposeClient) GetProjectStatuses() map[string]StackStatus {
	statuses := make(map[string]StackStatus)

	args := []string{"compose", "ls", "--all", "--format", "json"}
	output, err := c.executor.RunQuiet(args)
	if err != nil {
		c.logger.Debug("Failed to get compose statuses: %v", err)
		return statuses
	}

	var projects []struct {
		Name   string `json:"Name"`
		Status string `json:"Status"`
	}

	if err := json.Unmarshal(output, &projects); err != nil {
		c.logger.Debug("Failed to parse compose statuses: %v", err)
		return statuses
	}

	for _, project := range projects {
		statuses[project.Name] = normalizeStatus(project.Status)
	}

	return statuses
}

func (c *ComposeClient) buildArgs(stack *Stack, config *Config, action string) []string {
	args := []string{"compose"}
	args = append(args, config.CommonArgs...)
	args = append(args, "--project-directory", stack.Dir, "--project-name", stack.Name, action)
	return args
}

// normalizeStatus converts docker compose status to StackStatus.
func normalizeStatus(status string) StackStatus {
	// Remove count in parentheses: "running(1)" -> "running"
	if idx := strings.Index(status, "("); idx != -1 {
		status = status[:idx]
	}

	switch status {
	case "running":
		return StackStatusRunning
	case "exited", "stopped":
		return StackStatusStopped
	default:
		return StackStatusDown
	}
}
