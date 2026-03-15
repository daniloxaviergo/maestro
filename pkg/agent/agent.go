package agent

import (
	"fmt"

	"maestro/pkg/config"
)

// Agent manages agent identity and configuration loading.
type Agent struct {
	name string
	path string
	cfg  config.AgentConfig
}

// NewAgent creates a new Agent with the given name and config path.
// If name is empty, it reads from AGENT_NAME environment variable.
// If configPath is empty, it builds the path from ConfigDirFromEnv() and AgentNameFromEnv().
func NewAgent(name, configPath string) *Agent {
	if name == "" {
		name = config.AgentNameFromEnv()
	}

	if configPath == "" {
		configDir := config.ConfigDirFromEnv()
		configPath = buildConfigPath(configDir, name)
	}

	return &Agent{
		name: name,
		path: configPath,
		cfg:  config.AgentConfig{},
	}
}

// buildConfigPath constructs the config file path from config directory and agent name.
// The path format is: {configDir}/{agentName}/config.yml
func buildConfigPath(configDir, agentName string) string {
	return fmt.Sprintf("%s/%s/config.yml", configDir, agentName)
}

// LoadConfig reads and parses the configuration file from the configured path.
// Returns the loaded configuration or a default config if the file is missing or invalid.
func (a *Agent) LoadConfig() config.AgentConfig {
	a.cfg = config.LoadConfig(a.path)
	return a.cfg
}

// GetConfig returns the currently loaded configuration.
func (a *Agent) GetConfig() config.AgentConfig {
	return a.cfg
}

// GetName returns the agent name.
func (a *Agent) GetName() string {
	return a.name
}

// GetConfigPath returns the configured config file path.
func (a *Agent) GetConfigPath() string {
	return a.path
}

// NewTestAgent creates a new Agent for testing purposes.
// This should only be used in tests as it bypasses normal config loading.
func NewTestAgent(name string, cfg config.AgentConfig) *Agent {
	return &Agent{
		name: name,
		path: "/test/path/" + name,
		cfg:  cfg,
	}
}
