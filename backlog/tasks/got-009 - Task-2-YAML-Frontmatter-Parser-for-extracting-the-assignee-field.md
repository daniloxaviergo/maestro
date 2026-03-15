---
id: GOT-009
title: 'Task 2: YAML Frontmatter Parser for extracting the assignee field'
status: To Do
assignee: []
created_date: '2026-03-15 00:52'
updated_date: '2026-03-15 00:53'
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
ordinal: 4000
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
- [ ] #1 Successfully parse YAML frontmatter from valid markdown files
- [ ] #2 Extract assignee field as a slice of strings (array)
- [ ] #3 Handle files without frontmatter (treat as empty assignee array)
- [ ] #4 Handle malformed YAML gracefully with error logging
- [ ] #5 Support empty assignee arrays (assignee: [] or assignee:)
<!-- AC:END -->
