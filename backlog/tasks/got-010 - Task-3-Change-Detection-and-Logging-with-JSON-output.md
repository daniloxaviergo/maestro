---
id: GOT-010
title: 'Task 3: Change Detection and Logging with JSON output'
status: To Do
assignee: []
created_date: '2026-03-15 00:52'
updated_date: '2026-03-15 00:53'
labels:
  - logging
  - json
  - go
dependencies:
  - GOT-008
  - GOT-009
references:
  - >-
    backlog/docs/doc-002 -
    PRD-Monitor-File-Changes-in-.-backlog-tasks-When-assignee-Field-Is-Modified.md
priority: high
ordinal: 5000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement change detection logic to compare cached assignee values with current values and log changes to `./backlog/logs/assignee_changes.log` in JSON format.

### Implementation Notes

**Goal**: Detect assignee field modifications and log events with full context.

**Key Components**:
1. Create `pkg/cache/` package for storing previous assignee values
2. Implement `CompareAndLog()` function to detect changes
3. Write JSON-formatted log entries to `./backlog/logs/assignee_changes.log`
4. Handle first-run scenario (cache empty, treat all as "new assignees")

**Log Entry Format**:
```json
{
  "timestamp": "2026-03-14T10:30:00Z",
  "file": "backlog/tasks/task-001.md",
  "old_assignee": ["alice"],
  "new_assignee": ["bob"]
}
```

**Cache Implementation**:
- Simple in-memory map: `map[string][]string`
- No persistence required (start fresh on each run)
- Thread-safe access using `sync.RWMutex`
- Handle file deletions (remove from cache)

**Change Detection Logic**:
1. Parse current assignee value
2. Look up cached previous value
3. Compare slices (order-insensitive comparison)
4. If different, trigger log write

**Log Writing**:
- Non-blocking or async I/O to avoid blocking watcher
- Append mode to log file
- Handle log directory creation if missing
- Error handling for write failures

**Integration Points**:
- Receives parsed data from parser module
- Outputs log entries to `./backlog/logs/assignee_changes.log`

**Dependencies**:
- Go standard library: `encoding/json`, `os`, `io/ioutil`

**Acceptance Criteria**:
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Log event includes: file path, timestamp, old assignee (array), new assignee (array)
- [ ] #2 Log format is human-readable and machine-parsable (JSON or structured text)
- [ ] #3 Log output goes to ./backlog/logs/assignee_changes.log
- [ ] #4 Handle assignee additions, removals, and replacements correctly
- [ ] #5 Handle case where file previously had no cached value (treat as new assignee array)
<!-- AC:END -->
