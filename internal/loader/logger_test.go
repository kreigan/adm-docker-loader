package loader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("creates log file", func(t *testing.T) {
		dir := t.TempDir()
		logPath := filepath.Join(dir, "test.log")

		logger, err := NewLogger(logPath, false)
		if err != nil {
			t.Fatalf("NewLogger failed: %v", err)
		}
		t.Cleanup(func() {
			//nolint:errcheck // Test cleanup
			logger.Close()
		})

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			t.Error("Log file was not created")
		}
	})

	t.Run("writes to log file", func(t *testing.T) {
		dir := t.TempDir()
		logPath := filepath.Join(dir, "test.log")

		logger, err := NewLogger(logPath, false)
		if err != nil {
			t.Fatalf("NewLogger failed: %v", err)
		}

		logger.Info("test message")
		logger.Warning("warning message")
		logger.Error("error message")
		//nolint:errcheck // Test cleanup
		logger.Close()

		content, err := os.ReadFile(logPath)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		logStr := string(content)
		if !strings.Contains(logStr, "INFO") {
			t.Error("Log file missing INFO entry")
		}
		if !strings.Contains(logStr, "WARN") {
			t.Error("Log file missing WARN entry")
		}
		if !strings.Contains(logStr, "ERROR") {
			t.Error("Log file missing ERROR entry")
		}
	})

	t.Run("debug only in verbose mode", func(t *testing.T) {
		dir := t.TempDir()

		// Non-verbose logger
		logPath1 := filepath.Join(dir, "nonverbose.log")
		logger1, err := NewLogger(logPath1, false)
		if err != nil {
			t.Fatalf("NewLogger failed: %v", err)
		}
		logger1.Debug("debug message")
		//nolint:errcheck // Test cleanup
		logger1.Close()

		content1, err := os.ReadFile(logPath1)
		if err != nil {
			t.Fatalf("Failed to read log: %v", err)
		}
		if strings.Contains(string(content1), "DEBUG") {
			t.Error("Non-verbose logger should not write DEBUG")
		}

		// Verbose logger
		logPath2 := filepath.Join(dir, "verbose.log")
		logger2, err := NewLogger(logPath2, true)
		if err != nil {
			t.Fatalf("NewLogger failed: %v", err)
		}
		logger2.Debug("debug message")
		//nolint:errcheck // Test cleanup
		logger2.Close()

		content2, err := os.ReadFile(logPath2)
		if err != nil {
			t.Fatalf("Failed to read log: %v", err)
		}
		if !strings.Contains(string(content2), "DEBUG") {
			t.Error("Verbose logger should write DEBUG")
		}
	})

	t.Run("creates parent directories", func(t *testing.T) {
		dir := t.TempDir()
		logPath := filepath.Join(dir, "subdir", "nested", "test.log")

		logger, err := NewLogger(logPath, false)
		if err != nil {
			t.Fatalf("NewLogger failed to create parent dirs: %v", err)
		}
		//nolint:errcheck // Test cleanup
		logger.Close()

		if _, err := os.Stat(filepath.Dir(logPath)); os.IsNotExist(err) {
			t.Error("Parent directories were not created")
		}
	})
}
