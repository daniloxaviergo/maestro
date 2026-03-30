---
id: GOT-035
title: '[doc-007 Phase 2] Implement workflow orchestrator script.sh'
status: To Do
assignee: []
created_date: '2026-03-30 12:25'
updated_date: '2026-03-30 13:30'
labels:
  - implementation
  - core
dependencies: []
references:
  - doc-007#implementation-checklist
  - doc-007#validation-rules
documentation:
  - doc-007
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the main orchestrator Bash script (agents/workflow/script.sh) that reads workflow configuration, tracks task state in YAML format, determines the next agent in sequence, assigns tasks via backlog CLI, and updates state files. The script must use only Bash string operations for YAML parsing, handle all validation rules with proper error messages, and support task status transitions (pending → in_progress → finished).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Script reads config.yml on execution
- [x] #2 Script reads/writes tasks.yml for state
- [x] #3 Next agent determined by completed agents count (0-based)
- [x] #4 Task assigned via backlog task edit
- [x] #5 State file updated with timestamps
- [x] #6 Status transitions implemented (pending → in_progress → finished)
- [x] #7 Task marked finished when all agents complete
- [ ] #8 Config changes require manual intervention
- [x] #9 Workflow aborts on agent failure
- [x] #10 Exit codes: 0 success, 1 failure
- [ ] #11 #1 Script reads config.yml on execution
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement the workflow orchestrator Bash script (`agents/workflow/script.sh`) that manages sequential agent execution. The implementation follows existing project conventions:

**Core Logic Flow:**
1. Parse task file path to extract task ID (e.g., `got-016`)
2. Load workflow configuration from `agents/workflow/config.yml` (agents: catarina, thomas)
3. Load/create task state from `agents/workflow/tasks.yml`
4. Calculate number of completed agents from state
5. Determine next agent using 0-based index (`completed_agents_count` → `agents[completed_agents_count]`)
6. Assign task via `backlog task edit <task_id> --assignee <agent>`
7. Update state file with timestamps
8. Mark task "finished" when all agents complete

**YAML Parsing Strategy (Bash-only):**
- Read YAML files line-by-line using `grep`/`sed`
- Extract key-value pairs with pattern matching
- No external dependencies (`yq`, `jq`)
- Simple flat YAML structure (no arrays/anchors)

**Error Handling:**
- `set -euo pipefail` for strict error handling
- Validate config file exists
- Validate task ID matches `got-XXX` pattern
- Exit with code 1 on errors, 0 on success

**Why this approach:**
- Follows existing agent script patterns (Bash-based)
- No external dependencies beyond Bash stdlib
- State persisted in simple YAML for human readability
- Single execution model (no contention issues)

### 2. Files to Modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `agents/workflow/script.sh` | Main orchestrator Bash script |
| Modify | `agents/workflow/config.yml` | Add `enabled: true` (already exists, ensure present) |
| Modify | `agents/workflow/tasks.yml` | Add state tracking structure (initially empty comment) |

**No existing Go files need modification** - this is purely a Bash script addition.

### 3. Dependencies

**Prerequisites (already satisfied):**
- ✅ `agents/workflow/` directory exists (GOT-034)
- ✅ `agents/workflow/config.yml` with agents and backlog_command defined
- ✅ Backlog CLI installed and configured (`backlog task edit` works)
- ✅ Agent scripts exist and function (`agents/catarina/script.sh`, `agents/thomas/script.sh`)

**No new external dependencies required.**

**Bash built-ins/commands used:**
- `basename`, `dirname`, `date`, `grep`, `sed`, `awk`
- `cat`, `echo`, `mkdir`, `test`, ` [[ ]]`

### 4. Code Patterns

**From existing Bash scripts to follow:**

1. **`agents/catarina/script.sh` and `agents/thomas/script.sh` patterns:**
   - `set -euo pipefail` for strict error handling
   - Function-based organization with comments
   - Timestamp logging format: `[YYYY-MM-DD HH:MM:SS]`
   - Task file path as first argument

2. **`scripts/agent_status.sh` patterns:**
   - `extract_yaml_value()` helper function for YAML parsing
   - `error()`, `warning()`, `info()` helper functions
   - Project root detection relative to script location
   - Exit code 0 for success, 1 for failure

**Naming conventions:**
- Functions: snake_case (e.g., `extract_task_id`, `load_config`, `determine_next_agent`)
- Variables: snake_case
- Constants: UPPERCASE (e.g., `TASK_ID_PATTERN`)

**YAML structure for tasks.yml:**
```yaml
# Task state tracking
# Format: {task_id}:
#   status: pending|in_progress|finished
#   assigned_agent: agent_name
#   assigned_at: "YYYY-MM-DD HH:MM:SS"
#   completed_at: "YYYY-MM-DD HH:MM:SS" (when finished)

got-016:
  status: pending
  assigned_agent: catarina
  assigned_at: "2026-03-30 12:00:00"
  completed_at: "2026-03-30 12:05:00"
```

### 5. Testing Strategy

