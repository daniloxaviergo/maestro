---
id: doc-007
title: Workflow Agent Implementation
type: other
created_date: '2026-03-30 12:18'
---


# PRD: Workflow Agent Implementation

## Executive Summary

The Workflow Agent orchestrates sequential execution of downstream agents (Catarina, Thomas) in a task-handling pipeline. It manages task state persistently in YAML format and triggers the next agent in the sequence when the current one completes. The agent is implemented as a Bash script (`agents/workflow/script.sh`) that reads configuration, updates state on disk, and assigns tasks to agents via the Backlog CLI. No modifications to existing agent scripts are required—the workflow only assigns the next agent via `backlog task edit`.

---

## Key Requirements

| Requirement | Description | Status |
|-------------|-------------|--------|
| **R1** | Read workflow configuration from `agents/workflow/config.yml` | ✅ Defined |
| **R2** | Track task status in `agents/workflow/tasks.yml` | ✅ Defined |
| **R3** | Dynamically determine next agent based on completed agents | ✅ Defined |
| **R4** | Assign task to next agent via `backlog task edit` | ✅ Defined |
| **R5** | Update state file with assigned/completed timestamps | ✅ Defined |
| **R6** | Mark task as "finished" when all agents complete | ✅ Defined |
| **R7** | Handle config changes gracefully (manual intervention required) | ✅ Defined |
| **R8** | Abort immediately on downstream agent script errors | ✅ Defined |
| **R9** | Use simple YAML (flat key-value, no anchors/aliases) | ✅ Defined |
| **R10** | Single sequential execution (no file locking) | ✅ Defined |

---

## Technical Decisions

| Decision | Rationale |
|----------|-----------|
| **Language: Bash** | Existing agents use Bash; simple to implement with standard tools |
| **State format: Flat YAML** | Simple to parse and edit with Bash string operations |
| **No file locking** | Single execution model (no concurrent workflows per task) |
| **Backlog CLI only** | CLI always available, no REST fallback needed |
| **Bash YAML parsing/writing** | No external dependencies (`yq`, `jq`), simple string operations |
| **Task ID from file basename** | Simple, deterministic extraction from task file path |
| **State file: Direct write** | Atomic write not needed for single-execution workflow |
| **Partial state updates** | Preserve existing fields not being updated |
| **YAML format: Single-line per entry** | Simple and readable |
| **Agent assignment only** | Workflow assigns agent; agent's own script handles execution |

---

## Acceptance Criteria

### Functional

- [ ] **AC1**: Script reads `agents/workflow/config.yml` on execution
- [ ] **AC2**: Script reads/writes `agents/workflow/tasks.yml` for state
- [ ] **AC3**: Next agent is determined by counting completed agents (0-based index)
- [ ] **AC4**: Task is assigned via `backlog task edit <task_id> --assignee <agent>`
- [ ] **AC5**: State file is updated with agent assignment/completion timestamps
- [ ] **AC6**: Task status transitions: `pending` → `in_progress` → `finished`
- [ ] **AC7**: Task is marked "finished" when all config agents complete
- [ ] **AC8**: Config changes (agent added/removed) require manual intervention
- [ ] **AC9**: Workflow aborts immediately on downstream agent failure
- [ ] **AC10**: Script exits with code 0 on success, 1 on failure

### Non-Functional

- [ ] **NFC1**: State file uses simple, flat YAML structure
- [ ] **NFC2**: YAML parsing/writing uses only Bash string operations
- [ ] **NFC3**: No external dependencies (`yq`, `jq`, Python)
- [ ] **NFC4**: Single execution model (no concurrent workflow runs)
- [ ] **NFC5**: Agent scripts unchanged (no integration needed)

---

## Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `agents/workflow/config.yml` | **Create** | Workflow configuration (agent sequence, backlog command) |
| `agents/workflow/tasks.yml` | **Create** | Persistent task state tracking |
| `agents/workflow/script.sh` | **Create** | Main orchestrator Bash script |
| `agents/workflow/` | **Create** | Directory for workflow files |

---

## Files Created

