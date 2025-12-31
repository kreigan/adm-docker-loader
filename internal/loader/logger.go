package loader

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger handles logging to file and optionally to console.
type Logger struct {
	file    *os.File
	logger  *log.Logger
	verbose bool
}

// NewLogger creates a new logger instance.
func NewLogger(logPath string, verbose bool) (*Logger, error) {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0o750); err != nil {
		return nil, fmt.Errorf("creating log directory: %w", err)
	}

	// Open log file in append mode
	//nolint:gosec // Log path is from trusted base directory
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, fmt.Errorf("opening log file: %w", err)
	}

	// Create multi-writer for file and optionally stdout
	var writers []io.Writer
	writers = append(writers, file)

	// Only write to console in verbose mode
	if verbose {
		writers = append(writers, os.Stdout)
	}

	multiWriter := io.MultiWriter(writers...)
	logger := log.New(multiWriter, "", 0)

	return &Logger{
		file:    file,
		logger:  logger,
		verbose: verbose,
	}, nil
}

// Close closes the log file.
func (l *Logger) Close() error {
	return l.file.Close()
}

// Info logs an informational message.
func (l *Logger) Info(format string, args ...any) {
	l.log("INFO", format, args...)
}

// Warning logs a warning message.
func (l *Logger) Warning(format string, args ...any) {
	l.log("WARN", format, args...)
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...any) {
	l.log("ERROR", format, args...)
}

// Debug logs a debug message (only in verbose mode).
func (l *Logger) Debug(format string, args ...any) {
	if l.verbose {
		l.log("DEBUG", format, args...)
	}
}

// Console prints a clean message to stdout and logs it to file with timestamp.
// This is for user-facing messages that should always be visible.
func (l *Logger) Console(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	// Print to stdout without timestamp
	fmt.Println(message)
	// Log to file with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(l.file, "[%s] INFO: %s\n", timestamp, message)
}

// GetWriter returns an io.Writer for command output.
// Docker compose output always goes to both file and stdout.
func (l *Logger) GetWriter() io.Writer {
	return io.MultiWriter(l.file, os.Stdout)
}

func (l *Logger) log(level, format string, args ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	l.logger.Printf("[%s] %s: %s", timestamp, level, message)
}
