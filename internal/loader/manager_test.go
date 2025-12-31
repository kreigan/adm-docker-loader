package loader

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckDuplicates(t *testing.T) {
	t.Run("no duplicates returns nil", func(t *testing.T) {
		stacks := []*Stack{
			{Name: "web", Dir: "/stacks/01-web"},
			{Name: "database", Dir: "/stacks/02-database"},
		}
		if err := CheckDuplicates(stacks); err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("single duplicate returns error", func(t *testing.T) {
		stacks := []*Stack{
			{Name: "web", Dir: "/stacks/01-web"},
			{Name: "web", Dir: "/stacks/03-web"},
		}
		err := CheckDuplicates(stacks)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "duplicate") {
			t.Errorf("Expected 'duplicate' in error, got: %s", err.Error())
		}
	})

	t.Run("empty list returns nil", func(t *testing.T) {
		if err := CheckDuplicates([]*Stack{}); err != nil {
			t.Errorf("Expected no error for empty list, got: %v", err)
		}
	})
}

func TestWarnDuplicates(_ *testing.T) {
	// Just verify it doesn't panic
	WarnDuplicates([]*Stack{{Name: "web", Dir: "/01-web"}, {Name: "web", Dir: "/02-web"}})
	WarnDuplicates([]*Stack{})
}

func TestStackManagerExecuteAction(t *testing.T) {
	t.Run("invalid action returns error", func(t *testing.T) {
		dir := t.TempDir()
		config := &Config{}
		manager := NewStackManager(dir, config, newTestLogger(t), true) // dry-run

		err := manager.ExecuteAction("invalid", "")
		if err == nil {
			t.Error("Expected error for invalid action")
		}
		if !strings.Contains(err.Error(), "unrecognized") {
			t.Errorf("Expected 'unrecognized' in error, got: %s", err.Error())
		}
	})

	t.Run("start action executes", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		config := &Config{UpArgs: []string{"--detach"}}
		manager := NewStackManager(dir, config, newTestLogger(t), true) // dry-run

		err := manager.ExecuteAction("start", "web")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("stop action executes", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		config := &Config{Timeout: 10}
		manager := NewStackManager(dir, config, newTestLogger(t), true) // dry-run

		err := manager.ExecuteAction("stop", "web")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("down action executes", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		config := &Config{}
		manager := NewStackManager(dir, config, newTestLogger(t), true)

		err := manager.ExecuteAction("down", "web")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("restart action executes", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		config := &Config{Timeout: 10}
		manager := NewStackManager(dir, config, newTestLogger(t), true)

		err := manager.ExecuteAction("restart", "web")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("list action executes", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		config := &Config{}
		manager := NewStackManager(dir, config, newTestLogger(t), false)

		err := manager.ExecuteAction("list", "")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("unknown target stack returns error", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		config := &Config{}
		manager := NewStackManager(dir, config, newTestLogger(t), true)

		err := manager.ExecuteAction("start", "unknown")
		if err == nil {
			t.Error("Expected error for unknown stack")
		}
	})

	t.Run("fails on duplicate stacks", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))
		mustMkdir(t, filepath.Join(stacksDir, "02-web"))

		config := &Config{}
		manager := NewStackManager(dir, config, newTestLogger(t), true)

		err := manager.ExecuteAction("start", "")
		if err == nil {
			t.Error("Expected error for duplicate stacks")
		}
		if !strings.Contains(err.Error(), "duplicate") {
			t.Errorf("Expected 'duplicate' in error, got: %s", err.Error())
		}
	})

	t.Run("no stacks logs warning", func(t *testing.T) {
		dir := t.TempDir()
		mustMkdir(t, filepath.Join(dir, "stacks")) // Empty stacks dir

		config := &Config{}
		manager := NewStackManager(dir, config, newTestLogger(t), true)

		err := manager.ExecuteAction("start", "")
		if err != nil {
			t.Errorf("Expected nil for empty stacks, got: %v", err)
		}
	})
}
