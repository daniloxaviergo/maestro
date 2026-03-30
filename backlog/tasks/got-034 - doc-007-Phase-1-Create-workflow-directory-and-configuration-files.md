---
id: GOT-034
title: '[doc-007 Phase 1] Create workflow directory and configuration files'
status: Done
assignee: []
created_date: '2026-03-30 12:25'
updated_date: '2026-03-30 12:36'
labels:
  - setup
  - infrastructure
dependencies: []
references:
  - doc-007#files-to-modify
  - doc-007#technical-decisions
documentation:
  - doc-007
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the agents/workflow/ directory structure and initialize configuration files for the workflow agent system. This includes creating an empty agents/workflow/tasks.yml file and setting up the default agents/workflow/config.yml with the agent sequence (Catarina, Thomas) and Backlog CLI command configuration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 agents/workflow/ directory exists
- [x] #2 agents/workflow/config.yml created with agent sequence
- [x] #3 agents/workflow/tasks.yml created (empty initial state)
- [x] #4 Both files use flat YAML format (single-line entries)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task establishes the foundational directory structure and configuration files for the workflow agent system as specified in doc-007 Phase 1.

**Implementation approach:**
- Create `agents/workflow/` directory using standard shell commands
- Create `agents/workflow/config.yml` with agent sequence configuration (Catarina, Thomas) and Backlog CLI command
- Create empty `agents/workflow/tasks.yml` for task state persistence
- Use flat YAML format (single-line entries) consistent with existing agent configurations
- Follow existing project conventions for configuration file structure

**Configuration structure:**
- `config.yml`: Defines agent sequence (order of downstream agents) and Backlog CLI command template
- `tasks.yml`: Empty initial state file that will hold task tracking data in Phase 2

**Why this approach:**
- Matches the flat YAML structure used in existing `agents/*/config.yml` files
- Simple directory creation ensures clean starting point for Phase 2 implementation
- Empty `tasks.yml` provides placeholder for task state persistence
- No complex validation needed at this stage (handled in Phase 2)

### 2. Files to Modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `agents/workflow/` | New directory for workflow agent files |
| Create | `agents/workflow/config.yml` | Workflow configuration with agent sequence and Backlog CLI command |
| Create | `agents/workflow/tasks.yml` | Empty task state file (initial placeholder) |

**No existing files need modification** - this is a pure addition of new infrastructure.

### 3. Dependencies

**Prerequisites (already satisfied):**
- ✅ `agents/` directory exists with other agent configurations
- ✅ Existing agent `config.yml` files serve as reference patterns
- ✅ Backlog CLI available for task assignment (used in agent scripts)
- ✅ Bash scripting environment ready

**No new external dependencies required.**

**Related tasks (not blocking):**
- GOT-035 (Phase 2): Requires these files to exist before implementing script.sh
- GOT-036 (Phase 3): Requires these files for testing workflow

### 4. Code Patterns

**From existing agent configurations to follow:**

1. **Agent configuration pattern** (`agents/catarina/config.yml`, `agents/thomas/config.yml`):
   - Use flat YAML with single-line entries
   - Key-value pairs for configuration
   - Comments for documentation

2. **Directory structure pattern:**
   - Each agent has its own subdirectory under `agents/`
   - Configuration stored as `config.yml`
   - Agent scripts stored as `script.sh`

3. **YAML format conventions:**
   - `key: value` format (no anchors/aliases)
   - Single-line per entry
   - No indentation for top-level keys

**Example configuration structure to replicate:**
```yaml
# Agent configuration
script_path: "./agents/agent-name/script.sh"
tmux_session: "agent-name"
enabled: true
```

**New configuration needs:**
- `agents` array for agent sequence
- `backlog_command` for task assignment
- `enabled` flag for workflow activation

### 5. Testing Strategy

**Validation checks (manual verification):**

1. **Directory creation verification:**
   ```bash
   # Verify directory exists
   test -d agents/workflow && echo "Directory exists" || echo "Directory missing"
   ```

2. **Configuration file verification:**
   ```bash
   # Check config.yml exists and contains expected keys
   test -f agents/workflow/config.yml && grep -q "agents:" agents/workflow/config.yml && echo "Config valid" || echo "Config missing keys"
   
   # Verify YAML syntax (basic check)
   grep -E "^[a-z_]+:" agents/workflow/config.yml
   ```

3. **Tasks file verification:**
   ```bash
   # Verify tasks.yml exists
   test -f agents/workflow/tasks.yml && echo "Tasks file exists" || echo "Tasks file missing"
   ```

4. **Format validation:**
   ```bash
   # Check for flat YAML (no complex indentation)
   head -n 20 agents/workflow/config.yml | grep -v "^  " | grep -v "^$"
   ```

**Test cases to verify:**
- `Directory exists`: `agents/workflow/` directory created successfully
- `Config file created`: `agents/workflow/config.yml` exists with correct content
- `Tasks file created`: `agents/workflow/tasks.yml` exists (can be empty)
- `Format correctness`: Both files use flat YAML format (no multi-line values)
- `Expected keys present`: Config contains `agents`, `backlog_command`, `enabled` keys

### 6. Risks and Considerations

**No blocking issues.** This is a straightforward setup task.

**Design considerations:**
1. **Empty tasks.yml**: The initial empty file is a placeholder. In Phase 2, the script.sh will initialize it with default structure on first run if it doesn't exist.
2. **Flat YAML format**: Chosen for simplicity in Bash string manipulation (no YAML parsing libraries available in Phase 2).
3. **Agent sequence order**: Catarina first, then Thomas. Reordering requires updating both config and potentially existing state files.
4. **Backlog CLI dependency**: The workflow depends on Backlog CLI being installed and configured. This is already installed per existing agent setup.
5. **Future extensibility**: The flat YAML format may require rework if complex state transitions need to be stored, but this is out of scope for doc-007.

**Potential pitfalls:**
- YAML syntax errors could break parsing in Phase 2 - ensure no tabs, consistent indentation
- Missing `enabled` flag may cause ambiguity - always include it for consistency with agent configs
- Hardcoded paths in config may need adjustment for different environments

**Deployment considerations:**
- No runtime impact - just creates configuration files
- Existing agent scripts remain unchanged
- No database or migration script required
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implement doc-007 Phase 1: Create workflow directory and configuration files

Changes:
- Created agents/workflow/ directory
- Created agents/workflow/config.yml with flat YAML format (agents: catarina, thomas; backlog_command: backlog task edit; enabled: true)
- Created agents/workflow/tasks.yml with empty initial state and comment documentation

Verification:
- All 4 acceptance criteria passed
- go build ./... - SUCCESS
- go vet ./... - SUCCESS (no warnings)
- make build - SUCCESS

Design decisions:
- Used flat YAML format with comma-separated agents list to satisfy single-line entries requirement
- Minimal config structure matching existing agent patterns
- Empty tasks.yml with comment for future state tracking

Risks:
- None - straightforward file creation task

Follow-ups:
- Phase 2 will implement script.sh for workflow orchestration
- Phase 3 will test workflow integration
<!-- SECTION:FINAL_SUMMARY:END -->

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
- [ ] #10 Errors are logged not silently ignored
- [ ] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->
