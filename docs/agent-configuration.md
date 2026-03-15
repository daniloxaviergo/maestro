# Agent Configuration Documentation

## Overview

The Maestro agent orchestration system allows you to configure agents that automatically execute scripts when tasks are assigned to them. Each agent has its own configuration and script files.

## Directory Structure

```
agents/
├── agent-foo/
│   ├── config.yml      # Agent configuration
│   └── script.sh       # Script to execute (must be executable)
├── agent-bar/
│   ├── config.yml
│   └── script.sh
└── ...
```

## Configuration Format

Each agent's `config.yml` file uses YAML format with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `script_path` | string | No | Path to the bash script to execute. Can be relative (to project root) or absolute. If empty, script execution is skipped. |
| `tmux_session` | string | No | Tmux session name where the script will run. Defaults to "default" if not specified. |
| `enabled` | boolean | No | Whether the agent is active. Defaults to `false` if not specified. |

### Example Configuration

```yaml
script_path: "./agents/agent-foo/script.sh"
tmux_session: "agent-foo"
enabled: true
```

### Minimal Configuration (using defaults)

```yaml
script_path: "./agents/agent-foo/script.sh"
enabled: true
# tmux_session defaults to "default"
```

## Script Requirements

Scripts must:
1. Be executable (`chmod +x script.sh`)
2. Accept the task file path as the first argument (`$1`)
3. Be bash-compatible scripts

### Script Environment

- The script runs in the tmux session specified in `tmux_session`
- The working directory is the project root
- The task file path is passed as the first argument

### Example Script

```bash
#!/bin/bash
TASK_FILE="$1"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

echo "[$TIMESTAMP] Processing task: $TASK_FILE"
# Add your custom logic here
```

## Agent Activation

For an agent to execute its script:

1. The agent's `config.yml` must exist at `{config_dir}/{agent_name}/config.yml`
2. The `enabled` field must be `true` (or omitted, defaults to `false`)
3. The `script_path` must be set (if empty, execution is skipped)
4. The script file must exist and be executable

## Environment Variables

- `AGENTS_CONFIG_DIR`: Set the agents configuration directory (default: `./agents`)
- `AGENT_NAME`: Used to identify the agent when loading its config (loaded at runtime)

## Error Handling

- Missing config file: Agent is skipped with a warning log
- Invalid YAML: Default values are used with a warning
- Missing script: Script execution is skipped with a warning
- Script execution errors: Logged but do not stop the monitor

## Example Agents

See the example agents in `agents/agent-foo/` and `agents/agent-bar/` for working configurations.

## Configuration Validation

The system validates:
- YAML syntax of `config.yml`
- Script path existence (only logs warning if missing)
- File permissions (scripts should be executable)

## Next Steps

- See [Agent Orchestration Quickstart](agent-orchestration-quickstart.md) for step-by-step setup
- See [PRD-Agent-Orchestration-System](../backlog/docs/PRD-Agent-Orchestration-System.md) for technical details
