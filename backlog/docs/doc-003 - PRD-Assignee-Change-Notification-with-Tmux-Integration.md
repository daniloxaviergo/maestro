---
id: doc-003
title: 'PRD: Assignee Change Notification with Tmux Integration'
type: other
created_date: '2026-03-15 10:52'
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
- Task files in `./backlog/tasks/` store metadata in YAML frontmatter format
- The `assignee` field contains an array of user identifiers (e.g., `assignee: ["username"]`)
- File changes are monitored by the existing Maestro watcher
- YAML frontmatter parsing is already implemented (GOT-009, Done)

### Proposed Solution
Extend the existing file watcher to:
1. Parse the YAML frontmatter on WRITE events to extract the `assignee` field
2. Compare the new assignee value with the cached previous value
3. If the assignee changed, trigger a tmux notification message
4. Log the change event to `./backlog/logs/assignee_changes.log`

## Requirements

### User Stories

- **Role**: Project Manager
  - *As a project manager, I want to receive immediate notification when task assignees are changed so that I can track ownership transitions in real-time*

- **Role**: Developer
  - *As a developer working on a task, I want to see when someone else is assigned to my task so I can coordinate or transfer work*

- **Role**: Team Lead
  - *As a team lead, I want a complete log of all assignee changes with timestamps so I can review assignment history and investigate issues*

### Functional Requirements

#### FR1: File Watch Integration

Integrate assignee change detection into the existing Maestro file watcher.

##### Acceptance Criteria
- [ ] Reuse existing `pkg/watcher/` package for file monitoring
- [ ] Listen to WRITE events on `.md` files in `./backlog/tasks/`
- [ ] Handle recursive monitoring of subdirectories
- [ ] Properly debounce rapid file writes (50ms cooldown per file)
- [ ] Gracefully handle permission errors and missing files

#### FR2: Assignee Field Parsing

Extract and compare the `assignee` field from YAML frontmatter.

##### Acceptance Criteria
- [ ] Parse YAML frontmatter to extract `assignee` field
- [ ] Handle `assignee: []` (empty array) gracefully
- [ ] Handle `assignee:` (empty value) gracefully
- [ ] Handle files without frontmatter (treat as empty assignee)
- [ ] Handle malformed YAML with error logging
- [ ] Support single assignee: `assignee: ["alice"]`
- [ ] Support multiple assignees: `assignee: ["alice", "bob"]`

#### FR3: Change Detection

Compare current assignee with cached previous value and detect changes.

##### Acceptance Criteria
- [ ] Cache assignee value per file in memory (map: filepath -> assignee array)
- [ ] Compare old vs new assignee arrays for equality
- [ ] Handle first-time parsing (no cached value = treat as all new assignees)
- [ ] Handle empty assignee (no one assigned)
- [ ] Handle removal of assignees (e.g., "alice" removed)
- [ ] Handle addition of assignees (e.g., "bob" added)
- [ ] Handle replacement (e.g., "alice" replaced by "bob")

#### FR4: Tmux Notification

Trigger tmux status-line message on assignee changes.

##### Acceptance Criteria
- [ ] Execute `tmux display-message` command with notification
- [ ] Format: `Assignee changed to "[new_assignees]" for [filename]`
- [ ] Handle multiple assignees: `Assignee changed to "alice, bob" for task-001.md`
- [ ] Handle empty assignee: `Assignee changed to "none" for task-001.md`
- [ ] Command runs asynchronously (non-blocking)
- [ ] Command errors are logged but don't crash the watcher
- [ ] Support missing tmux gracefully (log warning, continue)

#### FR5: Event Logging

Log all assignee change events to a dedicated log file.

##### Acceptance Criteria
- [ ] Log to `./backlog/logs/assignee_changes.log`
- [ ] Log format: JSON with timestamp, file, old_assignee, new_assignee
- [ ] Timestamp in ISO 8601 format with timezone
- [ ] Array values serialized as JSON arrays
- [ ] Log write is non-blocking or buffered
- [ ] Handle log file creation if directory/file doesn't exist
- [ ] Handle log write errors gracefully

##### Log Format Example
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
  - Change detection latency: <500ms (95th percentile)
  - File parse time: <100ms per file
  - Log write time: <50ms per entry
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
- File watcher integration (reusing existing `pkg/watcher/`)
- YAML frontmatter parsing (reusing existing `pkg/parser/`)
- Assignee change detection with caching
- Tmux notification via `display-message` command
- JSON log output to `./backlog/logs/assignee_changes.log`
- Support for existing and newly created markdown files
- Debouncing (50ms cooldown per file)
- Graceful error handling and recovery

### Out of Scope
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
│              └───────────────┬───────────────┘              │
│                              │                              │
│           ┌──────────────────┼────────────────┐            │
│           ▼                  ▼                ▼            │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│   │   Tmux Notify│  │  JSON Logger │  │ Update Cache │     │
│   │              │  │              │  │              │     │
│   └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Package Structure

