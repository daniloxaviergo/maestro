---
id: GOT-008
title: 'Task 1: File Watcher Implementation using Go''s fsnotify'
status: To Do
assignee: []
created_date: '2026-03-15 00:52'
updated_date: '2026-03-15 00:53'
labels:
  - monitoring
  - filesystem
  - go
dependencies: []
references:
  - >-
    backlog/docs/doc-002 -
    PRD-Monitor-File-Changes-in-.-backlog-tasks-When-assignee-Field-Is-Modified.md
priority: high
ordinal: 3000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement a file watcher in Go that monitors the `./backlog/tasks` directory for real-time changes to markdown files.

### Implementation Notes

**Goal**: Detect file write/create/rename/remove events for markdown files in `./backlog/tasks`.

**Key Components**:
1. Create `main.go` with entry point that initializes the watcher
2. Use `fsnotify` library to watch `./backlog/tasks` directory (recursively)
3. Filter events to only process `.md` files
4. Handle concurrent file changes (debounce or queue mechanism)
5. Implement graceful shutdown on SIGINT/SIGTERM signals
6. Add error handling for file permission issues and I/O errors

**Architecture**:
```
cmd/
└── monitor/
    └── main.go           # Entry point, signal handling
pkg/
├── watcher/              # fsnotify wrapper
│   ├── watcher.go
│   └── events.go
└── cache/                # File state cache
    ├── cache.go
    └── types.go
```

**File Watcher Requirements**:
- Watch `./backlog/tasks` recursively for `.md` files
- Handle events: Create, Write, Rename, Remove
- Detect file write events (for assignee comparison)
- Implement rate limiting/debouncing for rapid changes
- Handle "file not found" errors when files are deleted

**Integration Points**:
- Output parsed file data to channels for processing by parser module
- Emit change events when assignee differences detected

**Dependencies**:
- `github.com/fsnotify/fsnotify`
- Go standard library: `os/signal`, `path/filepath`, `sync`

**Acceptance Criteria**:
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Watcher detects file write events in ./backlog/tasks
- [ ] #2 Watcher handles recursive monitoring of subdirectories
- [ ] #3 Watcher properly handles concurrent file changes without race conditions
- [ ] #4 Watcher gracefully handles file permission errors and other I/O issues
- [ ] #5 Watcher stops cleanly on interrupt signals (SIGINT, SIGTERM)
<!-- AC:END -->
