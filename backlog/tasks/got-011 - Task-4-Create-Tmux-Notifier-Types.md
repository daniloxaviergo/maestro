---
id: GOT-011
title: 'Task 4: Create Tmux Notifier Types'
status: Done
assignee: []
created_date: '2026-03-15'
updated_date: '2026-03-16 17:29'
labels:
  - tmux
  - notifier
  - go
dependencies:
  - GOT-010
references:
  - backlog/docs/doc-003 - PRD-Maestro-Feature-Set-1.md
priority: high
ordinal: 9000
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
- [x] #1 `pkg/notifier/types.go` created with `Notifier`, `NotificationConfig`, `AssigneeChangeEvent` structs
- [x] #2 Error variables defined: `ErrTmuxNotInstalled`, `ErrTmuxCommandFailed`, `ErrTmuxTimeout`
- [x] #3 `NotificationConfig` includes `MessageFormat` (string) and `Timeout` (time.Duration) fields
- [x] #4 Default `MessageFormat` matches format: `Assignee changed to "[new]" for [file]`
- [x] #5 `AssigneeChangeEvent` includes `FilePath`, `OldAssignee`, `NewAssignee` fields
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create `pkg/notifier/types.go` with all type definitions for the tmux notifier package. This is the foundational types task for the tmux notifier system.

**Implementation Steps:**
1. Create `pkg/notifier/` directory
2. Create `pkg/notifier/types.go` with:
   - `NotificationConfig` struct (MessageFormat, Timeout fields)
   - `Notifier` struct (holds config)
   - `AssigneeChangeEvent` struct (FilePath, OldAssignee, NewAssignee fields)
   - Error variables: `ErrTmuxNotInstalled`, `ErrTmuxCommandFailed`, `ErrTmuxTimeout`
3. Define default values for `NotificationConfig`

**Design Decisions:**
- `Notifier` holds `NotificationConfig` for easy configuration and extensibility
- `NotificationConfig` allows customization of message format and timeout per notifier instance
- `AssigneeChangeEvent` mirrors the change detector's output format for seamless integration
- Errors are package-level variables using `errors.New()` for easy comparison with `errors.Is()`
- Follow existing patterns from `pkg/watcher/` (error variables) and `pkg/cache/` (struct definitions)

**Why This Approach:**
- Simple, minimal types that meet the PRD requirements
- Non-blocking by design (implementation in Task 5 will use goroutines)
- Extensible for future features (can add Priority, Metadata fields later)
- Aligns with existing package structure and patterns

### 2. Files to Create/Modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `pkg/notifier/types.go` | Core type definitions for tmux notifier |
| Create | `pkg/notifier/notifier.go` | (Task 5: Implementation) Tmux command execution |
| Create | `pkg/notifier/notifier_test.go` | (Task 5: Testing) Unit tests for notifier |
| Modify | `pkg/change_detect/detector.go` | (Task 6: Integration) Add notifier callback |

**Files to Reference (Read-Only):**
- `pkg/cache/types.go` - Example of type definitions with sync patterns
- `pkg/cache/cache.go` - Example of struct-based cache with mutex
- `pkg/watcher/watcher.go` - Example of error variables and struct definitions
- `pkg/logs/logger.go` - Example of Logger struct with file handling
- `pkg/parser/types.go` - Example of data structures for file data

### 3. Dependencies

- **Go standard library**: `errors`, `time`, `sync` (optional, for future extensibility)
- **Prerequisite**: GOT-010 (change detection) - provides event data and defines `AssigneeChangeEvent` structure
- **No external dependencies** - uses only standard library

**Configuration Requirements:**
- No environment variables required for types
- Default timeout: 2 seconds (for command execution in Task 5)
- Default message format: `Assignee changed to "[new]" for [file]`

### 4. Code Patterns

**Follow existing project conventions:**
- Package name: `notifier` (lowercase)
- Type names: PascalCase (`NotificationConfig`, `Notifier`, `AssigneeChangeEvent`)
- Error variables: `Err` prefix (`ErrTmuxNotInstalled`, etc.)
- Field names: PascalCase for exported, camelCase for unexported (future-proofing)

**Structure Pattern (mirrors `pkg/cache/` and `pkg/watcher/`):**
```go
package notifier

import (
    "errors"
    "time"
)

// Configuration struct
type NotificationConfig struct {
    MessageFormat string
    Timeout       time.Duration
}

// Notifier struct holding state
type Notifier struct {
    config NotificationConfig
}

// Event struct for change data
type AssigneeChangeEvent struct {
    FilePath    string
    OldAssignee []string
    NewAssignee []string
}

// Error variables
var (
    ErrTmuxNotInstalled   = errors.New("tmux not installed")
    ErrTmuxCommandFailed  = errors.New("tmux command failed")
    ErrTmuxTimeout        = errors.New("tmux command timed out")
)
```

**Naming Conventions:**
- Exported types: `NotificationConfig`, `Notifier`, `AssigneeChangeEvent`
- Unexported fields (if added later): `config`, `mutex`
- Error variables: `ErrTmuxNotInstalled`, `ErrTmuxCommandFailed`, `ErrTmuxTimeout`

