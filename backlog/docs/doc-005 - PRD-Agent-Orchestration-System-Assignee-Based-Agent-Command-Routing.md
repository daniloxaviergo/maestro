---
id: doc-005
title: 'PRD: Agent Orchestration System - Assignee-Based Agent Command Routing'
type: other
created_date: '2026-03-15 18:50'
---
# PRD: Agent Orchestration System - Assignee-Based Agent Command Routing

## Overview

### Purpose
This system automatically detects when a backlog task is assigned to a specific agent and sends commands to that agent's tmux session to execute its designated bash script, enabling hands-free task delegation and execution across multiple agents.

### Goals
1. **Automated Agent Assignment**: When a task's `assignee` field matches an agent name, automatically trigger that agent's script execution
2. **Multi-Agent Support**: Support multiple agents, each with their own configuration and tmux session
3. **Real-Time Processing**: Detect assignee changes in real-time using file system monitoring
4. **Error Resilience**: Continue monitoring and processing even when individual agent configurations are missing or scripts fail

## Background

### Problem Statement
Currently, the Maestro project can detect assignee changes in backlog.md tasks and send notifications, but it does not:
- Route commands to specific agents based on their name
- Load agent-specific configurations from config.yml files
- Execute agent-specific bash scripts in their tmux sessions

The existing `notifier.ExecuteScript()` method exists but is not integrated into the assignee change detection flow.

### Current State
- File watcher detects markdown file changes
- Change detector identifies assignee changes in YAML frontmatter
- Notifier can send tmux messages and execute scripts
- Agent configuration loading exists (`pkg/agent`, `pkg/config`)
- Scripts can be executed via tmux in a session

**Gaps:**
1. No automatic routing of assignee changes to specific agents
2. No agent name matching logic
3. No orchestration layer to connect detection → agent lookup → script execution

### Proposed Solution
Add an orchestration layer that:
1. Monitors assignee changes in backlog tasks
2. Matches the assignee name to an agent's configured name
3. Loads the agent's config.yml from `{config_dir}/{agent_name}/config.yml`
4. Executes the agent's script in their tmux session

## Requirements

### User Stories

- **Developer**:
  - *As a developer, I want to assign a task to an agent by name so that the agent's script automatically executes in their tmux session*
  
- **Agent Operator**:
  - *As an agent operator, I want to configure which bash script runs for my agent so that I can customize my task processing pipeline*
  
- **System Administrator**:
  - *As a system administrator, I want the orchestration to handle missing agent configurations gracefully so that other agents continue to work*

### Functional Requirements

#### Task 1: Agent Matching Engine
Create a module that matches assignee names to configured agents.

##### User Flows
1. Assignee change detected in task file (e.g., `assignee: ["agent-foo"]`)
2. System searches for agent with matching name in agent configs
3. If found, retrieve agent's script path and tmux session
4. If not found, log warning and skip

##### Acceptance Criteria
- [ ] Assignee name (case-insensitive) matches agent config's `name` field
- [ ] Agent configuration loaded from `{config_dir}/{agent_name}/config.yml`
- [ ] Warning logged when no matching agent found
- [ ] Support multiple assignees; execute script for all matching agents

#### Task 2: Agent Configuration Loading
Implement agent configuration loading with fallback defaults.

##### User Flows
1. System reads `AGENTS_CONFIG_DIR` environment variable (default: `./agents`)
2. For each agent, checks `{config_dir}/{agent_name}/config.yml`
3. Load configuration with defaults:
   - `script_path`: empty string (skip execution if not set)
   - `tmux_session`: "default"
   - `enabled`: false

##### Acceptance Criteria
- [ ] Agent config loaded from `{config_dir}/{agent_name}/config.yml`
- [ ] Defaults applied for missing config fields
- [ ] Config file parse errors handled gracefully
- [ ] Missing config file returns default config

#### Task 3: Script Execution Routing
Integrate agent matching with script execution in the notify chain.

##### User Flows
1. Assignee change detected → Agent matched
2. Agent config loaded
3. If agent enabled and script_path configured:
   - Execute `bash {script_path}` in agent's tmux session
4. Errors logged, processing continues

##### Acceptance Criteria
- [ ] Script executes only when agent enabled and script_path configured
- [ ] tmux session created if it doesn't exist
- [ ] Script execution is non-blocking (goroutine)
- [ ] Errors logged but don't crash the system
- [ ] Timeout handling for long-running scripts

#### Task 4: Agent Configuration Directory Structure
Define and enforce the agent config directory structure.

##### User Flows
1. Agent configs stored in `{config_dir}/{agent_name}/config.yml`
2. Example:
   ```
   ./agents/
   ├── agent-foo/
   │   └── config.yml
   ├── agent-bar/
   │   └── config.yml
   ```

##### Acceptance Criteria
- [ ] Directory structure created if missing
- [ ] Config file validation on load
- [ ] Clear error messages for invalid structure

### Non-Functional Requirements

- **Performance**: Assignee detection and agent lookup should add < 100ms per file event
- **Concurrency**: Support multiple simultaneous agent script executions
- **Reliability**: Agent missing or script failing should not stop monitoring
- **Maintainability**: Clear separation between detection, matching, and execution
- **Observability**: All agent matching decisions logged with debug level

## Scope

### In Scope
1. Agent matching logic (assignee name → agent config)
2. Agent config loading with defaults
3. Script execution routing in the event flow
4. Error handling and logging
5. Documentation for agent configuration

### Out of Scope
1. Agent registration/dynamic agent discovery
2. Agent health monitoring
3. Script execution result reporting
4. Multi-step workflow orchestration
5. Resource usage limits for agent scripts

## Technical Considerations

### Existing System Impact
- Modify `pkg/change_detect/detector.go` to trigger agent orchestration
- Extend `pkg/notifier/notifier.go` to route to specific agents
- No changes to watcher or parser

### Dependencies
- **tmux**: Required for script execution
- **fsnotify**: Already in use for file watching
- **YAML parser**: Already in use for config files

### Constraints
- Agent names must match exactly (case-insensitive) with assignee field
- Script path is relative to agent config directory or absolute
- No inter-agent communication or coordination

## Success Metrics

### Quantitative
- Agent matching latency: < 100ms per event
- Script execution success rate: > 99%
- No performance degradation with > 10 agents

### Qualitative
- Developers can assign tasks to agents without manual intervention
- Agent operators can configure their scripts without code changes
- System continues operating when individual agents misconfigured

## Timeline & Milestones

### Key Dates
- [2026-03-16]: Agent matching engine design complete
- [2026-03-17]: Agent config loading implementation
- [2026-03-18]: Script execution routing integration
- [2026-03-19]: Testing and documentation
- [2026-03-20]: Launch

## Stakeholders

### Decision Makers
- Product Manager: Approve agent naming conventions and config structure

### Contributors
- Backend Developer: Implement orchestration logic
- DevOps: Configure agent environments and test

## Appendix

### Glossary
- **Agent**: A configured entity that executes bash scripts when assigned tasks
- **Tmux Session**: A terminal multiplexer session where agent scripts run
- **Assignee**: The agent name specified in a task's YAML frontmatter

### References
- [pkg/agent/agent.go](./pkg/agent/agent.go): Agent identity and config loading
- [pkg/notifier/notifier.go](./pkg/notifier/notifier.go): Tmux notification and script execution
- [cmd/monitor/main.go](./cmd/monitor/main.go): File monitoring entry point
- [pkg/config/config.go](./pkg/config/config.go): Configuration loading
- [pkg/change_detect/detector.go](./pkg/change_detect/detector.go): Assignee change detection
