---
id: doc-002
title: 'PRD: Monitor File Changes in ./backlog/tasks When assignee Field Is Modified'
type: other
created_date: '2026-03-15 00:50'
---
# PRD: Monitor File Changes in `./backlog/tasks` When the `assignee` Field Is Modified

## Overview

### Purpose
Monitor real-time changes to markdown task files in `./backlog/tasks` and trigger a log action when the `assignee` field is modified, enabling audit trail and notification capabilities for task assignments.

### Goals
- **G1**: Detect assignee field modifications in real-time with <500ms latency
- **G2**: Log all assignee change events with full context (file path, old value, new value, timestamp)
- **G3**: Support both existing files and newly created files in the tasks directory
- **G4**: Parse YAML frontmatter format to extract and compare assignee values

## Background

### Problem Statement
Currently, there is no automated way to track when task assignees are modified in the backlog system. Team members must manually check file change history or rely on git logs to audit who was assigned to a task and when changes occurred. This creates visibility gaps and makes it difficult to maintain an accurate audit trail.

### Current State
- Task files in `./backlog/tasks` store metadata in YAML frontmatter format
- The `assignee` field is an array (e.g., `assignee: []` or `assignee: ["username"]`)
- Changes to task files are tracked by git, but real-time visibility into assignee changes is not available
- No mechanism exists to trigger actions when specific field modifications occur

### Proposed Solution
Implement a file watcher in Go that monitors the `./backlog/tasks` directory for changes to markdown files. When a change is detected, parse the YAML frontmatter to extract the `assignee` field and compare it against the previous value. If the assignee field was modified, log the change event with full context.

## Requirements

### User Stories

- **Role**: Project Manager / Team Lead
  - *As a project manager, I want to know when task assignees are changed so that I can track ownership transitions and investigate if needed*

- **Role**: Developer
  - *As a developer, I want the system to log assignee changes automatically so I don't need to manually track who made changes*

- **Role**: Auditor
  - *As an auditor, I want a complete log of all assignee changes with timestamps so I can review assignment history*

### Functional Requirements

#### Task 1: File Watcher Implementation

Implement a file watching mechanism using Go's `fsnotify` or `otify` library to detect real-time changes to markdown files in `./backlog/tasks`.

##### User Flows
1. Program starts and initializes file watcher on `./backlog/tasks`
2. File watcher detects any change (create, write, rename, remove) to markdown files
3. For write events, the system compares the new assignee value with the cached previous value
4. If assignee changed, the log action is triggered

##### Acceptance Criteria
- [ ] Watcher detects file write events in `./backlog/tasks`
- [ ] Watcher handles recursive monitoring of subdirectories
- [ ] Watcher properly handles concurrent file changes without race conditions
- [ ] Watcher gracefully handles file permission errors and other I/O issues
- [ ] Watcher stops cleanly on interrupt signals (SIGINT, SIGTERM)

#### Task 2: YAML Frontmatter Parser

Parse the YAML frontmatter section of markdown files to extract the `assignee` field value.

##### User Flows
1. When a markdown file is detected, read the file content
2. Extract the YAML frontmatter (content between `---` delimiters)
3. Parse the YAML to extract the `assignee` field
4. Handle cases where frontmatter is missing or malformed

##### Acceptance Criteria
- [ ] Successfully parse YAML frontmatter from valid markdown files
- [ ] Extract `assignee` field as a slice of strings (array)
- [ ] Handle files without frontmatter (treat as empty assignee array)
- [ ] Handle malformed YAML gracefully with error logging
- [ ] Support empty assignee arrays (`assignee: []` or `assignee:`)

#### Task 3: Change Detection and Logging

Compare the current assignee value with the cached previous value and log changes.

##### User Flows
1. Cache the assignee value for each file after parsing
2. On subsequent file change events, re-parse and compare
3. If assignee differs, log the change event
4. Update the cached value

