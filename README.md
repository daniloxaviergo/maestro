# Maestro

Go-based file monitoring and agent orchestration system for Backlog.md task management.

## Overview

Maestro watches the `./backlog/tasks` directory for markdown file changes, detects assignee field modifications, logs changes to JSON, and can trigger tmux-based notifications or script execution.

## Features

- **Real-time file monitoring** using fsnotify
- **Assignee change detection** with order-insensitive comparison
- **JSON logging** with RFC3339 timestamps
- **Tmux notifications** via `display-message`
- **Agent script execution** in tmux sessions
- **Debouncing** to coalesce rapid successive writes (50ms cooldown)
- **Graceful shutdown** on SIGINT/SIGTERM

## Technology Stack

- **Language**: Go 1.25.7
- **File Watching**: fsnotify v1.9.0
- **YAML Parsing**: gopkg.in/yaml.v3 v3.0.1
- **Task Management**: Backlog.md (Markdown with YAML frontmatter)
- **Notifications**: Tmux

## Project Structure

```
maestro/
├── cmd/                 # Command-line applications
│   └── monitor/         # Main file watcher CLI
│       └── main.go
├── pkg/                 # Library code
│   ├── agent/           # Agent identity and configuration management
│   │   └── agent.go
│   ├── cache/           # File state caching with debouncing
│   │   ├── types.go
│   │   └── cache.go
│   ├── change_detect/   # Assignee change detection and logging
│   │   └── detector.go
│   ├── config/          # Configuration loading
│   │   └── config.go
│   ├── logs/            # JSON logging for assignee changes
│   │   └── logger.go
│   ├── matcher/         # Agent-assignee matching logic
│   │   └── matcher.go
│   ├── notifier/        # Tmux notification and script execution
│   │   ├── notifier.go
│   │   └── types.go
│   ├── parser/          # YAML frontmatter extraction
│   │   ├── parser.go
│   │   └── types.go
│   └── watcher/         # fsnotify wrapper for file monitoring
│       ├── watcher.go
│       └── events.go
├── bin/                 # Compiled binaries
├── docs/                # Project documentation
├── agents/              # Agent configuration directory
├── backlog/             # Backlog.md task files
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
├── Makefile             # Build and run commands
└── README.md            # This file
```

## Installation

### Prerequisites

- Go 1.20+
- Tmux (for notifications)

### Build

```bash
make build
```

Or directly with Go:

```bash
go build -o bin/monitor cmd/monitor/main.go
```

## Usage

### Basic Run

```bash
make run
```

The monitor will:
1. Watch `./backlog/tasks` recursively
2. Output file events to stdout
3. Log assignee changes to `assignee_changes.log`
4. Send tmux notifications (if configured)
5. Execute agent scripts (if configured)

### With Tmux Notifications

```bash
# Start tmux session
make tmux-start

# Run the monitor
make run

# Attach to see notifications
make tmux-attach

# Stop session
make tmux-stop
```

### Testing

```bash
# Run tests
go test ./...

# Static analysis
go vet ./...

# Full build with analysis
make build && go vet ./...
```

### Manual Testing

Terminal 1:
```bash
make run
```

Terminal 2:
```bash
# Create a task
echo '---
assignee: [qwen-code]
---
Task content' > backlog/tasks/test-task.md

# Modify assignee to trigger detection
echo '---
assignee: [some-other-user]
---
Task content' > backlog/tasks/test-task.md

# Check the log
cat assignee_changes.log
```

## Configuration

### Agent Configuration

Create agent directories in `./agents/{agent_name}/config.yml`:

```yaml
script_path: "/path/to/script.sh"
tmux_session: "agent-session"
enabled: true
```

**Fields:**
- `script_path`: Path to bash script to execute when agent is assigned
- `tmux_session`: Session name for script execution (default: "default")
- `enabled`: Whether script execution is active (default: false)

### Example Agent Setup

```bash
# Create agent directory
mkdir -p agents/qwen-code

# Create config
cat > agents/qwen-code/config.yml <<EOF
script_path: "./agents/qwen-code/script.sh"
tmux_session: "maestro-qwen"
enabled: true
EOF

# Create script
cat > agents/qwen-code/script.sh <<'EOF'
#!/bin/bash
echo "Task assigned to $AGENT_NAME at $(date)"
EOF

chmod +x agents/qwen-code/script.sh
```

## File Event Types

| Type | Description |
|------|-------------|
| CREATE | New markdown file created |
| WRITE | Existing file modified |
| REMOVE | File deleted |
| RENAME | File renamed (rarely triggered) |

Output format: `[timestamp] TYPE: /path/to/file.md`

## Assignee Change Detection Flow

1. File watcher detects CREATE or WRITE event
2. Parser extracts YAML frontmatter
3. Cache retrieves previous assignee value
4. Detector compares new vs cached assignee (order-insensitive)
5. If changed:
   - Log entry written to `assignee_changes.log`
   - Tmux notification sent (if configured)
   - Agent matcher runs to find matching agents
   - Agent scripts execute (if configured)
   - Cache updated with new assignee
6. If unchanged:
   - Cache simply updated with current assignee

## Development

### Code Style

- Package names: lowercase, short (e.g., `cache`, `watcher`, `parser`)
- Function names: CamelCase (e.g., `NewWatcher`, `ProcessFile`)
- Variables: camelCase (e.g., `fileWatcher`, `eventQueue`)
- Error variables: `Err` prefix (e.g., `ErrWatcherStopped`)

### Key Patterns

- Explicit error handling: `if err != nil { return err }`
- Buffered channels (capacity 100) to prevent goroutine leaks
- Mutex protection for concurrent access (`sync.RWMutex`)
- Non-blocking notifications via goroutines
- Debouncing (50ms cooldown) for rapid successive writes

### Build Commands

```bash
make build      # Build binary
make clean      # Remove binary
make run        # Run monitor
go vet ./...    # Static analysis
go mod tidy     # Update dependencies
```

## Backlog.md Integration

This project uses Backlog.md for task management:

- Task files: Markdown with YAML frontmatter
- Task IDs: `GOT-XXX` or prefixed (e.g., `AGENT-001`)
- Statuses: To Do, In Progress, Done
- Task directory: `backlog/tasks/`

**Note:** Task files are read-only by the monitor; only assignee field changes are detected.

## Environment Setup

```bash
# Create necessary directories
mkdir -p backlog/tasks agents

# Create a sample agent config
mkdir -p agents/qwen-code
cat > agents/qwen-code/config.yml <<EOF
script_path: "/path/to/script.sh"
tmux_session: "maestro-agent"
enabled: true
EOF
```

## Documentation

- [QWEN.md](./QWEN.md) - Detailed project context and architecture
- [AGENTS.md](./AGENTS.md) - Backlog.md MCP workflow guidelines

## License

MIT
