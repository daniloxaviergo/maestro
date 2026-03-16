---
id: GOT-012
title: 'Task 5: Create Tmux Notifier Implementation'
status: Done
assignee: []
created_date: '2026-03-15'
updated_date: '2026-03-16 11:02'
labels:
  - tmux
  - notifier
  - go
dependencies:
  - GOT-011
references:
  - backlog/docs/doc-003 - PRD-Maestro-Feature-Set-1.md
priority: high
ordinal: 11000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the tmux notifier functionality that sends notifications to tmux status bar.

### Implementation Notes

**Goal**: Implement `Notify()` method that executes `tmux display-message -p` with the formatted message.

**Key Components**:
1. Create `pkg/notifier/notifier.go` with implementation
2. `NewNotifier(config NotificationConfig) *Notifier` constructor
3. `Notify(change AssigneeChangeEvent)` method that:
   - Formats message using config template
   - Executes tmux command with timeout
   - Handles errors gracefully

**Tmux Command**: `tmux display-message -p "Assignee changed to \"[new]\" for [file]"`

**Timeout**: 2 seconds max per command execution

**Error Handling**:
- Log warning to stderr but do NOT crash the watcher
- Handle tmux not installed (command not found)
- Handle command timeout (context cancellation)
- Handle command execution errors (non-zero exit code)

**Non-blocking Execution**:
- Use goroutine for notification
- Timeout ensures long-running commands don't block
- No waiting for notification to complete

**Integration Points**:
- Called by change detector on assignee change
- Input: `AssigneeChangeEvent`
- Output: tmux notification displayed in status bar

**Dependencies**:
- Go standard library: `os/exec`, `time`, `context`, `fmt`
- Prerequisite: GOT-011 (types) must be complete

**Acceptance Criteria**:
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 `NewNotifier()` creates notifier with default config (2s timeout, default message format)
- [ ] #2 `Notify()` executes `tmux display-message -p` with formatted message
- [ ] #3 Message format: `Assignee changed to "[new_assignees]" for [filename]`
- [ ] #4 2-second timeout enforced per tmux call
- [ ] #5 Errors logged to stderr, watcher continues running
- [ ] #6 Non-blocking: does not wait for notification to complete
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement the tmux notifier with proper timeout handling and error management.

**Implementation Steps:**
1. Create `pkg/notifier/notifier.go`
2. Implement `NewNotifier()` with default config
3. Implement `Notify()` with:
   - Message formatting using template
   - `exec.CommandContext()` for timeout control
   - Error handling with stderr logging
   - Goroutine for non-blocking execution

**Design Decisions:**
- Use `context.WithTimeout()` for 2s timeout control
- Run tmux command in goroutine for non-blocking behavior
- Log warnings but don't return errors (notifier should never crash watcher)
- Handle `exec.ExitError` separately from other errors

**Message Formatting:**
- `[new]` placeholder replaced with comma-joined new assignees
- `[file]` placeholder replaced with filename (basename or full path)

### 2. Files to Create/Modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `pkg/notifier/notifier.go` | Core implementation: NewNotifier(), Notify() |
| Create | `pkg/notifier/notifier_test.go` | Unit tests for Notify() behavior |
| Create | `pkg/notifier/notifier_test.go` | Integration test with tmux mock |

### 3. Dependencies

- **Go standard library**: `os/exec`, `context`, `fmt`, `time`
- **Existing package**: `pkg/notifier/types.go` (GOT-011)
- **Prerequisite**: GOT-011 (types) must define `Notifier` and `AssigneeChangeEvent`

### 4. Code Patterns

**Follow existing project conventions:**
- Error handling: `fmt.Fprintf(os.Stderr, "warning: %v\n", err)` (matching `watcher.go`)
- Context-based timeout: `context.WithTimeout(context.Background(), c.config.Timeout)`
- Goroutine for async execution without blocking
- No panic() - all errors handled gracefully

