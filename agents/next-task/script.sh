#!/bin/bash
# next-task script
# This script is executed when a task is assigned to next-task agent
# It reads the config, removes the current task, and echoes the next task

TASK_FILE="$1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/config.yml"
LOG_FILE="/home/danilo/scripts/github/maestro/agents/next-task/execution.log"

# Extract task ID from TASK_FILE
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

echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: Processing task: $TASK_FILE" >> "$LOG_FILE"
echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: Task ID: $TASK_ID" >> "$LOG_FILE"

# Check if config file exists
if [[ ! -f "$CONFIG_FILE" ]]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: ERROR: Config file not found: $CONFIG_FILE" >> "$LOG_FILE"
    exit 0
fi

# Read tasks from config (comma-separated values after 'tasks:')
TASKS_LINE=$(grep "^tasks:" "$CONFIG_FILE" | sed 's/^tasks:[[:space:]]*//')
if [[ -z "$TASKS_LINE" ]]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: ERROR: No tasks found in config"
    exit 0
fi

# Convert comma-separated list to array
IFS=',' read -ra TASKS_ARRAY <<< "$TASKS_LINE"
# Trim whitespace from each task ID
for i in "${!TASKS_ARRAY[@]}"; do
    TASKS_ARRAY[$i]=$(echo "${TASKS_ARRAY[$i]}" | xargs)
done

echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: Current tasks: ${TASKS_ARRAY[*]}" >> "$LOG_FILE"

# Find the current task in the array and remove it
TASK_FOUND=false
NEW_TASKS=()
for task in "${TASKS_ARRAY[@]}"; do
    if [[ "$task" == "$TASK_ID" ]]; then
        TASK_FOUND=true
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: Found current task '$TASK_ID', removing from list" >> "$LOG_FILE"
    else
        NEW_TASKS+=("$task")
    fi
done

# If task not found, exit
if [[ "$TASK_FOUND" != "true" ]]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: Task '$TASK_ID' not found in config tasks list, exiting" >> "$LOG_FILE"
    exit 0
fi

# Update config file with remaining tasks
if [[ ${#NEW_TASKS[@]} -eq 0 ]]; then
    # No more tasks, clear the tasks line
    sed -i "s/^tasks:.*$/tasks: /" "$CONFIG_FILE"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: All tasks completed, config cleared" >> "$LOG_FILE"
else
    # Build new tasks string
    NEW_TASKS_STRING=$(IFS=', '; echo "${NEW_TASKS[*]}")
    sed -i "s/^tasks:.*$/tasks: $NEW_TASKS_STRING/" "$CONFIG_FILE"
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: Updated tasks: ${NEW_TASKS[*]}" >> "$LOG_FILE"
fi

# Check if there's a next task
if [[ ${#NEW_TASKS[@]} -gt 0 ]]; then
    NEXT_TASK="${NEW_TASKS[0]}"
    echo "NEXT_TASK: $NEXT_TASK - $PROJECT_PATH" >> "$LOG_FILE"
    cd $PROJECT_PATH
    sleep 3
    echo "NEXT_TASK: $NEXT_TASK - book" >> "$LOG_FILE"
    backlog task edit "$NEXT_TASK" --assignee "book" >> "$LOG_FILE"
    sleep 3
    echo "NEXT_TASK: $NEXT_TASK - workflow" >> "$LOG_FILE"
    backlog task edit "$NEXT_TASK" --assignee "workflow" >> "$LOG_FILE"
    sleep 3
else
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: No more tasks in queue" >> "$LOG_FILE"
fi

echo "[$(date '+%Y-%m-%d %H:%M:%S')] next-task: Task processing complete" >> "$LOG_FILE"
