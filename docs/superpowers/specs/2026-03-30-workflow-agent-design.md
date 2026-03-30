# Workflow Agent Design Spec

**Date:** 2026-03-30  
**Status:** Approved  
**Implementation:** Bash script in `agents/workflow/`

---

## Overview

The Workflow Agent orchestrates sequential execution of downstream agents in a task-handling pipeline. It reads a configurable agent sequence and manages task state persistently, triggering the next agent in the flow when the current one completes.

---

## Architecture

### Components

| Component | Purpose |
|---------|---------|
| `agents/workflow/script.sh` | Main orchestrator bash script |
| `agents/workflow/config.yml` | Workflow configuration (agent sequence, backlog command) |
| `agents/workflow/tasks.yml` | Persistent task state tracking |

### Data Flow

```
1. Task assigned to "workflow" agent
   ↓
2. workflow/script.sh runs with task file path
   ↓
3. Load config (agents sequence, backlog command)
   ↓
4. Load state (agents/workflow/tasks.yml)
   ↓
5. Determine next agent (based on completed agents)
   ↓
6. If more agents: backlog task edit --assignee <next_agent>
   ↓
7. Update state file
   ↓
8. Script exits
   ↓
9. Monitor detects assignee change, triggers next agent's script
   ↓
10. Repeat from step 2 until all agents complete
```

---

## Configuration

### `agents/workflow/config.yml`

```yaml
# Agent sequence for task processing (order matters)
agents: ["Catarina", "Thomas"]

# Backlog CLI command (full path or in PATH)
backlog_command: "/home/danilo/.local/bin/backlog"
```

**Fields:**
- `agents`: Array of agent names in execution order
- `backlog_command`: Path to backlog CLI (used for `backlog task edit`)

---

## State Management

### `agents/workflow/tasks.yml`

```yaml
got-016:
  status: in_progress  # pending | in_progress | finished
  agents:
    Catarina:
      assigned_at: "2024-06-01T10:00:00Z"
      status: completed
      completed_at: "2024-06-01T11:00:00Z"
    Thomas:
      assigned_at: "2024-06-01T11:00:00Z"
      status: assigned
got-017:
  status: pending
  agents: {}
```

**Fields:**
- `status`: Task workflow status
  - `pending`: Task added to workflow, not yet started
  - `in_progress`: At least one agent has worked on task
  - `finished`: All agents in sequence have completed
- `agents`: Map of agent name → agent state
  - `assigned_at`: ISO 8601 timestamp when agent was assigned
  - `status`: `assigned` or `completed`
  - `completed_at`: ISO 8601 timestamp when agent completed (if completed)

**Position tracking:** Current agent index is calculated dynamically from the `agents` map (not stored).

---

## Workflow Logic

### Pseudocode

```bash
workflow/script.sh TASK_FILE:
1. TASK_ID = basename(TASK_FILE, ".md")
2. CONFIG = load_config("agents/workflow/config.yml")
3. STATE = load_state("agents/workflow/tasks.yml")
4. 
5. # Get task state or initialize new task
6. TASK_STATE = get_task(STATE, TASK_ID) or init_task()
7. 
8. # Determine next agent
9. COMPLETED_AGENTS = get_completed_agents(TASK_STATE.agents)
10. NEXT_INDEX = length(COMPLETED_AGENTS)
11. 
12. IF NEXT_INDEX >= length(CONFIG.agents):
13.     # All agents completed
14.     TASK_STATE.status = "finished"
15.     save_state(STATE)
16.     exit 0
17. 
18. NEXT_AGENT = CONFIG.agents[NEXT_INDEX]
19. 
20. # Assign task to next agent
21. RESULT = run("backlog task edit ${TASK_ID} --assignee ${NEXT_AGENT}")
22. IF FAILED(RESULT):
23.     log_error("Failed to assign task to ${NEXT_AGENT}")
24.     exit 1
25. 
26. # Update state
27. TASK_STATE.agents[NEXT_AGENT].assigned_at = now()
28. TASK_STATE.agents[NEXT_AGENT].status = "assigned"
29. TASK_STATE.status = "in_progress"
30. save_state(STATE)
31. 
32. exit 0
```

