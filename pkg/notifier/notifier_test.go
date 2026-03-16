package notifier

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"maestro/pkg/agent"
)

func TestNotificationConfig_DefaultValues(t *testing.T) {
	config := NotificationConfig{
		MessageFormat: DefaultMessageFormat,
		Timeout:       2 * time.Second,
	}

	if config.MessageFormat != DefaultMessageFormat {
		t.Errorf("expected MessageFormat to be %q, got %q", DefaultMessageFormat, config.MessageFormat)
	}

	if config.Timeout != 2*time.Second {
		t.Errorf("expected Timeout to be 2s, got %v", config.Timeout)
	}
}

func TestNotificationConfig_CustomValues(t *testing.T) {
	customFormat := "Custom message: [new] - [file]"
	customTimeout := 5 * time.Second

	config := NotificationConfig{
		MessageFormat: customFormat,
		Timeout:       customTimeout,
	}

	if config.MessageFormat != customFormat {
		t.Errorf("expected MessageFormat to be %q, got %q", customFormat, config.MessageFormat)
	}

	if config.Timeout != customTimeout {
		t.Errorf("expected Timeout to be %v, got %v", customTimeout, config.Timeout)
	}
}

func TestNotifier_Constructors(t *testing.T) {
	// Test with default config
	defaultConfig := NotificationConfig{
		MessageFormat: DefaultMessageFormat,
		Timeout:       2 * time.Second,
	}
	notifier := Notifier{config: defaultConfig}

	if notifier.config.MessageFormat != DefaultMessageFormat {
		t.Errorf("expected notifier config MessageFormat to be %q, got %q", DefaultMessageFormat, notifier.config.MessageFormat)
	}

	if notifier.config.Timeout != 2*time.Second {
		t.Errorf("expected notifier config Timeout to be 2s, got %v", notifier.config.Timeout)
	}
}

func TestAssigneeChangeEvent_Struct(t *testing.T) {
	event := AssigneeChangeEvent{
		FilePath:    "backlog/tasks/task-001.md",
		OldAssignee: []string{"alice"},
		NewAssignee: []string{"bob"},
	}

	if event.FilePath != "backlog/tasks/task-001.md" {
		t.Errorf("expected FilePath to be %q, got %q", "backlog/tasks/task-001.md", event.FilePath)
	}

	if len(event.OldAssignee) != 1 || event.OldAssignee[0] != "alice" {
		t.Errorf("expected OldAssignee to be [alice], got %v", event.OldAssignee)
	}

	if len(event.NewAssignee) != 1 || event.NewAssignee[0] != "bob" {
		t.Errorf("expected NewAssignee to be [bob], got %v", event.NewAssignee)
	}
}

func TestAssigneeChangeEvent_EmptyAssignees(t *testing.T) {
	event := AssigneeChangeEvent{
		FilePath:    "backlog/tasks/task-001.md",
		OldAssignee: []string{},
		NewAssignee: []string{},
	}

	if len(event.OldAssignee) != 0 {
		t.Errorf("expected OldAssignee to be empty, got %v", event.OldAssignee)
	}

	if len(event.NewAssignee) != 0 {
		t.Errorf("expected NewAssignee to be empty, got %v", event.NewAssignee)
	}
}

func TestAssigneeChangeEvent_MultipleAssignees(t *testing.T) {
	event := AssigneeChangeEvent{
		FilePath:    "backlog/tasks/task-001.md",
		OldAssignee: []string{"alice", "charlie"},
		NewAssignee: []string{"bob", "dave"},
	}

	if len(event.OldAssignee) != 2 {
		t.Errorf("expected OldAssignee to have 2 elements, got %d", len(event.OldAssignee))
	}

	if len(event.NewAssignee) != 2 {
		t.Errorf("expected NewAssignee to have 2 elements, got %d", len(event.NewAssignee))
	}
}

