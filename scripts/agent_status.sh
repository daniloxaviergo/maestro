#!/usr/bin/env bash
#
# scripts/agent_status.sh - Agent status checker
#
# This script checks the status of agents by examining their execution logs.
# An agent is considered:
#   - "running" if the last log line is NOT "Task processing complete"
#   - "idle" if the last log line IS "Task processing complete" or no log exists
#
# Usage:
#   ./scripts/agent_status.sh          # Show human-readable status
#   ./scripts/agent_status.sh --json   # Output status as JSON
#   ./scripts/agent_status.sh --help   # Show this help message
#
# Dependencies:
#   - bash 4.0+
#   - Standard Unix tools: grep, sed, awk, tail
#

set -euo pipefail

# Project root (script is in scripts/ so we need to go up one level)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
AGENTS_DIR="${PROJECT_ROOT}/agents"

# =============================================================================
# Helper Functions
# =============================================================================

error() {
    >&2 echo "ERROR: $*"
}

warning() {
    >&2 echo "WARNING: $*"
}

info() {
    echo "$*"
}

# Extract value from YAML config file using grep/sed
# Usage: extract_yaml_value "config.yml" "key_name"
extract_yaml_value() {
    local config_file="$1"
    local key_name="$2"
    local value

    # Check if config file exists
    if [[ ! -f "${config_file}" ]]; then
        return 1
    fi

    # Extract the value using grep and sed
    # Handles: key: "value", key: value, key: "value with spaces"
    # This pattern captures everything after the colon, removing surrounding quotes
    value=$(grep -E "^[[:space:]]*${key_name}:" "${config_file}" 2>/dev/null | head -n1 | sed -E 's/^[^:]*:[[:space:]]*//; s/^["'"'"']//' | sed "s/[\"']$//" || true)

    if [[ -z "${value}" ]]; then
        return 1
    fi

    echo "${value}"
}

# Check if an agent is running or idle based on its execution log
# Returns: "running", "idle", or "unknown"
# Arguments: agent_name, log_file_path
check_agent_status() {
    local agent_name="$1"
    local log_file="$2"

    # If no log file exists, agent has never run (idle)
    if [[ ! -f "${log_file}" ]]; then
        echo "idle"
        return 0
    fi

    # Check if log file is empty
    if [[ ! -s "${log_file}" ]]; then
        echo "idle"
        return 0
    fi

    # Get the last line of the log file
    local last_line
    last_line=$(tail -n1 "${log_file}" 2>/dev/null || true)

    # Check if the last line contains "Task processing complete"
    if echo "${last_line}" | grep -q "Task processing complete"; then
        echo "idle"
    else
        echo "running"
    fi
}

# Get the last log timestamp from an agent's log file
# Returns: timestamp string or empty if no log exists
get_last_log_timestamp() {
    local log_file="$1"

    if [[ ! -f "${log_file}" ]]; then
        echo ""
        return 0
    fi

    # Extract timestamp from first line of log (format: [YYYY-MM-DD HH:MM:SS])
    tail -n1 "${log_file}" 2>/dev/null | grep -oE '\[[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\]' | tr -d '[]' || true
}

# Count task assignments in an agent's log file
# Arguments: log_file_path
# Returns: count of task assignments
count_task_assignments() {
    local log_file="$1"

    if [[ ! -f "${log_file}" ]] || [[ ! -s "${log_file}" ]]; then
        echo "0"
        return 0
    fi

    grep -c "Task assigned:" "${log_file}" 2>/dev/null || echo "0"
}