| File | Purpose |
|------|---------|
| `agents/workflow/config.yml` | Workflow configuration (agent sequence, backlog command) |
| `agents/workflow/tasks.yml` | Persistent task state tracking |
| `agents/workflow/script.sh` | Main orchestrator Bash script |
| `docs/workflow-agent-quickstart.md` | Setup and configuration guide (out of scope) |

---

## Validation Rules

| Context | Rule | Error Message |
|---------|------|---------------|
| **Config loading** | Config file must exist | `Error: Config file not found: <path>` |
| **State loading** | Missing state file = empty state | `Warning: Initializing new state file` |
| **Task ID** | Must match `got-XXX` pattern | `Error: Invalid task ID: <id>` |
| **Agent assignment** | Agent must exist in config | `Error: Agent <name> not in config` |
| **State update** | YAML must be valid | `Error: Failed to write state: <error>` |
| **Backlog CLI** | Command must succeed | `Error: Failed to assign task: <error>` |

---

## Out of Scope

- **Concurrent task processing**
- **Retry logic for failed assignments**
- **Timeout-based escalation**
- **Workflow metrics/logging**
- **Web UI for state visualization**
- **Agent script modifications**
- **Dynamic agent discovery/registering**
- **Task branching/forking/merging**

---

## Implementation Checklist

- [ ] **Step 1**: Create `agents/workflow/` directory
- [ ] **Step 2**: Create `agents/workflow/config.yml` with default configuration
- [ ] **Step 3**: Create empty `agents/workflow/tasks.yml`
- [ ] **Step 4**: Implement `agents/workflow/script.sh`:
  - [ ] Parse task ID from input file path
  - [ ] Load configuration from `config.yml`
  - [ ] Load state from `tasks.yml`
  - [ ] Calculate completed agents count
  - [ ] Determine next agent index
  - [ ] Assign task via `backlog task edit`
  - [ ] Update state file
  - [ ] Handle errors appropriately
- [ ] **Step 5**: Test with sample task (`got-016`)
- [ ] **Step 6**: Verify state transitions (pending → in_progress → finished)

---

## Stakeholder Alignment

| Stakeholder | Responsibility | Verification |
|-------------|----------------|----------------|
| **Product Owner** | Approve agent sequence, config structure | Review acceptance criteria |
| **Backend Developer** | Implement script, YAML handling | Code review, test coverage |
| **QA Engineer** | Test workflow transitions, edge cases | Manual testing checklist |
| **DevOps** | Configure Backlog CLI, test tmux integration | Integration testing |

---

## Traceability Matrix

| Requirement | User Story | Acceptance Criteria | Test Case |
|-------------|------------|---------------------|-----------|
| **R1** | As a developer, I want to configure agent sequence | AC1 | Config file parsed correctly |
| **R2** | As a developer, I want to track task state | AC2 | State file updated after each assignment |
| **R3** | As a developer, I want agents to execute in order | AC3 | Next agent = completed_count (0-based) |
| **R4** | As a developer, I want tasks assigned automatically | AC4 | `backlog task edit` executed for each step |
| **R5** | As a developer, I want timestamps for all events | AC5 | `assigned_at`, `completed_at` fields populated |
| **R6** | As a developer, I want tasks marked finished | AC6, AC7 | Status = "finished" when all agents complete |
| **R7** | As a developer, I want clear error on config change | AC8 | State machine blocks on config mismatch |
| **R8** | As a developer, I want clean failure handling | AC9 | Script exits with code 1 on agent failure |
| **R9** | As a developer, I want simple YAML | NFC1, NFC2 | YAML uses flat structure, Bash string ops only |
| **R10** | As a developer, I want single execution | NFC4 | No file locking, No concurrent runs |

---

## Validation

| Check | Status |
|-------|--------|
| **Code Quality Standards** | ✅ Bash script follows existing patterns |
| **Technical Feasibility** | ✅ No external dependencies required |
| **Alignment with User Needs** | ✅ Addresses sequential workflow orchestration |
| **Backlog Integration** | ✅ Uses existing `backlog task edit` CLI |
| **State Persistence** | ✅ Simple YAML file on disk |

---

## Ready for Implementation

The PRD is unambiguous, technically feasible, and stakeholder-aligned. Ready for developers to begin implementation of `agents/workflow/`.
