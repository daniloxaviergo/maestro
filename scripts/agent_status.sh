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

        # Store agent data
        agents_data+=("${agent_name}|${status}|${last_timestamp}|${log_path}|${script_path}")
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
    local agent_name status last_timestamp log_path script_path

    info "Agent Status"
    info "============"
    info ""

    if [[ ${#agents_data[@]} -eq 0 ]]; then
        info "No enabled agents found."
        return 0
    fi

    # Find maximum agent name length for alignment
    local max_name_len=0
    for entry in "${agents_data[@]}"; do
        agent_name="${entry%%|*}"
        if [[ ${#agent_name} -gt ${max_name_len} ]]; then
            max_name_len=${#agent_name}
        fi
    done

    # Display each agent
    for entry in "${agents_data[@]}"; do
        IFS='|' read -r agent_name status last_timestamp log_path script_path <<< "${entry}"

        # Add status indicator
        local status_indicator
        case "${status}" in
            running)
                status_indicator="RUNNING"
                ;;
            idle)
                status_indicator="IDLE"
                ;;
            *)
                status_indicator="UNKNOWN"
                ;;
        esac

        # Format output with alignment
        printf "  %-$(printf '%d' ${max_name_len})s  [%s]" "${agent_name}" "${status_indicator}"

        # Add extra info if available
        if [[ -n "${last_timestamp}" ]]; then
            printf "  (last: %s)" "${last_timestamp}"
        fi
        if [[ -n "${script_path}" ]]; then
            printf "  script: %s" "${script_path}"
        fi
        printf "\n"
    done
}

# Output status in JSON format
output_json() {
    local -a agents_data=("$@")
    local agent_name status last_timestamp log_path script_path
    local first=true

    echo "{"
    echo '  "agents": ['

    for entry in "${agents_data[@]}"; do
        IFS='|' read -r agent_name status last_timestamp log_path script_path <<< "${entry}"

        if [[ "${first}" == "true" ]]; then
            first=false
        else
            echo ","
        fi

        # Output JSON object for this agent
        printf '    {\n'
        printf '      "name": "%s",\n' "${agent_name}"
        printf '      "status": "%s",\n' "${status}"
        printf '      "last_timestamp": "%s",\n' "${last_timestamp}"
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
