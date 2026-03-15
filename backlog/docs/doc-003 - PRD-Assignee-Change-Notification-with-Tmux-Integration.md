---
id: doc-003
title: 'PRD: Assignee Change Notification with Tmux Integration'
type: other
created_date: '2026-03-15 10:52'
updated_date: '2026-03-15 11:49'
---
# PRD: Assignee Change Notification with Tmux Integration

## Overview

### Purpose
Implement real-time notification system that detects when the `assignee` field changes in markdown task files and triggers a tmux status-line message to notify the current tmux session user.

### Goals
- **G1**: Detect assignee field modifications within 500ms of file write events
- **G2**: Display clear, actionable notifications in tmux status line
- **G3**: Log all assignee change events with full context for audit trail
- **G4**: Support both single and multiple assignee scenarios
- **G5**: Handle file creation, modification, and deletion events gracefully

## Background

### Problem Statement
Currently, when task assignees are changed in the backlog system, there is no automated notification to alert users in real-time. Team members must manually check files or rely on external notifications (email, chat) to know when ownership changes occur. This creates visibility gaps and reduces responsiveness.

### Current State
- **FR1 IMPLEMENTED**: File watcher using `pkg/watcher/` with fsnotify, 50ms debouncing, recursive directory watching for `.md` files
- **FR2 IMPLEMENTED**: YAML frontmatter parsing in `pkg/parser/` extracts `assignee` field, handles missing frontmatter gracefully
- **FR5 IMPLEMENTED**: Event logging in `pkg/logs/` writes JSON to `./assignee_changes.log` (proven in logs)
- **FR3 PARTIAL**: `pkg/change_detect/` exists and compares cached vs new assignees
- **FR4 MISSING**: NO tmux integration exists - this is the only remaining component

### Proposed Solution
Extend the existing file watcher flow to add tmux notification:
1. Parser extracts `assignee` field from YAML frontmatter on WRITE/CREATE events
2. Change detector compares new assignee with cached previous value
3. If assignee changed, trigger tmux notification via `tmux display-message`
4. Event is logged to `./assignee_changes.log`

**Existing Flow (Partial):**
```
File Change → Watcher Event → Parser → Change Detector → Logger (JSON file)
```

**Target Flow (Complete):**
```
File Change → Watcher Event → Parser → Change Detector → Tmux Notifier → Logger
```

## Requirements

### User Stories

- **Role**: Project Manager
  - *As a project manager, I want to receive immediate notification when task assignees are changed so that I can track ownership transitions in real-time*

- **Role**: Developer
  - *As a developer working on a task, I want to see when someone else is assigned to my task so I can coordinate or transfer work*

- **Role**: Team Lead
  - *As a team lead, I want a complete log of all assignee changes with timestamps so I can review assignment history and investigate issues*

### Functional Requirements

#### FR1: File Watch Integration (ALREADY IMPLEMENTED)

The existing `pkg/watcher/` package handles file monitoring.

##### Acceptance Criteria (Already Met)
- [x] Reuses existing `pkg/watcher/` package for file monitoring
- [x] Listens to WRITE events on `.md` files in `./backlog/tasks/`
- [x] Handles recursive monitoring of subdirectories
- [x] Properly debounces rapid file writes (50ms cooldown per file)
- [x] Gracefully handles permission errors and missing files

##### Implementation Reference
- `pkg/watcher/watcher.go` - File watcher with debouncing
- `pkg/watcher/events.go` - Event types: CREATE/WRITE/REMOVE/RENAME

#### FR2: Assignee Field Parsing (ALREADY IMPLEMENTED)

YAML frontmatter parsing is already implemented in `pkg/parser/`.

##### Acceptance Criteria (Already Met)
- [x] Parses YAML frontmatter to extract `assignee` field
- [x] Handles `assignee: []` (empty array) gracefully
- [x] Handles `assignee:` (empty value) gracefully
- [x] Handles files without frontmatter (treats as empty assignee)
- [x] Handles malformed YAML with error logging
- [x] Supports single assignee: `assignee: ["alice"]`
- [x] Supports multiple assignees: `assignee: ["alice", "bob"]`

##### Implementation Reference
- `pkg/parser/parser.go` - YAML frontmatter extraction
- `pkg/parser/types.go` - `Frontmatter{Assignee []string}`, `FileData{FilePath, Frontmatter, Error, ParseTime}`

