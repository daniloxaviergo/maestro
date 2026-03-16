---
id: GOT-032
title: Improve agent_status.sh with more information about agents
status: To Do
assignee:
  - Catarina
created_date: '2026-03-16 15:28'
updated_date: '2026-03-16 15:28'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
in `scripts/agent_status.sh`  returns a table like:

```text
Name          | Status          | Processing In          | Task Count          | Avg Duration
catarina       | RUNNING   | 20m                         | 50                        | 10m
```

name - name of agent
status - RUNNING, IDLE
processing in - if status RUNNING how long processing the task
task count - how many tasks was processed
avg duration - average time spend to process task
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
