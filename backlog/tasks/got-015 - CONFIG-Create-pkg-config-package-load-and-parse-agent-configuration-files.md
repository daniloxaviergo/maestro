---
id: GOT-015
title: '[CONFIG] Create pkg/config package - load and parse agent configuration files'
status: To Do
assignee: []
created_date: '2026-03-15 17:16'
updated_date: '2026-03-15 17:36'
labels: []
dependencies: []
references:
  - backlog/docs/doc-004-per-agent-configuration.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create configuration loading and parsing package for agent YAML files
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 pkg/config/config.go with LoadConfig function
- [x] #2 pkg/config/types.go with AgentConfig struct (script_path, tmux_session, enabled)
- [x] #3 LoadConfig reads YAML file from path and returns AgentConfig
- [x] #4 Missing config file logs warning and returns default config
- [x] #5 YAML parsing errors are caught and logged
- [x] #6 AgentNameFromEnv() function to read AGENT_NAME environment variable
- [x] #7 ConfigDirFromEnv() function to read AGENTS_CONFIG_DIR environment variable
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The `pkg/config` package will implement YAML-based configuration loading for agent configuration files. The design follows the existing project patterns seen in `pkg/parser` and `pkg/cache`:

- **Package structure**: Two files (`types.go` for types, `config.go` for functions) mirroring `pkg/parser` conventions
- **YAML parsing**: Use `gopkg.in/yaml.v3` (already in `go.mod` dependencies)
- **Default values**: Return a default `AgentConfig` when config file is missing or invalid
- **Environment variables**: Read optional paths via `os.Getenv` with sensible defaults
- **Error handling**: Log warnings but never fail; return default config on any error
- **Thread safety**: No mutex needed - config loading is a read-only operation that returns new structs

**Architecture decisions:**
- Use `string` fields for YAML mapping (script_path → ScriptPath) matching Go snake_case conventions
- Keep struct simple with only the 3 required fields from PRD
- No config validation beyond type checking - config is agent-local responsibility

### 2. Files to Modify

| Action | File | Description |
|--------|------|-------------|
| Create | `pkg/config/types.go` | Define `AgentConfig` struct with YAML tags |
| Create | `pkg/config/config.go` | Implement `LoadConfig`, `AgentNameFromEnv`, `ConfigDirFromEnv` |
| Create | `pkg/config/config_test.go` | Unit tests for all public functions |
| Create | `pkg/config/fixtures/` | Test fixture YAML files |
| Modify | `go.mod` | (likely no change - yaml.v3 already present) |

### 3. Dependencies

- **Existing**: `gopkg.in/yaml.v3 v3.0.1` (already in `go.mod`)
- **Environment variables**:
  - `AGENT_NAME`: Optional agent identifier (default: empty string)
  - `AGENTS_CONFIG_DIR`: Optional config directory (default: `./agents`)
- **No blocking tasks**: Can be implemented independently

### 4. Code Patterns

Follow existing patterns from `pkg/parser` and `pkg/cache`:

```go
// Type definitions in types.go
type AgentConfig struct {
    ScriptPath  string `yaml:"script_path"`
    TmuxSession string `yaml:"tmux_session"`
    Enabled     bool   `yaml:"enabled"`
}

// Function returns in config.go
func LoadConfig(path string) AgentConfig {
    // Return default on any error
}

// Error handling pattern (from parser.go)
result := AgentConfig{ /* default values */ }
if err != nil {
    log.Printf("Warning: %v", err)
    return result
}
```

**Naming conventions:**
- Public functions: CamelCase (e.g., `LoadConfig`, `AgentNameFromEnv`)
- Struct fields: CamelCase with YAML tags
- Error messages: Start with capital, no trailing period in log

### 5. Testing Strategy

**Test cases:**
1. `TestLoadConfig_ValidYAML` - Parse valid config file
2. `TestLoadConfig_MissingFile` - Return default config with warning
3. `TestLoadConfig_InvalidYAML` - Parse errors logged, default returned
4. `TestLoadConfig_PartialConfig` - Missing fields use defaults (enabled=false, empty strings)
5. `TestAgentNameFromEnv_Present` - Read from environment
6. `TestAgentNameFromEnv_Empty` - Return empty when not set
7. `TestConfigDirFromEnv_Present` - Read from environment
8. `TestConfigDirFromEnv_Empty` - Return default when not set

**Fixture files:**
- `valid-config.yml` - Complete valid config
- `partial-config.yml` - Missing fields
- `invalid-yaml.yml` - Malformed YAML

### 6. Risks and Considerations

- **No blocking issues**: All acceptance criteria are straightforward
- **YAML library**: Using existing dependency `yaml.v3` (same as parser package)
- **Default values**: `Enabled` defaults to `false` (safe disabled state)
- **Path handling**: Use `os.Getenv` directly, no path resolution in config package (caller's responsibility)
- **Future extensibility**: Struct is simple but can be extended with new fields without breaking existing code
- **Test fixtures**: Need to create test data directory structure
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed successfully. Created pkg/config package with config loading and parsing for agent YAML files.
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 Code follows existing project conventions package structure naming error handling
- [ ] #2 go vet passes with no warnings
- [ ] #3 go build succeeds without errors
- [ ] #4 Unit tests added or updated for new or changed functionality
- [ ] #5 go test ... passes with no failures
- [ ] #6 Code comments added for non-obvious logic
- [ ] #7 README or docs updated if public behavior changes
- [ ] #8 make build succeeds
- [ ] #9 make run works as expected
- [ ] #10 Errors are logged not silently ignored
- [ ] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->
