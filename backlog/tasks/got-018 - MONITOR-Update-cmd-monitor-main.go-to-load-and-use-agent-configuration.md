---
id: GOT-018
title: '[MONITOR] Update cmd/monitor/main.go to load and use agent configuration'
status: To Do
assignee: []
created_date: '2026-03-15 17:17'
updated_date: '2026-03-15 18:01'
labels: []
dependencies: []
references:
  - backlog/docs/doc-004-per-agent-configuration.md
priority: high
ordinal: 8000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update monitor entry point to use agent configuration
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 cmd/monitor/main.go modified to load agent configuration at startup
- [ ] #2 Configuration loaded from AGENT_NAME and AGENTS_CONFIG_DIR environment variables
- [ ] #3 AgentConfig passed to notifier for script execution
- [ ] #4 On assignee change, notifier.ExecuteScript called with configured script path
- [ ] #5 Script path and tmux session from agent configuration
- [ ] #6 Missing config file logs warning but monitor continues
- [ ] #7 Config parsing errors handled gracefully
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
