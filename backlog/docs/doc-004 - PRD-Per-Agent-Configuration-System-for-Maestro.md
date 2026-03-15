---
id: doc-004
title: 'PRD: Per-Agent Configuration System for Maestro'
type: other
created_date: '2026-03-15 16:30'
updated_date: '2026-03-15 17:17'
---
### Example Configuration

```yaml
# agents/bob/config.yml
script_path: "agents/bob/config/default.sh"
tmux_session: "bob"
enabled: true

# agents/alice/config.yml
script_path: "agents/alice/config/default.sh"
tmux_session: "alice"
enabled: true
```

## Implementation Tasks

The following tasks have been created in Backlog.md to implement this feature:

| Task ID | Title | Status | Priority |
|---------|-------|--------|----------|
| [GOT-015](backlog://task/got-015) | [CONFIG] Create pkg/config package | To Do | High |
| [GOT-016](backlog://task/got-016) | [AGENT] Create pkg/agent package | To Do | High |
| [GOT-017](backlog://task/got-017) | [NOTIFY] Modify pkg/notifier to execute bash scripts | To Do | High |
| [GOT-018](backlog://task/got-018) | [MONITOR] Update cmd/monitor/main.go | To Do | High |
| [GOT-019](backlog://task/got-019) | [DOCS] Create example configs and documentation | To Do | Medium |