### State Calculation (No Stored Index)

```bash
# Find last completed agent by checking status field
# Next agent = agent at index (completed_count) in config
# Example: If Catarina is completed, Thomas is next (index 1)
# If Thomas is completed, task is finished (no more agents)
```

---

## Error Handling

| Scenario | Action |
|--------|--------|
| Config file missing | Log error, exit 1 |
| State file missing | Initialize empty state |
| Task not in state | Initialize new task with `status: pending` |
| `backlog task edit` fails | Log error, keep state unchanged, exit 1 |
| Agent not in config | Log warning, skip (or exit if critical) |
| Task file not found in backlog | Log warning, continue with state update |

---

## Execution Example

### Initial State

**Task:** `got-016.md` assigned to "workflow"

**Config:**
```yaml
agents: ["Catarina", "Thomas"]
backlog_command: "/home/danilo/.local/bin/backlog"
```

**State (initial):**
```yaml
got-016:
  status: pending
  agents: {}
```

### Run 1: Workflow assigns Catarina

1. Script runs with `got-016.md`
2. No completed agents → `NEXT_INDEX = 0`
3. `NEXT_AGENT = "Catarina"`
4. `backlog task edit got-016 --assignee Catarina`
5. State updated:
   ```yaml
   got-016:
     status: in_progress
     agents:
       Catarina:
         assigned_at: "2024-06-01T10:00:00Z"
         status: assigned
   ```

### Run 2: Catarina completes, assigns to workflow

1. Script runs with `got-016.md`
2. Catarina completed → `NEXT_INDEX = 1`
3. `NEXT_AGENT = "Thomas"`
4. `backlog task edit got-016 --assignee Thomas`
5. State updated:
   ```yaml
   got-016:
     status: in_progress
     agents:
       Catarina:
         assigned_at: "2024-06-01T10:00:00Z"
         status: completed
         completed_at: "2024-06-01T11:00:00Z"
       Thomas:
         assigned_at: "2024-06-01T11:00:00Z"
         status: assigned
   ```

### Run 3: Thomas completes, assigns to workflow

1. Script runs with `got-016.md`
2. All agents completed → `NEXT_INDEX = 2` (equals config length)
3. `TASK_STATE.status = "finished"`
4. State updated:
   ```yaml
   got-016:
     status: finished
     agents:
       Catarina:
         assigned_at: "2024-06-01T10:00:00Z"
         status: completed
         completed_at: "2024-06-01T11:00:00Z"
       Thomas:
         assigned_at: "2024-06-01T11:00:00Z"
         status: completed
         completed_at: "2024-06-01T12:00:00Z"
   ```

---

## Testing

### Manual Testing

```bash
# 1. Create workflow directory and config
mkdir -p agents/workflow
cat > agents/workflow/config.yml <<'EOF'
agents: ["Catarina", "Thomas"]
backlog_command: "/home/danilo/.local/bin/backlog"
EOF

# 2. Initialize state file
cat > agents/workflow/tasks.yml <<'EOF'
got-016:
  status: pending
  agents: {}
EOF

# 3. Assign task to workflow agent
backlog task edit got-016 --assignee workflow

# 4. Monitor will trigger workflow script
# 5. Check state file for updates
cat agents/workflow/tasks.yml

# 6. Repeat until task is finished
```

### Verification Checklist

- [ ] Script creates config file if missing
- [ ] Script creates state file if missing
- [ ] Task state updates correctly after each agent assignment
- [ ] Task marked as "finished" when all agents complete
- [ ] Errors are logged appropriately
- [ ] State file is properly formatted YAML

---

## Implementation Tasks

1. Create `agents/workflow/config.yml` with example configuration
2. Create empty `agents/workflow/tasks.yml` file
3. Write `agents/workflow/script.sh` implementing workflow logic
4. Test with a sample task
5. Update documentation

---

## Future Enhancements (Out of Scope)

- Concurrent task processing
- Retry logic for failed assignments
- Timeout-based escalation
- Workflow metrics/logging
- Web UI for workflow state visualization
