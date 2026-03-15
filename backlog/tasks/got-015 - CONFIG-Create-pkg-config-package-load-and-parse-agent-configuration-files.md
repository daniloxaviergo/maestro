---
id: GOT-015
title: '[CONFIG] Create pkg/config package - load and parse agent configuration files'
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
Create configuration loading and parsing package for agent YAML files
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 pkg/config/config.go with LoadConfig function
- [ ] #2 pkg/config/types.go with AgentConfig struct (script_path, tmux_session, enabled)
- [ ] #3 LoadConfig reads YAML file from path and returns AgentConfig
- [ ] #4 Missing config file logs warning and returns default config
- [ ] #5 YAML parsing errors are caught and logged
- [ ] #6 AgentNameFromEnv() function to read AGENT_NAME environment variable
- [ ] #7 ConfigDirFromEnv() function to read AGENTS_CONFIG_DIR environment variable
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
