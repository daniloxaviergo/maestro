---
id: GOT-029
title: Check if tmux session exists before creating new session in notifier
status: Done
assignee:
  - thomas
created_date: '2026-03-16 11:47'
updated_date: '2026-03-31 00:12'
labels:
  - bug
  - tmux
  - notifier
dependencies: []
priority: high
ordinal: 22000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement a session existence check before creating new tmux sessions in the notifier package. Currently, the ExecuteScript and ExecuteScriptsForAgents methods attempt to create a new tmux session unconditionally with tmux new-session -d -s name, which fails with exit status 1 when the session already exists (e.g., when a user is attached to the session). The check should use tmux list-sessions to verify session existence before attempting creation.
<!-- SECTION:DESCRIPTION:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Code follows existing project conventions package structure naming error handling
- [x] #2 go vet passes with no warnings
- [x] #3 go build succeeds without errors
- [x] #4 Unit tests added or updated for new or changed functionality
- [x] #5 go test ... passes with no failures
- [x] #6 Code comments added for non-obvious logic
- [x] #7 README or docs updated if public behavior changes
- [x] #8 make build succeeds
- [x] #9 make run works as expected
- [x] #10 Errors are logged not silently ignored
- [x] #11 Graceful degradation monitor continues if individual file processing fails
- [x] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Given the tmux session does not exist, ExecuteScript creates the session successfully
- [x] #2 Given the tmux session already exists, ExecuteScript does not attempt to create a new session and executes the script in the existing session
- [x] #3 Given the tmux session does not exist, ExecuteScriptsForAgents creates the session successfully
- [x] #4 Given the tmux session already exists, ExecuteScriptsForAgents does not attempt to create a new session and executes scripts in the existing session
- [x] #5 When a user is attached to a tmux session and the monitor runs, no error is logged about session creation failure
- [x] #6 The session existence check uses 'tmux list-sessions' command
- [x] #7 All existing acceptance criteria for ExecuteScript and ExecuteScriptsForAgents remain valid
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan

### 1. Technical Approach

Implement a session existence check before creating new tmux sessions in the notifier package. The approach uses `tmux list-sessions` to check if a session already exists before attempting to create it.

**Key Changes:**
1. Create a `sessionExists(sessionName string) (bool, error)` helper method that:
   - Executes `tmux list-sessions` to retrieve all sessions
   - Parses output to check if the session name exists (matches start of line format `sessionName:`)
   - Returns `true` if found, `false` otherwise
   - Returns error if command fails (for graceful fallback)

2. Update `ExecuteScript` to:
   - Call `sessionExists()` before creating session
   - Only execute `tmux new-session` if session does not exist
   - Log when session already exists (debug level)

3. Update `executeScriptForAgent` to apply the same pattern

**Design Decision:** Use `tmux list-sessions` parsing instead of `tmux has-session` for broader tmux version compatibility and more predictable output format.

### 2. Files to Modify

**Modified Files:**
- `pkg/notifier/notifier.go`:
  - Add `sessionExists(sessionName string) (bool, error)` method
  - Update `ExecuteScript` to check session existence before creation
  - Update `executeScriptForAgent` to check session existence before creation
  - Add debug log when session already exists

**No Changes Required:**
- `pkg/notifier/types.go` - No type changes needed
- `pkg/notifier/notifier_test.go` - Add unit tests for `sessionExists()` function

### 3. Dependencies

**No New Dependencies:**
- Uses existing `os/exec` package for command execution
- Uses `strings` package for output parsing
- No external YAML or configuration changes

**Prerequisites:**
- Verify `tmux` is installed (existing check via `ErrTmuxNotInstalled`)
- Verify `tmux list-sessions` command works (graceful degradation if command fails)

### 4. Code Patterns

**Error Handling:**
- If `tmux list-sessions` fails, log warning and skip session creation check (fallback to current behavior)
- Log when session already exists using `log.Printf` (consistent with `executeScriptForAgent`)
- Continue execution if session exists (do not return error)

**Session Name:**
- Use `cfg.TmuxSession` with fallback to `"default"`
- Consistent with existing code pattern

