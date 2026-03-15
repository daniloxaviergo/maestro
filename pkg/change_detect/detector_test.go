package change_detect

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"maestro/pkg/logs"
	"maestro/pkg/parser"
)

// createTempLogger creates a temporary log file for testing
func createTempLogger(t *testing.T) (*logs.Logger, string) {
	t.Helper()
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return logger, logPath
}

// createTestFileData creates a FileData for testing
func createTestFileData(filePath string, assignee []string) parser.FileData {
	return parser.FileData{
		FilePath: filePath,
		Frontmatter: parser.Frontmatter{
			Assignee: assignee,
		},
		Error: nil,
	}
}

func TestProcessFile_NewFileNoChange(t *testing.T) {
	logger, logPath := createTempLogger(t)
	d := NewDetector(logger)

	fileData := createTestFileData("test.md", []string{"alice"})

	changed, err := d.ProcessFile(fileData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if changed {
		t.Errorf("Expected no change detected for new file, got true")
	}

	// Process again with same assignee - should still not log
	changed, err = d.ProcessFile(fileData)
	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if changed {
		t.Errorf("Expected no change detected for same assignee on second run, got true")
	}

	// Verify no log file was created (first run should not log)
	// Since NewLogger creates an empty file, check if it's empty
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() > 0 {
		t.Errorf("Expected empty log file for first run, got %d bytes", info.Size())
	}
}

func TestProcessFile_SameAssigneeNoChange(t *testing.T) {
	logger, logPath := createTempLogger(t)
	d := NewDetector(logger)

	// First process - populate cache
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err := d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process with same assignee
	secondData := createTestFileData("test.md", []string{"alice"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if changed {
		t.Errorf("Expected no change detected for same assignee, got true")
	}

	// Verify no log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() > 0 {
		t.Errorf("Expected empty log file when assignee unchanged, got %d bytes", info.Size())
	}
}

func TestProcessFile_DifferentAssigneeChange(t *testing.T) {
	logger, logPath := createTempLogger(t)
	d := NewDetector(logger)

	// First process - add to cache
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err := d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process with different assignee
	secondData := createTestFileData("test.md", []string{"bob"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected for different assignee, got false")
	}

	// Verify log file was created and contains correct data
	entry := parseLogEntry(t, logPath)
	if len(entry.OldAssignee) != 1 || entry.OldAssignee[0] != "alice" {
		t.Errorf("Expected old_assignee to be ['alice'], got %v", entry.OldAssignee)
	}
	if len(entry.NewAssignee) != 1 || entry.NewAssignee[0] != "bob" {
		t.Errorf("Expected new_assignee to be ['bob'], got %v", entry.NewAssignee)
	}
}

func TestProcessFile_EmptyToNonEmpty(t *testing.T) {
	logger, logPath := createTempLogger(t)
	d := NewDetector(logger)

	// First process - empty assignee
	firstData := createTestFileData("test.md", []string{})
	_, err := d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - non-empty assignee
	secondData := createTestFileData("test.md", []string{"bob"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected for empty to non-empty, got false")
	}

	// Verify log file was created and contains correct data
	entry := parseLogEntry(t, logPath)
	if len(entry.OldAssignee) != 0 {
		t.Errorf("Expected old_assignee to be empty, got %v", entry.OldAssignee)
	}
	if len(entry.NewAssignee) != 1 || entry.NewAssignee[0] != "bob" {
		t.Errorf("Expected new_assignee to be ['bob'], got %v", entry.NewAssignee)
	}
}

func TestProcessFile_NonEmptyToEmpty(t *testing.T) {
	logger, logPath := createTempLogger(t)
	d := NewDetector(logger)

	// First process - non-empty assignee
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err := d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - empty assignee
	secondData := createTestFileData("test.md", []string{})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected for non-empty to empty, got false")
	}

	// Verify log file was created and contains correct data
	entry := parseLogEntry(t, logPath)
	if len(entry.OldAssignee) != 1 || entry.OldAssignee[0] != "alice" {
		t.Errorf("Expected old_assignee to be ['alice'], got %v", entry.OldAssignee)
	}
	if len(entry.NewAssignee) != 0 {
		t.Errorf("Expected new_assignee to be empty, got %v", entry.NewAssignee)
	}
}

func TestProcessFile_MultipleAssigneesOrderInsensitive(t *testing.T) {
	logger, logPath := createTempLogger(t)
	d := NewDetector(logger)

	// First process - assignees in one order
	firstData := createTestFileData("test.md", []string{"alice", "bob"})
	_, err := d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - same assignees in different order
	secondData := createTestFileData("test.md", []string{"bob", "alice"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if changed {
		t.Errorf("Expected no change detected for same assignees (different order), got true")
	}

	// Verify no log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() > 0 {
		t.Errorf("Expected empty log file for same assignees with different order, got %d bytes", info.Size())
	}
}

func TestProcessFile_AdditionOfAssignee(t *testing.T) {
	logger, logPath := createTempLogger(t)
	d := NewDetector(logger)

	// First process - single assignee
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err := d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - added assignee
	secondData := createTestFileData("test.md", []string{"alice", "bob"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected for added assignee, got false")
	}

	// Verify log file was created and contains correct data
	entry := parseLogEntry(t, logPath)
	if len(entry.OldAssignee) != 1 || entry.OldAssignee[0] != "alice" {
		t.Errorf("Expected old_assignee to be ['alice'], got %v", entry.OldAssignee)
	}
	if len(entry.NewAssignee) != 2 || entry.NewAssignee[0] != "alice" || entry.NewAssignee[1] != "bob" {
		t.Errorf("Expected new_assignee to be ['alice','bob'], got %v", entry.NewAssignee)
	}
}

func TestRemoveFile(t *testing.T) {
	logger, _ := createTempLogger(t)
	d := NewDetector(logger)

	// Add to cache
	d.cache.SetAssignee("test.md", []string{"alice"})
	d.processed["test.md"] = true

	// Remove file
	d.RemoveFile("test.md")

	// Verify cache entry was removed
	_, exists := d.cache.GetAssignee("test.md")
	if exists {
		t.Errorf("Expected cache entry to be removed")
	}
}

// contains returns true if substr is found in s
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// parseLogEntry parses a JSON log entry from a log file
func parseLogEntry(t *testing.T, logPath string) logs.AssigneeChange {
	t.Helper()
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	var entry logs.AssigneeChange
	if err := json.Unmarshal(content, &entry); err != nil {
		t.Fatalf("Failed to parse log entry: %v", err)
	}

	return entry
}