#### FR3: Change Detection (PARTIALLY IMPLEMENTED)

`pkg/change_detect/` exists and compares cached vs new assignees.

##### Acceptance Criteria (Partial Implementation)
- [x] Caches assignee value per file in memory (map: filepath -> assignee array)
- [x] Compares old vs new assignee arrays for equality
- [ ] Handle first-time parsing (no cached value = treat as all new assignees)
- [x] Handles empty assignee (no one assigned)
- [x] Handles removal of assignees (e.g., "alice" removed)
- [x] Handles addition of assignees (e.g., "bob" added)
- [x] Handles replacement (e.g., "alice" replaced by "bob")

##### Implementation Reference
- `pkg/change_detect/detector.go` - Detects changes, returns true if log was written

#### FR4: Tmux Notification (MISSING - MUST ADD)

Trigger tmux status-line message on assignee changes.

##### Acceptance Criteria (NEW - TO BE IMPLEMENTED)
- [ ] Execute `tmux display-message` command with notification
- [ ] Format: `Assignee changed to "[new_assignees]" for [filename]`
- [ ] Handle multiple assignees: `Assignee changed to "alice, bob" for task-001.md`
- [ ] Handle empty assignee: `Assignee changed to "none" for task-001.md`
- [ ] Command runs asynchronously (non-blocking)
- [ ] Command errors are logged but don't crash the watcher
- [ ] Support missing tmux gracefully (log warning, continue)

##### Implementation Notes
- Create `pkg/notifier/` or `pkg/assignee/` package for tmux integration
- Integrate with existing `change_detect.Detector` or watcher event flow
- Use `os/exec` to run `tmux display-message -p "message"` command
- Add timeout (e.g., 2s) to prevent hanging if tmux is unresponsive

#### FR5: Event Logging (ALREADY IMPLEMENTED)

JSON logging to `./assignee_changes.log` is already working.

##### Acceptance Criteria (Already Met)
- [x] Logs to `./assignee_changes.log`
- [x] Log format: JSON with timestamp, file, old_assignee, new_assignee
- [x] Timestamp in ISO 8601 format with timezone
- [x] Array values serialized as JSON arrays
- [x] Log write is non-blocking or buffered
- [x] Handles log file creation if directory/file doesn't exist
- [x] Handles log write errors gracefully

##### Implementation Reference
- `pkg/logs/logger.go` - JSON log writer
- Log format example:
```json
{
  "timestamp": "2026-03-15T10:30:00Z",
  "file": "backlog/tasks/GOT-015.md",
  "old_assignee": ["alice"],
  "new_assignee": ["bob"]
}
```

### Non-Functional Requirements

- **Performance**:
  - Change detection latency: <500ms (95th percentile) - **already achieved**
  - File parse time: <100ms per file - **already achieved**
  - Log write time: <50ms per entry - **already achieved**
  - Tmux command execution: <100ms (to be verified)
  - Memory usage: <50MB for typical task count (100-500 files)

- **Reliability**:
  - System should recover from file system errors without crashing
  - Cache starts fresh on each run (no persistence across runs)
  - Handle file system events that may coalesce
  - Handle tmux command failures gracefully

- **Maintainability**:
  - Code follows Go 1.20+ best practices
  - Package structure follows existing Maestro conventions
  - Logging includes appropriate debug/error levels
  - Clear separation of concerns (parsing, detection, notification, logging)

- **Compatibility**:
  - Go 1.20 or later
  - Linux, macOS, and Windows platforms
  - tmux 2.0+ (minimum for `display-message` command)
  - Markdown files in YAML frontmatter format

## Scope

### In Scope
- Tmux notification integration only (FR4 - the missing component)
- Integration with existing `change_detect` package
- Graceful error handling when tmux is unavailable
- Support for existing file watcher and parser components
- All existing functionality preserved

### Out of Scope
- Refactoring of existing FR1-FR5 components
- Email/SMS notifications (tmux only)
- Integration with external systems (webhooks, databases)
- Real-time notifications to users via other channels
- History replay (starting fresh on each run)
- Monitoring of other fields beyond `assignee`
- User interface for viewing logs or history
- Persistent cache across runs
- File rotation or cleanup for logs

