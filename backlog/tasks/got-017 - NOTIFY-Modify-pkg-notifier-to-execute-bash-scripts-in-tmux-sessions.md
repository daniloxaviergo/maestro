---
id: GOT-017
title: '[NOTIFY] Modify pkg/notifier to execute bash scripts in tmux sessions'
status: In Progress
assignee: []
created_date: '2026-03-15 17:17'
updated_date: '2026-03-15 18:38'
labels: []
dependencies: []
references:
  - backlog/docs/doc-004-per-agent-configuration.md
priority: high
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify notifier to execute bash scripts in tmux sessions
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 pkg/notifier/notifier.go modified to execute bash scripts in tmux sessions
- [x] #2 New ExecuteScript method to run script in tmux session
- [x] #3 ExecuteScript uses tmux send-keys to run script
- [x] #4 ExecuteScript handles missing script file with error logging
- [x] #5 ExecuteScript handles script execution failures with error logging
- [x] #6 ExecuteScript is non-blocking (runs in goroutine)
- [x] #7 Script output captured but not displayed
- [x] #8 tmux session name from configuration
- [x] #9 tmux session creation if it doesn't exist
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The `ExecuteScript` method will execute bash scripts within tmux sessions using `tmux send-keys`. The implementation will:

- **Script execution via tmux**: Use `tmux send-keys -t <session> <command> Enter` to run scripts in the configured session
- **Session management**: If the session doesn't exist, create it first using `tmux new-session -d -s <session>`
- **Non-blocking execution**: Run the entire sequence in a goroutine, similar to the existing `Notify()` method
- **Output handling**: Use `tmux capture-pane` to capture output after execution, but discard it (per acceptance criteria #7)
- **Error handling**: Check for script file existence before execution and handle command failures with detailed logging
- **Configuration integration**: Read `TmuxSession` from `AgentConfig` via the agent package

**Why this approach:**
- tmux `send-keys` is the standard way to execute commands in existing sessions
- Session-per-agent pattern (from doc-004) allows isolation
- Goroutine execution matches the existing `Notify()` pattern for non-blocking behavior
- Output capture but discard aligns with "script output captured but not displayed" acceptance criterion

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `pkg/notifier/types.go` | Modify | Add `Agent` field to `NotificationConfig` for accessing agent configuration |
| `pkg/notifier/notifier.go` | Modify | Add `ExecuteScript` method with tmux session management |
| `pkg/notifier/notifier_test.go` | Modify | Add unit tests for `ExecuteScript` method |

### 3. Dependencies

- **pkg/config**: Already created (GOT-015) - provides `AgentConfig` with `TmuxSession` field
- **pkg/agent**: Already created (GOT-016) - provides `Agent` struct with `LoadConfig()` and `GetConfig()`
- **tmux**: Must be installed and available in PATH (verified via which tmux)
- **bash**: Required for script execution (standard on most systems)

**Prerequisites:**
- Agent must have `enabled: true` in their config
- Script file at `script_path` must exist and be executable
- tmux session name from config (e.g., "bob" or "alice")

### 4. Code Patterns

Follow existing patterns in `pkg/notifier`:

- **Error variables**: Define distinct error types for different failure modes:
  - `ErrScriptNotFound` - when script file doesn't exist
  - `ErrScriptExecutionFailed` - when script returns non-zero exit code
  - `ErrSessionCreationFailed` - when tmux session cannot be created
  - `ErrTmuxCommandFailed` - when tmux command itself fails
  - `ErrTmuxTimeout` - when tmux command times out

- **Non-blocking execution**: All script execution runs in goroutines with context timeout
- **Error logging**: Use `fmt.Fprintf(os.Stderr, ...)` for warnings, consistent with `Notify()` method
- **Context usage**: Use `context.WithTimeout` for tmux command execution

**Implementation structure:**
```go
func (n *Notifier) ExecuteScript(agent *agent.Agent) {
    go func() {
        // 1. Validate script path from config
        // 2. Check session exists, create if not
        // 3. Execute script via tmux send-keys
        // 4. Handle errors with appropriate logging
    }()
}
```

### 5. Testing Strategy

Add comprehensive unit tests:

- **TestExecuteScript_ScriptNotFound**: Verify error when script file doesn't exist
- **TestExecuteScript_MissingSession**: Verify session creation works
- **TestExecuteScript_Success**: Verify script executes successfully
- **TestExecuteScript_NonBlocking**: Verify method returns immediately (goroutine)
- **TestExecuteScript_ExitCodeNonZero**: Verify error handling for script failures
- **TestExecuteScript_Timeout**: Verify context timeout behavior

**Test approach:**
- Use `t.Setenv("AGENT_NAME", "test-agent")` to set environment variables
- Create temporary script files for testing (clean up with `os.Remove`)
- Mock or test against actual tmux (tmux is available in test environment)
- Use `time.After` or context cancellation to verify non-blocking behavior

### 6. Risks and Considerations

**Blocking issues:**
- None identified - all prerequisites (config, agent packages) are already implemented

**Trade-offs:**
1. **Output capture vs silence**: The acceptance criterion says "script output captured but not displayed" - we'll capture it internally but discard it (not log it) to avoid polluting stderr
2. **Session reuse vs isolation**: Using existing session means concurrent script runs may interleave; this is acceptable per the PRD
3. **Timeout policy**: Will use the existing `config.Timeout` (default 2s) - may need tuning based on actual script execution times

**Implementation considerations:**
- Script path should be absolute or relative to agent config directory
- Need to ensure script is executable (`chmod +x`) or call bash explicitly
- Consider adding `WorkingDirectory` field to `AgentConfig` if scripts need to run from specific directories
- Session name from config allows per-agent isolation
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Code follows existing project conventions package structure naming error handling
- [x] #2 go vet passes with no warnings
- [x] #3 go build succeeds without errors
- [x] #4 Unit tests added or updated for new or changed functionality
- [x] #5 go test ... passes with no failures
- [x] #6 Code comments added for non-obvious logic
- [ ] #7 README or docs updated if public behavior changes
- [ ] #8 make build succeeds
- [ ] #9 make run works as expected
- [ ] #10 Errors are logged not silently ignored
- [ ] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->