# Calculate total elapsed time from log file and compute average duration
# Arguments: log_file_path
# Returns: average duration in minutes (integer) or 0
calculate_avg_duration() {
    local log_file="$1"
    local total_minutes=0
    local count=0

    if [[ ! -f "${log_file}" ]] || [[ ! -s "${log_file}" ]]; then
        echo "0"
        return 0
    fi

    # Extract all "Total time elapsed: Xm" entries and sum them
    while IFS= read -r line; do
        # Extract minutes from "Total time elapsed: Xm"
        if [[ "$line" =~ Total\ time\ elapsed:\ ([0-9]+)m ]]; then
            total_minutes=$((total_minutes + ${BASH_REMATCH[1]}))
            ((count++)) || true
        fi
    done < <(grep "Total time elapsed:" "${log_file}" 2>/dev/null || true)

    # Calculate average
    if [[ ${count} -gt 0 ]]; then
        echo $((total_minutes / count))
    else
        echo "0"
    fi
}

# Get the timestamp of the last "Task assigned" entry for processing duration
# Arguments: log_file_path
# Returns: epoch timestamp of last task assignment, or current time if not found
get_last_task_assigned_epoch() {
    local log_file="$1"

    if [[ ! -f "${log_file}" ]] || [[ ! -s "${log_file}" ]]; then
        date +%s
        return 0
    fi

    # Extract timestamp from last "Task assigned" line
    local timestamp
    timestamp=$(grep "Task assigned:" "${log_file}" 2>/dev/null | tail -n1 | grep -oE '\[[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\]' | tr -d '[]' || true)

    if [[ -z "${timestamp}" ]]; then
        echo "$(date +%s)"
        return 0
    fi

    # Convert timestamp to epoch
    date -d "${timestamp}" +%s 2>/dev/null || date +%s
}

# Format processing duration for display
# Converts seconds to human-readable format (e.g., "20m" or "1h15m")
# Arguments: epoch_duration (seconds)
# Returns: formatted string like "20m" or "1h15m"
format_duration() {
    local seconds=$1

    if [[ ${seconds} -lt 0 ]]; then
        echo "N/A"
        return 0
    fi

    # Convert seconds to hours and minutes using bash arithmetic
    # minutes % 60 gives remaining minutes after removing full hours
    local minutes=$((seconds / 60))
    local hours=$((minutes / 60))
    local remaining_minutes=$((minutes % 60))

    # If >= 1 hour, show "XhYm"; otherwise show "Xm"
    if [[ ${hours} -gt 0 ]]; then
        echo "${hours}h${remaining_minutes}m"
    else
        echo "${minutes}m"
    fi
}

