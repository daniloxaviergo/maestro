package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AssigneeChange represents a single log entry for an assignee change
type AssigneeChange struct {
	Timestamp  string   `json:"timestamp"`
	File       string   `json:"file"`
	OldAssignee []string `json:"old_assignee"`
	NewAssignee []string `json:"new_assignee"`
}

// Logger handles JSON logging of assignee changes
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	logPath  string
}

// NewLogger creates a new logger that writes to the specified log file
func NewLogger(logPath string) (*Logger, error) {
	// Ensure log directory exists
	dir := filepath.Dir(logPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file in append mode, create if doesn't exist
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		file:    file,
		logPath: logPath,
	}, nil
}

// LogAssigneeChange writes a JSON log entry for an assignee change
func (l *Logger) LogAssigneeChange(file string, oldAssignee, newAssignee []string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	change := AssigneeChange{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		File:        file,
		OldAssignee: oldAssignee,
		NewAssignee: newAssignee,
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(change, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Write to file with newline
	if _, err := l.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	return nil
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.file.Close()
}