```
pkg/
├── assignee/
│   ├── types.go          # Data types for assignee events
│   ├── parser.go         # YAML frontmatter parsing
│   ├── cache.go          # In-memory assignee cache
│   ├── detector.go       # Change detection logic
│   ├── notifier.go       # Tmux notification
│   └── logger.go         # JSON log file writing
└── watcher/
    └── events.go         # Extended event types
```

### Key Data Structures

```go
// AssigneeCache caches assignee values per file path
type AssigneeCache struct {
    mu      sync.RWMutex
    entries map[string]*AssigneeState
}

type AssigneeState struct {
    Filepath   string
    Assignees  []string
    LastChange time.Time
}

// AssigneeChangeEvent represents a detected change
type AssigneeChangeEvent struct {
    Timestamp   time.Time
    Filepath    string
    OldAssignees []string
    NewAssignees []string
}
```

### Sequence Diagram

```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  File Write │  │   Watcher   │  │  Parser     │  │  Detector   │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │               │               │               │
       │ WRITE event   │               │               │
       │──────────────>│               │               │
       │               │               │               │
       │               │ Parse YAML    │               │
       │               │──────────────>│               │
       │               │               │               │
       │               │ YAML data     │               │
       │               │<──────────────│               │
       │               │               │               │
       │               │ Cache miss    │               │
       │               │──────────────>│               │
       │               │               │               │
       │               │ Cache updated │               │
       │               │<──────────────│               │
       │               │               │               │
       │               │               │ Detect change │
       │               │               │──────────────>│
       │               │               │               │
       │               │               │ Change found  │
       │               │               │<──────────────│
       │               │               │               │
       │               │ Notify        │               │
       │               │──────────────>│               │
       │               │               │               │
       │               │               │ Notify tmux   │
       │               │               │──────────────>│
       │               │               │               │
       │               │               │ tmux output   │
       │               │               │<──────────────│
       │               │               │               │
```

## Implementation Plan

### Phase 1: Core Parsing and Detection

1. **Create `pkg/assignee/types.go`**
   - Define `AssigneeState`, `AssigneeChangeEvent` structs
   - Define constants (log file path, debounce duration)

2. **Update `pkg/assignee/parser.go`**
   - Implement YAML frontmatter parsing
   - Extract and normalize `assignee` field
   - Handle edge cases (missing frontmatter, empty arrays)

3. **Create `pkg/assignee/cache.go`**
   - Implement `AssigneeCache` with thread-safe operations
   - Methods: `Get()`, `Set()`, `HasChanged()`

4. **Create `pkg/assignee/detector.go`**
   - Integrate with watcher events
   - Compare cached vs new assignee values
   - Generate change events

### Phase 2: Notification and Logging

5. **Create `pkg/assignee/notifier.go`**
   - Execute tmux display-message command
   - Format message with new assignee(s)
   - Handle errors gracefully

6. **Create `pkg/assignee/logger.go`**
   - JSON log file writer
   - Append-only writes with buffering
   - Error handling and recovery

7. **Create `pkg/assignee/manager.go`**
   - Orchestrate all components
   - Handle event flow
   - Coordinate between watcher and subsystems

### Phase 3: Integration and Testing

8. **Update `cmd/monitor/main.go`**
   - Initialize assignee manager
   - Register callbacks for assignee events
   - Start assignee processing

9. **Write unit tests**
   - Parser tests (valid/invalid YAML, edge cases)
   - Cache tests (thread safety, comparison)
   - Detector tests (change detection logic)
   - Notifier tests (tmux command execution)

10. **Integration testing**
    - Test with real task files
    - Verify tmux notifications appear
    - Verify log file format and content
    - Test with rapid writes (debouncing)

## Success Metrics

### Quantitative
- Change detection latency: <500ms (95th percentile)
- File parse time: <100ms per file
- Log write time: <50ms per entry
- Memory overhead: <5MB for 100 files
- Tmux command execution: <100ms

### Qualitative
- Notifications appear clearly in tmux status line
- Log entries are easily searchable and parseable
- No false positives or missed changes
- System recovers gracefully from errors
- Code is maintainable and extensible

## Timeline & Milestones

### Key Dates
- **Design complete**: PRD approved, implementation plan reviewed
- **Implementation complete**: All phases 1-2 done, tests passing
- **Integration complete**: Phase 3 done, end-to-end tested
- **Testing complete**: All acceptance criteria verified
- **Launch/Release**: Deploy and run in production environment

## Stakeholders

### Decision Makers
- Product Owner: Approval of PRD scope and requirements

### Contributors
- Backend Engineer: Implementation of assignee detection system
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
- **GOT-010**: Change detection and JSON logging (Existing task - to be updated)
