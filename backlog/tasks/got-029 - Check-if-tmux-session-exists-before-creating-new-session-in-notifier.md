---
id: GOT-029
title: Check if tmux session exists before creating new session in notifier
status: To Do
assignee: []
created_date: '2026-03-16 11:47'
labels:
  - bug
  - tmux
  - notifier
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement a session existence check before creating new tmux sessions in the notifier package. Currently, the ExecuteScript and ExecuteScriptsForAgents methods attempt to create a new tmux session unconditionally with tmux new-session -d -s name, which fails with exit status 1 when the session already exists (e.g., when a user is attached to the session). The check should use tmux list-sessions to verify session existence before attempting creation.
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
