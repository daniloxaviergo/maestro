---
id: GOT-016
title: >-
  [AGENT] Create pkg/agent package - manage agent identity and configuration
  loading
status: In Progress
assignee: []
created_date: '2026-03-15 17:16'
updated_date: '2026-03-15 18:01'
labels: []
dependencies: []
references:
  - backlog/docs/doc-004-per-agent-configuration.md
priority: high
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create agent package to manage agent identity and configuration loading
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 pkg/agent/agent.go with Agent struct to manage agent identity and configuration
- [x] #2 Agent.LoadConfig() method to load config from configured path
- [x] #3 Agent.GetConfig() method to return loaded configuration
- [x] #4 Agent.GetName() method to return agent name
- [x] #5 Default config directory is ./agents/ configurable via AGENTS_CONFIG_DIR
- [x] #6 Agent name from AGENT_NAME environment variable
- [x] #7 Missing config file logs warning but doesn't crash agent
- [x] #8 Config file path is {config_dir}/{agent_name}/config.yml
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create `pkg/agent` package with an `Agent` struct that manages agent identity and configuration loading. The implementation follows existing project patterns:

- **Agent struct**: Holds agent name (from `AGENT_NAME` env var) and configuration state
- **Delegation to pkg/config**: Use existing `ConfigDirFromEnv()` and `AgentNameFromEnv()` functions
- **Path composition**: Build config path as `{config_dir}/{agent_name}/config.yml`
- **Graceful degradation**: Missing or invalid config returns default config with warning log
- **Constructor pattern**: `NewAgent()` creates agent with default or configured name

**Why this approach:**
- Follows existing `pkg/config` patterns for consistency
- Separates concerns: `pkg/config` handles parsing, `pkg/agent` manages agent state
- Simple, testable API with clear methods

### 2. Files to Modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `pkg/agent/agent.go` | Main implementation with Agent struct and methods |
| Create | `pkg/agent/agent_test.go` | Unit tests for all public methods |
| Create | `pkg/agent/fixtures/valid-config.yml` | Test fixture for valid config |
| Create | `pkg/agent/fixtures/invalid-yaml.yml` | Test fixture for invalid YAML |

**No existing files need modification** - this is a pure addition.

### 3. Dependencies

**Prerequisites (already satisfied):**
- ✅ `pkg/config` package exists with `ConfigDirFromEnv()` and `AgentNameFromEnv()` functions
- ✅ `pkg/config/types.go` defines `AgentConfig` struct with `ScriptPath`, `TmuxSession`, `Enabled` fields
- ✅ YAML parser (`gopkg.in/yaml.v3`) available in `go.mod`
- ✅ Existing project structure and conventions understood

**No new external dependencies required.**

### 4. Code Patterns

**From existing packages to follow:**

1. **pkg/config patterns:**
   - Environment variable handling with defaults
   - Warning logs via `log.Printf("Warning: ...")`
   - Return zero/default values on errors (no crash)

2. **pkg/cache patterns:**
   - Thread-safe operations using `sync.RWMutex`
   - Commented public functions
   - Clear error handling with `fmt.Errorf` wrapping

3. **pkg/notifier patterns:**
   - NewConstructor pattern (`NewNotifier()`, `NewAgent()`)
   - Configuration structs for options

**Naming conventions:**
- Struct: `Agent`
- Constructor: `NewAgent()`
- Methods: `LoadConfig()`, `GetConfig()`, `GetName()`, `GetConfigPath()`
- Variables: camelCase

### 5. Testing Strategy

**Test cases:**
- `TestNewAgent_NameFromEnv` - Agent name from `AGENT_NAME` env var
- `TestNewAgent_ConfigDirFromEnv` - Config dir from `AGENTS_CONFIG_DIR` env var
- `TestNewAgent_Defaults` - Default values when env vars not set
- `TestAgent_LoadConfig_ValidFile` - Load valid YAML config
- `TestAgent_LoadConfig_MissingFile` - Missing config logs warning, returns default
- `TestAgent_LoadConfig_InvalidYAML` - Invalid YAML logs warning, returns default
- `TestAgent_GetConfig` - Returns loaded config
- `TestAgent_GetName` - Returns agent name
- `TestAgent_GetConfigPath` - Builds correct config path

### 6. Risks and Considerations

**No blocking issues.** This is a straightforward addition.

**Design considerations:**
1. **Config reloading**: Current design loads config once. If config needs dynamic reloading, a `ReloadConfig()` method can be added later.
2. **Error handling**: Warning logs are sufficient per existing pattern in `pkg/config`
3. **Thread safety**: Agent will be read-only after LoadConfig; no mutex needed unless config changes at runtime (not required by acceptance criteria)
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Code follows existing project conventions package structure naming error handling
- [x] #2 go vet passes with no warnings
- [x] #3 go build succeeds without errors
- [x] #4 Unit tests added or updated for new or changed functionality
- [x] #5 go test ... passes with no failures
- [ ] #6 Code comments added for non-obvious logic
- [ ] #7 README or docs updated if public behavior changes
- [x] #8 make build succeeds
- [ ] #9 make run works as expected
- [ ] #10 Errors are logged not silently ignored
- [ ] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->