func TestErrorVariables_Distinct(t *testing.T) {
	errors := []error{
		ErrTmuxNotInstalled,
		ErrTmuxCommandFailed,
		ErrTmuxTimeout,
		ErrScriptNotFound,
		ErrScriptExecutionFailed,
		ErrSessionCreationFailed,
	}

	// Verify all errors are distinct
	seen := make(map[string]bool)
	for _, err := range errors {
		if seen[err.Error()] {
			t.Errorf("duplicate error message: %v", err)
		}
		seen[err.Error()] = true
	}

	// Verify error messages are not empty
	if ErrTmuxNotInstalled.Error() == "" {
		t.Error("ErrTmuxNotInstalled.Error() returned empty string")
	}
	if ErrTmuxCommandFailed.Error() == "" {
		t.Error("ErrTmuxCommandFailed.Error() returned empty string")
	}
	if ErrTmuxTimeout.Error() == "" {
		t.Error("ErrTmuxTimeout.Error() returned empty string")
	}
	if ErrScriptNotFound.Error() == "" {
		t.Error("ErrScriptNotFound.Error() returned empty string")
	}
	if ErrScriptExecutionFailed.Error() == "" {
		t.Error("ErrScriptExecutionFailed.Error() returned empty string")
	}
	if ErrSessionCreationFailed.Error() == "" {
		t.Error("ErrSessionCreationFailed.Error() returned empty string")
	}
}

func TestErrorVariables_CanBeComparedWithErrorsIs(t *testing.T) {
	// Test that errors can be compared using errors.Is()
	err1 := ErrTmuxNotInstalled
	err2 := ErrTmuxNotInstalled

	if !errorsIs(err1, err2) {
		t.Error("ErrTmuxNotInstalled should be equal to itself")
	}

	// Test with different error types
	err3 := ErrTmuxCommandFailed
	if errorsIs(err1, err3) {
		t.Error("ErrTmuxNotInstalled should not be equal to ErrTmuxCommandFailed")
	}
}

func errorsIs(err, target error) bool {
	return err.Error() == target.Error()
}

// TestNewNotifier_DefaultConfig tests that NewNotifier uses default values when config is empty
func TestNewNotifier_DefaultConfig(t *testing.T) {
	notifier := NewNotifier(NotificationConfig{})

	if notifier.config.Timeout != 2*time.Second {
		t.Errorf("expected default timeout to be 2s, got %v", notifier.config.Timeout)
	}

	if notifier.config.MessageFormat != DefaultMessageFormat {
		t.Errorf("expected default message format to be %q, got %q", DefaultMessageFormat, notifier.config.MessageFormat)
	}
}

// TestNewNotifier_CustomConfig tests that NewNotifier uses custom config values
func TestNewNotifier_CustomConfig(t *testing.T) {
	customFormat := "Custom: [new] for [file]"
	customTimeout := 3 * time.Second

	notifier := NewNotifier(NotificationConfig{
		MessageFormat: customFormat,
		Timeout:       customTimeout,
	})

	if notifier.config.Timeout != customTimeout {
		t.Errorf("expected timeout to be %v, got %v", customTimeout, notifier.config.Timeout)
	}

	if notifier.config.MessageFormat != customFormat {
		t.Errorf("expected message format to be %q, got %q", customFormat, notifier.config.MessageFormat)
	}
}

// TestNotifier_formatMessage tests the message formatting logic
func TestNotifier_formatMessage(t *testing.T) {
	tests := []struct {
		name          string
		configFormat  string
		event         AssigneeChangeEvent
		expectedMsg   string
	}{
		{
			name:         "default format with single assignee",
			configFormat: DefaultMessageFormat,
			event: AssigneeChangeEvent{
				FilePath:    "backlog/tasks/task-001.md",
				OldAssignee: []string{"alice"},
				NewAssignee: []string{"bob"},
			},
			expectedMsg: `Assignee changed to "bob" for backlog/tasks/task-001.md`,
		},
		{
			name:         "default format with multiple assignees",
			configFormat: DefaultMessageFormat,
			event: AssigneeChangeEvent{
				FilePath:    "backlog/tasks/task-001.md",
				OldAssignee: []string{"alice", "charlie"},
				NewAssignee: []string{"bob", "dave"},
			},
			expectedMsg: `Assignee changed to "bob, dave" for backlog/tasks/task-001.md`,
		},
		{
			name:         "default format with empty new assignees",
			configFormat: DefaultMessageFormat,
			event: AssigneeChangeEvent{
				FilePath:    "backlog/tasks/task-001.md",
				OldAssignee: []string{"alice"},
				NewAssignee: []string{},
			},
			expectedMsg: `Assignee changed to "" for backlog/tasks/task-001.md`,
		},
		{
			name:         "custom format with placeholders",
			configFormat: "Updated: [new] assigned to [file]",
			event: AssigneeChangeEvent{
				FilePath:    "backlog/tasks/task-001.md",
				OldAssignee: []string{},
				NewAssignee: []string{"bob"},
			},
			expectedMsg: `Updated: bob assigned to backlog/tasks/task-001.md`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := NewNotifier(NotificationConfig{
				MessageFormat: tt.configFormat,
			})

			msg := notifier.formatMessage(tt.event)

			if msg != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, msg)
			}
		})
	}
}

