# Maestro Project Context

## Project Overview

**Maestro** is a Go-based project that implements a file monitoring system for tracking changes to markdown files in a backlog management system. The project uses `fsnotify` to provide real-time file system event detection.

### Technology Stack

- **Language**: Go 1.25.7
- **File Watching**: `fsnotify` v1.9.0
- **Task Management**: Backlog.md MCP (Markdown-based task tracking system)
- **Build Tools**: Standard Go toolchain (`go build`, `go test`, `go vet`, `go mod`)

### Project Structure

```
maestro/
├── cmd/                 # Command-line applications
│   └── monitor/         # File watcher CLI
│       └── main.go
├── pkg/                 # Library code
│   ├── cache/           # File state caching with debouncing
│   │   ├── types.go
│   │   └── cache.go
│   └── watcher/         # fsnotify wrapper
│       ├── watcher.go
│       └── events.go
├── docs/                # Project documentation
│   └── setup-monitor.md
├── bin/                 # Compiled binaries
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
└── .env                 # Environment variables (API keys)
```

### Architecture

The file watcher uses a layered architecture:

1. **cmd/monitor** - CLI entry point with signal handling (SIGINT/SIGTERM)
2. **pkg/watcher** - fsnotify wrapper that:
   - Recursively watches `./backlog/tasks` directory
   - Filters events to only `.md` files
   - Converts fsnotify events to normalized `FileEvent` types
3. **pkg/cache** - File state management:
   - Caches file content hashes and metadata
   - Implements 50ms debouncing for rapid writes
   - Thread-safe with mutex protection

### Event Types

The watcher detects four event types:
- **CREATE**: New markdown file created
- **WRITE**: Existing file modified
- **REMOVE**: File deleted
- **RENAME**: File renamed

Events are output in format: `[timestamp] TYPE: /absolute/path/to/file.md`

## Building and Running

### Prerequisites

- Go 1.20+ installed
- `./backlog/tasks` directory must exist
- Write permissions to project directory

### Building

```bash
# Build the monitor binary
go build -o bin/monitor cmd/monitor/main.go

# Build all packages
go build ./...
```

### Running

```bash
# Run directly
go run cmd/monitor/main.go

# Or run the built binary
./bin/monitor
```

The monitor will:
1. Start watching `./backlog/tasks` recursively
2. Output events to stdout in real-time
3. Continue until SIGINT (Ctrl+C) or SIGTERM

### Testing

```bash
# Run tests
go test ./...

# Static analysis
go vet ./...

# Dependency management
go mod tidy
```

### Manual Testing

```bash
# Terminal 1: Start monitor
go run cmd/monitor/main.go

# Terminal 2: Trigger events
touch backlog/tasks/test-task.md          # CREATE
echo "content" >> backlog/tasks/test.md   # WRITE
rm backlog/tasks/test-task.md             # REMOVE
```

## Development Conventions

### Code Style

- **Package names**: Lowercase, short (e.g., `cache`, `watcher`)
- **Function names**: CamelCase (e.g., `NewWatcher`, `ProcessEvent`)
- **Variables**: camelCase (e.g., `fileWatcher`, `eventQueue`)
- **Error variables**: Prefix with `Err` (e.g., `ErrWatcherStopped`)

### Package Structure

- `cmd/` - CLI entry points only; initialize and run
- `pkg/` - Library code with separation of concerns:
  - `watcher/` - Filesystem watching logic
  - `cache/` - State management and debouncing

### Key Patterns

- Explicit error handling: `if err != nil { return err }`
- Buffered channels to prevent goroutine leaks
- Mutex protection for concurrent access
- Context for cancellation (optional, available for future extension)

### Debouncing

File write events are debounced with a 50ms cooldown per file to prevent flooding from rapid successive writes (e.g., during file saves).

### Error Handling

- Errors are logged but do not crash the watcher
- File not found errors on removal are handled gracefully
- Permission errors are logged and monitoring continues

## Backlog.md Task Management

This project uses Backlog.md for all task management:

- **Task file format**: Markdown with YAML frontmatter
- **Task IDs**: `GOT-XXX` (e.g., GOT-008)
- **Config**: `backlog/config.yml` defines project settings
- **Statuses**: To Do, In Progress, Done

### Task Workflow

1. Tasks are created in `backlog/tasks/` with frontmatter
2. Tasks progress through statuses (To Do → In Progress → Done)
3. Completed tasks may be moved to `backlog/completed/`
4. All work should be tracked in Backlog.md

## Current Implementation Status

### Completed Tasks

- **GOT-008** (Done): File watcher implementation using fsnotify
  - Recursive directory watching
  - Event filtering for `.md` files
  - Debouncing mechanism
  - Signal handling for graceful shutdown

### Related Upcoming Tasks (in backlog)

- **GOT-009**: YAML frontmatter parser
- **GOT-010**: Change detection and logging with JSON output

## Environment Variables

`.env` file contains API keys (do not commit secrets):
- `OLLAMA_API_KEY` - Ollama API key
- `OPENAI_API_KEY` - OpenAI API key

## Future Extensions (Out of Scope for Current Implementation)

- Configuration file for watch paths
- File output vs stdout via flags
- Integration with parser module (GOT-009)
- Performance metrics/health checks
- Systemd service files
