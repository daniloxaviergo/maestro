---
id: GOT-008
title: 'Task 1: File Watcher Implementation using Go''s fsnotify'
status: To Do
assignee: []
created_date: '2026-03-15 00:52'
updated_date: '2026-03-15 00:55'
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

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement a file watcher using Go's `fsnotify` library to monitor `./backlog/tasks` for markdown file changes.

**Architecture Overview:**
- Create a Go workspace module (`maestro`) at project root
- Structure code with separation of concerns: `cmd/`, `pkg/watcher/`, `pkg/cache/`
- Use `fsnotify.Watcher` for filesystem events with recursive directory watching
- Implement event debouncing via a map tracking last event time per file (50ms cooldown)
- Use `os/signal` package to catch SIGINT/SIGTERM for graceful shutdown
- Use channels for event propagation between watcher and handler goroutines

**Key Decisions:**
- **No persistence**: Start fresh on each run (per PRD requirements)
- **Debouncing**: Simple time-based cooldown to handle rapid file changes
- **Event filtering**: Only `.md` files processed; other events logged but ignored
- **Error handling**: Log errors without crashing; watcher continues monitoring

### 2. Files to Modify

**New Files to Create:**
- `go.mod` - Go module definition (`module maestro`)
- `cmd/monitor/main.go` - Entry point with signal handling and initialization
- `pkg/watcher/watcher.go` - fsnotify wrapper with recursive watching
- `pkg/watcher/events.go` - Event types and processing logic
- `pkg/cache/cache.go` - File state cache with debouncing
- `pkg/cache/types.go` - Cache data structures
- `docs/setup-monitor.md` - Documentation for running the monitor

**No Existing Files Modified:**
- Task files remain unchanged (only read, never written)
- No modifications to backlog structure or config

### 3. Dependencies

**Go Standard Library:**
- `os/signal` - Signal handling (SIGINT, SIGTERM)
- `path/filepath` - Path manipulation for filtering `.md` files
- `sync` - Mutex protection for concurrent access
- `time` - Timestamps for debouncing
- `errors` - Error wrapping and handling

**External Dependencies:**
- `github.com/fsnotify/fsnotify@v1.7.0` - Filesystem notifications
- `gopkg.in/yaml.v3` - YAML parsing (for future assignee comparison)

**Prerequisites:**
- Go 1.20+ installed (per GOT-001, GOT-003)
- `./backlog/tasks` directory must exist (created at project setup)
- Write permissions to project directory for log output

### 4. Code Patterns

**Go Conventions:**
```go
// Use explicit error handling
if err != nil {
    log.Printf("error: %v", err)
    return err
}

// Use context for cancellation (optional for future extension)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Use buffered channels to prevent goroutine leaks
eventChan := make(chan fsnotify.Event, 100)
```

**Package Structure:**
- `cmd/monitor/main.go` - CLI entry point only; initialize and run watcher
- `pkg/watcher/` - Filesystem watching logic, event filtering
- `pkg/cache/` - File state tracking, debouncing, caching

**Naming Conventions:**
- Package names: `watcher`, `cache` (lowercase)
- Functions: `NewWatcher`, `StartWatching`, `Stop`
- Variables: `fileWatcher`, `eventQueue` (camelCase)
- Errors: `ErrWatcherStopped`, `ErrInvalidPath`

**Integration Points:**
- Watcher emits events to channel
- Handler reads from channel, processes events
- Cache stores previous state for future assignee comparison (Task 2 integration)

### 5. Testing Strategy

**Unit Tests:**
- `pkg/watcher/watcher_test.go` - Test event detection, filtering, error handling
- `pkg/cache/cache_test.go` - Test caching, debouncing, concurrent access
- Mock filesystem events to test edge cases

**Integration Tests:**
- Create temporary directory, add/remove `.md` files
- Verify events are received and filtered correctly
- Test concurrent file modifications

**Manual Testing:**
```bash
# Run monitor and verify output
go run cmd/monitor/main.go

# Create a new task file and verify event logged
touch backlog/tasks/test-task.md

# Modify an existing task file and verify event logged
echo "test" >> backlog/tasks/got-008\ -\ Task-1-File-Watcher-Implementation-using-Gos-fsnotify.md
```

**Edge Cases to Cover:**
- Non-existent directory (graceful error)
- Permission denied on watch directory
- Rapid successive writes to same file (debouncing)
- File deleted during processing
- Non-`.md` file changes (ignored)

### 6. Risks and Considerations

**Known Risks:**
- **File system limits**: Linux has inotify limits (`fs.inotify.max_user_watches`); may need tuning for large directories
- **Event coalescing**: fsnotify may coalesce rapid writes into single events; debouncing mitigates this
- **Cross-platform differences**: Windows/macOS may handle events differently; testing required

**Trade-offs:**
- **Simplicity over robustness**: Using simple time-based debouncing instead of queue-based approach
- **No persistence**: Starting fresh each run (per PRD); could add persistence later
- **Blocking logger**: Using standard library `log` package; could be made async if needed

**Deployment Considerations:**
- Monitor runs as background process (not integrated with other tools yet)
- No systemd/service file included in scope ( Task 6 may cover)
- Log output to stdout/stderr (piping to file handled by shell)

**Future Extensions (Out of Scope):**
- Configuration file for watch paths
- Output to file vs stdout via flag
- Integration with parser module (Task 2)
- Performance metrics/health checks
<!-- SECTION:PLAN:END -->
