package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigDir = "./agents"
)

// LoadConfig reads and parses an agent configuration file from the given path.
// Returns a default AgentConfig if the file is missing or cannot be parsed.
func LoadConfig(path string) AgentConfig {
	result := AgentConfig{
		ScriptPath:  "",
		TmuxSession: "",
		Enabled:     false,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		logWarningf("config file not found at %q: %v", path, err)
		return result
	}

	if err := yaml.Unmarshal(data, &result); err != nil {
		logWarningf("failed to parse config file %q: %v", path, err)
		return result
	}

	return result
}

// AgentNameFromEnv reads the AGENT_NAME environment variable.
// Returns an empty string if the variable is not set.
func AgentNameFromEnv() string {
	return os.Getenv("AGENT_NAME")
}

// ConfigDirFromEnv reads the AGENTS_CONFIG_DIR environment variable.
// Returns the default config directory path if the variable is not set.
func ConfigDirFromEnv() string {
	dir := os.Getenv("AGENTS_CONFIG_DIR")
	if dir == "" {
		return defaultConfigDir
	}
	return dir
}

// logWarningf logs a warning message using log.Printf
func logWarningf(format string, args ...interface{}) {
	log.Printf("Warning: "+format, args...)
}
