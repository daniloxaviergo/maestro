---
id: GOT-033
title: Add config YML to maestro
status: In Progress
assignee:
  - Thomas
created_date: '2026-03-16 17:36'
updated_date: '2026-03-16 17:39'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Configuration file for watch paths (currently hardcoded `./backlog/tasks`)
the file should be in ./
<!-- SECTION:DESCRIPTION:END -->

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

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Config file maestro.yml created at project root with watch_paths field
- [ ] #2 pkg/config package exports MaestroConfig struct and LoadMaestroConfig function
- [ ] #3 cmd/monitor/main.go loads config and uses configured watch paths
- [ ] #4 Default behavior preserved when config file is missing
- [ ] #5 go vet passes with no warnings
- [ ] #6 go build succeeds without errors
- [ ] #7 Unit tests added for config loading
- [ ] #8 make build and make run work correctly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create a YAML configuration file (`maestro.yml`) at the project root to externalize hardcoded paths and add flexible configuration loading. The approach:

- **Define configuration structure**: Create `pkg/config/types.go` extension with a `MaestroConfig` struct containing watch paths, logging options, and debounce settings
- **Add config loading**: Extend `pkg/config/config.go` with `LoadMaestroConfig()` function that reads `./maestro.yml` and returns default values if missing
- **Update watcher**: Modify `pkg/watcher/watcher.go` to accept watch paths via config instead of hardcoded `./backlog/tasks`
- **Update monitor**: Modify `cmd/monitor/main.go` to load the config and use configured paths
- **Graceful fallback**: If `maestro.yml` doesn't exist, default to current behavior (`./backlog/tasks`)

**Key design decisions**:
- YAML format for consistency with existing agent configs
- Optional file (missing = use defaults) to avoid breaking existing deployments
- Watch paths as array to support multiple directories in future

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `pkg/config/types.go` | Modify | Add `MaestroConfig` struct with fields for watch paths, debounce ms, log directory |
| `pkg/config/config.go` | Modify | Add `LoadMaestroConfig()` function with YAML parsing and defaults |
| `pkg/watcher/watcher.go` | Modify | Remove hardcoded watch path, accept paths from config |
| `cmd/monitor/main.go` | Modify | Load config, use configured watch paths |
| `maestro.yml` | Create | Example configuration file with documented defaults |
| `docs/configuration.md` | Create | Document the new config file format and options |

### 3. Dependencies

- **No new dependencies** - uses existing `gopkg.in/yaml.v3`
- **GOT-015** (Done) - `pkg/config` package provides foundation
- **GOT-016** (Done) - `pkg/agent` package shows similar config pattern
- Existing `make build` and `go build` tooling works as-is

### 4. Code Patterns

Follow existing project conventions:

- **Error handling**: Log warnings on missing config, return defaults
- **YAML tags**: Use snake_case (`watch_paths`, `debounce_ms`)
- **Config loading pattern**: Same as `LoadConfig()` in `pkg/config/config.go`
- **Naming**: `MaestroConfig` struct, `LoadMaestroConfig()` function
- **Debounce handling**: Keep existing 50ms default in `pkg/cache/cache.go`, config should override if needed

### 5. Testing Strategy

- **Unit tests**: Add `pkg/config/config_test.go` test for `LoadMaestroConfig()` with:
  - Missing config file (uses defaults)
  - Valid config file with custom paths
  - Invalid YAML (graceful degradation)
- **Integration test**: Verify `make run` works with custom `maestro.yml`
- **Edge cases**: Empty paths array, non-existent directories (handled by watcher)

### 6. Risks and Considerations

- **Breaking change**: None - missing config = current behavior
- **Backward compatibility**: Fully maintained - defaults match current hardcoded values
- **File naming**: `maestro.yml` follows existing pattern (like `backlog/config.yml`)
- **Debounce override**: Current hardcoded 50ms in `pkg/cache/cache.go` - config should either override or add config for it
- **Watch path validation**: Watcher already validates paths exist; config should validate paths are directories
<!-- SECTION:PLAN:END -->
