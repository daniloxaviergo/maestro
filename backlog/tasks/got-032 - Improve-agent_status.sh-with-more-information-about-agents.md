---
id: GOT-032
title: Improve agent_status.sh with more information about agents
status: To Do
assignee:
  - Catarina
created_date: '2026-03-16 15:28'
updated_date: '2026-03-16 15:46'
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
- [x] #1 Code follows existing project conventions package structure naming error handling
- [x] #2 go vet passes with no warnings
- [x] #3 go build succeeds without errors
- [ ] #4 Unit tests added or updated for new or changed functionality
- [x] #5 go test ... passes with no failures
- [ ] #6 Code comments added for non-obvious logic
- [x] #7 README or docs updated if public behavior changes
- [x] #8 make build succeeds
- [ ] #9 make run works as expected
- [x] #10 Errors are logged not silently ignored
- [x] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
<!-- DOD:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Output table has 5 columns: Name, Status, Processing In, Task Count, Avg Duration
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The goal is to enhance `scripts/agent_status.sh` to display a comprehensive table with agent statistics. The script will parse agent execution logs to calculate:

1. **Status**: RUNNING (last log line doesn't contain "Task processing complete") or IDLE
2. **Processing In**: Duration since the last "Task assigned" log entry for running agents
3. **Task Count**: Total number of task assignments in the log file
4. **Avg Duration**: Average time spent processing tasks (calculated from "Total time elapsed" entries)

The implementation will:
- Parse execution.log files for each enabled agent
- Extract timestamps from log entries using regex
- Calculate time differences for processing duration
- Track task counts and elapsed times for average calculations
- Format output as a table with aligned columns
- Support both human-readable and JSON output formats

Key design decisions:
- Parse logs from `agents/{agent}/execution.log` paths
- Use bash date parsing and arithmetic for time calculations
- Gracefully handle missing or malformed log entries
- Follow existing script conventions for YAML parsing and agent discovery

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `scripts/agent_status.sh` | Modify | Main implementation: add processing duration, task count, avg duration calculation and table display |

### 3. Dependencies

- No new dependencies required (uses standard bash tools)
- Existing prerequisites:
  - Agents must have execution.log files (created by agent scripts)
  - Agent scripts must log "Task assigned" and "Total time elapsed" entries
  - Bash 4.0+ for associative arrays and advanced features
  - Standard Unix tools: grep, sed, awk, tail, date

### 4. Code Patterns

Follow existing patterns in `agent_status.sh`:
- Use `set -euo pipefail` for error handling
- Helper functions for YAML parsing and timestamp extraction
- Array-based data collection (`agents_data`)
- Separate `output_human()` and `output_json()` for display formats
- Exit with code 1 on errors, 0 on success
- Log warnings to stderr, info to stdout

New patterns to add:
- Time calculation function using `date +%s` for epoch timestamps
- Loop-based log parsing for specific patterns
- Accumulator variables for count and duration aggregation
- Floating-point division using `awk` for average calculation

### 5. Testing Strategy

- Test with existing execution.log files (catarina, agent-bar)
- Create mock log files with various scenarios:
  - Running agent (last line doesn't contain "Task processing complete")
  - Idle agent (last line contains "Task processing complete")
  - Agent with multiple task assignments
  - Agent with missing elapsed time entries
- Verify output format matches specification:
  - Columns aligned correctly
  - Duration formatted as "Xm" for minutes
  - Empty cells for missing data
- Test JSON output structure
- Test edge cases: no agents, disabled agents, malformed logs

### 6. Risks and Considerations

**Potential Issues:**
- Log format variations: Some agents may not log "Total time elapsed" entries
- Time calculations: Bash doesn't have native floating-point math (use `awk`)
- Timezone handling: Timestamps should use consistent format (currently `YYYY-MM-DD HH:MM:SS`)
- Processing duration: For running agents, calculate from last "Task assigned" to current time

**Trade-offs:**
- Log parsing approach: Reading entire log files vs. parsing only relevant lines
- Average calculation: Simple average vs. weighted average (use simple for initial implementation)
- Edge case handling: Show "N/A" for missing data vs. skipping the field

**Future Enhancements:**
- Support for different time units (hours, seconds)
- Historical data aggregation across multiple log files
- Task-specific duration tracking
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed successfully
<!-- SECTION:NOTES:END -->
