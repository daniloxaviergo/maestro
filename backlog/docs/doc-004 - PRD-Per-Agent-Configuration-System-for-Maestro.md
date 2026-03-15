---
id: doc-004
title: 'PRD: Per-Agent Configuration System for Maestro'
type: other
created_date: '2026-03-15 16:30'
---
# PRD: Per-Agent Configuration System for Maestro

## Overview

### Purpose
Enable each agent (monitor instance) to have its own YAML configuration file specifying which bash script to execute and which tmux session to use when assignee changes are detected, allowing for customized agent behavior.

### Goals
- **Goal 1**: Allow per-agent configuration via YAML files (bob.yml, alice.yml) stored in a configurable config directory
- **Goal 2**: When assignee change is detected, the agent should execute the configured bash script in its designated tmux session
- **Goal 3**: Graceful error handling - if script fails or doesn't exist, log error and continue monitoring
- **Goal 4**: Support multiple agents running on the same machine, each in separate tmux sessions

## Background

### Problem Statement
Currently, the Maestro file watcher has a single, hard-coded notification mechanism using `tmux display-message`. There is no way to:
- Customize which bash script runs when an assignee changes
- Execute commands in specific tmux sessions per agent
- Have different behavior for different agents monitoring the same files

### Current State
- All agents use the same notification logic via `notifier.Notify()`
- No configuration file support for agent-specific settings
- Bash script execution and tmux session targeting not supported
- Single global behavior for all agent instances

### Proposed Solution
Introduce a per-agent configuration system where:
1. Each agent has a YAML config file (e.g., `bob.yml`, `alice.yml`) defining:
   - Path to bash script to execute on assignee change
   - Tmux session name for script execution
2. Agent reads its config at startup (via environment variable or CLI flag)
3. On assignee change detection, agent executes configured script in its tmux session
4. Configuration files are stored in a configurable directory (default: `./agents/`)

## Requirements

### User Stories

- **Role**: System Administrator
  - *As a system administrator, I want to define agent-specific configuration files so that different agents can have different behaviors*
  
- **Role**: Developer
  - *As a developer, I want to specify which bash script runs on assignee changes so that I can customize notification methods per agent*
  
- **Role**: DevOps Engineer
  - *As a DevOps engineer, I want to specify the tmux session for script execution so that scripts run in isolated sessions per agent*

### Functional Requirements

#### Task 1: Configuration File Structure and Format

Define a YAML configuration file structure for agent configurations.

##### User Flows
1. Administrator creates `bob.yml` in the agents config directory
2. Configuration contains:
   - `script_path`: Path to bash script to execute
   - `tmux_session`: Name of tmux session for script execution
   - `enabled`: Boolean to enable/disable the agent (optional, default: true)
3. Agent reads configuration at startup
4. Agent validates configuration before starting

##### Acceptance Criteria
- [ ] Configuration file uses YAML format
- [ ] Must specify `script_path` field (string)
- [ ] Must specify `tmux_session` field (string)
- [ ] Optional `enabled` field with default value of `true`
- [ ] Agent reads config file at startup
- [ ] Agent validates config before starting

#### Task 2: Configuration File Loading and Discovery

Implement logic to discover and load agent configuration files.

##### User Flows
1. Agent starts and determines its identity (via environment variable `AGENT_NAME`)
2. Agent looks for config file at `{config_dir}/{agent_name}.yml`
3. Agent loads and parses the YAML configuration
4. If config file not found, agent uses defaults or logs warning and continues
5. Configuration directory path is configurable via environment variable `AGENTS_CONFIG_DIR`

##### Acceptance Criteria
- [ ] Config directory path configurable via `AGENTS_CONFIG_DIR` environment variable
- [ ] Default config directory is `./agents/`
- [ ] Agent name determined by `AGENT_NAME` environment variable
- [ ] Config file path is `{config_dir}/{agent_name}.yml`
- [ ] Missing config file logs warning but doesn't crash agent
- [ ] YAML parsing errors are caught and logged

#### Task 3: Bash Script Execution on Assignee Change

Execute the configured bash script when an assignee change is detected.

