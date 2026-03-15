package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"maestro/pkg/agent"
	"maestro/pkg/cache"
	"maestro/pkg/change_detect"
	"maestro/pkg/config"
	"maestro/pkg/logs"
	"maestro/pkg/matcher"
	"maestro/pkg/notifier"
	"maestro/pkg/parser"
	"maestro/pkg/watcher"
)

const (
	watchPath = "./backlog/tasks"
)

func main() {
	log.Println("Starting file monitor...")

	// Create log directory and logger
	logDir := "."
	logPath := filepath.Join(logDir, "assignee_changes.log")
	logger, err := logs.NewLogger(logPath)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Create change detector
	detector := change_detect.NewDetector(logger)

	// Load agent configurations from agents directory
	agentsDir := config.ConfigDirFromEnv()
	agentConfigs, err := loadAgentConfigs(agentsDir)
	if err != nil {
		log.Printf("Warning: failed to load agent configs: %v\n", err)
	}

	// Create agents and matcher
	var agents []*agent.Agent
	for _, cfgPath := range agentConfigs {
		agentName := extractAgentNameFromPath(cfgPath)
		a := agent.NewAgent(agentName, cfgPath)
		a.LoadConfig()
		agents = append(agents, a)
	}

	matcher := matcher.NewMatcher(agents)
	detector.SetMatcher(matcher)

	// Create and wire notifier for tmux notifications
	notifier := notifier.NewNotifier(notifier.NotificationConfig{})
	detector.SetNotifier(notifier)

	// Create watcher
	w, err := watcher.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}

	// Add watch path
	if err := w.AddWatch(watchPath); err != nil {
		log.Fatalf("Failed to add watch path %s: %v", watchPath, err)
	}

	log.Printf("Watching directory: %s", watchPath)

	// Start watching in background
	if err := w.Watch(); err != nil {
		log.Fatalf("Failed to start watching: %v", err)
	}

	// Set up signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Handle events
	go func() {
		for event := range w.Events() {
			fmt.Printf("[%s] %s: %s\n",
				event.Timestamp.Format(time.RFC3339Nano),
				event.Type,
				event.Path)

			// Handle remove events - clean up cache
			if event.Type == cache.EventRemove {
				detector.RemoveFile(event.Path)
				continue
			}

			// Skip non-WRITE and non-CREATE events for assignee change detection
			if event.Type != cache.EventWrite && event.Type != cache.EventCreate {
				continue
			}

			// Parse file to extract assignee
			fileData := parser.ParseFile(event.Path)
			if fileData.Error != nil {
				log.Printf("Warning: failed to parse %s: %v\n", event.Path, fileData.Error)
				continue
			}

			// Detect assignee change
			changed, err := detector.ProcessFile(fileData)
			if err != nil {
				log.Printf("Warning: failed to process %s: %v\n", event.Path, err)
				continue
			}

			if changed {
				log.Printf("Assignee change logged for %s\n", event.Path)
			}
		}
	}()

	// Wait for shutdown signal
	log.Println("Monitor running. Press Ctrl+C to stop.")
	select {
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down...", sig)
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	// Stop watcher
	w.Stop()
	log.Println("Monitor stopped")
}

// loadAgentConfigs loads all agent config file paths from the agents directory.
// Returns a list of config file paths or an error if the directory cannot be read.
func loadAgentConfigs(agentsDir string) ([]string, error) {
	entries, err := os.ReadDir(agentsDir)
	if err != nil {
		return nil, err
	}

	var configPaths []string
	for _, entry := range entries {
		if entry.IsDir() {
			cfgPath := filepath.Join(agentsDir, entry.Name(), "config.yml")
			configPaths = append(configPaths, cfgPath)
		}
	}

	return configPaths, nil
}

// extractAgentNameFromPath extracts the agent name from a config file path.
// The agent name is the parent directory name of the config file.
func extractAgentNameFromPath(configPath string) string {
	dir := filepath.Dir(configPath)
	return filepath.Base(dir)
}
