package config

// AgentConfig represents the configuration for an agent
type AgentConfig struct {
	ScriptPath  string `yaml:"script_path"`
	TmuxSession string `yaml:"tmux_session"`
	Enabled     bool   `yaml:"enabled"`
}

// MaestroConfig represents the main configuration for the Maestro system
type MaestroConfig struct {
	WatchPaths  []string `yaml:"watch_paths"`
	DebounceMs  int      `yaml:"debounce_ms"`
	LogDir      string   `yaml:"log_dir"`
}

// DefaultMaestroConfig returns the default configuration values
func DefaultMaestroConfig() MaestroConfig {
	return MaestroConfig{
		WatchPaths:  []string{"./backlog/tasks"},
		DebounceMs:  50,
		LogDir:      ".",
	}
}
