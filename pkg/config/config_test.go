package config

import (
	"os"
	"path/filepath"
	"testing"
)

const fixturesDir = "./fixtures"

func TestLoadConfig_ValidYAML(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "valid-config.yml")
	result := LoadConfig(filePath)

	// Check no error and correct values
	if result.ScriptPath != "/home/user/scripts/agent.sh" {
		t.Errorf("Expected ScriptPath to be '/home/user/scripts/agent.sh', got %q", result.ScriptPath)
	}
	if result.TmuxSession != "agent-workspace" {
		t.Errorf("Expected TmuxSession to be 'agent-workspace', got %q", result.TmuxSession)
	}
	if !result.Enabled {
		t.Errorf("Expected Enabled to be true, got false")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "non-existent.yml")
	result := LoadConfig(filePath)

	// Should return default config (all zero values)
	if result.ScriptPath != "" {
		t.Errorf("Expected empty ScriptPath for missing file, got %q", result.ScriptPath)
	}
	if result.TmuxSession != "" {
		t.Errorf("Expected empty TmuxSession for missing file, got %q", result.TmuxSession)
	}
	if result.Enabled {
		t.Errorf("Expected Enabled to be false for missing file, got true")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "invalid-yaml.yml")
	result := LoadConfig(filePath)

	// Should return default config on YAML parse error
	if result.ScriptPath != "" {
		t.Errorf("Expected empty ScriptPath for invalid YAML, got %q", result.ScriptPath)
	}
	if result.TmuxSession != "" {
		t.Errorf("Expected empty TmuxSession for invalid YAML, got %q", result.TmuxSession)
	}
	if result.Enabled {
		t.Errorf("Expected Enabled to be false for invalid YAML, got true")
	}
}

func TestLoadConfig_PartialConfig(t *testing.T) {
	filePath := filepath.Join(fixturesDir, "partial-config.yml")
	result := LoadConfig(filePath)

	// Should return config with defaults for missing fields
	if result.ScriptPath != "/home/user/scripts/agent.sh" {
		t.Errorf("Expected ScriptPath to be '/home/user/scripts/agent.sh', got %q", result.ScriptPath)
	}
	if result.TmuxSession != "" {
		t.Errorf("Expected empty TmuxSession for partial config, got %q", result.TmuxSession)
	}
	if result.Enabled {
		t.Errorf("Expected Enabled to be false (default) for partial config, got true")
	}
}

func TestAgentNameFromEnv_Present(t *testing.T) {
	// Set environment variable
	os.Setenv("AGENT_NAME", "test-agent")
	defer os.Unsetenv("AGENT_NAME")

	result := AgentNameFromEnv()
	if result != "test-agent" {
		t.Errorf("Expected 'test-agent', got %q", result)
	}
}

func TestAgentNameFromEnv_Empty(t *testing.T) {
	// Ensure variable is not set
	os.Unsetenv("AGENT_NAME")

	result := AgentNameFromEnv()
	if result != "" {
		t.Errorf("Expected empty string when not set, got %q", result)
	}
}

func TestConfigDirFromEnv_Present(t *testing.T) {
	// Set environment variable
	os.Setenv("AGENTS_CONFIG_DIR", "/custom/config/dir")
	defer os.Unsetenv("AGENTS_CONFIG_DIR")

	result := ConfigDirFromEnv()
	if result != "/custom/config/dir" {
		t.Errorf("Expected '/custom/config/dir', got %q", result)
	}
}

func TestConfigDirFromEnv_Default(t *testing.T) {
	// Ensure variable is not set
	os.Unsetenv("AGENTS_CONFIG_DIR")

	result := ConfigDirFromEnv()
	if result != defaultConfigDir {
		t.Errorf("Expected default config dir %q, got %q", defaultConfigDir, result)
	}
}
