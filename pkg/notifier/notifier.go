package notifier

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"maestro/pkg/agent"
	"maestro/pkg/config"
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

// ExecuteScript executes a bash script in the configured tmux session.
// This method is non-blocking - it executes the script in a goroutine.
// The script is run via tmux send-keys and output is captured but not displayed.
func (n *Notifier) ExecuteScript(filePath string) {
	go func() {
		// Check if agent is configured
		if n.config.Agent == nil {
			fmt.Fprintf(os.Stderr, "warning: agent not configured for ExecuteScript\n")
			return
		}

		agent := n.config.Agent
		cfg := agent.GetConfig()

		// Check if script path is configured
		if cfg.ScriptPath == "" {
			fmt.Fprintf(os.Stderr, "warning: script_path not configured for agent %s\n", agent.GetName())
			return
		}

		// Check if script file exists
		if _, err := os.Stat(cfg.ScriptPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: %v: %s\n", ErrScriptNotFound, cfg.ScriptPath)
			return
		}

		// Ensure tmux session exists
		sessionName := cfg.TmuxSession
		if sessionName == "" {
			sessionName = "default"
		}

		ctx, cancel := context.WithTimeout(context.Background(), n.config.Timeout)
		defer cancel()

		// Create session if it doesn't exist
		createCmd := exec.CommandContext(ctx, "tmux", "new-session", "-d", "-s", sessionName)
		if err := createCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v: %v\n", ErrSessionCreationFailed, err)
			return
		}

		// Execute script via tmux send-keys with file path as argument
		execCmd := exec.CommandContext(ctx, "tmux", "send-keys", "-t", sessionName, fmt.Sprintf("bash %s %s", cfg.ScriptPath, filePath), "Enter")

		if err := execCmd.Run(); err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Fprintf(os.Stderr, "warning: tmux script execution timed out\n")
			} else if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Fprintf(os.Stderr, "warning: %v: exit code %d\n", ErrScriptExecutionFailed, exitErr.ExitCode())
			} else {
				fmt.Fprintf(os.Stderr, "warning: %v: %v\n", ErrScriptExecutionFailed, err)
			}
			return
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

// ExecuteScriptsForAgents executes scripts for the given agents concurrently.
// This method is non-blocking - each script is executed in a goroutine.
// Disabled agents are skipped. Missing scripts or execution errors are logged as warnings.
func (n *Notifier) ExecuteScriptsForAgents(agents []*agent.Agent, filePath string) {
	for _, a := range agents {
		// Skip disabled agents
		cfg := a.GetConfig()
		if !cfg.Enabled {
			log.Printf("Warning: agent %s is disabled, skipping script execution", a.GetName())
			continue
		}

		// Skip if no script path configured
		if cfg.ScriptPath == "" {
			log.Printf("Warning: agent %s has no script_path configured, skipping", a.GetName())
			continue
		}

		// Execute script for this agent (non-blocking)
		go n.executeScriptForAgent(a, cfg, filePath)
	}
}

// executeScriptForAgent executes a script for a specific agent.
// This is a helper method that contains the actual script execution logic.
func (n *Notifier) executeScriptForAgent(agent *agent.Agent, cfg config.AgentConfig, filePath string) {
	// Check if script file exists
	if _, err := os.Stat(cfg.ScriptPath); os.IsNotExist(err) {
		log.Printf("Warning: %v: %s", ErrScriptNotFound, cfg.ScriptPath)
		return
	}

	// Ensure tmux session exists
	sessionName := cfg.TmuxSession
	if sessionName == "" {
		sessionName = "default"
	}

	ctx, cancel := context.WithTimeout(context.Background(), n.config.Timeout)
	defer cancel()

	// Create session if it doesn't exist
	createCmd := exec.CommandContext(ctx, "tmux", "new-session", "-d", "-s", sessionName)
	if err := createCmd.Run(); err != nil {
		log.Printf("Warning: %v: %v", ErrSessionCreationFailed, err)
		return
	}

	// Execute script via tmux send-keys with file path as argument
	execCmd := exec.CommandContext(ctx, "tmux", "send-keys", "-t", sessionName, fmt.Sprintf("bash %s %s", cfg.ScriptPath, filePath), "Enter")

	if err := execCmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Warning: tmux script execution timed out for agent %s", agent.GetName())
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Warning: %v: exit code %d for agent %s", ErrScriptExecutionFailed, exitErr.ExitCode(), agent.GetName())
		} else {
			log.Printf("Warning: %v: %v for agent %s", ErrScriptExecutionFailed, err, agent.GetName())
		}
	}
}
