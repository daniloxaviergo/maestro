# File Monitor Setup

This document describes how to set up and run the file monitor for the Maestro project.

## Overview

The file monitor watches the `./backlog/tasks` directory for changes to markdown files and outputs events in real-time.

## Agent Orchestration

Maestro supports an agent orchestration system that automatically executes scripts when tasks are assigned to specific agents.

### Prerequisites

- tmux installed (`tmux --version` to verify)
- Agent configurations in `agents/` directory (see [Agent Configuration](agent-configuration.md))
- At least one agent configured with `enabled: true` and `script_path` set

### How It Works

1. The monitor detects assignee changes in task files
2. Agent names are matched against configured agents (case-insensitive)
3. Matching agents with `enabled: true` and `script_path` configured have their scripts executed
4. Scripts run in their configured tmux session

### Quick Start

1. Create an agent configuration (see [Agent Orchestration Quickstart](agent-orchestration-quickstart.md))
2. Start the monitor: `make run`
3. Update a task file's `assignee` field to match your agent name
4. Watch tmux session for script execution

### Configuration

See [Agent Configuration](agent-configuration.md) for detailed configuration options.

### Environment Variables

- `AGENTS_CONFIG_DIR`: Override agents directory (default: `./agents`)

For full agent orchestration setup, see [Agent Orchestration Quickstart](agent-orchestration-quickstart.md).

## Prerequisites

- Go 1.20 or higher installed
- The `./backlog/tasks` directory must exist

## Installation

No additional installation is required. The monitor uses only standard library packages and the `fsnotify` library which is included in the module dependencies.

## Running the Monitor

To run the monitor:

```bash
go run cmd/monitor/main.go
```

The monitor will start and display:

```
Starting file monitor...
Watching directory: ./backlog/tasks
Monitor running. Press Ctrl+C to stop.
```

## Event Output

When file changes are detected, the monitor outputs events in the following format:

```
[2026-03-15T00:52:00.000000000Z] CREATE: backlog/tasks/got-008 - Task-1-File-Watcher-Implementation-using-Gos-fsnotify.md
[2026-03-15T00:52:01.000000000Z] WRITE: backlog/tasks/got-008 - Task-1-File-Watcher-Implementation-using-Gos-fsnotify.md
[2026-03-15T00:52:02.000000000Z] REMOVE: backlog/tasks/test-task.md
```

Event types:
- **CREATE**: A new markdown file was created
- **WRITE**: An existing markdown file was modified
- **REMOVE**: A markdown file was deleted
- **RENAME**: A markdown file was renamed

## Stopping the Monitor

Press `Ctrl+C` to gracefully stop the monitor. The application will:

1. Stop the file watcher
2. Close the event channel
3. Exit cleanly

## Debouncing

The monitor includes debouncing to handle rapid file changes. When a file is modified multiple times in quick succession (within 50ms), only the first event is processed. This prevents event flooding during file saves or edits.

## Error Handling

The monitor handles errors gracefully:
- If a file is deleted between the event notification and reading, the error is logged and monitoring continues
- File permission errors are logged but do not stop the watcher

## Manual Testing

To test the monitor:

1. Start the monitor in one terminal
2. Create a new task file in another terminal:
   ```bash
   touch backlog/tasks/test-task.md
   ```
3. Modify an existing file:
   ```bash
   echo "test" >> backlog/tasks/got-008\ -\ Task-1-File-Watcher-Implementation-using-Gos-fsnotify.md
   ```
4. Delete a file:
   ```bash
   rm backlog/tasks/test-task.md
   ```

You should see events logged in the monitor terminal.
