# Maestro Project Context

## Project Overview

**Maestro** is a Go-based file monitoring and agent orchestration system that tracks changes to markdown task files in a Backlog.md workflow. The project monitors the `./backlog/tasks` directory for file changes, detects assignee field modifications, logs changes to JSON, and can trigger tmux-based notifications or script execution.

### Technology Stack

- **Language**: Go 1.25.7
- **File Watching**: `fsnotify` v1.9.0 for real-time filesystem event detection
- **YAML Parsing**: `gopkg.in/yaml.v3` v3.0.1 for frontmatter extraction
- **Task Management**: Backlog.md (Markdown-based task tracking system with YAML frontmatter)
- **Build Tools**: Standard Go toolchain (`go build`, `go test`, `go vet`, `go mod`)
- **Notification**: Tmux for display-message and script execution

### Project Structure

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
│   ├── config/          # Configuration loading (agents, agents)
│   │   └── config.go
│   ├── logs/            # JSON logging for assignee changes
│   │   └── logger.go
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
├── agents/              # Agent configuration directory (default)
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
├── Makefile             # Build and run commands
└── .env                 # Environment variables (API keys)
```

### Architecture

The system uses a layered, composable architecture:

1. **cmd/monitor** - CLI entry point with signal handling (SIGINT/SIGTERM)
2. **pkg/watcher** - fsnotify wrapper that:
   - Recursively watches `./backlog/tasks` directory
   - Filters events to only `.md` files
   - Converts fsnotify events to normalized `FileEvent` types
   - Implements 50ms debouncing to handle rapid successive writes
3. **pkg/cache** - File state management:
   - Caches file content hashes and metadata
   - Tracks assignee field values across file events
   - Thread-safe with mutex protection
4. **pkg/parser** - YAML frontmatter extraction:
   - Extracts YAML frontmatter from markdown files
   - Parses `assignee`, `status`, `id`, `title` fields
5. **pkg/change_detect** - Assignee change detection:
   - Compares current assignee with cached assignee
   - Logs changes to JSON with timestamp
   - Triggers tmux notifications via pkg/notifier
6. **pkg/notifier** - Notification system:
   - Sends tmux `display-message` notifications
   - Executes bash scripts in tmux sessions for complex actions
   - Non-blocking (async) execution via goroutines

### File Event Types

The watcher detects four event types:
- **CREATE**: New markdown file created
- **WRITE**: Existing file modified
- **REMOVE**: File deleted
- **RENAME**: File renamed (rarely triggered due to fsnotify behavior)

Events are output in format: `[timestamp] TYPE: /absolute/path/to/file.md`

### Assignee Change Detection Flow

1. File watcher detects CREATE or WRITE event
2. Parser extracts YAML frontmatter from markdown file
3. Cache retrieves previous assignee value (or uses empty if first run)
4. Detector compares new vs. cached assignee (order-insensitive)
5. If changed:
   - Log entry written to `assignee_changes.log` (JSON format)
   - Tmux notification sent (if configured)
   - Cache updated with new assignee
6. If agent is configured, associated script may execute

### Configuration

#### Agent Configuration

Agents are defined in YAML config files (default path: `./agents/{agent_name}/config.yml`):

```yaml
script_path: "/path/to/script.sh"
tmux_session: "maestro-agent"
enabled: true
```

## Building and Running

### Prerequisites

- Go 1.20+ installed
- `./backlog/tasks` directory must exist
- Write permissions to project directory
- Tmux installed (for notifications)

### Building

```bash
# Build the monitor binary
make build

# Or directly with go
go build -o bin/monitor cmd/monitor/main.go

# Build all packages
go build ./...

# Run static analysis
go vet ./...
```

### Running

```bash
# Run the monitor (watch for assignee changes)
make run

# Or directly
go run cmd/monitor/main.go
```

The monitor will:
1. Start watching `./backlog/tasks` recursively
2. Output file events to stdout in real-time
3. Log assignee changes to `assignee_changes.log` (JSON)
4. Send tmux notifications for assignee changes (if configured)
5. Continue until SIGINT (Ctrl+C) or SIGTERM

### Tmux Notifications

To use tmux notifications, start a tmux session first:

```bash
# Start tmux session for notifications
make tmux-start

# Then run the monitor
make run

# Attach to see notifications
make tmux-attach

# Stop tmux session
make tmux-stop
```

### Testing

```bash
# Run tests
go test ./...

# Static analysis
go vet ./...

# Dependency management
go mod tidy

# Full build with analysis
make build && go vet ./...
```

### Manual Testing

```bash
# Terminal 1: Start monitor
make run

# Terminal 2: Trigger events
touch backlog/tasks/test-task.md          # CREATE
echo "content" >> backlog/tasks/test.md   # WRITE
rm backlog/tasks/test-task.md             # REMOVE

