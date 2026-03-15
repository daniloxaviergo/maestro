---
id: GOT-024
title: 'Task 5: Create example agent configurations and documentation'
status: In Progress
assignee: []
created_date: '2026-03-15 18:54'
updated_date: '2026-03-15 23:08'
labels:
  - task
  - docs
  - agent
dependencies:
  - GOT-020
  - GOT-021
references:
  - >-
    /home/danilo/scripts/github/maestro/backlog/docs/PRD-Agent-Orchestration-System.md
priority: low
ordinal: 12437.5
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Task 5: Create agent example configuration files and documentation
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
<!-- DOD:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Example agent configurations created in agents/ directory with at least two example agents
- [ ] #2 Each agent has config.yml with script_path, tmux_session, and enabled fields
- [ ] #3 Example bash scripts created that can execute via tmux
- [ ] #4 Documentation added to docs/ explaining agent configuration format
- [ ] #5 Documentation added explaining how to create new agents
- [ ] #6 Documentation references examples from PRD
- [ ] #7 go vet passes with no warnings on all new files
- [ ] #8 go build succeeds without errors
- [ ] #9 make build succeeds
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task creates example agent configurations and documentation to support the agent orchestration system. The approach will be:

1. **Create example agent configuration directory structure** - Set up a working example in the `agents/` directory with at least two example agents (e.g., `agent-foo`, `agent-bar`)

2. **Create example configuration files** - Each agent gets a `config.yml` with realistic configuration values showing all available options

3. **Create example bash scripts** - Simple sample scripts that agents can execute, demonstrating the integration flow

4. **Update documentation** - Add comprehensive documentation in `docs/` explaining:
   - Agent configuration file format
   - How to create new agents
   - Directory structure expectations
   - Configuration options and defaults
   - How to test agent configurations

5. **Add examples to README/PRD** - Reference example configs in the PRD and project README

**Rationale for approach:**
- Examples follow existing project conventions (YAML config, tmux integration)
- Minimal but functional examples that work out of the box
- Documentation integrated with existing project structure

### 2. Files to Modify

**New files to create:**

| Path | Purpose |
|------|---------|
| `agents/agent-foo/config.yml` | Example agent configuration for "agent-foo" |
| `agents/agent-foo/script.sh` | Example bash script for agent-foo |
| `agents/agent-bar/config.yml` | Example agent configuration for "agent-bar" |
| `agents/agent-bar/script.sh` | Example bash script for agent-bar |
| `docs/agent-configuration.md` | Comprehensive documentation for agent configuration |
| `docs/agent-orchestration-quickstart.md` | Quick start guide for setting up agents |

**Files to update:**

| Path | Changes |
|------|---------|
| `docs/setup-monitor.md` | Add section on agent configuration and testing |
| `backlog/docs/PRD-Agent-Orchestration-System.md` | Add "Example Configurations" section linking to docs |

**No existing files modified** - This is purely additive work.

### 3. Dependencies

**Prerequisites (already in place):**
- [x] GOT-020: Agent Matching Engine - Matcher package exists and tested
- [x] GOT-021: Script Execution Routing - Notifier has `ExecuteScriptsForAgents()` method
- [x] GOT-015 to GOT-018: Agent configuration system - All packages implemented

**No external dependencies** - Uses standard Go YAML parser already in `go.mod`

**Setup required before implementation:**
- None - codebase is ready for documentation/examples

**Blocking issues:** None

### 4. Code Patterns

**Follow existing patterns:**

| Pattern | Example |
|---------|---------|
| YAML config format | Use `script_path`, `tmux_session`, `enabled` fields |
| Agent directory structure | `{config_dir}/{agent_name}/config.yml` |
| Script path | Relative to agent directory or absolute path |
| Tmux session | Default to "default" if not specified |
| Error handling | Log warnings, don't crash on misconfiguration |
| File permissions | Scripts should be executable (`chmod +x`) |

**Naming conventions:**
- Agent names: lowercase, hyphen-separated (e.g., `agent-foo`, `agent-bar`)
- Config files: `config.yml` (always)
- Script files: `script.sh` (convention, not enforced)

**Error handling:**
- Missing config: Return defaults, log warning
- Missing script: Log warning, skip execution
- Invalid YAML: Log warning, return defaults
- Tmux errors: Log warnings, continue processing

### 5. Testing Strategy

**Configuration validation:**

1. **YAML syntax validation** - Ensure config files are valid YAML using Go's yaml.v3 parser
2. **Path validation** - Verify script paths exist (or are empty for optional scripts)
3. **Directory structure validation** - Confirm `{agent}/config.yml` exists

**Integration testing:**

1. **Manual test with monitor:**
   ```bash
   # Start tmux session
   make tmux-start
   
   # Run monitor
   make run
   
   # Update a task file's assignee to "agent-foo"
   # Verify tmux notification appears and script runs
   ```

2. **Edge cases to test:**
   - Missing config file (should use defaults)
   - Missing script path (should skip execution)
   - Disabled agent (should log warning, skip)
   - Invalid YAML (should log warning, use defaults)

**Testing commands:**
```bash
# Build to verify no compile errors
make build

# Run static analysis
go vet ./...

# Run tests for related packages
go test ./pkg/matcher/...
go test ./pkg/agent/...
go test ./pkg/config/...
```

### 6. Risks and Considerations

**Known risks:**
1. **Permissions** - Scripts need execute permissions (`chmod +x`); documentation must clarify this
2. **Tmux availability** - Agent orchestration requires tmux; if not installed, scripts will fail
3. **Path resolution** - Relative paths in config are resolved from current working directory, not config file location

**Design trade-offs:**
1. **Simple examples** - Example scripts will be minimal (just `echo` statements) to avoid complexity
2. **No error recovery** - If script fails, agent continues; this is intentional (graceful degradation)
3. **No versioning** - Example configs won't have version fields; may be added in future

**Deployment considerations:**
- Example configs should use relative paths that work from project root
- Scripts should be documented as needing `chmod +x`
- Documentation should mention tmux is required for script execution

**Future enhancements (out of scope):**
- Agent health monitoring
- Script execution result reporting
- Multiple scripts per agent
- Script dependencies
<!-- SECTION:PLAN:END -->
