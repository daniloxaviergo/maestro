---
id: GOT-016
title: >-
  [AGENT] Create pkg/agent package - manage agent identity and configuration
  loading
status: To Do
assignee: []
created_date: '2026-03-15 17:16'
labels: []
dependencies: []
references:
  - backlog/docs/doc-004-per-agent-configuration.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create agent package to manage agent identity and configuration loading
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 pkg/agent/agent.go with Agent struct to manage agent identity and configuration
- [ ] #2 Agent.LoadConfig() method to load config from configured path
- [ ] #3 Agent.GetConfig() method to return loaded configuration
- [ ] #4 Agent.GetName() method to return agent name
- [ ] #5 Default config directory is ./agents/ configurable via AGENTS_CONFIG_DIR
- [ ] #6 Agent name from AGENT_NAME environment variable
- [ ] #7 Missing config file logs warning but doesn't crash agent
- [ ] #8 Config file path is {config_dir}/{agent_name}/config.yml
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
