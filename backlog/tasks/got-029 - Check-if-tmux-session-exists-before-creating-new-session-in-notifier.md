---
id: GOT-029
title: Check if tmux session exists before creating new session in notifier
status: To Do
assignee: []
created_date: '2026-03-16 11:47'
updated_date: '2026-03-16 11:48'
labels:
  - bug
  - tmux
  - notifier
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement a session existence check before creating new tmux sessions in the notifier package. Currently, the ExecuteScript and ExecuteScriptsForAgents methods attempt to create a new tmux session unconditionally with tmux new-session -d -s name, which fails with exit status 1 when the session already exists (e.g., when a user is attached to the session). The check should use tmux list-sessions to verify session existence before attempting creation.
<!-- SECTION:DESCRIPTION:END -->

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
- [ ] #13 - [ ] Code follows existing project conventions (package structure, naming, error handling)
- [ ] #14 - [ ] go vet ./... passes with no warnings
- [ ] #15 - [ ] go build ./... succeeds without errors
- [ ] #16 - [ ] Manual testing: session exists scenario (attach to session, run monitor, verify no error)
- [ ] #17 - [ ] Manual testing: session does not exist scenario (no session, run monitor, verify session created and script executes)
- [ ] #18 - [ ] Error handling: if 'tmux list-sessions' fails, log warning and skip session creation check (fallback to current behavior)
- [ ] #19 - [ ] Code comments added for non-obvious logic
- [ ] #20 - [ ] No breaking changes to existing functionality
<!-- DOD:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Given the tmux session does not exist, ExecuteScript creates the session successfully
- [ ] #2 Given the tmux session already exists, ExecuteScript does not attempt to create a new session and executes the script in the existing session
- [ ] #3 Given the tmux session does not exist, ExecuteScriptsForAgents creates the session successfully
- [ ] #4 Given the tmux session already exists, ExecuteScriptsForAgents does not attempt to create a new session and executes scripts in the existing session
- [ ] #5 When a user is attached to a tmux session and the monitor runs, no error is logged about session creation failure
- [ ] #6 The session existence check uses 'tmux list-sessions' command
- [ ] #7 All existing acceptance criteria for ExecuteScript and ExecuteScriptsForAgents remain valid
<!-- AC:END -->

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
<!-- SECTION:NOTES:END -->