**Logging:**
- Debug level when session already exists: `log.Printf("Session %s already exists, skipping creation", sessionName)`
- Warning level when `tmux list-sessions` fails: `log.Printf("Warning: failed to list tmux sessions: %v", err)`

**Code Style:**
- Match existing function naming: `executeScriptForAgent` (lowercase helper)
- Use `log.Printf` for non-critical warnings (consistent with agent orchestration code)
- No changes to error types (existing `ErrSessionCreationFailed` still used)

### 5. Testing Strategy

**Unit Tests to Add:**
- `TestSessionExists_SessionDoesNotExist` - Verify returns `false, nil` for non-existent session
- `TestSessionExists_SessionExists` - Verify returns `true, nil` for existing session
- `TestSessionExists_TmuxNotInstalled` - Verify graceful error handling

**Integration Tests (already covered by acceptance criteria):**
1. Start tmux session manually, run monitor, verify no session creation error
2. No tmux session, run monitor, verify session is created successfully

**Edge Cases:**
- Session does not exist (should create)
- Session exists and user is attached (should not fail)
- Session exists but is detached (should not fail)
- `tmux list-sessions` command fails (should fallback gracefully)
- Multiple agents with same session name (check once per session, not per script)

**Verification Commands:**
```bash
# Test with session exists
tmux new-session -d -s default
make build && ./monitor
# Check no session creation error in output

# Test with session does not exist
tmux kill-session -t default 2>/dev/null || true
make build && ./monitor
# Verify session is created and script executes
```

### 6. Risks and Considerations

**Known Risks:**

1. **tmux list-sessions parsing**: Output format may vary across tmux versions
   - *Mitigation*: Parse using `strings.HasPrefix` to match `sessionName:` at line start
   - *Alternative fallback*: If parsing fails, log warning and skip check

2. **Race condition**: Session could be destroyed between check and creation
   - *Mitigation*: Accept as edge case; existing error handling catches this

3. **Performance**: Extra command execution for every script run
   - *Mitigation*: Minimal overhead; session existence is checked once per agent (not per script in batch)

**No Blocking Issues**: Implementation is straightforward and follows existing error handling patterns in the codebase.

**Testing Considerations:**
- Tests require tmux to be installed (existing test infrastructure handles this)
- Unit tests may need to mock `os/exec` for isolation
- Integration tests should verify actual tmux behavior
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
### 1. Technical Approach

Modify `pkg/notifier/notifier.go` to implement a session existence check before creating new tmux sessions:

**Session Check Function:**
- Create a helper method `sessionExists(sessionName string) (bool, error)`
- Use `tmux list-sessions` to list all sessions
- Parse output to check if session with given name exists
- Return `true` if found, `false` otherwise
- Handle errors gracefully (e.g., tmux not installed, command fails)

**Update ExecuteScript:**
- Before calling `tmux new-session -d -s sessionName`, check if session exists
- Only create session if it does not exist
- Log when session already exists (debug level)

**Update executeScriptForAgent:**
- Apply same session existence check pattern
- Ensure consistency between both methods

**Architecture Decision:** Use `tmux list-sessions` with string parsing instead of `tmux has-session` for compatibility across tmux versions.

### 2. Files to Modify

**Modified Files:**
- `pkg/notifier/notifier.go` - Add session existence check logic
  - Add `sessionExists(sessionName string) (bool, error)` helper method
  - Update `ExecuteScript` to check before creating session
  - Update `executeScriptForAgent` to check before creating session
  - Update error message to clarify session already exists scenario

**No Changes Required:**
- `pkg/notifier/types.go` - No type changes needed
- `pkg/notifier/notifier_test.go` - Add tests for session existence check if tests exist
- No changes to agent configuration format
- No changes to CLI or Makefile

### 3. Dependencies

**No New Dependencies:**
- Uses existing `os/exec` package
- No external YAML or configuration changes
- Uses standard tmux commands (`list-sessions`, `send-keys`)

**Prerequisites Check:**
- Verify `tmux` is installed (existing check in `types.go` with `ErrTmuxNotInstalled`)
- Verify `tmux list-sessions` command works (graceful degradation if command fails)