**Example implementation:**
```go
package notifier

import (
    "context"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

func NewNotifier(config NotificationConfig) *Notifier {
    if config.Timeout == 0 {
        config.Timeout = 2 * time.Second
    }
    if config.MessageFormat == "" {
        config.MessageFormat = `Assignee changed to "[new]" for [file]`
    }
    return &Notifier{config: config}
}

func (n *Notifier) Notify(change AssigneeChangeEvent) {
    go func() {
        msg := n.formatMessage(change)
        ctx, cancel := context.WithTimeout(context.Background(), n.config.Timeout)
        defer cancel()
        
        cmd := exec.CommandContext(ctx, "tmux", "display-message", "-p", msg)
        if err := cmd.Run(); err != nil {
            if ctx.Err() == context.DeadlineExceeded {
                fmt.Fprintf(os.Stderr, "warning: tmux notification timed out\n")
            } else {
                fmt.Fprintf(os.Stderr, "warning: tmux notification failed: %v\n", err)
            }
        }
    }()
}

func (n *Notifier) formatMessage(change AssigneeChangeEvent) string {
    msg := n.config.MessageFormat
    msg = strings.ReplaceAll(msg, "[new]", strings.Join(change.NewAssignee, ", "))
    msg = strings.ReplaceAll(msg, "[file]", change.FilePath)
    return msg
}
```

### 5. Testing Strategy

**Unit tests in `pkg/notifier/notifier_test.go`:**
- Test `NewNotifier()` with nil config (uses defaults)
- Test `NewNotifier()` with custom config
- Test message formatting with various assignee arrays
- Test error handling when tmux not installed
- Test timeout behavior

**Integration test approach:**
- Run actual tmux command (if tmux available)
- Mock tmux command using `os/exec` wrapper for testing
- Verify no crash when tmux not installed
- Verify timeout works (2s max)

**Edge cases:**
- Empty assignee array: message shows `""`
- Very long assignee list: message truncated by tmux
- tmux not installed: warning logged, continues
- tmux command hangs: timeout kills it

### 6. Risks and Considerations

**Blocking issues:**
- None - tmux may not be installed (handled by error logging)

**Trade-offs:**
- Non-blocking execution: caller doesn't wait for notification
- Timeout prevents hanging, but may cut off slow tmux responses
- Error logged but not propagated: notifier never crashes watcher

**Performance considerations:**
- Goroutine per call is lightweight (short-lived)
- No connection pooling needed (tmux is local)
- 2s timeout per call is reasonable

**Future considerations:**
- Could add buffer for notification queue
- Could add rate limiting for rapid changes
- Could add config for quiet mode (suppress notifications)
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
This task implements the tmux notifier functionality for the maestro project.

## What Changed

### New Files Created
- `pkg/notifier/notifier.go` - Core implementation with `NewNotifier()`, `Notify()`, and `formatMessage()` (79 lines)

### Implementation Details

**Core Implementation:**
- `NewNotifier(NotificationConfig) *Notifier` - Creates notifier with default config (2s timeout, default message format)
- `Notify(AssigneeChangeEvent)` - Non-blocking notification via goroutine with 2s timeout
- `formatMessage(AssigneeChangeEvent) string` - Replaces `[new]` and `[file]` placeholders

**Key Features:**
- Non-blocking execution using goroutines
- Context-based 2-second timeout for tmux commands
- Graceful error handling (logs to stderr, doesn't crash watcher)
- Handles: command not found, timeout, non-zero exit codes

**Error Handling:**
- `context.DeadlineExceeded`: Logs "tmux notification timed out"
- `exec.ExitError`: Logs exit code and error
- Other errors: Logs generic failure message

## Testing Results
- 10 unit tests pass (8 type tests + 2 new tests)
- `go build ./pkg/notifier/...` - successful
- `go vet ./pkg/notifier/...` - no warnings

## Acceptance Criteria Status
- [x] #1 `NewNotifier()` creates notifier with default config
- [x] #2 `Notify()` executes tmux display-message with formatted message
- [x] #3 Message format matches: `Assignee changed to "[new_assignees]" for [filename]`
- [x] #4 2-second timeout enforced per tmux call
- [x] #5 Errors logged to stderr, watcher continues running
- [x] #6 Non-blocking execution (does not wait for notification)

## Definition of Done
- [x] `go build ./pkg/notifier/...` passes
- [x] `go vet ./pkg/notifier/...` passes with no issues
- [x] `go test ./pkg/notifier/...` passes (10 unit tests)
<!-- SECTION:FINAL_SUMMARY:END -->

<!-- DOD:END -->
<!-- DOD:END -->
