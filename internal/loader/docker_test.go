package loader

import (
	"errors"
	"testing"
)

func TestComposeClientHasContainers(t *testing.T) {
	stack := &Stack{Name: "web", Dir: "/stacks/01-web"}

	t.Run("returns true when containers exist", func(t *testing.T) {
		mock := &MockDockerExecutor{RunQuietOut: []byte("abc123\n")}
		client := NewComposeClient(mock, newTestLogger(t), &Config{})

		if !client.HasContainers(stack) {
			t.Error("Expected true when containers exist")
		}
		if len(mock.RunQuietCalls) != 1 {
			t.Errorf("Expected 1 call, got %d", len(mock.RunQuietCalls))
		}
	})

	t.Run("returns false when no containers", func(t *testing.T) {
		mock := &MockDockerExecutor{RunQuietOut: []byte("")}
		client := NewComposeClient(mock, newTestLogger(t), &Config{})

		if client.HasContainers(stack) {
			t.Error("Expected false when no containers")
		}
	})

	t.Run("returns false on error", func(t *testing.T) {
		mock := &MockDockerExecutor{RunQuietError: errors.New("docker error")}
		client := NewComposeClient(mock, newTestLogger(t), &Config{})

		if client.HasContainers(stack) {
			t.Error("Expected false on error")
		}
	})
}

func TestComposeClientGetProjectStatuses(t *testing.T) {
	t.Run("parses json output", func(t *testing.T) {
		mock := &MockDockerExecutor{
			RunQuietOut: []byte(`[{"Name":"web","Status":"running(1)"},{"Name":"db","Status":"exited"}]`),
		}
		client := NewComposeClient(mock, newTestLogger(t), &Config{})

		statuses := client.GetProjectStatuses()
		if statuses["web"] != StackStatusRunning {
			t.Errorf("Expected web=running, got %s", statuses["web"])
		}
		if statuses["db"] != StackStatusStopped {
			t.Errorf("Expected db=stopped, got %s", statuses["db"])
		}
	})

	t.Run("returns empty map on error", func(t *testing.T) {
		mock := &MockDockerExecutor{RunQuietError: errors.New("error")}
		client := NewComposeClient(mock, newTestLogger(t), &Config{})

		statuses := client.GetProjectStatuses()
		if len(statuses) != 0 {
			t.Error("Expected empty map on error")
		}
	})
}

func TestComposeClientOperations(t *testing.T) {
	stack := &Stack{Name: "web", Dir: "/stacks/01-web"}
	config := &Config{Timeout: 10}
	stackConfig := &StackConfig{}

	tests := []struct {
		name      string
		operation func(*ComposeClient) error
		wantCmd   string
	}{
		{"Up", func(c *ComposeClient) error { return c.Up(stack, stackConfig) }, "up"},
		{"Start", func(c *ComposeClient) error { return c.Start(stack, stackConfig) }, "start"},
		{"Stop", func(c *ComposeClient) error { return c.Stop(stack, stackConfig) }, "stop"},
		{"Down", func(c *ComposeClient) error { return c.Down(stack, stackConfig) }, "down"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockDockerExecutor{}
			client := NewComposeClient(mock, newTestLogger(t), config)

			if err := tt.operation(client); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(mock.RunCalls) != 1 {
				t.Fatalf("Expected 1 call, got %d", len(mock.RunCalls))
			}
			if !sliceContains(mock.RunCalls[0], tt.wantCmd) {
				t.Errorf("Expected %q in args, got %v", tt.wantCmd, mock.RunCalls[0])
			}
		})
	}
}

func TestDefaultDockerExecutor(t *testing.T) {
	t.Run("dry run does not execute", func(t *testing.T) {
		executor := NewDockerExecutor(newTestLogger(t), true)
		err := executor.Run([]string{"version"})
		if err != nil {
			t.Errorf("Dry run should not error: %v", err)
		}
	})

	t.Run("RunQuiet returns nil in dry run", func(t *testing.T) {
		executor := NewDockerExecutor(newTestLogger(t), true)
		output, err := executor.RunQuiet([]string{"version"})
		if err != nil {
			t.Errorf("Dry run should not error: %v", err)
		}
		if output != nil {
			t.Errorf("Dry run should return nil output, got: %v", output)
		}
	})
}
