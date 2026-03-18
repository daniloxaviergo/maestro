#!/bin/bash
# catarina script
# This script is executed when a task is assigned to catarina
# The script receives the task file path as the first argument

TASK_FILE="$1"
START_TIME=$(date +%s)

# Extract task ID and project path from TASK_FILE
# Task file format: /path/to/project/backlog/tasks/got-XXX
# Task ID: e.g., got-028
# Project path: /path/to/project

# Extract the filename (task ID) from the path
TASK_ID=$(basename "$TASK_FILE")

# Find the project path by locating the backlog/tasks directory
# Walk up the directory tree until we find the project root
PROJECT_PATH=$(dirname "$TASK_FILE")
while [[ "$PROJECT_PATH" != "/" && "$PROJECT_PATH" != "." ]]; do
    if [[ -d "$PROJECT_PATH/backlog/tasks" ]]; then
        # Check if this is actually the project root by verifying backlog/config.yml exists
        if [[ -f "$PROJECT_PATH/backlog/config.yml" ]]; then
            break
        fi
    fi
    PROJECT_PATH=$(dirname "$PROJECT_PATH")
done

echo "[$(date '+%Y-%m-%d %H:%M:%S')] catarina: Processing task: $TASK_FILE"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] catarina: Task ID: $TASK_ID"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] catarina: Project path: $PROJECT_PATH"

# Log the assignment to a file
LOG_FILE="/home/danilo/scripts/github/maestro/agents/catarina/execution.log"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
echo "[$TIMESTAMP] Task assigned: $TASK_FILE (Task ID: $TASK_ID, Project: $PROJECT_PATH)" >> "$LOG_FILE"

echo "#############################################"
echo "#############################################"
echo "#############################################"

notify-send \
  -i /home/danilo/scripts/github/maestro/agents/catarina/icon.png \
  -a "Maestro" \
  "Catarina" \
  "Plan the task: $TASK_ID"

echo $PROJECT_PATH
cd $PROJECT_PATH
qwen "/plan $TASK_ID" --yolo --output-format stream-json --include-partial-messages | jq 'select(.type? == "assistant") | .message.content[]? | select(.type? == "text") | .text?'

END_TIME=$(date +%s)
ELAPSED=$((END_TIME - START_TIME))
ELAPSED_MINUTES=$((ELAPSED / 60))

echo "#############################################"
echo "#############################################"
echo "#############################################"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] catarina: Total time elapsed: ${ELAPSED_MINUTES}m" >> "$LOG_FILE"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] catarina: Task processing complete" >> "$LOG_FILE"

notify-send \
  -i /home/danilo/scripts/github/maestro/agents/catarina/icon.png \
  -w \
  -a "Maestro" \
  "Catarina" \
  "Finished task: $TASK_ID ${ELAPSED_MINUTES}m"
