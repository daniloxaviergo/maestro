---
id: GOT-021
title: >-
  Task 2: Script Execution Routing - Integrate agent matching with script
  execution
status: Done
assignee: []
created_date: '2026-03-15 18:52'
updated_date: '2026-03-16 11:02'
labels:
  - task
  - orchestration
dependencies:
  - GOT-018
references:
  - >-
    /home/danilo/scripts/github/maestro/backlog/docs/PRD-Agent-Orchestration-System.md
priority: high
ordinal: 16000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update GOT-021 to In Progress with implementation plan
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The Script Execution Routing will integrate agent matching with script execution in the change detection flow. The implementation follows the existing architecture patterns:

- **Extend `pkg/notifier/` to support multi-agent script execution**: Add a new method `ExecuteScriptsForAssignees()` that takes a list of matching agents and executes their scripts concurrently
- **Update `pkg/change_detect/detector.go`**: Modify `ProcessFile()` to trigger script execution after assignee change detection using the matcher
- **No changes to matcher**: The `pkg/matcher/matcher.go` already provides the matching logic via `MatchAssignees()`

**Architecture decisions:**
- Keep script execution in `notifier.go` since it handles tmux interactions
- Use goroutines for concurrent script execution (non-blocking, per existing pattern)
- Pass `*Matcher` to detector during initialization for agent matching
- Script execution is conditional on agent config `Enabled` field

