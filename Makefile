.PHONY: build run clean tmux

# Build the monitor binary
build:
	go build -o bin/monitor cmd/monitor/main.go

# Run the monitor directly
run:
	go run cmd/monitor/main.go

# Clean build artifacts
clean:
	rm -rf bin/monitor

# --- Tmux commands for testing notifications ---

# Start a detached tmux session for notifications
tmux-start:
	tmux new-session -d -s maestro-test

# List tmux sessions
tmux-list:
	tmux ls

# Attach to the tmux session
tmux-attach:
	tmux attach -t maestro-test

# Kill the tmux session
tmux-stop:
	tmux kill-session -t maestro-test

# Run monitor with tmux session pre-started (all-in-one)
run-tmux:
	$(MAKE) tmux-start
	go run cmd/monitor/main.go

# Test tmux notifications manually
tmux-test:
	@echo "=== Testing tmux notification ==="
	$(MAKE) tmux-start
	@echo "Session started. To see notification:"
	@echo "  1. Run: make run (in another terminal)"
	@echo "  2. Update a task file's assignee"
	@echo "  3. Attach with: make tmux-attach"
	@echo ""
	@echo "To stop session: make tmux-stop"
	$(MAKE) tmux-list
