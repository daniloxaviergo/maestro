package notifier

import (
	"testing"
	"time"
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
