---
id: GOT-030
title: Get correct task id and workdir
status: To Do
assignee: []
created_date: '2026-03-16 14:31'
updated_date: '2026-03-16 14:39'
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