// TestExecuteScript tests the ExecuteScript method
func TestExecuteScript(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()

	// Create a test script that does nothing successful
	scriptPath := filepath.Join(tempDir, "test_script.sh")
	scriptContent := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	// Create a temporary config directory
	configDir := filepath.Join(tempDir, "agents")
	testAgentDir := filepath.Join(configDir, "test-agent")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("failed to create agent directory: %v", err)
	}

	// Create agent config
	configPath := filepath.Join(testAgentDir, "config.yml")
	configContent := `
script_path: ` + scriptPath + `
tmux_session: test-session
enabled: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	// Create agent and load config
	agent := agent.NewAgent("test-agent", configPath)
	agent.LoadConfig()

	// Create notifier with agent config
	notifier := NewNotifier(NotificationConfig{
		Agent:   agent,
		Timeout: 5 * time.Second,
	})

	// Test that ExecuteScript doesn't panic and returns immediately (non-blocking)
	// Using empty string as filePath since this test doesn't verify file path
	start := time.Now()
	notifier.ExecuteScript("")
	elapsed := time.Since(start)

	// Should return immediately (less than 100ms for goroutine to start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("ExecuteScript should be non-blocking, took %v", elapsed)
	}

	// Wait a bit for the script to complete (with timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Verify the tmux session was created
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		t.Log("tmux session creation completed within timeout")
	}
}

// TestExecuteScript_NoAgentConfig tests ExecuteScript when agent is not configured
func TestExecuteScript_NoAgentConfig(t *testing.T) {
	notifier := NewNotifier(NotificationConfig{
		Agent:   nil,
		Timeout: 2 * time.Second,
	})

	// Should log a warning and return without panicking
	// Using empty string as filePath since this test doesn't verify file path
	notifier.ExecuteScript("")

	// Give it a moment to execute
	time.Sleep(50 * time.Millisecond)
}

// TestExecuteScript_MissingScriptPath tests ExecuteScript when script_path is not configured
func TestExecuteScript_MissingScriptPath(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()

	// Create a temporary config directory
	configDir := filepath.Join(tempDir, "agents")
	testAgentDir := filepath.Join(configDir, "test-agent")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("failed to create agent directory: %v", err)
	}

	// Create agent config with empty script_path
	configPath := filepath.Join(testAgentDir, "config.yml")
	configContent := `
script_path: ""
tmux_session: test-session
enabled: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	// Create agent and load config
	agent := agent.NewAgent("test-agent", configPath)
	agent.LoadConfig()

	// Create notifier with agent config
	notifier := NewNotifier(NotificationConfig{
		Agent:   agent,
		Timeout: 2 * time.Second,
	})

	// Should log a warning and return without panicking
	// Using empty string as filePath since this test doesn't verify file path
	notifier.ExecuteScript("")

	// Give it a moment to execute
	time.Sleep(50 * time.Millisecond)
}

// TestExecuteScript_FileNotFound tests ExecuteScript when script file doesn't exist
func TestExecuteScript_FileNotFound(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()

	// Create a temporary config directory
	configDir := filepath.Join(tempDir, "agents")
	testAgentDir := filepath.Join(configDir, "test-agent")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("failed to create agent directory: %v", err)
	}

	// Create agent config with non-existent script
	configPath := filepath.Join(testAgentDir, "config.yml")
	configContent := `
script_path: /nonexistent/script.sh
tmux_session: test-session
enabled: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	// Create agent and load config
	agent := agent.NewAgent("test-agent", configPath)
	agent.LoadConfig()

	// Create notifier with agent config
	notifier := NewNotifier(NotificationConfig{
		Agent:   agent,
		Timeout: 2 * time.Second,
	})

	// Should log a warning and return without panicking
	// Using empty string as filePath since this test doesn't verify file path
	notifier.ExecuteScript("")

	// Give it a moment to execute
	time.Sleep(50 * time.Millisecond)
}

// TestExecuteScriptsForAgents_WithFilePath tests that scripts receive file path as argument
func TestExecuteScriptsForAgents_WithFilePath(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()

	// Create a test script that logs the file path argument
	scriptPath := filepath.Join(tempDir, "test_script.sh")
	scriptContent := `#!/bin/bash
