---
id: GOT-026
title: 'Bug: Script execution missing task file path argument'
status: To Do
assignee: []
created_date: '2026-03-16 00:30'
updated_date: '2026-03-16 00:31'
labels:
  - bug
  - script-execution
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix the script execution to pass the task file path as an argument to agent scripts when they are invoked. Currently, scripts are executed without any arguments, so the task file path ($1) is empty even though the script is triggered.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 - [ ] When assignee changes, agent scripts are invoked with the task file path as the first argument
- [ ] #2 - [ ] The script receives the full absolute path to the task file
- [ ] #3 - [ ] Agent scripts can access and process the task file content via the passed argument
- [ ] #4 - [ ] Manual test: Create a task with assignee `[agent-bar]`, verify script receives the file path in the log output
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
