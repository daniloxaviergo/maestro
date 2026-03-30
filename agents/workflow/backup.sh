#!/bin/bash
# Workflow orchestrator script
# Manages sequential agent execution for Backlog.md tasks
# Reads workflow configuration, tracks task state, assigns tasks to next agent

set -euo pipefail

# ============================================================================
# Constants
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_PATH="${SCRIPT_DIR}/config.yml"
STATE_PATH="${SCRIPT_DIR}/tasks.yml"

TASK_FILE="$1"

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

# Log the assignment to a file
LOG_FILE="/home/danilo/scripts/github/maestro/agents/workflow/execution.log"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

# ============================================================================
# Helper Functions
# ============================================================================

error() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] ERROR: $*" >&2
}

warning() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] WARNING: $*" >&2
}

info() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] INFO: $*"
}

# Extract YAML value from a file (simple key: value parsing)
# Usage: extract_yaml_value <file> <key>
extract_yaml_value() {
    local file="$1"
    local key="$2"
    grep -E "^${key}:" "$file" 2>/dev/null | sed "s/^${key}:[[:space:]]*//" | sed 's/^"//;s/"$//'
}

# Extract agents list from config (comma-separated)
# Usage: get_agents <config_path>
get_agents() {
    local config_file="$1"
    if [[ ! -f "$config_file" ]]; then
        error "Config file not found: $config_file"
        return 1
    fi
    extract_yaml_value "$config_file" "agents"
}

# Extract backlog command from config
# Usage: get_backlog_command <config_path>
get_backlog_command() {
    local config_file="$1"
    if [[ ! -f "$config_file" ]]; then
        error "Config file not found: $config_file"
        return 1
    fi
    extract_yaml_value "$config_file" "backlog_command"
}

# Load task state from state file
# Usage: load_task_state <state_file> <task_id>
# Outputs: status assigned_agent assigned_at completed_at (one per line)
load_task_state() {
    local state_file="$1"
    local task_id="$2"
    
    if [[ ! -f "$state_file" ]]; then
        # State file doesn't exist, return empty state
        echo "pending"
        echo ""
        echo ""
        echo ""
        return 0
    fi
    
    # Check if task exists in state file
    if ! grep -q "^${task_id}:" "$state_file" 2>/dev/null; then
        # Task not in state file, return empty state
        echo "pending"
        echo ""
        echo ""
        echo ""
        return 0
    fi
    
    # Extract task block - find lines starting with 2 spaces (indented values)
    # Stop when we hit a line that starts at column 0 (next task or comment)
    local start_line end_line
    start_line=$(grep -n "^${task_id}:" "$state_file" | head -1 | cut -d: -f1)
    
    # Find end of task block (next line starting at column 0, or end of file)
    local next_task_line
    next_task_line=$(tail -n +$((start_line + 1)) "$state_file" | grep -n "^[^[:space:]]" | head -1 | cut -d: -f1) || true
    
    if [[ -z "$next_task_line" ]]; then
        # No next task found, use end of file
        end_line=$(wc -l < "$state_file")
    else
        # Calculate actual end line (line after the last indented line)
        end_line=$((start_line + next_task_line - 1))
    fi
    
    # Extract the task block
    local task_block
    if [[ "$end_line" -ge "$start_line" ]]; then
        task_block=$(sed -n "${start_line},${end_line}p" "$state_file")
    else
        task_block=$(sed -n "${start_line},\$p" "$state_file")
    fi
    
    # Extract values - handle empty values properly
    local status assigned_agent assigned_at completed_at
    status=$(echo "$task_block" | grep -E "^[[:space:]]*status:" | head -1 | sed 's/.*status:[[:space:]]*//' | sed 's/^"//;s/"$//') || true
    assigned_agent=$(echo "$task_block" | grep -E "^[[:space:]]*assigned_agent:" | head -1 | sed 's/.*assigned_agent:[[:space:]]*//' | sed 's/^"//;s/"$//') || true
    assigned_at=$(echo "$task_block" | grep -E "^[[:space:]]*assigned_at:" | head -1 | sed 's/.*assigned_at:[[:space:]]*//' | sed 's/^"//;s/"$//') || true
    completed_at=$(echo "$task_block" | grep -E "^[[:space:]]*completed_at:" | head -1 | sed 's/.*completed_at:[[:space:]]*//' | sed 's/^"//;s/"$//') || true
    
    echo "${status:-pending}"
    echo "${assigned_agent:-}"
    echo "${assigned_at:-}"
    echo "${completed_at:-}"
}