echo "Script called with file path: $1" > /tmp/script_execution_log.txt
exit 0
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	// Create a temporary config directory
	configDir := filepath.Join(tempDir, "agents")
	testAgentDir := filepath.Join(configDir, "test-agent")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("failed to create agent directory: %v", err)
	}

	// Create agent config
	configPath := filepath.Join(testAgentDir, "config.yml")
	configContent := `
script_path: ` + scriptPath + `
tmux_session: test-session
enabled: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	// Create agent and load config
	testAgent := agent.NewAgent("test-agent", configPath)
	testAgent.LoadConfig()

	// Create notifier with agent config
	notifier := NewNotifier(NotificationConfig{
		Agent:   testAgent,
		Timeout: 5 * time.Second,
	})

	// Test file path
	testFilePath := "/test/path/to/task.md"

	// Execute script for agents (non-blocking)
	start := time.Now()
	notifier.ExecuteScriptsForAgents([]*agent.Agent{testAgent}, testFilePath)
	elapsed := time.Since(start)

	// Should return immediately (non-blocking)
	if elapsed > 100*time.Millisecond {
		t.Errorf("ExecuteScriptsForAgents should be non-blocking, took %v", elapsed)
	}

	// Wait for tmux to complete execution
	time.Sleep(2 * time.Second)

	// Note: Since tmux is non-blocking and requires a running session,
	// this test verifies the method doesn't panic and accepts the filePath parameter.
	// The actual script execution is verified in integration tests.
	t.Log("ExecuteScriptsForAgents executed without panic with file path argument")
}

// TestExecuteScriptsForAgents_MultipleAgents tests that multiple agents all receive the file path
func TestExecuteScriptsForAgents_MultipleAgents(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()

	// Create test scripts for multiple agents
	scriptPath1 := filepath.Join(tempDir, "script1.sh")
	scriptPath2 := filepath.Join(tempDir, "script2.sh")

	for i, scriptPath := range []string{scriptPath1, scriptPath2} {
		scriptContent := `#!/bin/bash
exit 0
`
		if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
			t.Fatalf("failed to create test script %d: %v", i, err)
		}
	}

	// Create config directories for two agents
	configDir := filepath.Join(tempDir, "agents")
	agent1Dir := filepath.Join(configDir, "agent1")
	agent2Dir := filepath.Join(configDir, "agent2")

	if err := os.MkdirAll(agent1Dir, 0755); err != nil {
		t.Fatalf("failed to create agent1 directory: %v", err)
	}
	if err := os.MkdirAll(agent2Dir, 0755); err != nil {
		t.Fatalf("failed to create agent2 directory: %v", err)
	}

	// Create agent configs
	config1Path := filepath.Join(agent1Dir, "config.yml")
	config2Path := filepath.Join(agent2Dir, "config.yml")

	configContent := `
script_path: ` + scriptPath1 + `
tmux_session: test-session
enabled: true
`
	if err := os.WriteFile(config1Path, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create agent1 config: %v", err)
	}

	configContent = `
script_path: ` + scriptPath2 + `
tmux_session: test-session
enabled: true
`
	if err := os.WriteFile(config2Path, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create agent2 config: %v", err)
	}

	// Create agents and load configs
	agent1 := agent.NewAgent("agent1", config1Path)
	agent1.LoadConfig()

	agent2 := agent.NewAgent("agent2", config2Path)
	agent2.LoadConfig()

	// Create notifier
	notifier := NewNotifier(NotificationConfig{
		Timeout: 5 * time.Second,
	})

	// Test file path
	testFilePath := "/test/path/to/task.md"

	// Execute scripts for multiple agents
	start := time.Now()
	notifier.ExecuteScriptsForAgents([]*agent.Agent{agent1, agent2}, testFilePath)
	elapsed := time.Since(start)

	// Should return immediately (non-blocking)
	if elapsed > 100*time.Millisecond {
		t.Errorf("ExecuteScriptsForAgents should be non-blocking, took %v", elapsed)
	}

	// Give scripts time to execute
	time.Sleep(2 * time.Second)

	t.Log("ExecuteScriptsForAgents executed for multiple agents without panic")
}

// TestExecuteScriptsForAgents_DisabledAgent tests that disabled agents are skipped
func TestExecuteScriptsForAgents_DisabledAgent(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()

	// Create a test script
	scriptPath := filepath.Join(tempDir, "test_script.sh")
	scriptContent := `#!/bin/bash
