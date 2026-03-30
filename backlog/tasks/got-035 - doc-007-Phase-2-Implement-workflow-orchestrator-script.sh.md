---
id: GOT-035
title: '[doc-007 Phase 2] Implement workflow orchestrator script.sh'
status: To Do
assignee: []
created_date: '2026-03-30 12:25'
labels:
  - implementation
  - core
dependencies: []
references:
  - doc-007#implementation-checklist
  - doc-007#validation-rules
documentation:
  - doc-007
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the main orchestrator Bash script (agents/workflow/script.sh) that reads workflow configuration, tracks task state in YAML format, determines the next agent in sequence, assigns tasks via backlog CLI, and updates state files. The script must use only Bash string operations for YAML parsing, handle all validation rules with proper error messages, and support task status transitions (pending → in_progress → finished).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Script reads config.yml on execution
- [ ] #2 Script reads/writes tasks.yml for state
- [ ] #3 Next agent determined by completed agents count (0-based)
- [ ] #4 Task assigned via backlog task edit
- [ ] #5 State file updated with timestamps
- [ ] #6 Status transitions implemented (pending → in_progress → finished)
- [ ] #7 Task marked finished when all agents complete
- [ ] #8 Config changes require manual intervention
- [ ] #9 Workflow aborts on agent failure
- [ ] #10 Exit codes: 0 success, 1 failure
<!-- AC:END -->

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
