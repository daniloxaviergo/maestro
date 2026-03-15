---
id: GOT-020
title: 'Task 1: Agent Matching Engine - Match assignee names to configured agents'
status: In Progress
assignee: []
created_date: '2026-03-15 18:52'
updated_date: '2026-03-15 20:12'
labels:
  - task
  - agent
  - orchestration
dependencies:
  - GOT-015
  - GOT-016
references:
  - >-
    /home/danilo/scripts/github/maestro/backlog/docs/PRD-Agent-Orchestration-System.md
priority: high
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Task 1: Agent Matching Engine - Match assignee names to configured agents
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The Agent Matching Engine will match assignee names from task YAML frontmatter to configured agents in the system. The implementation follows the existing architecture patterns in the codebase:

- **New package `pkg/matcher`**: Dedicated module for agent matching logic, following the pattern of `pkg/change_detect`, `pkg/cache`, etc.
- **Matching algorithm**: Case-insensitive comparison of assignee names against agent names from their configurations
- **Multi-assignee support**: Handle multiple assignees in a task; match each against all configured agents
- **Graceful degradation**: If no matching agent found, log warning and continue; if no agent configs exist, continue monitoring

**Architecture decisions:**
- Use a simple map-based lookup (`agentName -> Agent`) for O(1) matching
- Delegate config loading to existing `pkg/agent` and `pkg/config` packages
- Return list of matching agents (can be empty) rather than failing when no match found
- No caching needed - agent configs are loaded fresh on each match request per the PRD

**Why this approach:**
- Minimal code changes - existing packages handle loading/parsing
- Clear separation: `pkg/matcher` handles matching, `pkg/agent` manages agent state
- Follows existing error handling pattern (log warnings, don't crash)
- Supports future extensibility (multiple matches, priority matching, etc.)

### 2. Files to Modify

| Action | File | Description |
|--------|------|-------------|
| Create | `pkg/matcher/matcher.go` | Main implementation with `Matcher` struct and `MatchAssignees` method |
| Create | `pkg/matcher/matcher_test.go` | Unit tests for matching logic |

**No existing files need modification** - this is a pure addition that uses existing packages.

### 3. Dependencies

**Prerequisites (already satisfied):**
- ✅ `pkg/agent` package exists with `Agent` struct and `GetName()`, `GetConfig()` methods
- ✅ `pkg/config` package exists with `LoadConfig()` and `ConfigDirFromEnv()` functions
- ✅ Agent config files stored in `{config_dir}/{agent_name}/config.yml` structure

**External dependencies:**
- No new dependencies required (uses existing `pkg/agent`, `pkg/config`)

**No blocking tasks** - GOT-015 and GOT-016 are prerequisites and already completed.

### 4. Code Patterns

**From existing packages to follow:**

1. **pkg/agent patterns:**
   - Constructor pattern: `NewMatcher(agents []*Agent) *Matcher`
   - Public methods: `MatchAssignees(assignees []string) []*Agent`
   - Error handling: Log warnings via `log.Printf("Warning: ...")`

2. **pkg/change_detect patterns:**
   - Method names in CamelCase
   - Clear return types (`[]*Agent` for matches, empty slice if none)
   - Comments for public functions

3. **Error handling patterns:**
   - Missing agent config: log warning, continue with other agents
   - No matching agents: log debug/warning message
   - Never return errors from match function - always return list of matches

**Naming conventions:**
- Struct: `Matcher`
- Constructor: `NewMatcher`
- Method: `MatchAssignees(assignees []string) []*Agent`
- Variables: `agentMap`, `matchedAgents`, `agentName`

### 5. Testing Strategy

**Test cases:**
1. `TestNewMatcher_EmptyAgents` - Creates matcher with no agents
2. `TestNewMatcher_SingleAgent` - Creates matcher with one agent
3. `TestNewMatcher_MultipleAgents` - Creates matcher with multiple agents
4. `TestMatchAssignees_NoMatches` - No agents match the assignees (logs warning)
5. `TestMatchAssignees_SingleMatch` - One assignee matches one agent
6. `TestMatchAssignees_MultipleMatches` - Multiple assignees match multiple agents
7. `TestMatchAssignees_CaseInsensitive` - Match is case-insensitive
8. `TestMatchAssignees_PartialMatch` - Some match, some don't
9. `TestMatchAssignees_EmptyInput` - Empty assignee list returns empty matches

**Verification:**
- `go test ./pkg/matcher/...` - All tests pass
- `go vet ./pkg/matcher/...` - No warnings
- Test coverage target: ≥80% for `MatchAssignees` function

### 6. Risks and Considerations

**No blocking issues.** Implementation is straightforward.

**Design considerations:**
1. **Case sensitivity**: Match will be case-insensitive (convert agent names to lowercase in map key)
2. **Performance**: O(n) where n = number of agents (acceptable for typical < 10 agents)
3. **Multiple matches**: Same assignee can match multiple agents if names duplicated (logs warning per duplicate)
4. **Agent reloading**: Current design loads agents once on matcher creation; for dynamic agent config changes, a `ReloadAgents()` method can be added later
5. **Order preservation**: Matched agents returned in order of assignees in the task (deterministic)

**Future extensibility:**
- Priority matching (if agent has priority field, sort matches)
- Weighted matching (fuzzy match with confidence scores)
- Agent grouping (match assignee to agent group, then expand group to agents)

**Integration notes:**
- To be called from `pkg/change_detect/detector.go` after assignee change detection
- Integration task (GOT-022) will wire matcher into the change detection flow
- This task (GOT-020) only implements the matching engine; integration is separate
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed successfully. Created pkg/matcher/matcher.go with NewMatcher constructor and MatchAssignees method. All tests pass, go vet passes, and make build succeeds. The implementation follows project conventions: case-insensitive matching, proper error handling with log warnings, and O(1) map-based lookup for performance.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
### Summary

Implemented the Agent Matching Engine (`pkg/matcher`) to match assignee names from task YAML frontmatter to configured agents in the system.

### Changes

- **Created `pkg/matcher/matcher.go`**: Implements `Matcher` struct with:
  - `NewMatcher(agents []*Agent)`: Constructor that builds case-insensitive lookup map
  - `MatchAssignees(assignees []string) []*Agent`: Matches assignees to agents, returns empty slice if no matches
  
- **Created `pkg/matcher/matcher_test.go`**: 11 comprehensive test cases covering:
  - Empty agents/assignees handling
  - Single/multiple agent scenarios
  - Case-insensitive matching
  - Partial matches
  - Duplicate assignees
  - Order preservation

### Verification

- ✅ `go test ./...` - All tests pass (11/11 in matcher package)
- ✅ `go vet ./...` - No warnings
- ✅ `make build` - Build succeeds
- ✅ Follows project conventions (package structure, naming, error handling)

### Notes

- Case-insensitive matching implemented via lowercase map keys
- Logs warnings when assignees don't match any configured agent
- Returns empty slice when no matches (graceful degradation)
- Preserves order of matched agents to match input assignee order
- No caching (reloads agents fresh per task per PRD specification)
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Code follows existing project conventions package structure naming error handling
- [x] #2 go vet passes with no warnings
- [x] #3 go build succeeds without errors
- [x] #4 Unit tests added or updated for new or changed functionality
- [x] #5 go test ... passes with no failures
- [x] #6 Code comments added for non-obvious logic
- [ ] #7 README or docs updated if public behavior changes
- [x] #8 make build succeeds
- [ ] #9 make run works as expected
- [x] #10 Errors are logged not silently ignored
- [x] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->
