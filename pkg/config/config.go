package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigDir = "./agents"
)

// DefaultConfigPath is the default path to the maestro config file
var DefaultConfigPath = "./maestro.yml"

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

// LoadMaestroConfig reads and parses the main Maestro configuration file.
// Returns a default MaestroConfig if the file is missing or cannot be parsed.
func LoadMaestroConfig() MaestroConfig {
	result := DefaultMaestroConfig()

	data, err := os.ReadFile(DefaultConfigPath)
	if err != nil {
		logWarningf("maestro config file not found at %q: %v, using defaults", DefaultConfigPath, err)
		return result
	}

	if err := yaml.Unmarshal(data, &result); err != nil {
		logWarningf("failed to parse maestro config file %q: %v, using defaults", DefaultConfigPath, err)
		return DefaultMaestroConfig()
	}

	// Validate watch_paths - use default if empty
	if len(result.WatchPaths) == 0 {
		logWarningf("watch_paths is empty in config, using default")
		result.WatchPaths = DefaultMaestroConfig().WatchPaths
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