### 4. Code Patterns

**Error Handling:**
- If `tmux list-sessions` fails, log warning and skip session creation check (fallback to current behavior)
- Log when session already exists (use `log.Printf` for consistency)
- Continue execution if session exists (do not return error)

**Session Name:**
- Use `cfg.TmuxSession` with fallback to `"default"`
- Consistent with existing code pattern

**Logging:**
- Debug level when session already exists: `log.Printf("Session %s already exists, skipping creation", sessionName)`
- Warning level when `tmux list-sessions` fails: `log.Printf("Warning: failed to list tmux sessions: %v", err)`

### 5. Testing Strategy

**Manual Testing Steps:**
1. Start a tmux session: `tmux new-session -d -s test-session`
2. Run the monitor: `make run`
3. Modify a task file to change assignee
4. Verify no error is logged about session creation failure
5. Verify script executes in the existing session

**Edge Cases to Cover:**
- Session does not exist (should create)
- Session exists and user is attached (should not fail)
- Session exists but is detached (should not fail)
- `tmux list-sessions` command fails (should fallback gracefully)
- Multiple agents with same session name (should check once per agent, not once per script execution)

**Verification Commands:**
```bash
# Test with session exists
tmux new-session -d -s default
make run
# Check no session creation error in output

# Test with session does not exist
tmux kill-session -t default 2>/dev/null || true
make run
# Verify session is created and script executes
```

### 6. Risks and Considerations

**Known Risks:**

1. **tmux list-sessions parsing**: Parsing the output of `tmux list-sessions` may be fragile if output format changes
   - *Mitigation*: Use robust string matching (e.g., check if session name appears at start of line)
   - *Alternative*: Use `tmux has-session -t <name>` which returns exit code 0/1

2. **Race condition**: Session could be destroyed between check and creation
   - *Mitigation*: Accept this as edge case; error handling already covers this scenario

3. **Performance**: Extra command execution for every script run
   - *Mitigation*: Session existence is cached per script execution; minimal overhead

**No Blocking Issues**: Implementation is straightforward and follows existing error handling patterns in the codebase.

Implementation note: README/doc updates not required as this is an internal implementation change with no public API changes.

Verification: make build succeeds and all tests pass. make run works correctly with session existence check in place.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Implemented session existence check before creating tmux sessions in the notifier package to prevent errors when sessions already exist.

### What Changed

**Modified Files:**
- `pkg/notifier/notifier.go`: Added `sessionExists()` helper function and updated `ExecuteScript()` and `executeScriptForAgent()` to check session existence before creating new sessions
- `pkg/notifier/notifier_test.go`: Added unit tests `TestSessionExists` and `TestSessionExists_SessionNameParsing`

### Implementation Details

1. **sessionExists()**: A helper function that uses `tmux list-sessions` to check if a session exists by parsing the output format `sessionName:windows=...`

2. **ExecuteScript()**: Now calls `sessionExists()` before attempting to create a session. If the session exists, it logs a message and skips creation.

3. **executeScriptForAgent()**: Same pattern as ExecuteScript() - checks session existence first, logs "Session already exists" when no creation is needed.

4. **Error Handling**: If `tmux list-sessions` fails, logs a warning and continues with session creation (graceful degradation).

### Why This Fix

Previously, when a user was attached to a tmux session or a session already existed from a previous run, the monitor would attempt to create a new session with `tmux new-session -d -s name` which would fail with exit status 1, causing errors in the logs.

With this fix:
- If a session doesn't exist, it's created normally
- If a session exists (even with user attached), the script executes in the existing session without error
- No more "tmux session creation failed" errors when the monitor runs while a user is attached to a session

### Tests

All tests pass including:
- New unit tests for `sessionExists()` function
- Existing tests continue to pass, confirming backward compatibility
- Integration tests verify the behavior with existing sessions

### Verification

- `go vet ./...` - no warnings
- `go test ./...` - all tests pass
- `make build` - builds successfully
<!-- SECTION:FINAL_SUMMARY:END -->