**Why this approach:**
- Minimal code changes - existing packages handle matching and tmux execution
- Clear separation: matcher finds agents, notifier executes scripts
- Follows existing error handling pattern (log warnings, don't crash)
- Concurrent execution for multiple agents without blocking
- Reuses existing `ExecuteScript()` method with agent-specific config

### 2. Files to Modify

| Action | File | Description |
|--------|------|-------------|
| Modify | `pkg/notifier/notifier.go` | Add `ExecuteScripts()` method for multi-agent script execution |
| Modify | `pkg/change_detect/detector.go` | Add matcher field, import matcher package, update `ProcessFile()` to trigger scripts |
| Modify | `cmd/monitor/main.go` | Create matcher instance, wire matcher to detector |
| Create | `pkg/change_detect/detector_test.go` | Update tests for new matcher integration (if tests exist) |

### 3. Dependencies

**Prerequisites (already satisfied):**
- `pkg/matcher` package with `MatchAssignees(assignees []string) []*Agent` method (GOT-020 completed)
- `pkg/agent` package with `Agent` struct and `GetName()`, `GetConfig()` methods
- `pkg/config` package with `AgentConfig` struct including `ScriptPath`, `TmuxSession`, `Enabled`
- `pkg/notifier` package with `ExecuteScript()` method for single agent
- `pkg/change_detect` package with `ProcessFile()` and `SetNotifier()` methods

**Integration points:**
- `pkg/change_detect.Detector` needs a `*matcher.Matcher` field
- Detector calls `matcher.MatchAssignees(newAssignee)` after change detection
- For each matched agent, call `notifier.ExecuteScriptForAgent(agent)` (new method)
- Script execution is non-blocking (goroutine)

**No new external dependencies** - uses existing packages.

**Blocking tasks:**
- **GOT-018** must be complete (agent config loading) - prerequisite for agent configs to be valid
- **GOT-020** must be complete (matcher) - provides the matching logic

### 4. Code Patterns

**From existing packages to follow:**

1. **pkg/matcher patterns:**
   - Constructor: `NewMatcher(agents []*Agent) *Matcher`
   - Method: `MatchAssignees(assignees []string) []*Agent`
   - Case-insensitive matching via lowercase lookup map
   - Returns empty slice if no matches (not an error)

2. **pkg/notifier patterns:**
   - Method names in CamelCase
   - Non-blocking execution via goroutines
   - Context with timeout for tmux commands
   - Graceful error handling with `fmt.Fprintf(os.Stderr, "warning: ...")`

3. **pkg/change_detect patterns:**
   - `ProcessFile(fileData parser.FileData) (bool, error)` return changed status
   - Update cache before triggering side effects (notifiers/scripts)
   - Log warnings but don't return errors for non-critical issues

4. **New method signature:**
   ```go
   func (n *Notifier) ExecuteScriptsForAgents(agents []*agent.Agent)
   ```
   - Takes slice of matched agents
   - Iterates and calls existing `ExecuteScript()` for each (or inline logic)
   - Non-blocking via goroutines

**Naming conventions:**
- Method: `ExecuteScriptsForAgents(agents []*agent.Agent)` (plural for multi-agent)
- Field in Detector: `matcher *matcher.Matcher`
- Constructor parameter: `agents []*agent.Agent`

### 5. Testing Strategy

**Test cases to add/modify:**

1. **`TestDetector_ProcessFile_WithMatcher_ScriptExecuted`** - Assignee change triggers script for matched agent
2. **`TestDetector_ProcessFile_WithMatcher_NoMatch_WarningLogged`** - No matching agent, script not executed
3. **`TestDetector_ProcessFile_MultipleAgents`** - Multiple assignees trigger multiple scripts
4. **`TestDetector_ProcessFile_AgentDisabled`** - Disabled agent's script not executed
5. **`TestDetector_ProcessFile_ScriptNotFound_Warning`** - Missing script logs warning
6. **Unit tests for `ExecuteScriptsForAgents()`**:
   - Empty agents slice (no-op)
   - Single agent with script
   - Multiple agents with scripts (concurrent)
   - Agent with no script path (skip)
   - Agent with disabled config (skip)

**Verification:**
- `go test ./pkg/change_detect/...` - All tests pass
- `go test ./pkg/notifier/...` - All tests pass
- `go vet ./...` - No warnings
- `make build` - Build succeeds
- Test coverage target: ≥80% for new methods

### 6. Risks and Considerations

**Design considerations:**

1. **Concurrent execution**: Multiple agents' scripts will run concurrently when multiple assignees match. This is intentional per PRD goals but could cause tmux session contention if multiple agents share the same tmux session.

2. **Error isolation**: Each script execution is independent - failure of one agent's script should not affect others. Implemented via goroutines with individual error handling.

3. **Script path resolution**: The `ExecuteScript()` method already resolves script path from agent config. Ensure this works for both relative and absolute paths.

4. **Tmux session handling**: Each agent has its own tmux session per config. If multiple agents share a session, their scripts will queue in tmux (expected behavior).

5. **Performance**: Script execution is non-blocking but can take time. The 2-second timeout from notifier config applies per script. Consider if this is appropriate.

6. **Reentrancy**: If assignee changes rapidly, could trigger multiple script executions for same file before cache updates. This is acceptable per PRD (real-time processing).

**Potential pitfalls:**

1. **Order of script execution**: Not guaranteed when concurrent. If order matters, need to add sequential execution option.

2. **Resource cleanup**: No cleanup of goroutines on shutdown. This is consistent with existing codebase but could be improved.

3. **Debugging**: Multiple concurrent scripts could be hard to debug. Consider adding correlation IDs or file-specific identifiers to log messages.

**Deployment considerations:**

1. No database migrations or config changes required
2. Backward compatible - if no agents match, system behaves as before
3. Can be rolled out without downtime
4. Existing agents will automatically start executing scripts when assigned tasks
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Script execution routing has been integrated with agent matching in the change detection flow. The implementation adds support for executing agent scripts when assignee changes occur.

**Changes made:**

1. **`pkg/change_detect/detector.go`**: 
   - Already had `matcher *matcher.Matcher` field and `SetMatcher()` method
   - `ProcessFile()` now calls `matcher.MatchAssignees(newAssignee)` after assignee change detection
   - Triggers `notifier.ExecuteScriptsForAgents(matchedAgents)` for matched agents

2. **`pkg/notifier/notifier.go`**:
   - `ExecuteScriptsForAgents(agents []*agent.Agent)` - multi-agent script execution method
   - `executeScriptForAgent(agent *agent.Agent, cfg config.AgentConfig)` - helper for single agent
   - Skips disabled agents and agents without script path
   - Non-blocking execution via goroutines with per-agent error handling

3. **`cmd/monitor/main.go`**:
   - Creates `Matcher` with loaded agents
   - Wires matcher to detector via `detector.SetMatcher(matcher)`

4. **`pkg/change_detect/detector_test.go`**:
   - Fixed test calls to use `SetMatcher()` instead of constructor argument
   - Tests cover: no match, script execution, multiple agents, disabled agents

**Testing:**
- All tests pass: `go test ./...` - 26 tests pass
- Build succeeds: `make build` - no errors
- Static analysis: `go vet ./...` - no warnings

**Risks/Follow-ups:**
- Concurrent script execution when multiple agents match (intentional per PRD)
- No cleanup of goroutines on shutdown (consistent with existing codebase)
<!-- SECTION:FINAL_SUMMARY:END -->