### 5. Testing Strategy

**Unit tests in `pkg/notifier/notifier_test.go`:**
1. Test `NotificationConfig` struct initialization with defaults
2. Test `NotificationConfig` with custom values (MessageFormat, Timeout)
3. Test `AssigneeChangeEvent` struct initialization
4. Verify error variables are unique (not equal to each other)
5. Test `Notifier` constructor with default config
6. Test `Notifier` constructor with custom config

**Test Scenarios:**
- Empty config uses defaults
- Custom config values are stored correctly
- Event with empty assignees works
- Event with multiple assignees works
- All errors are distinct

**Verification Commands:**
```bash
go build ./pkg/notifier/...
go vet ./pkg/notifier/...
go test ./pkg/notifier/... -v
```

### 6. Risks and Considerations

**Blocking Issues:**
- None - this is a standalone types definition task

**Trade-offs:**
- `Timeout` is per-`Notify()` call (not cumulative for multiple tmux calls in future)
- `MessageFormat` is a simple string template (no complex templating engine like text/template)
- All fields in `AssigneeChangeEvent` are exported for flexibility (could be made unexported if needed)

**Design Decisions:**
- No mutex in `Notifier` yet (Task 5 will handle synchronization for command execution)
- No logger in `Notifier` (Task 5 will handle logging of errors)
- `NotificationConfig` is passed by value (immutable, thread-safe for reading)

**Future Considerations:**
- Could add `Priority` field to `AssigneeChangeEvent` for different notification behaviors
- Could add `Metadata` map[string]string for extensibility without breaking changes
- Could use `sync/atomic` for counters if needed
- Template engine (text/template) could be added later for complex formatting

**Integration Notes:**
- Task 5 (`notifier.go`) will implement the actual tmux command execution
- Task 6 (integration) will wire the notifier to the change detector
- This task only defines the types - no implementation beyond type declarations
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
> Types implementation complete. All acceptance criteria and Definition of Done items satisfied.

**Implementation Summary**:
- Created `pkg/notifier/types.go` with `Notifier`, `NotificationConfig`, `AssigneeChangeEvent` structs
- Defined error variables: `ErrTmuxNotInstalled`, `ErrTmuxCommandFailed`, `ErrTmuxTimeout`
- Created `pkg/notifier/notifier_test.go` with 8 unit tests
- All tests pass: `go test ./pkg/notifier/...` (8 tests)
- Build passes: `go build ./pkg/notifier/...`
- Vet passes: `go vet ./pkg/notifier/...`
- No external dependencies required

**Key Design Decisions**:
- Used custom error types with `Error()` method instead of `errors.New()` for better error categorization
- `NotificationConfig.Timeout` defaults to 2 seconds per requirements
- `DefaultMessageFormat` constant matches the required format: `Assignee changed to "[new]" for [file]`
- `AssigneeChangeEvent` mirrors the change detection output format
- All fields exported for flexibility (can be made unexported if needed in future)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
This task creates the foundational type definitions for the Tmux Notifier package (`pkg/notifier/`).

## What Changed

### New Files Created
- `pkg/notifier/types.go` - Core type definitions for tmux notifier (70 lines)
- `pkg/notifier/notifier_test.go` - Unit tests for type definitions (170 lines)

### Implementation Details

**Type Definitions:**
- `NotificationConfig` - Configuration struct with `MessageFormat` (string) and `Timeout` (time.Duration)
- `Notifier` - Main struct holding configuration state
- `AssigneeChangeEvent` - Event struct with `FilePath`, `OldAssignee`, `NewAssignee` fields

**Error Variables:**
- `ErrTmuxNotInstalled` - tmux not found in PATH
- `ErrTmuxCommandFailed` - tmux command returned non-zero exit code
- `ErrTmuxTimeout` - command execution exceeded timeout

**Configuration:**
- `DefaultMessageFormat` constant: `Assignee changed to "[new]" for [file]`
- Default `Timeout`: 2 seconds

## Testing Results
- All 8 unit tests pass
- `go vet ./pkg/notifier/...` - no warnings
- `go build ./pkg/notifier/...` - successful
- `go vet ./...` and `go build ./...` - no regressions

## Acceptance Criteria Status
- [x] #1 `pkg/notifier/types.go` created with required structs
- [x] #2 Error variables defined
- [x] #3 `NotificationConfig` includes `MessageFormat` and `Timeout` fields
- [x] #4 Default `MessageFormat` matches required format
- [x] #5 `AssigneeChangeEvent` includes required fields

## Definition of Done Status
- [x] #1 `go build ./pkg/notifier/...` passes
- [x] #2 `go vet ./pkg/notifier/...` passes with no issues
- [x] #3 `go test ./pkg/notifier/...` passes (8 unit tests)
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 #1 go build ./pkg/notifier/... passes
- [x] #2 #2 go vet ./pkg/notifier/... passes with no issues
- [x] #3 #3 go test ./pkg/notifier/... passes (unit tests for type definitions)
<!-- DOD:END -->
