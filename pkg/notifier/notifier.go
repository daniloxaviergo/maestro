package notifier

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// NewNotifier creates a new tmux notifier with the given config
func NewNotifier(config NotificationConfig) *Notifier {
	if config.Timeout == 0 {
		config.Timeout = 2 * time.Second
	}
	if config.MessageFormat == "" {
		config.MessageFormat = DefaultMessageFormat
	}
	return &Notifier{config: config}
}

// Notify sends a tmux notification for the given assignee change event.
// This method is non-blocking - it executes the tmux command in a goroutine.
func (n *Notifier) Notify(change AssigneeChangeEvent) {
	go func() {
		msg := n.formatMessage(change)
		ctx, cancel := context.WithTimeout(context.Background(), n.config.Timeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, "tmux", "display-message", "-p", msg)
		if err := cmd.Run(); err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Fprintf(os.Stderr, "warning: tmux notification timed out\n")
			} else if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Fprintf(os.Stderr, "warning: tmux notification failed with exit code %d: %v\n", exitErr.ExitCode(), err)
			} else {
				fmt.Fprintf(os.Stderr, "warning: tmux notification failed: %v\n", err)
			}
		}
	}()
}

// formatMessage formats the notification message using the configured template
func (n *Notifier) formatMessage(change AssigneeChangeEvent) string {
	msg := n.config.MessageFormat
	msg = strings.ReplaceAll(msg, "[new]", strings.Join(change.NewAssignee, ", "))
	msg = strings.ReplaceAll(msg, "[file]", change.FilePath)
	return msg
}
