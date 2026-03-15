---
id: GOT-011
title: 'Task 4: Create Tmux Notifier Types'
status: To Do
assignee: []
created_date: '2026-03-15'
updated_date: '2026-03-15'
labels:
  - tmux
  - notifier
  - go
dependencies:
  - GOT-010
references:
  - backlog/docs/doc-003 - PRD-Maestro-Feature-Set-1.md
priority: high
ordinal: 1100
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the data types and configuration structures for the Tmux Notifier package (`pkg/notifier/`).

### Implementation Notes

**Goal**: Define the core data structures for the tmux notifier system.

**Key Components**:
1. Create `pkg/notifier/types.go` with type definitions
2. Define `Notifier` struct to hold configuration and state
3. Define `NotificationConfig` struct for customizable message format
4. Define error variables for common failure cases

**Data Structures**:
```go
type Notifier struct {
    config NotificationConfig
}

type NotificationConfig struct {
    MessageFormat string  // Template for display message
    Timeout       time.Duration
}

type AssigneeChangeEvent struct {
    FilePath    string
    OldAssignee []string
    NewAssignee []string
}
```

**Error Variables**:
- `ErrTmuxNotInstalled` - tmux not found in PATH
- `ErrTmuxCommandFailed` - tmux command returned non-zero exit code
- `ErrTmuxTimeout` - command execution exceeded timeout

**Message Format**:
- Default: `Assignee changed to "[new]" for [file]`
- Placeholders: `[new]` for new assignees, `[file]` for filename

**Integration Points**:
- Used by `pkg/change_detect/` to notify on assignee changes
- Input: `AssigneeChangeEvent` from detector
- Output: tmux display-message command execution

**Dependencies**:
- Go standard library: `os/exec`, `time`, `context`
- Prerequisite: GOT-010 (change detection) provides the event data

**Acceptance Criteria**:
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 `pkg/notifier/types.go` created with `Notifier`, `NotificationConfig`, `AssigneeChangeEvent` structs
- [ ] #2 Error variables defined: `ErrTmuxNotInstalled`, `ErrTmuxCommandFailed`, `ErrTmuxTimeout`
- [ ] #3 `NotificationConfig` includes `MessageFormat` (string) and `Timeout` (time.Duration) fields
- [ ] #4 Default `MessageFormat` matches format: `Assignee changed to "[new]" for [file]`
- [ ] #5 `AssigneeChangeEvent` includes `FilePath`, `OldAssignee`, `NewAssignee` fields
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create `pkg/notifier/types.go` with all necessary type definitions for the tmux notifier package.

**Implementation Steps:**
1. Create `pkg/notifier/` directory
2. Create `types.go` with struct definitions
3. Define error variables using `errors.New()`
4. Define default configuration values

**Design Decisions:**
- `Notifier` holds `NotificationConfig` for easy configuration
- `NotificationConfig` allows customization of message format and timeout
- `AssigneeChangeEvent` mirrors the change detector's output format
- Errors are package-level variables for easy comparison with `errors.Is()`

### 2. Files to Create/Modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `pkg/notifier/types.go` | Type definitions: Notifier, NotificationConfig, AssigneeChangeEvent, errors |
| Create | `pkg/notifier/notifier.go` | (Task 2) Implementation of Notify() method |
| Create | `pkg/notifier/notifier_test.go` | Unit tests for type definitions |

### 3. Dependencies

- **Go standard library**: `errors`, `time`
- **Prerequisite**: GOT-010 (change detection) - provides `AssigneeChangeEvent` data

### 4. Code Patterns

**Follow existing project conventions:**
- Package structure mirrors `pkg/watcher/` and `pkg/cache/`
- Error variables use `Err` prefix (matching `watcher.go`)
- Exported types use PascalCase
- Unexported types use lowercase

**Example code:**
```go
package notifier

import (
    "errors"
    "time"
)

var (
    ErrTmuxNotInstalled   = errors.New("tmux not installed")
    ErrTmuxCommandFailed  = errors.New("tmux command failed")
    ErrTmuxTimeout        = errors.New("tmux command timed out")
)

type NotificationConfig struct {
    MessageFormat string
    Timeout       time.Duration
}

type AssigneeChangeEvent struct {
    FilePath    string
    OldAssignee []string
    NewAssignee []string
}

type Notifier struct {
    config NotificationConfig
}
```

### 5. Testing Strategy

**Unit tests in `pkg/notifier/notifier_test.go`:**
- Test `NotificationConfig` default values
- Test `AssigneeChangeEvent` struct initialization
- Verify error variables are unique
- Test config with custom message format

**Verification:**
- `go build ./pkg/notifier/...`
- `go test ./pkg/notifier/...`
- `go vet ./pkg/notifier/...`

### 6. Risks and Considerations

**Blocking issues:**
- None - standalone types task

**Trade-offs:**
- `Timeout` is per-call (not cumulative for multiple tmux calls)
- `MessageFormat` is simple string template (no complex templating engine)
- `AssigneeChangeEvent` fields are all exported for flexibility

**Future considerations:**
- Could add `Priority` field to event for different notification behaviors
- Could add `Metadata` map for extensibility
<!-- SECTION:PLAN:END -->
