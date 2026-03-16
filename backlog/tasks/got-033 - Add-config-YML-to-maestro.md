---
id: GOT-033
title: Add config YML to maestro
status: To Do
assignee:
  - Catarina
created_date: '2026-03-16 17:36'
updated_date: '2026-03-16 17:36'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Configuration file for watch paths (currently hardcoded `./backlog/tasks`)
the file should be in ./
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
