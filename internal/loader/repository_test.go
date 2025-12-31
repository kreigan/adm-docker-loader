package loader

import (
	"path/filepath"
	"strings"
	"testing"
)

const testStackName = "web" // Common test stack name

func TestStackRepositoryFindAll(t *testing.T) {
	t.Run("discovers stacks", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")

		// Create stack directories
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))
		mustMkdir(t, filepath.Join(stacksDir, "02-database"))
		mustMkdir(t, filepath.Join(stacksDir, "invalid")) // Should be skipped

		mock := &MockDockerExecutor{RunQuietOut: []byte("[]")}
		compose := NewComposeClient(mock, newTestLogger(t), &Config{})
		repo := NewStackRepository(dir, newTestLogger(t), compose)

		stacks, err := repo.FindAll()
		if err != nil {
			t.Fatalf("FindAll failed: %v", err)
		}

		if len(stacks) != 2 {
			t.Errorf("Expected 2 stacks, got %d", len(stacks))
		}
		if stacks[0].Name != testStackName {
			t.Errorf("Expected first stack %q, got %q", testStackName, stacks[0].Name)
		}
		if stacks[1].Name != "database" {
			t.Errorf("Expected second stack 'database', got %q", stacks[1].Name)
		}
	})

	t.Run("returns error for missing directory", func(t *testing.T) {
		dir := t.TempDir()
		mock := &MockDockerExecutor{}
		compose := NewComposeClient(mock, newTestLogger(t), &Config{})
		repo := NewStackRepository(dir, newTestLogger(t), compose)

		_, err := repo.FindAll()
		if err == nil {
			t.Error("Expected error for missing stacks directory")
		}
	})

	t.Run("sorts by directory name", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")

		mustMkdir(t, filepath.Join(stacksDir, "20-second"))
		mustMkdir(t, filepath.Join(stacksDir, "10-first"))
		mustMkdir(t, filepath.Join(stacksDir, "30-third"))

		mock := &MockDockerExecutor{RunQuietOut: []byte("[]")}
		compose := NewComposeClient(mock, newTestLogger(t), &Config{})
		repo := NewStackRepository(dir, newTestLogger(t), compose)

		stacks, err := repo.FindAll()
		if err != nil {
			t.Fatalf("FindAll failed: %v", err)
		}
		if stacks[0].Name != "first" || stacks[1].Name != "second" || stacks[2].Name != "third" {
			t.Errorf("Stacks not sorted correctly: %v", stacks)
		}
	})
}

func TestStackRepositoryFindByName(t *testing.T) {
	t.Run("finds by stack name", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		mock := &MockDockerExecutor{RunQuietOut: []byte("[]")}
		compose := NewComposeClient(mock, newTestLogger(t), &Config{})
		repo := NewStackRepository(dir, newTestLogger(t), compose)

		stack, err := repo.FindByName(testStackName)
		if err != nil {
			t.Fatalf("FindByName failed: %v", err)
		}
		if stack.Name != testStackName {
			t.Errorf("Expected %q, got %q", testStackName, stack.Name)
		}
	})

	t.Run("finds by directory name", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		mock := &MockDockerExecutor{RunQuietOut: []byte("[]")}
		compose := NewComposeClient(mock, newTestLogger(t), &Config{})
		repo := NewStackRepository(dir, newTestLogger(t), compose)

		stack, err := repo.FindByName("01-web")
		if err != nil {
			t.Fatalf("FindByName failed: %v", err)
		}
		if stack.Name != testStackName {
			t.Errorf("Expected %q, got %q", testStackName, stack.Name)
		}
	})

	t.Run("returns error for unknown stack", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		mock := &MockDockerExecutor{RunQuietOut: []byte("[]")}
		compose := NewComposeClient(mock, newTestLogger(t), &Config{})
		repo := NewStackRepository(dir, newTestLogger(t), compose)

		_, err := repo.FindByName("unknown")
		if err == nil {
			t.Error("Expected error for unknown stack")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' in error, got: %s", err.Error())
		}
	})
}

func TestStackRepositoryStatus(t *testing.T) {
	t.Run("populates status from compose", func(t *testing.T) {
		dir := t.TempDir()
		stacksDir := filepath.Join(dir, "stacks")
		mustMkdir(t, filepath.Join(stacksDir, "01-web"))

		mock := &MockDockerExecutor{
			RunQuietOut: []byte(`[{"Name":"web","Status":"running(1)"}]`),
		}
		compose := NewComposeClient(mock, newTestLogger(t), &Config{})
		repo := NewStackRepository(dir, newTestLogger(t), compose)

		stacks, err := repo.FindAll()
		if err != nil {
			t.Fatalf("FindAll failed: %v", err)
		}
		if stacks[0].Status != StackStatusRunning {
			t.Errorf("Expected running status, got %s", stacks[0].Status)
		}
	})
}