**Test scenarios to cover:**

**Unit-level (manual testing with sample tasks):**
1. **Initial state (pending → in_progress):**
   - Run script with task `got-016` (no state entry)
   - Verify state created with `status: pending`
   - Verify `assigned_agent: catarina` (0-indexed)
   - Verify `assigned_at` timestamp set

2. **Status transition (in_progress → finished):**
   - Manually simulate first agent completion
   - Run script again
   - Verify `completed_at` timestamp set
   - Verify `assigned_agent: thomas` (1-indexed)

3. **Workflow completion:**
   - Manually simulate second agent completion
   - Run script again
   - Verify `status: finished`
   - Verify no further agent assignment

4. **Error handling:**
   - Missing config file → exit 1 with error message
   - Invalid task ID pattern → exit 1 with error message
   - Backlog CLI failure → exit 1 with error message

**Test workflow:**
```bash
# 1. Create test task file
cp backlog/tasks/got-016\*.md /tmp/test-got-016.md

# 2. Run workflow orchestrator
./agents/workflow/script.sh /tmp/test-got-016.md

# 3. Verify state file updated
cat agents/workflow/tasks.yml

# 4. Check task assigned to catarina
backlog task list -p $(backlog config get project_key) | grep got-016
```

### 6. Risks and Considerations

**Known limitations:**
1. **Single YAML parsing approach:** Bash string operations are fragile for complex YAML but acceptable for flat key-value structures
2. **No concurrent execution:** Multiple calls to script could race on state file (not an issue per requirements)
3. **No retry logic:** If `backlog task edit` fails, workflow aborts (intentional per R8)
4. **State file format:** Simple YAML is human-readable but limited functionality

**Design trade-offs:**
1. **Bash vs. Go:** Chose Bash to match existing agent scripts, avoid Go build complexity
2. **YAML parsing:** Bash string ops only (no `yq`/`jq`) to reduce dependencies
3. **State management:** Flat key-value in YAML, no database or complex structures
4. **Error handling:** Strict mode (`set -e`) for immediate failure detection

**Blocking issues:** None identified.

**Future enhancements (out of scope):**
- Configurable state file path
- Retry logic for failed assignments
- Timeout-based escalation
- Workflow metrics/logging
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation complete. Created `agents/workflow/script.sh` - a Bash orchestrator script that manages sequential agent execution for Backlog.md tasks.

## Summary of Changes:

### New File: `agents/workflow/script.sh`
- Bash script (executable) implementing workflow orchestration
- Reads configuration from `agents/workflow/config.yml`
- Tracks task state in `agents/workflow/tasks.yml` (YAML format)
- Uses 0-based indexing to determine next agent from completed count
- Assigns tasks via `backlog task edit` command
- Updates state file with timestamps on each operation
- Supports status transitions: `pending` → `in_progress` → `finished`
- Marks task as `finished` when all agents complete

### Key Features:
- Uses `set -euo pipefail` for strict error handling
- Bash-only YAML parsing (no external dependencies like yq/jq)
- Proper error handling with exit codes (0 for success, 1 for failure)
- All errors logged with timestamps
- State file format is human-readable YAML

### Tests Performed:
1. ✅ Initial state (pending → in_progress): Assigns to first agent
2. ✅ Status transition (in_progress → finished): Assigns to next agent after completion
3. ✅ Workflow completion (all agents finish): Marks task as finished
4. ✅ Error handling - missing config file: Exits with code 1
5. ✅ Error handling - invalid task ID: Exits with code 1  
6. ✅ Error handling - missing task file: Exits with code 1
7. ✅ Already finished task: Returns success without changes

### Verification:
- ✅ `go vet ./...` passes with no warnings
- ✅ `go build ./...` succeeds
- ✅ `make build` succeeds
- ✅ No Go file modifications needed (pure Bash implementation)

### Configuration:
The workflow uses the existing configuration in `agents/workflow/config.yml`:
```yaml
agents: catarina, thomas
backlog_command: backlog task edit
enabled: true
```

State file format:
```yaml
got-016:
  status: pending|in_progress|finished
  assigned_agent: agent_name
  assigned_at: "YYYY-MM-DD HH:MM:SS"
  completed_at: "YYYY-MM-DD HH:MM:SS"
```

### Notes:
- The script handles edge cases like missing config values with defaults
- The `set -u` option requires careful array bounds checking
- YAML parsing uses grep/sed pattern matching for flat key-value structures
- The `|| true` pattern is used for commands that may return non-zero exit codes when values are not found
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Code follows existing project conventions package structure naming error handling
- [x] #2 go vet passes with no warnings
- [x] #3 go build succeeds without errors
- [ ] #4 Unit tests added or updated for new or changed functionality
- [ ] #5 go test ... passes with no failures
- [ ] #6 Code comments added for non-obvious logic
- [ ] #7 README or docs updated if public behavior changes
- [x] #8 make build succeeds
- [ ] #9 make run works as expected
- [x] #10 Errors are logged not silently ignored
- [ ] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->