# Update task state in state file
# Usage: update_task_state <state_file> <task_id> <status> <assigned_agent> <assigned_at> <completed_at>
update_task_state() {
    local state_file="$1"
    local task_id="$2"
    local status="$3"
    local assigned_agent="${4:-}"
    local assigned_at="${5:-}"
    local completed_at="${6:-}"

    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    # Create state file if it doesn't exist
    if [[ ! -f "$state_file" ]]; then
        echo "# Task state tracking for workflow agent" > "$state_file"
        echo "# Format: {task_id}:" >> "$state_file"
        echo "#   status: pending|in_progress|finished" >> "$state_file"
        echo "#   assigned_agent: agent_name" >> "$state_file"
        echo "#   assigned_at: \"YYYY-MM-DD HH:MM:SS\"" >> "$state_file"
        echo "#   completed_at: \"YYYY-MM-DD HH:MM:SS\" (when finished)" >> "$state_file"
        echo "" >> "$state_file"
    fi

    # If task doesn't exist, append new entry
    if ! grep -q "^${task_id}:" "$state_file" 2>/dev/null; then
        echo "${task_id}:" >> "$state_file"
        echo "  status: ${status}" >> "$state_file"
        if [[ -n "$assigned_agent" ]]; then
            echo "  assigned_agent: ${assigned_agent}" >> "$state_file"
        fi
        if [[ -n "$assigned_at" ]]; then
            echo "  assigned_at: \"${assigned_at}\"" >> "$state_file"
        fi
        if [[ -n "$completed_at" ]]; then
            echo "  completed_at: \"${completed_at}\"" >> "$state_file"
        fi
        echo "" >> "$state_file"
        return 0
    fi

    # Update existing task entry using temp file
    local temp_file="${state_file}.tmp"
    : > "$temp_file"  # Create empty temp file
    local in_task_block=false
    local task_block_done=false

    while IFS= read -r line || [[ -n "$line" ]]; do
        if [[ "$task_block_done" == "true" ]]; then
            # Already finished updating task block, just copy remaining lines
            echo "$line" >> "$temp_file"
            continue
        fi
        
        if [[ "$in_task_block" == "true" ]]; then
            # Currently in task block, check if we've reached the end
            if [[ "$line" =~ ^[^[:space:]] ]]; then
                # End of task block (next task starts at column 0)
                in_task_block=false
                task_block_done=true
                echo "$line" >> "$temp_file"
            fi
            continue
        fi
        
        # Not in task block yet
        if [[ "$line" =~ ^${task_id}: ]]; then
            in_task_block=true
            echo "$line" >> "$temp_file"
            echo "  status: ${status}" >> "$temp_file"

            if [[ -n "$assigned_agent" ]]; then
                echo "  assigned_agent: ${assigned_agent}" >> "$temp_file"
            fi
            if [[ -n "$assigned_at" ]]; then
                echo "  assigned_at: \"${assigned_at}\"" >> "$temp_file"
            fi
            if [[ -n "$completed_at" ]]; then
                echo "  completed_at: \"${completed_at}\"" >> "$temp_file"
            fi
            continue
        fi

        # Regular line, just copy it
        echo "$line" >> "$temp_file"
    done < "$state_file"

    # Copy temp file back
    mv "$temp_file" "$state_file"
}

# Get completed agents count from state file
# Usage: get_completed_agents_count <state_file> <task_id>
get_completed_agents_count() {
    local state_file="$1"
    local task_id="$2"
    local count=0

    if [[ ! -f "$state_file" ]]; then
        echo "0"
        return 0
    fi

    if ! grep -q "^${task_id}:" "$state_file" 2>/dev/null; then
        echo "0"
        return 0
    fi

    # Extract task block
    local start_line end_line
    start_line=$(grep -n "^${task_id}:" "$state_file" | head -1 | cut -d: -f1)
    local next_task_line
    next_task_line=$(tail -n +$((start_line + 1)) "$state_file" | grep -n "^[^[:space:]]" | head -1 | cut -d: -f1) || true
    
    if [[ -z "$next_task_line" ]]; then
        end_line=$(wc -l < "$state_file")
    else
        end_line=$((start_line + next_task_line - 1))
    fi

    local task_block
    if [[ "$end_line" -ge "$start_line" ]]; then
        task_block=$(sed -n "${start_line},${end_line}p" "$state_file")
    else
        task_block=$(sed -n "${start_line},\$p" "$state_file")
    fi

    # Count completed entries (where completed_at is set)
    if [[ -n "$task_block" ]] && echo "$task_block" | grep -q "completed_at:"; then
        # Count how many times completed_at appears
        count=$(echo "$task_block" | grep -c "completed_at:")
    fi

    echo "$count"
}

# Assign task to agent via backlog CLI
# Usage: assign_task <backlog_command> <task_id> <agent>
assign_task() {
    local backlog_command="$1"
    local task_id="$2"
    local agent="$3"
    
    info "Assigning task $task_id to $agent"
    
    if ! $backlog_command "$task_id" --assignee "$agent"; then
        error "Failed to assign task $task_id to $agent"
        return 1
    fi
    
    info "Task $task_id assigned to $agent successfully"
}