## Technical Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     File Watcher (fsnotify)                 │
│                     ┌───────────────────────┐               │
│                     │  WRITE event for .md  │               │
│                     └─────────────┬─────────┘               │
│                                   │                         │
│                                   ▼                         │
│                    ┌────────────────────────┐               │
│                    │     Parser             │               │
│                    │  Extract assignee      │               │
│                    └─────────────┬─────────┘               │
│                                   │                         │
│                                   ▼                         │
│              ┌───────────────────────────────┐              │
│              │     Change Detection          │              │
│              │  Compare with cached value    │              │
│              └───────▲─────┬──────▲──────────┘              │
│                      │     │     │                          │
│           ┌──────────┴─┐  │     └─┐  ┌────────────────┐     │
│           ▼            │  ▼       │  │  Update Cache  │     │
│   ┌──────────────┐     │  ▼       │  └────────────────┘     │
│   │   Tmux Notify│     │  Log     │                         │
│   │              │     │  Write   │                         │
│   └──────────────┘     │  ───────►│                         │
│           ▲            │          │                         │
│           └────────────┼──────────┘                         │
│                        │                                    │
└───────────────────────────────────────────────────────────────┘
```

### Package Structure (Current State)

```
pkg/
├── cache/                # File state caching with debouncing
│   ├── types.go
│   └── cache.go
├── change_detect/        # Assignee change detection
│   └── detector.go
├── logs/                 # JSON log writing (FR5 - DONE)
│   └── logger.go
├── parser/               # YAML frontmatter parsing (FR2 - DONE)
│   ├── parser.go
│   └── types.go
└── watcher/              # File watching with fsnotify (FR1 - DONE)
    ├── watcher.go
    └── events.go

# TO BE ADDED (FR4):
├── notifier/             # Tmux notification (NEW)
│   ├── notifier.go
│   └── types.go
```

### Key Data Structures (From Existing Code)

```go
// From pkg/parser/types.go
type Frontmatter struct {
    Assignee []string `yaml:"assignee,omitempty"`
    // ... other fields
}

type FileData struct {
    FilePath  string
    Frontmatter *Frontmatter
    Error     error
    ParseTime time.Time
}

// From pkg/cache/cache.go (simplified)
type AssigneeCache struct {
    mu      sync.RWMutex
    entries map[string]*AssigneeState
}

type AssigneeState struct {
    Assignees []string
    LastChange time.Time
}
```

### Sequence Diagram

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  File Write │  │   Watcher   │  │  Parser     │  │  Detector   │  │  Notifier   │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └───▲───┬─────┘
       │               │               │               │               │   │        │
       │ WRITE event   │               │               │               │   │        │
       │──────────────>│               │               │               │   │        │
       │               │               │               │               │   │        │
       │               │ Parse YAML    │               │               │   │        │
       │               │──────────────>│               │               │   │        │
       │               │               │               │               │   │        │
       │               │ YAML data     │               │               │   │        │
       │               │<──────────────│               │               │   │        │
       │               │               │               │               │   │        │
       │               │ Cache check   │               │               │   │        │
       │               │──────────────>│               │               │   │        │
       │               │               │               │               │   │        │
       │               │ Cache state   │               │               │   │        │
       │               │<──────────────│               │               │   │        │
       │               │               │ Detect change │               │   │        │
       │               │               │──────────────>│               │   │        │
       │               │               │               │               │   │        │
       │               │               │ Change found  │               │   │        │
       │               │               │<──────────────│               │   │        │
       │               │               │               │               │   │        │
       │               │ Notify        │               │               │   │        │
       │               │──────────────>│               │               │   │        │
       │               │               │               │               │   │        │
       │               │               │ Notify tmux   │               │   │        │
       │               │               │──────────────>│               │   │        │
       │               │               │               │               │   │        │
       │               │               │ tmux output   │               │   │        │
       │               │               │<──────────────│               │   │        │
       │               │               │               │               │   │        │
       │               │ Log           │               │               │   │        │
       │               │───────────────────────────────┴───┬──────────>│   │        │
       │               │                                   │           │   │        │
       │               │ Log entry written                 │           │   │        │
       │               │<──────────────────────────────────┴───────────┘   │        │
```

## Implementation Plan

### Phase 1: Tmux Notifier Package (NEW)