##### User Flows
1. Assignee change is detected in a markdown file
2. Agent reads its configuration for the configured script path
3. Agent executes the bash script in the configured tmux session
4. If script fails or doesn't exist, log error and continue monitoring

##### Acceptance Criteria
- [ ] Bash script executes when assignee change detected
- [ ] Script runs in the configured tmux session
- [ ] Script output is captured but not displayed (background execution)
- [ ] If script doesn't exist, log error and continue
- [ ] If script execution fails, log error and continue
- [ ] Script execution is non-blocking (doesn't delay monitoring)

#### Task 4: Tmux Session Targeting

Execute commands in the specified tmux session.

##### User Flows
1. Agent configures tmux session name from config file
2. When executing bash script, agent uses `tmux send-keys` to run in the session
3. Script output displayed in the tmux session

##### Acceptance Criteria
- [ ] Commands execute in the configured tmux session
- [ ] If tmux session doesn't exist, create it or log error
- [ ] `tmux send-keys` used for script execution in session
- [ ] Session name configurable via `tmux_session` field

### Non-Functional Requirements

- **Performance**: Configuration loading should not delay agent startup by more than 100ms
- **Security**: Bash script paths should be validated to prevent directory traversal attacks
- **Compatibility**: Works on Linux and macOS (tmux available)
- **Scalability**: Should support at least 10 agents running simultaneously on the same machine
- **Maintainability**: Configuration structure should be documented and easy to extend
- **Error Handling**: All configuration errors should be logged, never crash the agent

## Scope

### In Scope
- YAML configuration file structure (script_path, tmux_session, enabled)
- Configuration file discovery via environment variables
- Bash script execution on assignee change
- Tmux session targeting for script execution
- Graceful error handling (log and continue)
- Default values for optional configuration fields

### Out of Scope
- Configuration file validation with strict schema (future enhancement)
- Configuration file updates without agent restart (future enhancement)
- Web UI or CLI for managing agent configurations
- Centralized configuration management system
- Support for Windows (tmux not available)

## Technical Considerations

### Existing System Impact
- Requires modifications to `cmd/monitor/main.go` to load and use configuration
- Requires modifications to `pkg/change_detect/detector.go` to pass configuration to notifier
- Requires modifications to `pkg/notifier/notifier.go` to execute bash scripts in tmux sessions
- No changes to existing cache, parser, or watcher logic required

### Dependencies
- tmux (already available in existing system)
- YAML parser (gopkg.in/yaml.v3, already in use)

### Constraints
- Agent name must be unique per machine
- Tmux session names must be unique per machine
- Bash scripts must exist and be executable
- No dynamic configuration reloading (must restart agent to apply changes)

## Success Metrics

### Quantitative
- Configuration file parsing completes in < 100ms
- Script execution starts within 500ms of assignee change detection
- No agent crashes due to configuration errors

### Qualitative
- Administrators can easily create and modify agent configurations
- Scripts execute reliably in the correct tmux sessions
- Error messages clearly indicate configuration problems

## Timeline & Milestones

### Key Dates
- [Date]: Design complete - configuration structure finalized
- [Date]: Implementation complete - all tasks implemented
- [Date]: Testing complete - all acceptance criteria met
- [Date]: Documentation complete - user guide and examples provided

## Stakeholders

### Decision Makers
- [User]: Product Owner

### Contributors
- [Qwen Code]: Implementation
- [User]: Review and testing

## Appendix

### Glossary
- **Agent**: An instance of the Maestro file watcher with its own configuration
- **Tmux session**: A named tmux session where scripts are executed
- **Bash script**: Shell script executed when assignee changes are detected

### References
- [backlog://workflow/overview]: Backlog.md workflow overview
- [cmd/monitor/main.go]: Current monitor entry point
- [pkg/notifier/types.go]: Current notifier types
- [pkg/change_detect/detector.go]: Current change detection logic

### Example Configuration

```yaml
# agents/bob.yml
script_path: "/home/bob/scripts/notify.sh"
tmux_session: "bob-notifications"
enabled: true

# agents/alice.yml
script_path: "/home/alice/scripts/alert.sh"
tmux_session: "alice-alerts"
enabled: true
```