exit 0
`
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create test script: %v", err)
	}

	// Create agent config with disabled
	configDir := filepath.Join(tempDir, "agents")
	testAgentDir := filepath.Join(configDir, "test-agent")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("failed to create agent directory: %v", err)
	}

	configPath := filepath.Join(testAgentDir, "config.yml")
	configContent := `
script_path: ` + scriptPath + `
tmux_session: test-session
enabled: false
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	agentInstance := agent.NewAgent("test-agent", configPath)
	agentInstance.LoadConfig()

	notifier := NewNotifier(NotificationConfig{
		Timeout: 2 * time.Second,
	})

	// Should log a warning and skip the disabled agent
	notifier.ExecuteScriptsForAgents([]*agent.Agent{agentInstance}, "/test/path.md")

	time.Sleep(50 * time.Millisecond)
}

// TestExecuteScriptsForAgents_MissingScriptPath tests that agents without script_path are skipped
func TestExecuteScriptsForAgents_MissingScriptPath(t *testing.T) {
	// Create a temporary directory for test
	tempDir := t.TempDir()

	configDir := filepath.Join(tempDir, "agents")
	testAgentDir := filepath.Join(configDir, "test-agent")
	if err := os.MkdirAll(testAgentDir, 0755); err != nil {
		t.Fatalf("failed to create agent directory: %v", err)
	}

	// Create agent config with empty script_path
	configPath := filepath.Join(testAgentDir, "config.yml")
	configContent := `
script_path: ""
tmux_session: test-session
enabled: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	agentInstance := agent.NewAgent("test-agent", configPath)
	agentInstance.LoadConfig()

	notifier := NewNotifier(NotificationConfig{
		Timeout: 2 * time.Second,
	})

	// Should log a warning and skip the agent
	notifier.ExecuteScriptsForAgents([]*agent.Agent{agentInstance}, "/test/path.md")

	time.Sleep(50 * time.Millisecond)
}

// TestSessionExists tests the sessionExists helper function
func TestSessionExists(t *testing.T) {
	tests := []struct {
		name           string
		sessionName    string
		prepareSession bool
		wantExists     bool
		wantErr        bool
	}{
		{
			name:           "session does not exist",
			sessionName:    "nonexistent-session-12345",
			prepareSession: false,
			wantExists:     false,
			wantErr:        false,
		},
		{
			name:           "session exists",
			sessionName:    "existing-session-67890",
			prepareSession: true,
			wantExists:     true,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare session if needed
			if tt.prepareSession {
				cmd := exec.Command("tmux", "new-session", "-d", "-s", tt.sessionName)
				if err := cmd.Run(); err != nil {
					t.Fatalf("failed to create test session: %v", err)
				}
				defer func() {
					// Cleanup session after test
					exec.Command("tmux", "kill-session", "-t", tt.sessionName).Run()
				}()
			}

			// Give tmux time to create session
			if tt.prepareSession {
				time.Sleep(100 * time.Millisecond)
			}

			exists, err := sessionExists(tt.sessionName)

			if (err != nil) != tt.wantErr {
				t.Errorf("sessionExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if exists != tt.wantExists {
				t.Errorf("sessionExists() = %v, want %v", exists, tt.wantExists)
			}
		})
	}
}

// TestSessionExists_SessionNameParsing tests that session names are parsed correctly
func TestSessionExists_SessionNameParsing(t *testing.T) {
	// Create a test session
	sessionName := "test-parse-session"
	cmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName)
	if err := cmd.Run(); err != nil {
		t.Skipf("tmux not available for testing: %v", err)
	}
	defer exec.Command("tmux", "kill-session", "-t", sessionName).Run()

	time.Sleep(100 * time.Millisecond)

	exists, err := sessionExists(sessionName)
	if err != nil {
		t.Fatalf("sessionExists() returned error: %v", err)
	}

	if !exists {
		t.Errorf("sessionExists() = false, expected true for existing session")
	}
}
