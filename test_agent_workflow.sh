#!/bin/bash

set -euo pipefail

$task_file=$(echo "/home/danilo/scripts/github/maestro/backlog/tasks/got-035")

cat ./agents/workflow/tasks.yml
# returns empty because dont have any task

./agents/workflow/script.sh "$task_file"
# expected logs
# INFO: Workflow agents: catarina, thomas
# INFO: Backlog command: backlog task edit
# INFO: Task ID: got-035
# INFO: Current state: status=in_progress, assigned_agent=catarina

cat ./agents/workflow/tasks.yml
# return in_progress with catarina agent
# got-035:
#   - status: in_progress
#     assigned_agent: catarina
#     assigned_at: "2026-03-30 11:00:00"

# simulate catarina finished
backlog task edit got-035 --assignee "workflow"

./agents/workflow/script.sh "$task_file"
# expected logs
# INFO: Workflow agents: catarina, thomas
# INFO: Backlog command: backlog task edit
# INFO: Task ID: got-035
# INFO: Agent catarina completed work
# INFO: Current state: status=finished, assigned_agent=catarina
# INFO: Assigning to next agent: thomas
# INFO: Current state: status=in_progress, assigned_agent=thomas

cat ./agents/workflow/tasks.yml
# return in_progress with thomas agent
# got-035:
#   - status: finished
#     assigned_agent: catarina
#     assigned_at: "2026-03-30 11:00:00"
#     completed_at: "2026-03-30 11:01:00"
#   - status: in_progress
#     assigned_agent: thomas
#     assigned_at: "2026-03-30 11:01:00"

# simulate thomas finished
backlog task edit got-035 --assignee "workflow"

./agents/workflow/script.sh "$task_file"

cat ./agents/workflow/tasks.yml
# return in_progress with thomas agent
# got-035:
#   - status: finished
#     assigned_agent: catarina
#     assigned_at: "2026-03-30 11:00:00"
#     completed_at: "2026-03-30 11:01:00"
#   - status: finished
#     assigned_agent: thomas
#     assigned_at: "2026-03-30 11:01:00"
#     completed_at: "2026-03-30 11:02:00"