1. **Create `pkg/notifier/types.go`**
   - Define `Notifier` struct with tmux command configuration
   - Define `NotificationConfig` for customizable message format
   - Define error variables for tmux-related errors

2. **Create `pkg/notifier/notifier.go`**
   - Implement `NewNotifier()` constructor
   - Implement `Notify(change AssigneeChangeEvent)` method
   - Execute `tmux display-message -p "message"` command
   - Handle command errors gracefully (log warning, continue)
   - Support timeout (2s) to prevent hanging

3. **Integrate with existing detector**
   - Modify `pkg/change_detect/detector.go` to call notifier
   - Pass `AssigneeChangeEvent` to notifier on change detected
   - Ensure non-blocking notifier execution

### Phase 2: Integration and Testing

4. **Update `cmd/monitor/main.go`**
   - Initialize `notifier.Notifier` with default config
   - Pass notifier to change detector or event handler
   - Ensure notifier is called on assignee changes

5. **Write unit tests**
   - Notifier tests (tmux command execution, error handling)
   - Integration tests (full flow from file change to notification)
   - Test with tmux present/absent scenarios

6. **Integration testing**
   - Test with real task files
   - Verify tmux notifications appear in status line
   - Test with tmux not available (should log warning, continue)
   - Test with rapid writes (debouncing still applies)
   - Verify all logs still written correctly

### Implementation Notes

- **Error Handling**: Tmux command errors should NOT crash the watcher. Log warning and continue monitoring.
- **Timeout**: Use `exec.CommandContext` with 2s timeout to prevent hanging.
- **Asynchronous**: Run tmux command in goroutine or non-blocking manner.
- **Fallback**: If tmux is unavailable, log warning: "tmux not available for notification: [error]"

### Files to Modify

| File | Change |
|------|--------|
| `pkg/notifier/types.go` | CREATE - Notifier type definitions |
| `pkg/notifier/notifier.go` | CREATE - Tmux notification implementation |
| `pkg/change_detect/detector.go` | MODIFY - Add notifier callback |
| `cmd/monitor/main.go` | MODIFY - Initialize and wire notifier |
| `go.mod` | UPDATE - No new dependencies needed |

## Success Metrics

### Quantitative
- Change detection latency: <500ms (95th percentile) - **already achieved**
- File parse time: <100ms per file - **already achieved**
- Log write time: <50ms per entry - **already achieved**
- Tmux notification latency: <200ms (from change to tmux display)
- Memory overhead: <5MB for 100 files - **already achieved**

### Qualitative
- Notifications appear clearly in tmux status line
- No false positives or missed changes - **already achieved**
- System recovers gracefully from tmux unavailability
- No crash when tmux is not installed or not running

## Timeline & Milestones

### Key Dates
- **Design complete**: PRD approved, implementation plan reviewed
- **Implementation complete**: FR4 (tmux notifier) implemented, tests passing
- **Integration complete**: Notifier integrated, end-to-end tested
- **Testing complete**: All acceptance criteria verified
- **Launch/Release**: Deploy and run in production environment

## Stakeholders

### Decision Makers
- Product Owner: Approval of PRD scope and requirements

### Contributors
- Backend Engineer: Implementation of tmux notifier
- QA Engineer: Testing and validation of change detection and notification

## Appendix

### Glossary
- **YAML frontmatter**: Metadata section at the top of markdown files, delimited by `---`
- **Assignee field**: The `assignee` key in frontmatter containing an array of user identifiers
- **fsnotify**: A Go library for file system change notifications
- **Debouncing**: Method to delay processing until file writes stabilize (50ms cooldown)
- **Tmux**: Terminal multiplexer with window management and status line

### References
- Maestro project context: `QWEN.md`
- Backlog.md workflow: `backlog://workflow/overview`
- YAML spec: https://yaml.org/spec/
- fsnotify documentation: https://github.com/fsnotify/fsnotify
- Go yaml.v3: https://gopkg.in/yaml.v3
- tmux manual: https://man7.org/linux/man-pages/man1/tmux.1.html

### Related Tasks
- **GOT-008**: File watcher implementation (Done)
- **GOT-009**: YAML frontmatter parser (Done)
- **GOT-010**: Change detection and JSON logging (Done - partial)
- **doc-003**: This PRD - updated to reflect actual codebase state
