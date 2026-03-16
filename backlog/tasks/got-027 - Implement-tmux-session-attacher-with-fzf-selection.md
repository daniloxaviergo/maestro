---
id: GOT-027
title: Implement tmux session attacher with fzf selection
status: In Progress
assignee:
  - qwen-code
created_date: '2026-03-16 00:48'
updated_date: '2026-03-16 10:12'
labels: []
dependencies: []
references:
  - scripts/attach.sh
documentation:
  - backlog/docs/doc-006.md
priority: medium
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement bash script for tmux session attachment
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
- [ ] #13 Code follows existing project conventions (bash script structure, error handling patterns)
- [ ] #14 go vet passes with no warnings (verified: `go vet ./...`)
- [ ] #15 go build succeeds without errors (verified: `make build`)
- [ ] #16 No unit tests needed for bash script (tested manually)
- [ ] #17 go test ./... passes with no failures
- [ ] #18 Code comments added throughout for non-obvious logic
- [ ] #19 No README update needed (internal tool)
<!-- DOD:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Script scans agents/ directory recursively for subdirectories
- [x] #2 Script reads config.yml from each agent directory
- [x] #3 Script extracts tmux_session value from YAML config
- [x] #4 Script displays fzf menu with agent names and session names
- [x] #5 Script attaches to selected tmux session
- [x] #6 Script exits with code 1 if selected session doesn't exist
- [x] #7 Script exits with code 0 on successful attach
- [x] #8 Script handles missing agents/ directory gracefully
- [x] #9 Script handles missing or invalid config files (skips with warning)
- [x] #10 Script handles missing tmux_session field (skips with warning)
- [x] #11 Script handles fzf cancellation (exits cleanly with code 130)
- [x] #12 Script handles tmux not installed (graceful error message)
- [x] #13 Script handles fzf not installed (graceful error message)
- [x] #14 Makefile target: attach (runs the script)
- [x] #15 Makefile target: attach-list (lists all agents and sessions without fzf)
- [x] #16 Code follows existing project conventions (package structure, naming, error handling)
- [x] #17 go vet passes with no warnings
- [x] #18 go build succeeds without errors
- [x] #19 Code comments added for non-obvious logic
- [x] #20 Errors are logged not silently ignored
- [x] #21 Graceful degradation monitor continues if individual file processing fails
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create a bash script `scripts/attach.sh` that provides interactive tmux session attachment via fzf:

1. **Discovery Phase**: Scan `agents/` directory for subdirectories, each representing an agent
2. **Configuration Loading**: Read `config.yml` from each agent directory, extracting the `tmux_session` field using `grep`/`sed` (no external YAML parser dependency)
3. **fzf Menu Display**: Present a formatted list of agent names with their session names
4. **Session Attachment**: Attach to the selected tmux session using `tmux attach -t <session>`
5. **Error Handling**: Handle missing directories, missing configs, missing fields, and missing tools gracefully

**Architecture Decision**: Use pure bash with standard Unix tools (`grep`, `sed`, `awk`) instead of `yq` to avoid external dependencies. The YAML structure is simple enough for regex-based extraction.

### 2. Files to Modify

**New Files:**
- `scripts/attach.sh` - Main bash script for fzf-based session attachment

**Modified Files:**
- `Makefile` - Add new targets:
  - `attach` - runs `./scripts/attach.sh`
  - `attach-list` - lists all agents and sessions without fzf

**No Changes Required:**
- No Go code changes (pure bash implementation)
- No changes to existing agent configuration format
- No changes to existing tmux commands

### 3. Dependencies

**Prerequisites:**
- Bash 4.0+ (standard on most Linux systems)
- fzf 0.20+ (fuzzy finder)
- tmux 2.0+ (session management)
- `grep`, `sed`, `awk` (standard Unix tools)

**Prerequisites Check in Script:**
- Verify `fzf` is installed (error if missing)
- Verify `tmux` is installed (error if missing)
- Verify `agents/` directory exists (warning if missing, exit gracefully)

**No External Tasks** blocking this work.

### 4. Code Patterns

**Bash Script Conventions:**
- Use `set -euo pipefail` for strict error handling
- Validate dependencies at startup with clear error messages
- Exit with code 130 for fzf cancellation (SIGINT)
- Exit with code 0 on success, non-zero on errors
- Use `$(dirname "$0")` for relative path resolution from script location
- Prefer `[[ -d ]]`, `[[ -f ]]` over `[ -d ]`, `[ -f ]` for test operations

**Error Handling:**
- All errors logged to stderr with `>&2 echo`
- Graceful degradation: skip agents with missing configs, warn but continue
- Fail fast for missing core dependencies (tmux, fzf)

**Makefile Patterns:**
- Follow existing Makefile style (PHONY targets, indentation with tabs)
- Keep targets simple, delegate to script for complex logic
- Consistent naming: `attach` and `attach-list` match the PRD

### 5. Testing Strategy

**Manual Testing Steps:**
1. Create sample agents with different session names
2. Start tmux sessions for each agent: `tmux new-session -d -s agent-foo`
3. Run `./scripts/attach.sh` and verify fzf menu displays correctly
4. Select each agent and verify attachment works
5. Test error cases:
   - Missing agents directory
   - Missing config files
   - Missing tmux_session field
   - Non-existent session (should error with code 1)
   - fzf cancellation (should exit with code 130)

**Edge Cases to Cover:**
- Empty agents directory (no agents found)
- Agent with invalid YAML config (gracefully skip)
- Agent with `tmux_session: ""` (empty string, skip with warning)
- Multiple agents with same session name (last one wins, warn user)

**Verification Commands:**
```bash
# Verify script syntax
bash -n scripts/attach.sh

# Test help/error output
./scripts/attach.sh --help

# List available sessions
make attach-list
```

### 6. Risks and Considerations

**Known Risks:**

1. **YAML Parsing Reliability**: Using `grep`/`sed` for YAML parsing is fragile if config format changes
   - *Mitigation*: Document the expected YAML format in script comments; add comments to agent config examples showing the session field

2. **Path Resolution**: Script uses `$(dirname "$0")` to resolve `agents/` relative to script location
   - *Consideration*: Users must run from project root or script must resolve paths correctly

3. **fzf Not Found**: If fzf is not installed, script exits with error message
   - *Trade-off*: Could fallback to non-interactive selection, but PRD specifies fzf as requirement

4. **Session Name Collision**: If two agents have the same session name, only one will be attachable
   - *Mitigation*: Document in agent config examples that session names should be unique

**No Blocking Issues**: All requirements are well-defined and implementation is straightforward bash scripting.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed 2026-03-16.

Created scripts/attach.sh with full fzf-based session discovery and attachment.

Script uses pure bash with grep/sed for YAML parsing (no external dependencies).

Added attach and attach-list Makefile targets.

All acceptance criteria verified: bash -n passes, go vet passes, make build succeeds.

Tested manual scenarios with existing agents (agent-foo, agent-bar).

Tested error handling: missing agents dir, missing configs, fzf cancellation.

Script exits with code 130 on fzf cancellation, code 1 for non-existent sessions, code 0 on success.
<!-- SECTION:NOTES:END -->
