package agent

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	fixturesDir      = "./fixtures"
	testAgentName    = "test-agent"
	testConfigDir    = "./test-config"
	testConfigPath   = "./test-config/test-agent/config.yml"
)

// setupTestEnv sets up test environment variables and returns cleanup function
func setupTestEnv(agentName, configDir string) func() {
	// Store original values
	origAgentName := os.Getenv("AGENT_NAME")
	origConfigDir := os.Getenv("AGENTS_CONFIG_DIR")

	// Set test values
	if agentName != "" {
		os.Setenv("AGENT_NAME", agentName)
	} else {
		os.Unsetenv("AGENT_NAME")
	}

	if configDir != "" {
		os.Setenv("AGENTS_CONFIG_DIR", configDir)
	} else {
		os.Unsetenv("AGENTS_CONFIG_DIR")
	}

	return func() {
		// Restore original values
		os.Setenv("AGENT_NAME", origAgentName)
		os.Setenv("AGENTS_CONFIG_DIR", origConfigDir)
	}
}

func TestNewAgent_NameFromEnv(t *testing.T) {
	cleanup := setupTestEnv("env-agent", "")
	defer cleanup()

	agent := NewAgent("", "")
	if agent.GetName() != "env-agent" {
		t.Errorf("Expected agent name from env var 'env-agent', got %q", agent.GetName())
	}
}

func TestNewAgent_ConfigDirFromEnv(t *testing.T) {
	cleanup := setupTestEnv("test-agent", "/custom/config")
	defer cleanup()

	agent := NewAgent("", "")
	expectedPath := "/custom/config/test-agent/config.yml"
	if agent.GetConfigPath() != expectedPath {
		t.Errorf("Expected config path %q, got %q", expectedPath, agent.GetConfigPath())
	}
}

func TestNewAgent_ExplicitNameAndPath(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	agent := NewAgent("explicit-agent", "/explicit/path/config.yml")
	if agent.GetName() != "explicit-agent" {
		t.Errorf("Expected agent name 'explicit-agent', got %q", agent.GetName())
	}
	if agent.GetConfigPath() != "/explicit/path/config.yml" {
		t.Errorf("Expected config path '/explicit/path/config.yml', got %q", agent.GetConfigPath())
	}
}

func TestNewAgent_Defaults(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	agent := NewAgent("", "")
	// Agent name should be empty when env var not set
	if agent.GetName() != "" {
		t.Errorf("Expected empty agent name when not set, got %q", agent.GetName())
	}
	// Config path should use default config dir
	expectedDefaultPath := "./agents//config.yml"
	if agent.GetConfigPath() != expectedDefaultPath {
		t.Errorf("Expected default config path %q, got %q", expectedDefaultPath, agent.GetConfigPath())
	}
}

func TestAgent_LoadConfig_ValidFile(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	filePath := filepath.Join(fixturesDir, "valid-config.yml")
	agent := NewAgent("test", filePath)

	cfg := agent.LoadConfig()

	if cfg.ScriptPath != "/home/user/scripts/agent.sh" {
		t.Errorf("Expected ScriptPath '/home/user/scripts/agent.sh', got %q", cfg.ScriptPath)
	}
	if cfg.TmuxSession != "agent-workspace" {
		t.Errorf("Expected TmuxSession 'agent-workspace', got %q", cfg.TmuxSession)
	}
	if !cfg.Enabled {
		t.Errorf("Expected Enabled true, got false")
	}
}

func TestAgent_LoadConfig_MissingFile(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	filePath := filepath.Join(fixturesDir, "non-existent.yml")
	agent := NewAgent("test", filePath)

	cfg := agent.LoadConfig()

	// Should return default config (all zero values)
	if cfg.ScriptPath != "" {
		t.Errorf("Expected empty ScriptPath for missing file, got %q", cfg.ScriptPath)
	}
	if cfg.TmuxSession != "" {
		t.Errorf("Expected empty TmuxSession for missing file, got %q", cfg.TmuxSession)
	}
	if cfg.Enabled {
		t.Errorf("Expected Enabled false for missing file, got true")
	}
}

func TestAgent_LoadConfig_InvalidYAML(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	filePath := filepath.Join(fixturesDir, "invalid-yaml.yml")
	agent := NewAgent("test", filePath)

	cfg := agent.LoadConfig()

	// Should return default config on YAML parse error
	if cfg.ScriptPath != "" {
		t.Errorf("Expected empty ScriptPath for invalid YAML, got %q", cfg.ScriptPath)
	}
	if cfg.TmuxSession != "" {
		t.Errorf("Expected empty TmuxSession for invalid YAML, got %q", cfg.TmuxSession)
	}
	if cfg.Enabled {
		t.Errorf("Expected Enabled false for invalid YAML, got true")
	}
}

func TestAgent_GetConfig(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	filePath := filepath.Join(fixturesDir, "valid-config.yml")
	agent := NewAgent("test", filePath)

	// Load config first
	agent.LoadConfig()

	// Get config and verify
	cfg := agent.GetConfig()

	if cfg.ScriptPath != "/home/user/scripts/agent.sh" {
		t.Errorf("Expected ScriptPath '/home/user/scripts/agent.sh', got %q", cfg.ScriptPath)
	}
}

func TestAgent_GetName(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	agent := NewAgent("my-agent", "/path/config.yml")
	if agent.GetName() != "my-agent" {
		t.Errorf("Expected agent name 'my-agent', got %q", agent.GetName())
	}
}

func TestAgent_GetConfigPath(t *testing.T) {
	cleanup := setupTestEnv("", "")
	defer cleanup()

	agent := NewAgent("test", "/custom/path/config.yml")
	if agent.GetConfigPath() != "/custom/path/config.yml" {
		t.Errorf("Expected config path '/custom/path/config.yml', got %q", agent.GetConfigPath())
	}
}

func TestAgent_GetConfigPath_BuildsCorrectly(t *testing.T) {
	cleanup := setupTestEnv("my-agent", "/config/dir")
	defer cleanup()

	agent := NewAgent("", "")
	expectedPath := "/config/dir/my-agent/config.yml"
	if agent.GetConfigPath() != expectedPath {
		t.Errorf("Expected config path %q, got %q", expectedPath, agent.GetConfigPath())
	}
}
