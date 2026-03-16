package change_detect

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"maestro/pkg/agent"
	"maestro/pkg/logs"
	"maestro/pkg/matcher"
	"maestro/pkg/notifier"
	"maestro/pkg/parser"
)

// orchestrator_test.go contains integration tests for the agent orchestration system.
// These tests verify the full flow from assignee change detection through agent matching
// to script execution via the notifier. Tests cover error handling, graceful degradation,
// and edge cases like disabled agents, missing scripts, and concurrent processing.

// createTempDir creates a temporary directory and returns cleanup function
func createTempDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return tmpDir
}

// createTempConfig creates a temporary agent config file and returns path
func createTempConfig(t *testing.T, tmpDir, agentName, scriptPath string, enabled bool) string {
	t.Helper()
	agentDir := filepath.Join(tmpDir, "agents", agentName)
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		t.Fatalf("failed to create agent directory: %v", err)
	}

	configPath := filepath.Join(agentDir, "config.yml")
	configContent := `
script_path: ` + scriptPath + `
tmux_session: test-session
enabled: ` + boolToString(enabled) + `
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}
	return configPath
}

// boolToString converts bool to string for YAML
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// TestOrchestrator_FullFlow tests the complete orchestrator flow
func TestOrchestrator_FullFlow(t *testing.T) {
	tmpDir := createTempDir(t)

	// Create a test script
	scriptPath := filepath.Join(tmpDir, "script.sh")
	scriptContent := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	// Create test agent config
	configPath := createTempConfig(t, tmpDir, "alice", scriptPath, true)

	// Create test agent
	testAgent := agent.NewAgent("alice", configPath)
	testAgent.LoadConfig()

	// Create matcher
	matcher := matcher.NewMatcher([]*agent.Agent{testAgent})

	// Create detector with temp logger
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	d.SetMatcher(matcher)

	// Create a notifier with a short timeout for testing
	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// First process - add to cache
	firstData := createTestFileData("test.md", []string{"alice"})
	changed, err := d.ProcessFile(firstData)

	if err != nil {
		t.Errorf("First ProcessFile returned error: %v", err)
	}
	if changed {
		t.Errorf("Expected no change for first run, got true")
	}

	// Second process - add another assignee to trigger change and script execution
	secondData := createTestFileData("test.md", []string{"alice", "bob"})

	changed, err = d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("Second ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected for different assignee, got false")
	}

	// Verify log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Expected log file to have content, got empty")
	}
}

// TestOrchestrator_DisabledAgent tests that disabled agents are skipped
func TestOrchestrator_DisabledAgent(t *testing.T) {
	tmpDir := createTempDir(t)

	// Create script for enabled agent
	scriptPath1 := filepath.Join(tmpDir, "script1.sh")
	scriptContent := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(scriptPath1, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	// Create enabled agent config
	configPath1 := createTempConfig(t, tmpDir, "alice", scriptPath1, true)
	testAgent1 := agent.NewAgent("alice", configPath1)
	testAgent1.LoadConfig()

	// Create disabled agent config (script won't exist, but should not panic)
	scriptPath2 := "/nonexistent/script2.sh"
	configPath2 := createTempConfig(t, tmpDir, "bob", scriptPath2, false)
	testAgent2 := agent.NewAgent("bob", configPath2)
	testAgent2.LoadConfig()

	// Create matcher with both agents
	matcher := matcher.NewMatcher([]*agent.Agent{testAgent1, testAgent2})

	// Create detector with temp logger
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	d.SetMatcher(matcher)

	// Create a notifier
	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// First process - add to cache
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err = d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - add bob (disabled agent should be skipped gracefully)
	secondData := createTestFileData("test.md", []string{"alice", "bob"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected, got false")
	}

	// Give goroutines time to execute (disabled agent should not execute script)
	time.Sleep(100 * time.Millisecond)

	// Verify log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Expected log file to have content, got empty")
	}
}

// TestOrchestrator_MissingScript tests graceful handling of missing script
func TestOrchestrator_MissingScript(t *testing.T) {
	tmpDir := createTempDir(t)

	// Agent config with non-existent script
	configPath := createTempConfig(t, tmpDir, "alice", "/nonexistent/script.sh", true)

	testAgent := agent.NewAgent("alice", configPath)
	testAgent.LoadConfig()

	// Create matcher
	matcher := matcher.NewMatcher([]*agent.Agent{testAgent})

	// Create detector with temp logger
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	d.SetMatcher(matcher)

	// Create a notifier
	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// First process - add to cache
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err = d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - change assignee (script execution should fail gracefully)
	secondData := createTestFileData("test.md", []string{"bob"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected, got false")
	}

	// Give goroutines time to execute
	time.Sleep(100 * time.Millisecond)

	// Verify log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Expected log file to have content, got empty")
	}
}

// TestOrchestrator_MultipleAgentsSameAssignee tests when multiple agents match the same assignee
func TestOrchestrator_MultipleAgentsSameAssignee(t *testing.T) {
	tmpDir := createTempDir(t)

	// Create scripts for both agents
	scriptPath1 := filepath.Join(tmpDir, "script1.sh")
	scriptContent := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(scriptPath1, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	scriptPath2 := filepath.Join(tmpDir, "script2.sh")
	if err := os.WriteFile(scriptPath2, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	// Create two agents with same assignee name (unusual but valid)
	configPath1 := createTempConfig(t, tmpDir, "alice", scriptPath1, true)
	testAgent1 := agent.NewAgent("alice", configPath1)
	testAgent1.LoadConfig()

	configPath2 := createTempConfig(t, tmpDir, "alice-secondary", scriptPath2, true)
	testAgent2 := agent.NewAgent("alice-secondary", configPath2)
	testAgent2.LoadConfig()

	// Create matcher
	matcher := matcher.NewMatcher([]*agent.Agent{testAgent1, testAgent2})

	// Create detector with temp logger
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	d.SetMatcher(matcher)

	// Create a notifier
	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// First process - add to cache
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err = d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - change to both agents
	secondData := createTestFileData("test.md", []string{"alice", "alice-secondary"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected, got false")
	}

	// Give goroutines time to execute
	time.Sleep(100 * time.Millisecond)

	// Verify log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Expected log file to have content, got empty")
	}
}

// TestOrchestrator_ConcurrentFileProcessing tests concurrent file processing
func TestOrchestrator_ConcurrentFileProcessing(t *testing.T) {
	// Create script in temp directory
	tmpDir := createTempDir(t)
	scriptPath := filepath.Join(tmpDir, "script.sh")
	scriptContent := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	configPath := createTempConfig(t, tmpDir, "alice", scriptPath, true)
	testAgent := agent.NewAgent("alice", configPath)
	testAgent.LoadConfig()

	matcher := matcher.NewMatcher([]*agent.Agent{testAgent})

	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	d.SetMatcher(matcher)

	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// Process multiple files concurrently
	numFiles := 10
	done := make(chan bool, numFiles)

	for i := 0; i < numFiles; i++ {
		go func(index int) {
			filePath := filepath.Join(tmpDir, "test"+string(rune('0'+index))+".md")
			// Use a file data with alice as assignee
			fileData := parser.FileData{
				FilePath: filePath,
				Frontmatter: parser.Frontmatter{
					Assignee: []string{"alice"},
				},
				Error: nil,
			}
			_, err := d.ProcessFile(fileData)
			if err != nil {
				t.Errorf("ProcessFile for file %d returned error: %v", index, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numFiles; i++ {
		<-done
	}

	// Verify at least some log entries were created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}

	// First run should not log, but since we have 10 files, only first run per file
	// won't log. This is expected behavior.
	t.Logf("Log file size: %d bytes", info.Size())
}

// TestOrchestrator_NoNotifier tests that detector works without notifier
func TestOrchestrator_NoNotifier(t *testing.T) {
	tmpDir := createTempDir(t)

	configPath := createTempConfig(t, tmpDir, "alice", "/nonexistent/script.sh", true)
	testAgent := agent.NewAgent("alice", configPath)
	testAgent.LoadConfig()

	matcher := matcher.NewMatcher([]*agent.Agent{testAgent})

	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	d.SetMatcher(matcher)
	// Notifier is nil

	// First process - add to cache
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err = d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - change assignee (no notifier, but should still log)
	secondData := createTestFileData("test.md", []string{"bob"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected, got false")
	}

	// Verify log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Expected log file to have content, got empty")
	}
}

// TestOrchestrator_NoMatcher tests that detector works without matcher
func TestOrchestrator_NoMatcher(t *testing.T) {
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	// Matcher is nil
	d.SetMatcher(nil)

	// Create a notifier
	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// First process - add to cache
	firstData := createTestFileData("test.md", []string{"alice"})
	_, err = d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - change assignee (no matcher, but should still log and notify)
	secondData := createTestFileData("test.md", []string{"bob"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if !changed {
		t.Errorf("Expected change detected, got false")
	}

	// Give notifier time to execute
	time.Sleep(100 * time.Millisecond)

	// Verify log file was created
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Expected log file to have content, got empty")
	}
}

// TestOrchestrator_EmptyAgentList tests edge case with empty agent list
func TestOrchestrator_EmptyAgentList(t *testing.T) {
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create matcher with empty agent list
	matcher := matcher.NewMatcher(nil)
	d := NewDetector(logger)
	d.SetMatcher(matcher)

	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// Process a file with assignee
	fileData := createTestFileData("test.md", []string{"alice"})
	changed, err := d.ProcessFile(fileData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	// First run should not log
	if changed {
		t.Errorf("Expected no change for first run, got true")
	}
}

// TestProcessFile_WithMatcher_OrderInsensitiveWithAgents tests that order-insensitive
// assignee matching works correctly with agent orchestration
func TestProcessFile_WithMatcher_OrderInsensitiveWithAgents(t *testing.T) {
	tmpDir := createTempDir(t)

	scriptPath := filepath.Join(tmpDir, "script.sh")
	scriptContent := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	configPath := createTempConfig(t, tmpDir, "alice", scriptPath, true)
	testAgent := agent.NewAgent("alice", configPath)
	testAgent.LoadConfig()

	matcher := matcher.NewMatcher([]*agent.Agent{testAgent})
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)
	d.SetMatcher(matcher)

	notifier := notifier.NewNotifier(notifier.NotificationConfig{
		Timeout: 1 * time.Second,
	})
	d.SetNotifier(notifier)

	// First process - assignees in one order
	firstData := createTestFileData("test.md", []string{"alice", "bob"})
	_, err = d.ProcessFile(firstData)
	if err != nil {
		t.Fatalf("First ProcessFile failed: %v", err)
	}

	// Second process - same assignees in different order (should not trigger script)
	secondData := createTestFileData("test.md", []string{"bob", "alice"})

	changed, err := d.ProcessFile(secondData)

	if err != nil {
		t.Errorf("ProcessFile returned error: %v", err)
	}
	if changed {
		t.Errorf("Expected no change for same assignees (different order), got true")
	}
}

// TestProcessFile_NilFileData tests handling of empty FileData
func TestProcessFile_NilFileData(t *testing.T) {
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)

	// Process empty file data (no error)
	changed, err := d.ProcessFile(parser.FileData{})

	if err != nil {
		t.Errorf("ProcessFile returned error for empty FileData: %v", err)
	}
	if changed {
		t.Errorf("Expected no change for empty FileData, got true")
	}
}

// TestProcessFile_ErrorFileData tests handling of FileData with error
func TestProcessFile_ErrorFileData(t *testing.T) {
	tmpLoggerDir := t.TempDir()
	logPath := filepath.Join(tmpLoggerDir, "test_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	d := NewDetector(logger)

	// Process file data with error (using a real error)
	fileData := parser.FileData{
		FilePath: "test.md",
		Frontmatter: parser.Frontmatter{
			Assignee: []string{"alice"},
		},
		Error: os.ErrNotExist,
	}

	changed, err := d.ProcessFile(fileData)

	if err != nil {
		t.Errorf("ProcessFile returned error for error FileData: %v", err)
	}
	if changed {
		t.Errorf("Expected no change for error FileData, got true")
	}
}