# ============================================================================
# Main Workflow Logic
# ============================================================================

main() {
    local task_id="$1"
    
    # Validate config file
    if [[ ! -f "$CONFIG_PATH" ]]; then
        error "Config file not found: $CONFIG_PATH"
        exit 1
    fi
    
    # Get workflow configuration
    local agents_str backlog_cmd
    agents_str=$(get_agents "$CONFIG_PATH") || exit 1
    backlog_cmd=$(get_backlog_command "$CONFIG_PATH") || exit 1
    
    # Parse agents into array
    IFS=', ' read -ra agents <<< "$agents_str"
    if [[ ${#agents[@]} -eq 0 ]]; then
        error "No agents defined in config"
        exit 1
    fi
    
    info "Workflow agents: ${agents[*]}"
    info "Backlog command: $backlog_cmd"
    
    # Extract task ID
    # local task_id
    # task_id=$(extract_task_id "$task_file") || exit 1
    info "Task ID: $task_id"
    
    # Load current task state - use a temp file to avoid issues with empty values
    local state_file_tmp
    state_file_tmp=$(mktemp)
    load_task_state "$STATE_PATH" "$task_id" > "$state_file_tmp"
    
    local status assigned_agent assigned_at completed_at
    status=$(sed -n '1p' "$state_file_tmp")
    assigned_agent=$(sed -n '2p' "$state_file_tmp")
    assigned_at=$(sed -n '3p' "$state_file_tmp")
    completed_at=$(sed -n '4p' "$state_file_tmp")
    rm -f "$state_file_tmp"

    info "Current state: status=$status, assigned_agent=$assigned_agent"

    # Determine next agent based on completed count
    local completed_count next_agent
    completed_count=$(get_completed_agents_count "$STATE_PATH" "$task_id")
    
    if [[ $completed_count -ge ${#agents[@]} ]]; then
        error "No more agents available (completed: $completed_count, agents: ${#agents[@]})"
        exit 1
    fi
    next_agent="${agents[$completed_count]}"

    info "Next agent: $next_agent (completed count: $completed_count)"

    # Determine action based on status
    local timestamp
    timestamp=$(date '+%Y-%m-%d %H:%M:%S')

    case "$status" in
        pending)
            # First assignment - assign to first agent
            info "Task is pending, assigning to $next_agent"
            if ! assign_task "$backlog_cmd" "$task_id" "$next_agent"; then
                exit 1
            fi

            # Update state
            update_task_state "$STATE_PATH" "$task_id" "in_progress" "$next_agent" "$timestamp"
            info "State updated: status=in_progress, assigned_agent=$next_agent, assigned_at=$timestamp"
            ;;
        in_progress)
            # Check if current assigned agent needs to be assigned
            # If assigned_agent is empty or different from next_agent, we need to assign
            if [[ -z "$assigned_agent" || "$assigned_agent" != "$next_agent" ]]; then
                info "Assigning task to $next_agent (was: ${assigned_agent:-none})"
                if ! assign_task "$backlog_cmd" "$task_id" "$next_agent"; then
                    exit 1
                fi

                update_task_state "$STATE_PATH" "$task_id" "in_progress" "$next_agent" "$timestamp"
                info "State updated: assigned_agent=$next_agent, assigned_at=$timestamp"
            else
                # Current agent completed their work
                info "Agent $assigned_agent completed work"
                
                # Mark completion
                update_task_state "$STATE_PATH" "$task_id" "in_progress" "$assigned_agent" "$assigned_at" "$timestamp"
                
                # Increment completed count
                completed_count=$((completed_count + 1))
                
                if [[ $completed_count -ge ${#agents[@]} ]]; then
                    # All agents completed, mark task as finished
                    info "All agents completed, marking task as finished"
                    update_task_state "$STATE_PATH" "$task_id" "finished" "$assigned_agent" "$assigned_at" "$timestamp"
                    info "Task $task_id workflow complete"
                    exit 0
                fi
                
                next_agent="${agents[$completed_count]}"
                
                # Assign to next agent
                info "Assigning to next agent: $next_agent"
                if ! assign_task "$backlog_cmd" "$task_id" "$next_agent"; then
                    exit 1
                fi
                
                update_task_state "$STATE_PATH" "$task_id" "in_progress" "$next_agent" "$timestamp"
                info "State updated: assigned_agent=$next_agent, assigned_at=$timestamp"
            fi
            ;;
        finished)
            info "Task $task_id is already finished"
            exit 0
            ;;
        *)
            error "Unknown status: $status"
            exit 1
            ;;
    esac
    
    info "Workflow orchestration complete"
    exit 0
}

# Run main function
main "$TASK_ID"


# /home/danilo/scripts/github/maestro/backlog/tasks/got-035
# ./agents/workflow/script.sh /home/danilo/scripts/github/maestro/backlog/tasks/got-035