# Terminal 3: Check log file
cat assignee_changes.log
```

## Development Conventions

### Code Style

- **Package names**: Lowercase, short (e.g., `cache`, `watcher`, `parser`)
- **Function names**: CamelCase (e.g., `NewWatcher`, `ProcessFile`, `DetectChanges`)
- **Variables**: camelCase (e.g., `fileWatcher`, `eventQueue`, `logger`)
- **Error variables**: Prefix with `Err` (e.g., `ErrWatcherStopped`, `ErrScriptNotFound`)

### Package Structure

- `cmd/` - CLI entry points only; initialize and run components
- `pkg/` - Library code with separation of concerns:
  - `watcher/` - Filesystem watching logic (fsnotify wrapper)
  - `cache/` - State management and debouncing
  - `parser/` - YAML frontmatter extraction
  - `change_detect/` - Assignee change detection logic
  - `notifier/` - Tmux notification and script execution
  - `config/` - Configuration loading
  - `agent/` - Agent identity and config management
  - `logs/` - JSON logging

### Key Patterns

- **Explicit error handling**: `if err != nil { return err }`
- **Buffered channels** to prevent goroutine leaks (capacity 100 for events)
- **Mutex protection** for concurrent access (`sync.RWMutex`)
- **Context for cancellation** (available for future extension)
- **Non-blocking notifications** via goroutines

### Debouncing

File write events are debounced with a 50ms cooldown per file to prevent flooding from rapid successive writes (e.g., during file saves).

### Error Handling

- Errors are logged but do not crash the watcher
- File not found errors on removal are handled gracefully
- Permission errors are logged and monitoring continues
- Failed notifications are logged with exit codes

## Backlog.md Task Management

This project uses Backlog.md for ALL TASK MANAGEMENT VIA MCP tools:

- **Task file format**: Markdown with YAML frontmatter
- **Task IDs**: `GOT-XXX` (e.g., GOT-008) or prefixed with labels (e.g., AGENT-001)
- **Config**: `backlog/config.yml` defines project settings
- **Statuses**: To Do, In Progress, Done
- **Milestones**: Managed via `backlog/milestones/` directory

### Task Workflow

1. Tasks are created in `backlog/tasks/` with frontmatter
2. Tasks progress through statuses (To Do → In Progress → Done)
3. Completed tasks may be moved to `backlog/completed/`
4. All work should be tracked in Backlog.md
5. Task files are read-only by the monitor (only assignee field changes are detected)

## Current Implementation Status

### Completed Tasks

- **GOT-008** (Done): File watcher implementation using fsnotify
  - Recursive directory watching
  - Event filtering for `.md` files
  - Debouncing mechanism (50ms cooldown)
  - Signal handling for graceful shutdown

- **GOT-009** (Done): YAML frontmatter parser
  - Extracts YAML frontmatter from markdown
  - Parses `assignee`, `status`, `id`, `title` fields
  - Handles missing frontmatter gracefully

- **GOT-010** (Done): Change detection and JSON logging
  - Assignee change detection (order-insensitive comparison)
  - JSON logging to `assignee_changes.log`
  - Thread-safe logging with mutex protection

- **GOT-011** to **GOT-014** (Done): Tmux notifier
  - Notification system for assignee changes
  - tmux `display-message` integration
  - Script execution via tmux sessions

- **GOT-015** to **GOT-018** (Done): Agent configuration
  - Config package for loading agent configurations
  - Agent package for managing agent identity
  - Integration with detector for script execution

- **GOT-019** to **GOT-025** (Done): Agent orchestration
  - Agent matching and script execution routing
  - Monitor integration with agent orchestration
  - Example configurations and documentation

### Related Upcoming Tasks (in backlog)

- **GOT-026+**: Testing and documentation improvements
- **GOT-027+**: Performance optimizations and metrics
- **GOT-028+**: Systemd service files for production deployment

## Future Extensions (Out of Scope for Current Implementation)

- Configuration file for watch paths (currently hardcoded `./backlog/tasks`)
- File output vs stdout via CLI flags
- Additional notification backends (Slack, email, etc.)
- Performance metrics/health checks (prometheus metrics)
- Systemd service files for production deployment
- Web UI for monitoring and logs
- Database-backed state persistence
- Cross-platform event coalescing improvements
- File watching limits tuning (`fs.inotify.max_user_watches`)

## Important Notes

1. **Task files are read-only**: The monitor only reads task files; it never modifies them
2. **Assignee field tracking**: Changes to the `assignee` YAML field trigger notifications
3. **Debouncing**: Rapid successive writes are coalesced to prevent duplicate notifications
4. **Log format**: JSON with RFC3339 timestamps for structured logging
5. **Tmux non-blocking**: Notifications execute asynchronously to avoid blocking the monitor
