package loader

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// MockDockerExecutor implements DockerExecutor for testing.
//
//nolint:govet // Field order is logically grouped for readability
type MockDockerExecutor struct {
	RunCalls      [][]string
	RunQuietCalls [][]string
	RunQuietOut   []byte
	RunError      error
	RunQuietError error
}

func (m *MockDockerExecutor) Run(args []string) error {
	m.RunCalls = append(m.RunCalls, args)
	return m.RunError
}

func (m *MockDockerExecutor) RunQuiet(args []string) ([]byte, error) {
	m.RunQuietCalls = append(m.RunQuietCalls, args)
	return m.RunQuietOut, m.RunQuietError
}

// =============================================================================
// Test helpers
// =============================================================================

func newTestLogger(t *testing.T) *Logger {
	t.Helper()
	logger, err := NewLogger(filepath.Join(t.TempDir(), "test.log"), false)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	t.Cleanup(func() {
		//nolint:errcheck // Test cleanup, error not critical
		logger.Close()
	})
	return logger
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to write %s: %v", name, err)
	}
}

func assertSliceEqual(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("Slice length mismatch: got %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("Slice[%d] mismatch: got %q, want %q", i, got[i], want[i])
		}
	}
}

func sliceContains(slice []string, target string) bool {
	return slices.Contains(slice, target)
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", path, err)
	}
}
