# File Monitor Setup

This document describes how to set up and run the file monitor for the Maestro project.

## Overview

The file monitor watches the `./backlog/tasks` directory for changes to markdown files and outputs events in real-time.

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
