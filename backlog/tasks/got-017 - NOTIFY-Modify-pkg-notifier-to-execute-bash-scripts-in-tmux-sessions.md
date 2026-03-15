---
id: GOT-017
title: '[NOTIFY] Modify pkg/notifier to execute bash scripts in tmux sessions'
status: To Do
assignee: []
created_date: '2026-03-15 17:17'
updated_date: '2026-03-15 18:01'
labels: []
dependencies: []
references:
  - backlog/docs/doc-004-per-agent-configuration.md
priority: high
ordinal: 7000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify notifier to execute bash scripts in tmux sessions
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 pkg/notifier/notifier.go modified to execute bash scripts in tmux sessions
- [ ] #2 New ExecuteScript method to run script in tmux session
- [ ] #3 ExecuteScript uses tmux send-keys to run script
- [ ] #4 ExecuteScript handles missing script file with error logging
- [ ] #5 ExecuteScript handles script execution failures with error logging
- [ ] #6 ExecuteScript is non-blocking (runs in goroutine)
- [ ] #7 Script output captured but not displayed
- [ ] #8 tmux session name from configuration
- [ ] #9 tmux session creation if it doesn't exist
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
