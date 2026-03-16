---
id: GOT-030
title: Get correct task id and workdir
status: To Do
assignee:
  - Catarina
created_date: '2026-03-16 14:31'
updated_date: '2026-03-16 14:42'
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

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The goal is to create a helper function that parses a task file path and extracts two key pieces of information:
1. **Task ID** (e.g., `got-028`) - the base name of the task file without the markdown extension
2. **Project Workdir** (e.g., `/home/danilo/scripts/github/maestro`) - the root directory of the project containing the backlog

**Approach:**
- Create a new utility package `pkg/taskpath` to encapsulate task path parsing logic
- The function will take a full file path and return task_id and workdir
- Use path manipulation (filepath.SplitList, filepath.Dir, filepath.Base) to extract components
- The workdir detection will look for the `backlog/tasks` pattern to identify the project root
- The task ID is simply the filename (without `.md` extension)

**Why this approach:**
- Clean separation of concerns - task path utilities belong in their own package
- Reusable by any agent script that needs to understand its task context
- Simple path manipulation is reliable and doesn't require regex
- Backward compatible - existing scripts continue to work as before

**Trade-offs:**
- Workdir detection relies on the `backlog/tasks` directory structure
- If a project uses a different structure, the heuristic may need adjustment
- No validation that the extracted task_id actually corresponds to an existing task file

### 2. Files to Modify

| File | Action | Purpose |
|------|--------|---------|
| `pkg/taskpath/types.go` | **Create** | Type definitions for function signatures |
| `pkg/taskpath/taskpath.go` | **Create** | Main implementation of path parsing logic |
| `pkg/taskpath/taskpath_test.go` | **Create** | Unit tests for the parser |
| `pkg/notifier/notifier.go` | **Modify** | Import and use taskpath package for path parsing |
| `backlog/tasks/got-030 - Get-correct-task-id-and-workdir.md` | **Modify** | Add implementation plan to task |

**Files to reference:**
- `cmd/monitor/main.go` - For understanding how file paths are currently used
- `pkg/parser/parser.go` - Similar file parsing patterns

### 3. Dependencies

**Prerequisites:**
- GOT-030 task created and approved (this task)
- All existing tasks (GOT-008 through GOT-029) must be complete
- No external dependencies required - uses only standard library (`path/filepath`, `strings`)

**Setup steps:**
- Create `pkg/taskpath/` directory
- Run `go mod tidy` to ensure dependencies are up to date
- Run `go build ./...` to verify compilation

### 4. Code Patterns

**Conventions to follow:**
1. **Package naming**: lowercase, short (`taskpath`, not `taskpathutil` or `taskpathparser`)
2. **Function naming**: CamelCase exported, camelCase unexported (`ParseTaskPath`, `detectWorkdir`)
3. **Error handling**: Return error as last return value, handle gracefully in calling code
4. **Documentation**: Exported functions need comment lines
5. **Testing**: Unit tests follow Go convention (`_test.go` suffix)

**Example patterns from codebase:**
```go
// Parser pattern from pkg/parser
func ParseFile(filePath string) FileData {
    // ... implementation
    return result
}

// Error handling pattern from pkg/notifier
if err != nil {
    log.Printf("Warning: %v", err)
    return
}
```

### 5. Testing Strategy

**Unit tests will cover:**
1. **Normal cases**: Standard task paths like `/home/danilo/scripts/github/maestro/backlog/tasks/got-028.md`
2. **Different projects**: Paths in different project directories
3. **Task ID extraction**: Verify correct extraction of `got-028` from filename
4. **Workdir extraction**: Verify correct project root detection
5. **Edge cases**:
   - Paths without `.md` extension
   - Paths with special characters
   - Very deep nested paths

**Testing approach:**
- Test files in `pkg/taskpath/taskpath_test.go`
- Use table-driven tests for multiple input cases
- Verify both task_id and workdir are correctly extracted
- No mocking needed - pure path manipulation

**Verification:**
- `go test ./pkg/taskpath/...` should pass
- `go vet ./pkg/taskpath/...` should have no warnings
- Manual test with sample paths to verify output

### 6. Risks and Considerations

**Known risks:**
1. **Workdir detection heuristic**: Assumes `backlog/tasks` structure exists; may fail for alternative project layouts
2. **No validation**: Task ID extraction doesn't verify the task file exists
3. **Path separators**: May behave differently on Windows vs Linux (though backlog uses forward slashes)

**Mitigation strategies:**
- Use `filepath` package for cross-platform compatibility
- Document the expected path structure clearly
- Add error handling for unexpected path formats
- Consider adding a configuration option for workdir override

**Implementation considerations:**
- Workdir should be the directory containing `backlog/`, not `backlog/` itself
- Task ID should be the filename without `.md` extension
- Handle both absolute and relative paths correctly
- Return empty strings on error or malformed paths

### 7. Implementation Steps

**Phase 1: Create taskpath package**
1. Create `pkg/taskpath/types.go` with type definitions
2. Create `pkg/taskpath/taskpath.go` with `ParseTaskPath(filePath string) (taskID, workdir string, err error)`
3. Implement path parsing logic:
   - Extract base filename using `filepath.Base()`
   - Remove `.md` extension to get task ID
   - Walk up directory tree looking for `backlog/tasks` pattern
   - Return parent of `backlog/` as workdir

**Phase 2: Write unit tests**
1. Create `pkg/taskpath/taskpath_test.go`
2. Add test cases for:
   - Standard maestro path: `/home/danilo/scripts/github/maestro/backlog/tasks/got-028.md`
   - Different project: `/home/danilo/scripts/github/dca/backlog/tasks/got-1.md`
   - Task ID edge cases: `got-1`, `got-015`, etc.
3. Run tests with `go test ./pkg/taskpath/...`

**Phase 3: Update existing code (optional integration)**
1. Review `pkg/notifier/notifier.go` to see if taskpath can improve script execution
2. Consider adding task ID and workdir as environment variables or script arguments
3. Update monitor to pass parsed values if needed

**Phase 4: Verification**
1. Run `go build ./...` to verify no compilation errors
2. Run `go vet ./...` to check for issues
3. Run `make build` to verify Makefile still works
4. Test manually with sample script that uses the new package
<!-- SECTION:PLAN:END -->

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
