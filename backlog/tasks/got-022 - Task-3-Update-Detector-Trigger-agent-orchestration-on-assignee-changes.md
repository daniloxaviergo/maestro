---
id: GOT-022
title: 'Task 3: Update Detector - Trigger agent orchestration on assignee changes'
status: In Progress
assignee: []
created_date: '2026-03-15 18:53'
updated_date: '2026-03-15 22:12'
labels:
  - task
  - orchestration
dependencies:
  - GOT-020
references:
  - >-
    /home/danilo/scripts/github/maestro/backlog/docs/PRD-Agent-Orchestration-System.md
priority: high
ordinal: 12250
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Task 3: Update Detector to trigger agent orchestration on assignee changes
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task is to update the `Detector` to trigger agent orchestration when assignee changes occur. The system already has:

1. **Matcher** (`pkg/matcher`): Matches assignee names to configured agents (case-insensitive)
2. **Agent** (`pkg/agent`): Manages agent identity and configuration
3. **Config** (`pkg/config`): Loads agent configuration from YAML files
4. **Notifier** (`pkg/notifier`): Executes scripts via tmux sessions

The `Detector` already has infrastructure for agent orchestration but needs updates to fully integrate with the existing agent system:

- `Detector` already calls `d.matcher.MatchAssignees(newAssignee)` when assignee changes
- `Detector` already calls `d.notifier.ExecuteScriptsForAgents(matchedAgents)` 
- However, the `NewDetector` function signature has changed to require a matcher

**Approach:**
1. Verify current implementation in `detector.go` is complete and correct
2. Ensure proper error handling for agent script execution
3. Add comprehensive tests for agent orchestration flow
4. Update monitor main to properly wire detector with matcher
5. Verify graceful degradation (missing scripts, disabled agents, etc.)

### 2. Files to Modify

No files need modification. The implementation is already complete in:
- `pkg/change_detect/detector.go` - Core logic for agent orchestration is present
- `pkg/matcher/matcher.go` - Agent matching logic exists
- `pkg/notifier/notifier.go` - Script execution for multiple agents exists

**Files to verify (no changes needed if current state is correct):**
- `pkg/change_detect/detector_test.go` - Tests already cover agent orchestration
- `pkg/matcher/matcher_test.go` - Tests already comprehensive
- `cmd/monitor/main.go` - Integration with detector/matcher/notifier

### 3. Dependencies

- **GOT-020 (Agent Matching Engine)**: MUST be complete first - provides `matcher.Matcher` interface
- **GOT-015 (pkg/config)**: MUST be complete - loads agent configuration
- **GOT-016 (pkg/agent)**: MUST be complete - manages agent identity
- **GOT-017 (Tmux script execution)**: MUST be complete - `ExecuteScriptsForAgents` method
- **GOT-013 (Integrate Notifier with Detector)**: MUST be complete - detector has `SetNotifier`

**Prerequisites for testing:**
- `./agents/` directory with agent configs for integration testing
- Tmux installed for script execution tests

### 4. Code Patterns

Follow existing patterns in the codebase:

**Error Handling:**
- Log warnings for recoverable errors (missing scripts, disabled agents)
- Do not crash on individual agent failures
- Continue processing other agents if one fails

**Concurrency:**
- Script execution is non-blocking (goroutines)
- Detector uses mutex-protected cache
- Matcher uses read-only operations (thread-safe)

**Naming:**
- Go conventions: camelCase for functions/variables, PascalCase for types
- Package names: short lowercase (e.g., `agent`, `matcher`, `notifier`)

**Integration pattern:**
```go
matcher := matcher.NewMatcher(agents)
detector.SetMatcher(matcher)

notifier := notifier.NewNotifier(notifier.NotificationConfig{})
detector.SetNotifier(notifier)
```

### 5. Testing Strategy

**Unit Tests** (already exist in `detector_test.go`):
- `TestProcessFile_WithMatcher_NoMatch` - No agents match
- `TestProcessFile_WithMatcher_ScriptExecuted` - Script execution triggered
- `TestProcessFile_MultipleAssigneesWithMatcher` - Multiple agents matched
- `TestProcessFile_AgentDisabled` - Disabled agent skipped
- `TestDetector_SetMatcher` - Matcher can be set after construction

**Test Coverage:**
- ✓ New file creation (no change detection)
- ✓ Same assignee (no change)
- ✓ Different assignee (change detected)
- ✓ Multiple assignees order-insensitive
- ✓ Agent matching (match/no-match)
- ✓ Disabled agent handling
- ✓ Multiple agent orchestration

**Missing Test Scenarios:**
- Script execution with missing script file
- Script execution with timeout
- Multiple agents with mixed enabled/disabled states

**Verification Steps:**
1. `go test ./...` - All tests pass
2. `go vet ./...` - No warnings
3. `make build` - Binary compiles
4. Manual test with real agent configs

### 6. Risks and Considerations

**No blocking issues.** The implementation is already complete. Risks are minimal:

- **Graceful Degradation**: Already handled - disabled agents skipped, missing scripts logged
- **Concurrent Execution**: Already handled - each script in its own goroutine
- **Timeout Handling**: Already handled - context with timeout per script
- **Error Logging**: Already handled - all failures logged with agent name

**Trade-offs:**
- Script execution is non-blocking (may lose error output if tmux not attached)
- Agent matching is case-insensitive (allows flexibility but may cause confusion)
- No retry logic for failed script executions (intentional for simplicity)

**Deployment Considerations:**
- Ensure `./agents/` directory structure exists before runtime
- Agent configs must have valid `script_path` (absolute or relative to working directory)
- Tmux must be installed and accessible via PATH
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 Code follows existing project conventions package structure naming error handling
- [ ] #2 go vet passes with no warnings
- [ ] #3 go build succeeds without errors
- [ ] #4 Unit tests added or updated for new or changed functionality
- [ ] #5 go test ... passes with no failures
- [ ] #6 Code comments added for non-obvious logic
- [ ] #7 README or docs updated if public behavior changes
- [ ] #8 make build succeeds
- [ ] #9 make run works as expected
- [ ] #10 Errors are logged not silently ignored
- [ ] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
- [ ] #13 #1 Verify detector.go has complete agent orchestration integration (matcher + notifier)
- [ ] #14 #2 Verify all existing tests pass with current implementation
- [ ] #15 #3 Run `go vet ./...` and fix any warnings
- [ ] #16 #4 Run `make build` successfully
- [ ] #17 #5 Add integration test with real agent config if missing
<!-- DOD:END -->
