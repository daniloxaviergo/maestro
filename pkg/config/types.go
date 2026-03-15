package config

// AgentConfig represents the configuration for an agent
type AgentConfig struct {
	ScriptPath  string `yaml:"script_path"`
	TmuxSession string `yaml:"tmux_session"`
	Enabled     bool   `yaml:"enabled"`
}