# Discover all agents and their status
# Arguments: format (human, json)
# Returns: formatted output of agent statuses
discover_agents() {
    local format="${1:-human}"
    local agent_dir config_file script_path log_path agent_name status last_timestamp enabled
    local -a agents_data=()

    # Check if agents directory exists
    if [[ ! -d "${AGENTS_DIR}" ]]; then
        error "Agents directory not found: ${AGENTS_DIR}"
        exit 1
    fi

    # Scan agent directories
    for agent_dir in "${AGENTS_DIR}"/*/; do
        # Skip if no agent directories found
        [[ -d "${agent_dir}" ]] || continue

        agent_name=$(basename "${agent_dir}")
        config_file="${agent_dir}config.yml"

        # Get agent configuration
        enabled=$(extract_yaml_value "${config_file}" "enabled" || echo "true")
        script_path=$(extract_yaml_value "${config_file}" "script_path" || echo "")

        # Default log path: agents/{agent}/execution.log
        log_path="${agent_dir}execution.log"

        # Skip disabled agents (unless we want to show them)
        if [[ "${enabled}" == "false" ]]; then
            continue
        fi

        # Check agent status
        status=$(check_agent_status "${agent_name}" "${log_path}")

        # Get last log timestamp
        last_timestamp=$(get_last_log_timestamp "${log_path}")

        # Calculate processing duration (for running agents)
        local processing_seconds=0
        if [[ "${status}" == "running" ]]; then
            processing_seconds=$(($(date +%s) - $(get_last_task_assigned_epoch "${log_path}")))
        fi

        # Count task assignments
        local task_count
        task_count=$(count_task_assignments "${log_path}")

        # Calculate average duration
        local avg_duration
        avg_duration=$(calculate_avg_duration "${log_path}")

        # Store agent data: name|status|processing_seconds|task_count|avg_duration|log_path|script_path
        agents_data+=("${agent_name}|${status}|${processing_seconds}|${task_count}|${avg_duration}|${log_path}|${script_path}")
    done

    # Output based on format
    if [[ "${format}" == "json" ]]; then
        output_json "${agents_data[@]}"
    else
        output_human "${agents_data[@]}"
    fi
}

# Output status in human-readable format
output_human() {
    local -a agents_data=("$@")
    local agent_name status processing_seconds task_count avg_duration log_path script_path

    info "Agent Status"
    info "============"
    info ""

    if [[ ${#agents_data[@]} -eq 0 ]]; then
        info "No enabled agents found."
        return 0
    fi

    # Column headers
    local name_header="Name"
    local status_header="Status"
    local processing_header="Processing In"
    local count_header="Task Count"
    local duration_header="Avg Duration"

    # Find column widths
    local max_name_len=${#name_header}
    local max_status_len=${#status_header}
    local max_processing_len=${#processing_header}
    local max_count_len=${#count_header}
    local max_duration_len=${#duration_header}

    for entry in "${agents_data[@]}"; do
        IFS='|' read -r agent_name status processing_seconds task_count avg_duration log_path script_path <<< "${entry}"

        # Update name column width
        if [[ ${#agent_name} -gt ${max_name_len} ]]; then
            max_name_len=${#agent_name}
        fi

        # Update status column width
        local status_len
        case "${status}" in
            running)
                status_len=7  # RUNNING length
                ;;
            idle)
                status_len=4  # IDLE length
                ;;
            *)
                status_len=7  # UNKNOWN length
                ;;
        esac
        if [[ ${status_len} -gt ${max_status_len} ]]; then
            max_status_len=${status_len}
        fi

        # Update processing column width
        local processing_str
        if [[ "${status}" == "running" ]]; then
            processing_str=$(format_duration "${processing_seconds}")
        else
            processing_str="IDLE"
        fi
        if [[ ${#processing_str} -gt ${max_processing_len} ]]; then
            max_processing_len=${#processing_str}
        fi

        # Update count column width
        local count_str="${task_count}"
        if [[ ${#count_str} -gt ${max_count_len} ]]; then
            max_count_len=${#count_str}
        fi

        # Update duration column width
        local duration_str="${avg_duration}m"
        if [[ ${#duration_str} -gt ${max_duration_len} ]]; then
            max_duration_len=${#duration_str}
        fi
    done

    # Print header row
    printf "  %-$(printf '%d' ${max_name_len})s  " "${name_header}"
    printf "%-$(printf '%d' ${max_status_len})s  " "${status_header}"
    printf "%-$(printf '%d' ${max_processing_len})s  " "${processing_header}"
    printf "%-$(printf '%d' ${max_count_len})s  " "${count_header}"
    printf "%s\n" "${duration_header}"

    # Print separator row
    printf "  "
    printf "%-${max_name_len}s  " "$(printf '%0.s=' $(seq 1 ${max_name_len}) | tr '=' '-')"
    printf "%-${max_status_len}s  " "$(printf '%0.s=' $(seq 1 ${max_status_len}) | tr '=' '-')"
    printf "%-${max_processing_len}s  " "$(printf '%0.s=' $(seq 1 ${max_processing_len}) | tr '=' '-')"
    printf "%-${max_count_len}s  " "$(printf '%0.s=' $(seq 1 ${max_count_len}) | tr '=' '-')"
    printf "%s\n" "$(printf '%0.s=' $(seq 1 ${max_duration_len}) | tr '=' '-')"

    # Display each agent
    for entry in "${agents_data[@]}"; do
        IFS='|' read -r agent_name status processing_seconds task_count avg_duration log_path script_path <<< "${entry}"

        # Format status
        local status_str
        case "${status}" in
            running)
                status_str="RUNNING"
                ;;
            idle)
                status_str="IDLE"
                ;;
            *)
                status_str="UNKNOWN"
                ;;
        esac

        # Format processing duration
        local processing_str
        if [[ "${status}" == "running" ]]; then
            processing_str=$(format_duration "${processing_seconds}")
        else
            processing_str="IDLE"
        fi

        # Format duration with 'm' suffix
        local duration_str="${avg_duration}m"

        # Print row
        printf "  %-$(printf '%d' ${max_name_len})s  " "${agent_name}"
        printf "%-$(printf '%d' ${max_status_len})s  " "${status_str}"
        printf "%-$(printf '%d' ${max_processing_len})s  " "${processing_str}"
        printf "%-$(printf '%d' ${max_count_len})s  " "${task_count}"
        printf "%s\n" "${duration_str}"
    done
}

# Output status in JSON format
output_json() {
    local -a agents_data=("$@")
    local agent_name status processing_seconds task_count avg_duration log_path script_path
    local first=true

    echo "{"
    echo '  "agents": ['

    for entry in "${agents_data[@]}"; do
        IFS='|' read -r agent_name status processing_seconds task_count avg_duration log_path script_path <<< "${entry}"

        if [[ "${first}" == "true" ]]; then
            first=false
        else
            echo ","
        fi

        # Format processing duration as string
        local processing_str="IDLE"
        if [[ "${status}" == "running" ]]; then
            processing_str=$(format_duration "${processing_seconds}")
        fi

        # Output JSON object for this agent
        printf '    {\n'
        printf '      "name": "%s",\n' "${agent_name}"
        printf '      "status": "%s",\n' "${status}"
        printf '      "processing_in": "%s",\n' "${processing_str}"
        printf '      "task_count": %s,\n' "${task_count}"
        printf '      "avg_duration": %s,\n' "${avg_duration}"
        printf '      "log_path": "%s",\n' "${log_path}"
        printf '      "script_path": "%s"\n' "${script_path}"
        printf '    }'
    done

    if [[ ${#agents_data[@]} -gt 0 ]]; then
        echo ""
    fi
    echo "  ]"
    echo "}"
}

# Show help message
show_help() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Agent Status Checker

Options:
    --help, -h      Show this help message
    --json, -j      Output status as JSON
    --list, -l      Alias for --json (for compatibility)

Description:
    This script checks the status of agents by examining their execution logs.
    An agent is considered:
      - "running" if the last log line is NOT "Task processing complete"
      - "idle" if the last log line IS "Task processing complete" or no log exists

    If an agent has never run (no log file), it is considered idle.

    Human-readable output shows:
      - Name: Agent name
      - Status: RUNNING or IDLE
      - Processing In: Duration since last task assignment (for running agents)
      - Task Count: Number of tasks processed
      - Avg Duration: Average processing time per task

    JSON output includes:
      - processing_in: Formatted processing duration string
      - task_count: Total tasks processed
      - avg_duration: Average processing time (integer minutes)

Examples:
    $(basename "$0")              Show human-readable status
    $(basename "$0") --json       Output status as JSON
    $(basename "$0") --help       Show this help message

Dependencies:
    - bash 4.0+
    - Standard Unix tools: grep, sed, awk, tail

Project: Maestro
EOF
}

# =============================================================================
# Main Logic
# =============================================================================

# Parse command line arguments
FORMAT="human"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --help|-h)
            show_help
            exit 0
            ;;
        --json|-j)
            FORMAT="json"
            shift
            ;;
        --list|-l)
            FORMAT="json"
            shift
            ;;
        *)
            error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Check agents directory exists
if [[ ! -d "${AGENTS_DIR}" ]]; then
    error "Agents directory not found: ${AGENTS_DIR}"
    exit 1
fi

# Discover and display agent statuses
discover_agents "${FORMAT}"

exit 0
