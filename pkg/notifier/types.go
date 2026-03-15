package notifier

import "time"

// NotificationConfig holds configuration for the notifier
type NotificationConfig struct {
	MessageFormat string
	Timeout       time.Duration
}

// Notifier holds the configuration and state for the tmux notifier
type Notifier struct {
	config NotificationConfig
}

// AssigneeChangeEvent represents a change in file assignee
type AssigneeChangeEvent struct {
	FilePath    string
	OldAssignee []string
	NewAssignee []string
}

// DefaultMessageFormat is the default message format for notifications
const DefaultMessageFormat = `Assignee changed to "[new]" for [file]`

// Error variables for common failure cases
var (
	ErrTmuxNotInstalled   = errorTmuxNotInstalled{}
	ErrTmuxCommandFailed  = errorTmuxCommandFailed{}
	ErrTmuxTimeout        = errorTmuxTimeout{}
)

// errorTmuxNotInstalled is returned when tmux is not found in PATH
type errorTmuxNotInstalled struct{}

func (e errorTmuxNotInstalled) Error() string {
	return "tmux not installed"
}

// errorTmuxCommandFailed is returned when tmux command returns non-zero exit code
type errorTmuxCommandFailed struct{}

func (e errorTmuxCommandFailed) Error() string {
	return "tmux command failed"
}

// errorTmuxTimeout is returned when command execution exceeds timeout
type errorTmuxTimeout struct{}

func (e errorTmuxTimeout) Error() string {
	return "tmux command timed out"
}
