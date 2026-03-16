---
id: doc-006
title: 'PRD: Tmux Session Attacher with fzf'
type: other
created_date: '2026-03-16 00:41'
---
# PRD: Tmux Session Attacher with fzf

## Overview

### Purpose
Create a bash script that allows users to attach to tmux sessions for Maestro agents using fzf for interactive agent selection. This simplifies the workflow for developers and administrators who need to monitor agent activity.

### Goals
- Enable one-command agent session attachment with interactive selection
- Reduce manual effort in managing multiple agent tmux sessions
- Improve developer experience when troubleshooting or monitoring agents
- Create a reusable pattern for tmux session management

## Background

### Problem Statement
Currently, the Makefile provides basic tmux commands (`tmux-start`, `tmux-attach`) but only for a hardcoded `maestro-test` session. Users managing multiple agents must manually construct tmux commands with exact session names, which is error-prone and inconvenient.

### Current State
- Developers must know the exact tmux session name for each agent
- Session names are stored in agent config files (`tmux_session` field)
- No discoverable way to list available agent sessions
- Manual command construction: `tmux attach -t <session-name>`

### Proposed Solution
A bash script that:
1. Scans the `agents/` directory for configured agents
2. Extracts session names from each agent's config
3. Uses fzf to present an interactive list for selection
4. Attaches to the selected tmux session

## Requirements

### User Stories

- **Developer/Admin**: As a developer or system administrator, I want to attach to any agent's tmux session with a single command and fzf selection, so I don't need to remember or look up session names.

- **Multiple Agents**: As someone managing multiple agents, I want to see all available agent sessions in one place, so I can quickly switch between them.

### Functional Requirements

#### Task 1: fzf-based Session Selection Script

Create a bash script `scripts/attach.sh` that discovers and attaches to agent tmux sessions.

##### User Flows
1. User runs `./scripts/attach.sh` from project root
2. Script scans `agents/` directory for agent subdirectories
3. For each agent, script reads `config.yml` and extracts `tmux_session` value
4. Script presents fzf menu with agent names and their session names
5. User selects an agent via fzf
6. Script attaches to the corresponding tmux session
7. If session doesn't exist, script displays error message and exits with non-zero code

##### Acceptance Criteria
- [ ] Script scans `agents/` directory recursively for subdirectories
- [ ] Script reads `config.yml` from each agent directory
- [ ] Script extracts `tmux_session` value from YAML config
- [ ] Script displays fzf menu with agent names and session names
- [ ] Script attaches to selected tmux session
- [ ] Script exits with code 1 if selected session doesn't exist
- [ ] Script exits with code 0 on successful attach

#### Task 2: Error Handling

Implement robust error handling for edge cases.

##### Acceptance Criteria
- [ ] Script handles missing `agents/` directory gracefully
- [ ] Script handles missing or invalid config files (skips with warning)
- [ ] Script handles missing `tmux_session` field (skips with warning)
- [ ] Script handles fzf cancellation (exits cleanly with code 130)
- [ ] Script handles tmux not installed (graceful error message)
- [ ] Script handles fzf not installed (graceful error message)

#### Task 3: Integration with Makefile

Add convenient Makefile commands for the new script.

##### Acceptance Criteria
- [ ] Makefile targets: `attach` (runs the script)
- [ ] Makefile target: `attach-list` (lists all agents and sessions without fzf)

##### Acceptance Criteria for Makefile
- [ ] `make attach` runs `./scripts/attach.sh`
- [ ] `make attach-list` shows all agents with their session names

### Non-Functional Requirements

- **Performance**: Script should complete discovery and display fzf menu within 500ms
- **Compatibility**: Bash 4.0+, tmux 2.0+, fzf 0.20+
- **Maintainability**: Code should be well-commented, follow bash best practices
- **Error Messages**: All errors should be clear and actionable

## Scope

### In Scope
- Bash script `scripts/attach.sh` for tmux session attachment
- fzf-based interactive agent selection
- Discovery of agents from `agents/` directory
- Reading session names from agent config files
- Basic error handling and user feedback
- Makefile integration

### Out of Scope
- Creating tmux sessions (only attach to existing)
- Session management (create/kill/list sessions)
- Customizing fzf options (colors, layout, etc.)
- Persistent configuration for default session
- Multiple session attachment (multi-attach)

## Technical Considerations

### Existing System Impact
- New script in `scripts/` directory (parallel to `cmd/` and `pkg/`)
- No changes to Go code required
- Existing tmux commands in Makefile remain unchanged

### Dependencies
- **bash 4.0+**: Standard shell on most systems
- **fzf**: Required for interactive selection
- **tmux**: Required for session attachment
- **YAML parser**: Use `grep`/`sed` or `yq` for config parsing

### Constraints
- Must work in the Maestro project directory structure
- Session names must match those in agent config files
- No changes to existing agent configuration format

## Success Metrics

### Quantitative
- Time to attach to session: < 2 seconds (including fzf display)
- Script execution time: < 500ms for discovery

### Qualitative
- Users can attach to any agent session without reading docs
- Error messages clearly indicate what went wrong
- Fallback behavior is predictable and non-blocking

## Timeline & Milestones

- [ ] Design complete: PRD reviewed and approved
- [ ] Implementation complete: Script and Makefile integration
- [ ] Testing complete: Manual testing with multiple agents
- [ ] Documentation: Update README with new command

## Stakeholders

### Decision Makers
- Product Owner: Approve PRD scope and priorities

### Contributors
- Developer: Implement bash script
- QA/Testing: Verify error handling and edge cases

## Appendix

### Glossary
- **tmux**: Terminal multiplexer for managing multiple sessions
- **fzf**: Fuzzy finder for interactive selection
- **agent**: Maestro component that processes task assignments

### References
- [tmux man page](https://man7.org/linux/man-pages/man1/tmux.1.html): Tmux commands and options
- [fzf GitHub](https://github.com/junegunn/fzf): Fzf documentation
- Agent config format: `agents/*/config.yml`
- Current Makefile tmux commands: `Makefile`

### Example Agent Config
```yaml
script_path: "./agents/agent-bar/script.sh"
tmux_session: "agent-bar"
enabled: true
```

### Example Usage
```bash
# Attach to any agent session
./scripts/attach.sh

# Or via Makefile
make attach

# List all available sessions
make attach-list
```
