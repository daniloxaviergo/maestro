---
id: GOT-014
title: 'Task 7: Update Monitor Main to Use Notifier'
status: In Progress
assignee: []
created_date: '2026-03-15'
updated_date: '2026-03-15 17:24'
labels:
  - tmux
  - notifier
  - integration
  - go
dependencies:
  - GOT-011
  - GOT-012
  - GOT-013
references:
  - backlog/docs/doc-003 - PRD-Maestro-Feature-Set-1.md
priority: high
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the monitor main application to initialize and wire the tmux notifier to the change detector.

### Implementation Notes

**Goal**: Wire the complete notification flow from monitor entry point.

**Key Components**:
1. Modify `cmd/monitor/main.go` to initialize notifier
2. Wire notifier to change detector
3. Ensure notifier is called on assignee changes

**Implementation Steps**:
1. Import `github.com/yourorg/maestro/pkg/notifier`
2. Create `cmd/monitor/main.go`:
   - Initialize notifier with default config
   - Create change detector
   - Wire notifier to detector
   - Process file events through pipeline

**Default Configuration**:
- Timeout: 2 seconds (standard)
- Message format: `Assignee changed to "[new]" for [file]`

**Data Flow**:
1. Watcher detects file event → sends `FileEvent` to channel
2. Consumer receives event → calls parser
3. Parser returns `FileData` → cache lookup
4. Change detector compares → detects change
5. Change detector calls notifier → tmux notification

**Error Handling**:
- Monitor doesn't handle notifier errors (handled internally)
- All other error paths as per existing implementation

**Dependencies**:
- Existing: `cmd/monitor/main.go` (GOT-008)
- Existing: `pkg/change_detect/` (GOT-013)
- Existing: `pkg/notifier/` (GOT-011, GOT-012)
- Go standard library only

**Acceptance Criteria**:
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 `cmd/monitor/main.go` imports `pkg/notifier`
- [ ] #2 Notifier initialized with default config (2s timeout, default format)
- [ ] #3 Notifier wired to change detector via `SetNotifier()`
- [ ] #4 Notify() called on assignee change events
- [ ] #5 Integration test: tmux notification triggered when assignee changes
- [ ] #6 #1 `cmd/monitor/main.go` imports `pkg/notifier`
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Update the monitor main to initialize and wire the notifier to the change detector.

**Implementation Steps:**
1. Read current `cmd/monitor/main.go` structure
2. Add notifier import
3. Create notifier with default config
4. Create detector with cache
5. Wire notifier to detector
6. Monitor runs as before

**Design Decisions:**
- Notifier uses default config (no configuration file for this task)
- Detector created first, then notifier wired
- No configuration flags for notifier in this task (future enhancement)

### 2. Files to Create/Modify

| Action | File | Purpose |
|--------|------|---------|
| Modify | `cmd/monitor/main.go` | Initialize notifier, wire to detector |

### 3. Dependencies

- **Existing package**: `cmd/monitor/main.go` (GOT-008)
- **Existing package**: `pkg/change_detect/` (GOT-013)
- **Existing package**: `pkg/notifier/` (GOT-011, GOT-012)
- **Go standard library**: None new

### 4. Code Patterns

**Follow existing project conventions:**
- Error handling: `fmt.Fprintf(os.Stderr, "error: %v\n", err)` (matching existing code)
- Signal handling: SIGINT/SIGTERM for graceful shutdown
- Resource cleanup: close channels, cancel contexts

**Example main.go integration:**
```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/yourorg/maestro/pkg/cache"
    "github.com/yourorg/maestro/pkg/change_detect"
    "github.com/yourorg/maestro/pkg/notifier"
    "github.com/yourorg/maestro/pkg/parser"
    "github.com/yourorg/maestro/pkg/watcher"
)

func main() {
    // Initialize cache
    c := cache.New()
    
    // Initialize parser
    p := parser.NewParser()
    
    // Initialize change detector
    detector := change_detect.NewDetector(c)
    
    // Initialize notifier (optional, defaults to 2s timeout)
    notifier := notifier.NewNotifier(notifier.NotificationConfig{})
    
    // Wire notifier to detector
    detector.SetNotifier(notifier)
    
    // Initialize watcher
    w, err := watcher.New("./backlog/tasks", c)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
    }
    defer w.Close()
    
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Start watcher
    go func() {
        if err := w.Watch(); err != nil {
            fmt.Fprintf(os.Stderr, "error: %v\n", err)
        }
    }()
    
    // Process events (existing logic)
    // ... existing event processing code ...
    
    // Wait for signal
    <-sigChan
}
```

### 5. Testing Strategy

**Integration test approach:**
1. Start monitor with notifier and detector
2. Create task file with assignee
3. Update assignee field
4. Verify tmux notification displayed
5. Verify monitor continues running

**Manual verification:**
- Run monitor: `go run cmd/monitor/main.go`
- In another terminal, update task file's assignee
- Observe tmux status bar notification
- Verify monitor doesn't crash

**Edge cases:**
- tmux not installed: warning logged, monitor continues
- No notifier: detector works normally
- Multiple rapid changes: each triggers notification (async)

### 6. Risks and Considerations

**Blocking issues:**
- tmux may not be installed (handled by error logging in notifier)

**Trade-offs:**
- No configuration file for notifier (simple defaults)
- No way to disable notifier without code changes (future feature)

**Performance considerations:**
- Notifier async call doesn't block main loop
- Multiple notifications queue naturally (no backpressure yet)

**Future considerations:**
- Add CLI flag to disable notifier
- Add config file for notifier settings
- Add notification metrics/health check endpoint
<!-- SECTION:PLAN:END -->
