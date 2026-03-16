---
id: GOT-009
title: 'Task 2: YAML Frontmatter Parser for extracting the assignee field'
status: Done
assignee: []
created_date: '2026-03-15 00:52'
updated_date: '2026-03-16 17:29'
labels:
  - parser
  - yaml
  - go
dependencies:
  - GOT-008
references:
  - >-
    backlog/docs/doc-002 -
    PRD-Monitor-File-Changes-in-.-backlog-tasks-When-assignee-Field-Is-Modified.md
priority: high
ordinal: 16000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a YAML frontmatter parser in Go to extract the `assignee` field from markdown task files in `./backlog/tasks`.

### Implementation Notes

**Goal**: Parse YAML frontmatter from markdown files and extract the assignee array.

**Key Components**:
1. Create `pkg/parser/` package with frontmatter extraction logic
2. Parse content between `---` delimiters at start of file
3. Use `gopkg.in/yaml.v3` or `github.com/go-yaml/yaml` for YAML parsing
4. Extract `assignee` field as `[]string` (handle missing/empty cases)
5. Return structured data: `{filePath, oldAssignee, newAssignee, err}`

**YAML Frontmatter Format**:
```yaml
---
id: task-001
title: Example Task
assignee: ["username"]
---
# Content
```

**Error Handling**:
- Handle files without frontmatter (treat as empty assignee array)
- Handle malformed YAML with descriptive error logging
- Handle missing `assignee` field (treat as empty array)
- Handle invalid file paths or read errors

**Data Structures**:
```go
type FileData struct {
    FilePath   string
    Assignee   []string
    Error      error
    ParseTime  time.Duration
}
```

**Integration Points**:
- Receives file paths from watcher module via channels
- Outputs parsed data to cache comparison module
- Must be performant (<100ms per file as per NFR)

**Dependencies**:
- `gopkg.in/yaml.v3` or `github.com/go-yaml/yaml`

**Acceptance Criteria**:
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Successfully parse YAML frontmatter from valid markdown files
- [x] #2 Extract assignee field as a slice of strings (array)
- [x] #3 Handle files without frontmatter (treat as empty assignee array)
- [x] #4 Handle malformed YAML gracefully with error logging
- [x] #5 Support empty assignee arrays (assignee: [] or assignee:)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create a YAML frontmatter parser package (`pkg/parser/`) that extracts the `assignee` field from markdown files. The parser will:

- Read file content and detect YAML frontmatter delimited by `---`
- Use `gopkg.in/yaml.v3` (already available in go.mod) to parse YAML
- Extract `assignee` as `[]string` with graceful fallbacks for missing/malformed data
- Return structured `FileData` with path, assignee slice, error, and parse duration

**Design decisions:**
- Separate parser from cache/watcher to maintain single responsibility
- No frontmatter = empty assignee array (not an error)
- Malformed YAML returns error but doesn't crash the system
- Thread-safe parser instance (no shared state)

### 2. Files to Modify

| Action | File | Purpose |
|--------|------|---------|
| Create | `pkg/parser/parser.go` | Main parser logic: `ParseFile(path string) FileData` |
| Create | `pkg/parser/types.go` | Data structures: `FileData`, `Frontmatter` |
| Create | `pkg/parser/parser_test.go` | Unit tests for all edge cases |
| Create | `pkg/parser/fixtures/` | Test fixture files with various frontmatter formats |

### 3. Dependencies

- **Existing**: `gopkg.in/yaml.v3` (already in `go.mod`)
- **Existing**: `os`, `time`, `io/ioutil` from Go standard library
- **Dependency on other tasks**: GOT-008 (File Watcher) - parser receives file paths from watcher

### 4. Code Patterns

**Follow existing project conventions:**
- Package structure mirrors `pkg/watcher/` and `pkg/cache/`
- Error handling: return error in struct, don't panic
- Naming: camelCase for fields, PascalCase for exported types
- File naming: `parser.go` for main logic, `types.go` for structures
- Logging: use `fmt.Fprintf(os.Stderr, ...)` for errors (matching `watcher.go`)

**Data structures to define:**
```go
type Frontmatter struct {
    ID        string   `yaml:"id"`
    Title     string   `yaml:"title"`
    Assignee  []string `yaml:"assignee"`
    Status    string   `yaml:"status"`
    // other common fields as needed
}

type FileData struct {
    FilePath   string
    Assignee   []string
    Error      error
    ParseTime  time.Duration
}
```

### 5. Testing Strategy

**Test cases to cover:**
1. Valid frontmatter with assignee array: `assignee: ["alice", "bob"]`
2. Valid frontmatter with empty assignee: `assignee: []` or `assignee:`
3. Missing frontmatter (content starts with `#` or text)
4. Malformed YAML (syntax errors)
5. Missing `assignee` field (file has frontmatter but no assignee key)
6. Parse time validation (<100ms for typical files)

**Test approach:**
- Use subtests in `parser_test.go`
- Create fixture files in `fixtures/` directory
- Benchmark test for performance validation
- Error case tests verify error messages are descriptive

### 6. Risks and Considerations

**Blocking/dependencies:**
- The parser is useless without GOT-008 (file watcher) producing file paths
- GOT-010 (change detection) depends on this parser for input data

**Trade-offs:**
- Simple string search for `---` delimiters vs using a markdown parser (choosing simple search for performance and minimal dependencies)
- No caching of parsed results (stateless parser, caching done by cache package)
- YAMLLint-style validation not in scope (only parsing, not validation)

**Performance considerations:**
- Must meet <100ms NFR for typical task files (~100-500 lines)
- Avoid unnecessary allocations in the parse loop
- Read file once, parse once

**Future extensibility:**
- Parser returns `Frontmatter` struct that can grow with more fields
- `FileData` includes `Error` field for detailed error reporting
- Parser can be extended to return all frontmatter keys if needed
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Successfully implemented YAML frontmatter parser in Go. All 5 acceptance criteria verified through unit tests.

Created pkg/parser/ package with types.go (FileData, Frontmatter structs) and parser.go (ParseFile function using gopkg.in/yaml.v3).

Created comprehensive test suite with 9 test cases covering all edge cases: valid frontmatter, empty assignee, missing assignee, no frontmatter, malformed YAML, non-existent file, and performance validation.

All tests pass and application builds successfully. Parser meets <100ms NFR requirement.

Created test fixtures in pkg/parser/fixtures/ directory for reproducible test cases.

Integration point: Parser is ready to receive file paths from GOT-008 (file watcher) and output parsed data to cache comparison module.
<!-- SECTION:NOTES:END -->