##### Acceptance Criteria
- [ ] Log event includes: file path, timestamp, old assignee (array), new assignee (array)
- [ ] Log format is human-readable and machine-parsable (JSON or structured text)
- [ ] Log output goes to `./backlog/logs/assignee_changes.log`
- [ ] Handle assignee additions, removals, and replacements correctly
- [ ] Handle case where file previously had no cached value (treat as new assignee array)

##### Log Format Example
```json
{
  "timestamp": "2026-03-14T10:30:00Z",
  "file": "backlog/tasks/task-001.md",
  "old_assignee": ["alice"],
  "new_assignee": ["bob"]
}
```

### Non-Functional Requirements

- **Performance**: 
  - File watcher should detect changes within 500ms
  - File parsing should complete in <100ms for typical task files
  - Log writes should be non-blocking or use async I/O

- **Reliability**: 
  - System should recover from file system errors without crashing
  - Cache should persist or be rebuilt on restart (start fresh mode)
  - Handle file system events that may coalesce (multiple writes to same file)

- **Maintainability**: 
  - Code should follow Go best practices and conventions
  - Package structure should allow for future extension (e.g., different triggers)
  - Logging should include appropriate debug/error levels

- **Compatibility**: 
  - Go 1.20 or later
  - Linux, macOS, and Windows platforms
  - Markdown files in the specified YAML frontmatter format

## Scope

### In Scope
- File watcher implementation using Go file system monitoring
- YAML frontmatter parsing for extract `assignee` field
- Change detection with compare-and-cache logic
- Log output to `./backlog/logs/assignee_changes.log` in JSON format
- Support for both existing and newly created markdown files

### Out of Scope
- Email/SMS notifications (logging is the sole trigger action)
- Integration with external systems (webhooks, databases)
- Real-time notifications to users
- History replay (starting fresh on each run)
- Monitoring of other fields beyond `assignee`
- User interface for viewing logs

## Technical Considerations

### Existing System Impact
- This implementation adds a new background process that monitors the `./backlog/tasks` directory
- No modifications to existing markdown file format or structure required
- Logs are appended to a separate file, not modifying any task files

### Dependencies
- **Go standard library**: `os`, `io/ioutil`, `time`, `log`, `encoding/json`
- **External libraries**: 
  - `github.com/fsnotify/fsnotify` (file watching)
  - `gopkg.in/yaml.v3` or `github.com/go-yaml/yaml` (YAML parsing)

### Constraints
- Start fresh mode: no cache persistence across runs
- Assignee field is always treated as an array for consistency
- Log file grows indefinitely (no rotation or cleanup in scope)

## Success Metrics

### Quantitative
- Change detection latency: <500ms (95th percentile)
- File parse time: <100ms per file
- Log write time: <50ms per entry
- Memory usage: <50MB for typical task count (100-500 files)

### Qualitative
- Log entries should be easily searchable and parseable
- Errors should be logged with sufficient context for debugging
- System should be easy to run and monitor

## Timeline & Milestones

### Key Dates
- **Design complete**: Implementation plan reviewed and approved
- **Implementation complete**: Code passes all acceptance criteria
- **Testing complete**: Integration testing with real task files
- **Launch/Release**: Deploy and run in production environment

## Stakeholders

### Decision Makers
- Product Owner: Approval of PRD scope and requirements

### Contributors
- Backend Engineer: Implementation of file watcher and parser
- QA Engineer: Testing and validation of change detection

## Appendix

### Glossary
- **YAML frontmatter**: Metadata section at the top of markdown files, delimited by `---`
- **Assignee field**: The `assignee` key in frontmatter containing an array of user identifiers
- **fsnotify**: A Go library for file system change notifications

### References
- Backlog.md workflow: `backlog://workflow/overview`
- YAML spec: https://yaml.org/spec/
- fsnotify documentation: https://github.com/fsnotify/fsnotify
- Go yaml.v3: https://gopkg.in/yaml.v3
