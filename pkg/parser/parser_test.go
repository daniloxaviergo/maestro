package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const fixturesDir = "./fixtures"

func TestParseFile_ValidFrontmatter(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "valid-frontmatter.md")
	result := ParseFile(filePath)

	// Check no error
	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}

	// Check assignee extraction
	expectedAssignees := []string{"alice", "bob"}
	if len(result.Frontmatter.Assignee) != len(expectedAssignees) {
		t.Errorf("Expected %d assignees, got %d", len(expectedAssignees), len(result.Frontmatter.Assignee))
	}
	for i, expected := range expectedAssignees {
		if result.Frontmatter.Assignee[i] != expected {
			t.Errorf("Expected assignee[%d] to be %q, got %q", i, expected, result.Frontmatter.Assignee[i])
		}
	}

	// Check other fields
	if result.Frontmatter.ID != "task-001" {
		t.Errorf("Expected ID to be 'task-001', got %q", result.Frontmatter.ID)
	}
	if result.Frontmatter.Title != "Example Task" {
		t.Errorf("Expected Title to be 'Example Task', got %q", result.Frontmatter.Title)
	}

	// Check parse time is reasonable (< 1 second)
	if result.ParseTime >= time.Second {
		t.Errorf("Parse time too long: %v", result.ParseTime)
	}
}

func TestParseFile_EmptyAssignee(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "empty-assignee.md")
	result := ParseFile(filePath)

	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}

	// Check assignee is empty slice
	if len(result.Frontmatter.Assignee) != 0 {
		t.Errorf("Expected empty assignee array, got: %v", result.Frontmatter.Assignee)
	}
}

func TestParseFile_MissingAssignee(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "missing-assignee.md")
	result := ParseFile(filePath)

	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}

	// Check assignee defaults to empty slice
	if len(result.Frontmatter.Assignee) != 0 {
		t.Errorf("Expected empty assignee array (default), got: %v", result.Frontmatter.Assignee)
	}
}

func TestParseFile_NoFrontmatter(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "no-frontmatter.md")
	result := ParseFile(filePath)

	// Files without frontmatter should not be an error
	if result.Error != nil {
		t.Errorf("Expected no error for file without frontmatter, got: %v", result.Error)
	}

	// Assignee should be empty slice
	if len(result.Frontmatter.Assignee) != 0 {
		t.Errorf("Expected empty assignee array for file without frontmatter, got: %v", result.Frontmatter.Assignee)
	}
}

func TestParseFile_MalformedYAML(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "malformed-yaml.md")
	result := ParseFile(filePath)

	// This should return an error
	if result.Error == nil {
		t.Errorf("Expected error for malformed YAML, got nil")
	} else {
		// Verify error message is descriptive
		errMsg := result.Error.Error()
		if len(errMsg) < 10 {
			t.Errorf("Error message should be descriptive, got: %q", errMsg)
		}
	}

	// Assignee should still be empty slice on error
	if len(result.Frontmatter.Assignee) != 0 {
		t.Errorf("Expected empty assignee on error, got: %v", result.Frontmatter.Assignee)
	}
}

func TestParseFile_SingleAssignee(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "single-assignee.md")
	result := ParseFile(filePath)

	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}

	// Check single assignee
	if len(result.Frontmatter.Assignee) != 1 {
		t.Errorf("Expected 1 assignee, got %d", len(result.Frontmatter.Assignee))
	}
	if result.Frontmatter.Assignee[0] != "charlie" {
		t.Errorf("Expected assignee to be 'charlie', got %q", result.Frontmatter.Assignee[0])
	}
}

func TestParseFile_NonExistentFile(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "non-existent.md")
	result := ParseFile(filePath)

	if result.Error == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestParseFile_FilePathInResult(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "valid-frontmatter.md")
	result := ParseFile(filePath)

	if result.FilePath != filePath {
		t.Errorf("Expected FilePath to be %q, got %q", filePath, result.FilePath)
	}
}

// Test for performance requirement (<100ms)
func TestParseFile_Performance(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "valid-frontmatter.md")
	result := ParseFile(filePath)

	if result.Error != nil {
		t.Skipf("Skipping performance test due to error: %v", result.Error)
	}

	// Check parse time is under 100ms
	maxExpectedTime := 100 * time.Millisecond
	if result.ParseTime >= maxExpectedTime {
		t.Errorf("Parse time %v exceeded maximum expected %v", result.ParseTime, maxExpectedTime)
	}
}

// Benchmark test for performance tracking
func BenchmarkParseFile(b *testing.B) {
	filePath := filepath.Join(fixturesDir, "valid-frontmatter.md")
	
	// Ensure file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		b.Skipf("Fixture file not found: %s", filePath)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseFile(filePath)
	}
}
