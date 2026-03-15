.PHONY: build run clean

# Build the monitor binary
build:
	go build -o bin/monitor cmd/monitor/main.go

# Run the monitor directly
run:
	go run cmd/monitor/main.go

# Clean build artifacts
clean:
	rm -rf bin/monitor
