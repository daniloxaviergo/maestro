---
id: GOT-013
title: 'Task 6: Integrate Notifier with Detector'
status: To Do
assignee: []
created_date: '2026-03-15'
updated_date: '2026-03-15'
labels:
  - tmux
  - notifier
  - integration
  - go
dependencies:
  - GOT-011
  - GOT-012
  - GOT-010
references:
  - backlog/docs/doc-003 - PRD-Maestro-Feature-Set-1.md
priority: high
ordinal: 1300
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Integrate the tmux notifier with the change detector to notify on assignee changes.

### Implementation Notes

**Goal**: Wire the notifier to receive change events from the change detector and trigger notifications.

**Key Components**:
1. Modify `pkg/change_detect/detector.go` to accept notifier callback
2. Pass `AssigneeChangeEvent` to notifier on change detected
3. Ensure non-blocking execution

**Integration Points**:
- Change detector detects assignee change
- Detector calls notifier with `AssigneeChangeEvent`
- Notifier sends tmux notification (non-blocking)

**Non-blocking Execution**:
- Detector calls `Notify()` which returns immediately
- Notifier handles async execution internally
- Detector doesn't wait for tmux command completion

**Error Handling**:
- Notifier handles all tmux errors internally
- Detector doesn't see or handle notifier errors
- Watcher continues regardless of notification success

**Data Flow**:
1. Change detector: `CompareAndLog()` detects assignee change
2. Change detector: calls `notifier.Notify(event)` 
3. Notifier: formats message, executes tmux command asynchronously

**Dependencies**:
- Existing: `pkg/change_detect/detector.go` (GOT-010)
- Existing: `pkg/notifier/` (GOT-011, GOT-012)
- Go standard library only

**Acceptance Criteria**:
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 `pkg/change_detect/detector.go` accepts `*notifier.Notifier` parameter
- [ ] #2 Detector calls `notifier.Notify()` with `AssigneeChangeEvent` on assignee change
- [ ] #3 Non-blocking: detector doesn't wait for notification to complete
- [ ] #4 Notifier handles all errors internally (detector doesn't see errors)
- [ ] #5 Integration test: assignee change triggers tmux notification
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Integrate the notifier with the existing change detector by adding it as an optional callback.

**Implementation Steps:**
1. Read current `pkg/change_detect/detector.go` structure
2. Add `Notifier` field to `Detector` struct
3. Add `SetNotifier()` method or pass in constructor
4. In `ProcessFile()` or new method, call `notifier.Notify()` on change detected

**Design Decisions:**
- Make notifier optional (can be nil)
- Detector calls notifier without error checking
- Notifier handles all async logic internally

**Integration Points:**
- Detector detects assignee change via `CompareAndLog()`
- Detector creates `AssigneeChangeEvent`
- Detector calls `notifier.Notify(event)` if notifier is not nil

### 2. Files to Create/Modify

| Action | File | Purpose |
|--------|------|---------|
| Modify | `pkg/change_detect/detector.go` | Add notifier field, call Notify() on change |
| Create | `pkg/change_detect/detector_test.go` | Update tests with notifier integration |
| Create | `pkg/change_detect/integration_test.go` | Integration test with notifier |

### 3. Dependencies

- **Existing package**: `pkg/change_detect/detector.go` (GOT-010)
- **Existing package**: `pkg/notifier/` (GOT-011, GOT-012)
- **Go standard library**: None new

### 4. Code Patterns

**Follow existing project conventions:**
- Optional dependency: check `if d.notifier != nil` before calling
- Non-blocking: call method that returns immediately
- Error handling: notifier handles its own errors

**Example integration:**
```go
package changedetect

import (
    "github.com/yourorg/maestro/pkg/notifier"
)

type Detector struct {
    cache      *cache.Cache
    notifier   *notifier.Notifier  // Optional
}

func NewDetector(c *cache.Cache) *Detector {
    return &Detector{cache: c}
}

func (d *Detector) SetNotifier(n *notifier.Notifier) {
    d.notifier = n
}

func (d *Detector) CompareAndLog(filePath string, currentAssignee []string) error {
    prevAssignee := d.cache.GetAssignee(filePath)
    
    // First run: no previous value, just cache and return
    if prevAssignee == nil {
        d.cache.SetAssignee(filePath, currentAssignee)
        return nil
    }
    
    // Compare assignees
    if slicesEqual(prevAssignee, currentAssignee) {
        return nil
    }
    
    // Change detected: log and notify
    event := notifier.AssigneeChangeEvent{
        FilePath:    filePath,
        OldAssignee: prevAssignee,
        NewAssignee: currentAssignee,
    }
    
    // Non-blocking: notifier handles async
    if d.notifier != nil {
        d.notifier.Notify(event)
    }
    
    // Update cache and log
    d.cache.SetAssignee(filePath, currentAssignee)
    return d.logChange(event)
}
```

### 5. Testing Strategy

**Unit tests in `pkg/change_detect/detector_test.go`:**
- Test detector without notifier (no panic, no call)
- Test detector with notifier (Notify() called)
- Test with nil notifier (graceful handling)

**Integration test approach:**
- Run monitor with detector and notifier
- Update a task file's assignee field
- Verify tmux notification is triggered
- Verify detector continues working

**Edge cases:**
- Notifier is nil: detector works normally
- Notifier returns immediately: detector doesn't block
- Multiple changes: each triggers separate notification

### 6. Risks and Considerations

**Blocking issues:**
- None - notifier is optional, can be nil

**Trade-offs:**
- Notifier is optional: some deployments may not use tmux
- No feedback to detector about notification success
- Detector doesn't handle notifier errors (they're internal)

**Performance considerations:**
- Detector not affected by tmux latency (async call)
- No buffering or batching of notifications

**Future considerations:**
- Could add notification queue for high-frequency changes
- Could add metrics for notification success rate
<!-- SECTION:PLAN:END -->
