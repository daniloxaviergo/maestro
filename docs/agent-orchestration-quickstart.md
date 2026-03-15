# Agent Orchestration Quickstart

## Overview

This quickstart guide will help you create a new agent and get it working with the Maestro file monitor.

## Prerequisites

- Maestro monitor is running (see [setup-monitor.md](setup-monitor.md))
- tmux is installed (`tmux --version` to verify)
- Project directory writable

## Step 1: Create Agent Directory

Create a new directory for your agent in the `agents/` folder:

```bash
mkdir -p agents/my-agent
```

## Step 2: Create Configuration File

Create `agents/my-agent/config.yml` with your agent's settings:

```yaml
script_path: "./agents/my-agent/script.sh"
tmux_session: "my-agent"
enabled: true
```

### Configuration Reference

| Field | Default | Description |
|-------|- ------|-------------|
| `script_path` | (empty) | Path to the bash script to execute |
| `tmux_session` | "default" | Tmux session name for script execution |
| `enabled` | false | Whether the agent processes tasks |

## Step 3: Create Script

Create `agents/my-agent/script.sh`:

```bash
#!/bin/bash
# My-agent script
TASK_FILE="$1"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

echo "[$TIMESTAMP] my-agent: Processing task: $TASK_FILE"

# Add your custom logic here
# Example: Parse the task, run commands, send notifications

echo "[$TIMESTAMP] my-agent: Task processing complete"
```

Make it executable:

```bash
chmod +x agents/my-agent/script.sh
```

## Step 4: Start Tmux Session

Start a tmux session for your agent:

```bash
tmux new-session -d -s my-agent
```

## Step 5: Test Your Agent

1. Start the monitor in one terminal:
   ```bash
   make run
   ```

2. In another terminal, create a test task with your agent assigned:

   ```bash
   cat > backlog/tasks/test-task.md <<EOF
---
id: TEST-001
title: Test Task
assignee: ["my-agent"]
status: To Do
---

This is a test task for agent orchestration.
EOF
   ```

3. Watch the monitor output and your tmux session for the script execution

## Step 6: Verify

Check the tmux session output:

```bash
tmux capture-pane -p -t my-agent
```

Or attach to see live output:

```bash
tmux attach-session -t my-agent
```

## Agent Matching Rules

- Agent names are matched case-insensitively
- The assignee field in the task YAML must match the agent directory name
- Multiple agents can be assigned to a single task
- If no matching agent is found, a warning is logged and processing continues

## Troubleshooting

### Agent Not Responding

1. Check the monitor logs for warnings about missing agents
2. Verify the `config.yml` exists at `agents/{agent_name}/config.yml`
3. Ensure `enabled: true` is set in the config

### Script Not Executing

1. Verify `script_path` is set in `config.yml`
2. Ensure the script file exists and is executable (`chmod +x`)
3. Check tmux session exists: `tmux ls`

### Path Issues

- Relative paths in `script_path` are resolved from the project root
- Use absolute paths if you need guaranteed resolution

## Example Agents

For working examples, see:
- `agents/agent-foo/` - Simple echo-based agent
- `agents/agent-bar/` - Logging-based agent

## Environment Configuration

Set environment variables to customize agent behavior:

```bash
# Override default agents directory
export AGENTS_CONFIG_DIR="/path/to/custom/agents"

# Run monitor with custom config directory
AGENTS_CONFIG_DIR="/path/to/custom/agents" make run
```

## Next Steps

- See [agent-configuration.md](agent-configuration.md) for detailed configuration options
- See [PRD-Agent-Orchestration-System](../backlog/docs/PRD-Agent-Orchestration-System.md) for technical architecture
