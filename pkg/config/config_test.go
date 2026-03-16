package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMaestroConfig_MissingFile(t *testing.T) {
	// Ensure no config file exists
	os.Remove(DefaultConfigPath)

	cfg := LoadMaestroConfig()

	// Verify defaults
	if len(cfg.WatchPaths) != 1 {
		t.Errorf("Expected 1 watch path, got %d", len(cfg.WatchPaths))
	}
	if cfg.WatchPaths[0] != "./backlog/tasks" {
		t.Errorf("Expected default watch path './backlog/tasks', got '%s'", cfg.WatchPaths[0])
	}
	if cfg.DebounceMs != 50 {
		t.Errorf("Expected default debounce_ms 50, got %d", cfg.DebounceMs)
	}
	if cfg.LogDir != "." {
		t.Errorf("Expected default log_dir '.', got '%s'", cfg.LogDir)
	}
}

func TestLoadMaestroConfig_WithCustomConfig(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "maestro-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a custom config file
	configPath := filepath.Join(tmpDir, "maestro.yml")
	configData := `
watch_paths:
  - "/custom/path1"
  - "/custom/path2"
debounce_ms: 100
log_dir: "/var/log/maestro"
`
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Temporarily change DefaultConfigPath for testing
	originalPath := DefaultConfigPath
	DefaultConfigPath = configPath
	defer func() { DefaultConfigPath = originalPath }()

	cfg := LoadMaestroConfig()

	// Verify custom values
	if len(cfg.WatchPaths) != 2 {
		t.Errorf("Expected 2 watch paths, got %d", len(cfg.WatchPaths))
	}
	if cfg.WatchPaths[0] != "/custom/path1" {
		t.Errorf("Expected watch path '/custom/path1', got '%s'", cfg.WatchPaths[0])
	}
	if cfg.WatchPaths[1] != "/custom/path2" {
		t.Errorf("Expected watch path '/custom/path2', got '%s'", cfg.WatchPaths[1])
	}
	if cfg.DebounceMs != 100 {
		t.Errorf("Expected debounce_ms 100, got %d", cfg.DebounceMs)
	}
	if cfg.LogDir != "/var/log/maestro" {
		t.Errorf("Expected log_dir '/var/log/maestro', got '%s'", cfg.LogDir)
	}
}

func TestLoadMaestroConfig_EmptyWatchPaths(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "maestro-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a config file with empty watch_paths
	configPath := filepath.Join(tmpDir, "maestro.yml")
	configData := `
watch_paths: []
debounce_ms: 50
log_dir: "."
`
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Temporarily change DefaultConfigPath for testing
	originalPath := DefaultConfigPath
	DefaultConfigPath = configPath
	defer func() { DefaultConfigPath = originalPath }()

	cfg := LoadMaestroConfig()

	// Verify fallback to defaults when watch_paths is empty
	if len(cfg.WatchPaths) != 1 {
		t.Errorf("Expected 1 watch path after fallback, got %d", len(cfg.WatchPaths))
	}
	if cfg.WatchPaths[0] != "./backlog/tasks" {
		t.Errorf("Expected fallback watch path './backlog/tasks', got '%s'", cfg.WatchPaths[0])
	}
}

func TestDefaultMaestroConfig(t *testing.T) {
	cfg := DefaultMaestroConfig()

	if len(cfg.WatchPaths) != 1 {
		t.Errorf("Expected 1 watch path, got %d", len(cfg.WatchPaths))
	}
	if cfg.WatchPaths[0] != "./backlog/tasks" {
		t.Errorf("Expected default watch path './backlog/tasks', got '%s'", cfg.WatchPaths[0])
	}
	if cfg.DebounceMs != 50 {
		t.Errorf("Expected default debounce_ms 50, got %d", cfg.DebounceMs)
	}
	if cfg.LogDir != "." {
		t.Errorf("Expected default log_dir '.', got '%s'", cfg.LogDir)
	}
}

func TestLoadMaestroConfig_InvalidYAML(t *testing.T) {
	// Create a temporary directory for test
	tmpDir, err := os.MkdirTemp("", "maestro-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an invalid YAML file
	configPath := filepath.Join(tmpDir, "maestro.yml")
	configData := `
watch_paths:
  - "./valid/path"
debounce_ms: invalid_number
`
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Temporarily change DefaultConfigPath for testing
	originalPath := DefaultConfigPath
	DefaultConfigPath = configPath
	defer func() { DefaultConfigPath = originalPath }()

	cfg := LoadMaestroConfig()

	// Verify fallback to defaults when YAML is invalid
	if len(cfg.WatchPaths) != 1 {
		t.Errorf("Expected 1 watch path after fallback, got %d", len(cfg.WatchPaths))
	}
	if cfg.WatchPaths[0] != "./backlog/tasks" {
		t.Errorf("Expected fallback watch path './backlog/tasks', got '%s'", cfg.WatchPaths[0])
	}
}
