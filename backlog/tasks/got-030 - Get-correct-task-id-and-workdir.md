---
id: GOT-030
title: Get correct task id and workdir
status: To Do
assignee: []
created_date: '2026-03-16 14:31'
updated_date: '2026-03-16 17:28'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
To agent catarina script.sh receive this `/home/danilo/scripts/github/maestro/backlog/tasks/got-028` TASK_FILE

in script extract the task and the project path
in this sample task_id is got-028
and work dir is /home/danilo/scripts/github/maestro

some sample os TASK_FILE
/home/danilo/scripts/github/maestro/backlog/tasks/got-028
/home/danilo/scripts/github/maestro/backlog/tasks/got-001
/home/danilo/scripts/github/maestro/backlog/tasks/got-101
/home/danilo/scripts/github/maestro/backlog/tasks/got-1
/home/danilo/scripts/github/dca/backlog/tasks/got-1
/home/danilo/scripts/github/go-todotask/backlog/tasks/got-015
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
- [ ] #1 Task ID correctly extracted from task file paths with various formats
- [ ] #2 Project path correctly identified by detecting backlog/tasks structure
- [ ] #3 All existing agent scripts updated to use new extraction method
- [ ] #4 Unit tests added for task path extraction package
- [ ] #5 go vet passes with no warnings
- [ ] #6 go build succeeds without errors
- [ ] #7 make build succeeds
- [ ] #8 make run works with updated scripts
- [ ] #9 Documentation added for new extraction package
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

**Problem Statement:**
Agent scripts receive the task file path (e.g., `/home/danilo/scripts/github/maestro/backlog/tasks/got-028`) and must extract:
- `task_id`: The task identifier (e.g., `got-028` - the filename)
- `workdir`: The project root directory (e.g., `/home/danilo/scripts/github/maestro`)

**Current Implementation Issues:**
- Script uses directory traversal (walking up from `backlog/tasks`)
- Relies on detecting `backlog/tasks` directory + `backlog/config.yml` to find project root
- May fail if `config.yml` doesn't exist or if the directory structure varies

**Proposed Solution:**
Create a Go helper library that extracts task ID and project path reliably, then update all agent scripts to use this library.

**Implementation approach:**
1. Create `pkg/taskpath/` package with functions to extract task ID and project path
2. Use path parsing (basename for task ID, parent directory traversal for project path)
3. Add validation to ensure we found a valid project root (check for `backlog/` structure)
4. Update all agent scripts to use the Go helper via a shared script or function

**Why this approach:**
- Go provides robust path manipulation (more reliable than bash `dirname`/`basename` chains)
- Centralized logic = easier to fix if project structure changes
- Can add validation and error handling
- Works across different shell environments

### 2. Files to Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/taskpath/types.go` | **Create** | Type definitions for task path parsing results |
| `pkg/taskpath/taskpath.go` | **Create** | Core library for extracting task ID and project path |
| `agents/catarina/script.sh` | **Modify** | Update to use Go helper for task path extraction |
| `agents/agent-foo/script.sh` | **Modify** | Update to use Go helper for task path extraction |
| `agents/agent-bar/script.sh` | **Modify** | Update to use Go helper for task path extraction |
| `docs/task-path-extraction.md` | **Create** | Documentation for the new package |
| `cmd/extract-path/` | **Create** | (Optional) CLI utility for testing the extraction logic |

### 3. Dependencies

**Prerequisites:**
- All features from tasks GOT-008 through GOT-029 must be implemented (all marked as Done)
- Go 1.25.7+ available (verified in project context)
- `./backlog/tasks` directory structure must exist

**No blocking tasks** - This is a refactoring/enhancement task.

**Build verification:**
- `make build` must succeed after adding new package
- `go vet ./...` must pass with no warnings
- `go test ./taskpath/...` must pass

### 4. Code Patterns

**Conventions to Follow (from project QWEN.md):**

1. **Package naming**: Lowercase, short names (e.g., `taskpath`)
2. **Function naming**: CamelCase (e.g., `ParseTaskPath`, `ExtractTaskID`)
3. **Variables**: camelCase (e.g., `taskFilePath`, `projectPath`)
4. **Error handling**: Explicit `if err != nil` checks with `return err`
5. **Comments**: Add comments for non-obvious logic (path traversal strategy)

**Error handling pattern:**
```go
func ParseTaskPath(filePath string) (TaskPathResult, error) {
    if filePath == "" {
        return TaskPathResult{}, ErrEmptyPath
    }
    // ...
}
```

**Directory traversal strategy:**
- Start from `backlog/tasks/` parent directory
- Walk up directory tree
- Check for `backlog/config.yml` to identify project root
- Handle edge cases: absolute vs relative paths, symlinks

### 5. Testing Strategy

**Unit Tests:**
Test in `pkg/taskpath/taskpath_test.go`:

| Test Case | Description |
|-----------|-------------|
| `TestParseTaskID` | Extract task ID from various paths |
| `TestExtractProjectPath` | Find project root from various paths |
| `TestInvalidPath` | Handle empty or malformed paths |
| `TestNestedProjects` | Handle deeply nested project paths |
| `TestSymlinkHandling` | Verify symlink resolution |

**Test examples:**
```go
// Test case 1: Standard path
Input:  "/home/danilo/scripts/github/maestro/backlog/tasks/got-028"
Output: TaskID="got-028", ProjectPath="/home/danilo/scripts/github/maestro"

// Test case 2: Nested project
Input:  "/home/danilo/scripts/github/dca/backlog/tasks/got-001"
Output: TaskID="got-001", ProjectPath="/home/danilo/scripts/github/dca"

// Test case 3: Non-standard structure (should still work)
Input:  "/home/danilo/scripts/github/go-todotask/backlog/tasks/got-015"
Output: TaskID="got-015", ProjectPath="/home/danilo/scripts/github/go-todotask"
```

**Integration Test:**
- Create sample task files and test script execution
- Verify `TASK_ID` and `PROJECT_PATH` are correctly extracted
- Test with agent scripts to ensure end-to-end flow works

**Verification:**
- `go test ./pkg/taskpath/...` passes
- Sample script runs and outputs correct values
- `make build` succeeds
- `go vet ./...` passes

### 6. Risks and Considerations

**Known risks:**
1. **Backward compatibility**: Existing scripts may rely on current extraction logic
   - Mitigation: Provide migration path, keep fallback logic if needed
   
2. **Directory structure variations**: Projects may have different directory layouts
   - Mitigation: Make the extraction logic configurable via environment variable or config file
   
3. **Symlink issues**: Project roots may be symlinks
   - Mitigation: Use `filepath.EvalSymlinks()` to resolve symlinks before processing
   
4. **Performance**: Directory traversal could be slow for very deep trees
   - Mitigation: Limit traversal depth (e.g., max 20 levels up) - sufficient for any real project
   
5. **Edge cases**: Windows path separators, malformed paths
   - Mitigation: Use `filepath` package for cross-platform compatibility, validate inputs

**Design decisions:**
1. **Why Go library instead of pure bash?** 
   - More reliable path manipulation
   - Better error handling
   - Easier to test
   
2. **Why check for `backlog/config.yml`?**
   - This is the project root indicator used by the Maestro project
   - Distinguishes Maestro projects from other `backlog/` directories
   
3. **Why not use a fixed path length?**
   - Projects may have varying directory depths
   - Environment variable override allows customization

**Rollout considerations:**
- Update all agent scripts to use the new extraction method
- Document the change in each agent's `script.sh` file
- Test with existing projects before deploying to production
- No breaking changes - old behavior can be kept as fallback
<!-- SECTION:PLAN:END -->
