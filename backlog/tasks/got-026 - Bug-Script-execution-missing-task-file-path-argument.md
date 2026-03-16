---
id: GOT-026
title: 'Bug: Script execution missing task file path argument'
status: In Progress
assignee: []
created_date: '2026-03-16 00:30'
updated_date: '2026-03-16 02:05'
labels:
  - bug
  - script-execution
dependencies: []
documentation:
  - backlog/docs/doc-006.md
priority: high
ordinal: 2000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix the script execution to pass the task file path as an argument to agent scripts when they are invoked. Currently, scripts are executed without any arguments, so the task file path ($1) is empty even though the script is triggered.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 - [ ] When assignee changes, agent scripts are invoked with the task file path as the first argument
- [ ] #2 - [ ] The script receives the full absolute path to the task file
- [ ] #3 - [ ] Agent scripts can access and process the task file content via the passed argument
- [ ] #4 - [ ] Manual test: Create a task with assignee `[agent-bar]`, verify script receives the file path in the log output
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The bug is in `pkg/notifier/notifier.go` where agent scripts are executed via tmux without passing the task file path as an argument. Currently, the `tmux send-keys` command only executes `bash {script_path}` but should execute `bash {script_path} {task_file_path}`.

The fix requires:
1. Modify `ExecuteScriptsForAgents` to accept the task file path as a parameter
2. Modify `executeScriptForAgent` to include the file path in the tmux command
3. Update `change_detect/detector.go` to pass the file path when triggering script execution

**Key observation**: The `AssigneeChangeEvent` struct in `notifier/types.go` already has a `FilePath` field that captures the task file path, but it's only used for notifications, not script execution.

**Architecture decision**: Pass the file path through the `AssigneeChangeEvent` to `ExecuteScriptsForAgents`, ensuring scripts receive it as `$1`.

### 2. Files to Modify

| File | Change |
|------|--------|
| `pkg/notifier/types.go` | No changes needed (FilePath already exists in AssigneeChangeEvent) |
| `pkg/notifier/notifier.go` | Modify `ExecuteScriptsForAgents` and `executeScriptForAgent` to accept and use file path |
| `pkg/change_detect/detector.go` | Update the call to `ExecuteScriptsForAgents` to pass the file path |
| `pkg/notifier/notifier_test.go` | Add tests for script execution with file path argument |

### 3. Dependencies

- **No new dependencies** - Uses existing `exec` package for tmux commands
- **Prerequisites**: 
  - Task file must exist and be readable (handled by existing parser)
  - tmux must be installed (existing error handling covers this)
  - Agent script must exist and be executable (existing checks cover this)
- **Related tasks**: 
  - GOT-020 (Agent Matching Engine) - Done
  - GOT-022 (Detector trigger) - Done  
  - GOT-025 (Testing) - Already has notifier tests to extend

### 4. Code Patterns

**Follow existing patterns in the codebase:**
- Use `fmt.Sprintf` for command construction (already used in `executeScriptForAgent`)
- Error handling with `log.Printf` for warnings (existing pattern)
- Non-blocking execution via goroutines (existing pattern)
- Context timeouts for tmux commands (existing pattern)

**Naming conventions to follow:**
- Function names: camelCase (e.g., `executeScriptForAgent`)
- Error variables: prefix with `Err` (e.g., `ErrScriptExecutionFailed`)
- Struct fields: camelCase (e.g., `FilePath`, `NewAssignee`)

**Integration pattern:**
- The `AssigneeChangeEvent` is already constructed in `detector.go` and passed to `notifier.Notify()`
- Need to also pass it (or just the file path) to `ExecuteScriptsForAgents`

### 5. Testing Strategy

**Unit tests to add/modify:**

1. **`notifier_test.go`** - Add tests for script execution with file path:
   - `TestExecuteScriptForAgents_WithFilePath`: Verify scripts receive file path as $1
   - `TestExecuteScriptWithFilePath_Argument`: Create a test script that echoes $1 and verify it matches the expected path
   - Test with multiple agents to ensure each receives the same file path

2. **`detector_test.go`** - Add integration test:
   - Test full flow: change assignee → detector processes → scripts executed with file path
   - Verify the change is logged AND scripts receive the path

3. **Manual test** (from acceptance criteria):
   - Create task with assignee `[agent-bar]`
   - Run monitor
   - Check `agents/agent-bar/execution.log` for the file path

**Edge cases to cover:**
- Script with no arguments (current behavior) - should still work but won't receive path
- Multiple agents assigned to same task - all should receive same file path
- Script that doesn't use $1 - should still execute without error

### 6. Risks and Considerations

**Blocking issues: None** - This is a straightforward parameter addition.

**Trade-offs:**
1. **Approach A** (chosen): Modify `ExecuteScriptsForAgents` signature to accept file path
   - Pros: Clean API, clear intent
   - Cons: Slight API change (add parameter)

2. **Alternative**: Store file path in `NotificationConfig` or `Notifier` struct
   - Pros: No API change
   - Cons: Stateful notifier, harder to test, not thread-safe if multiple changes concurrent

3. **Alternative**: Extract file path from `AssigneeChangeEvent` inside `ExecuteScriptsForAgents`
   - Pros: No signature change
   - Cons: Less flexible, ties notifier to detector's data structure

**Implementation choice**: Approach A is cleanest. The method already has a `[]*agent.Agent` parameter; adding `filePath string` is a natural extension.

**Potential pitfalls:**
- Path sanitization: The file path comes from `fsnotify` which should provide absolute paths, but should verify with `filepath.Abs()` before passing
- Concurrent file changes: If multiple assignee changes happen rapidly, each script invocation should use the correct file path for that event
- Empty file path: Should handle gracefully (log warning, skip script execution)

**Rollout considerations:**
- No database/schema migrations needed
- No config file changes needed
- Scripts that don't use $1 will continue to work (backward compatible)
- Scripts that expect $1 will now receive the correct value

### 7. Implementation Steps

```markdown
1. Update `notifier.go`:
   - Modify `ExecuteScriptsForAgents(agents []*agent.Agent, filePath string)` signature
   - Update `executeScriptForAgent(agent *agent.Agent, cfg config.AgentConfig, filePath string)` signature
   - Change `fmt.Sprintf("bash %s", cfg.ScriptPath)` to `fmt.Sprintf("bash %s %s", cfg.ScriptPath, filePath)`
   - Update all call sites

2. Update `detector.go`:
   - In `ProcessFile()`, when calling `ExecuteScriptsForAgents`, pass `filePath` from `fileData.FilePath`

3. Update `notifier_test.go`:
   - Add test for script execution with file path argument
   - Verify the tmux command includes the file path

4. Run tests:
   - `go test ./...` - all tests pass
   - `go vet ./...` - passes with no warnings
   - `go build ./...` - succeeds without errors
   - `make build` - succeeds

5. Manual verification:
   - Start tmux session
   - Run monitor
   - Create task with assignee
   - Verify script log contains file path
```
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
<!-- DOD:END -->
