#!/usr/bin/env bash
#
# scripts/attach.sh - Tmux session attacher with fzf selection
#
# This script scans the agents/ directory for configured agents and presents
# an interactive fzf menu to attach to any agent's tmux session.
#
# Usage:
#   ./scripts/attach.sh          # Attach to an agent session via fzf
#   ./scripts/attach.sh --list   # List all available agent sessions
#   ./scripts/attach.sh --help   # Show this help message
#
# Dependencies:
#   - bash 4.0+
#   - fzf 0.20+
#   - tmux 2.0+
#   - grep, sed, awk (standard Unix tools)
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

# Check if a required command is installed
check_dependency() {
    local cmd="$1"
    if ! command -v "$cmd" &>/dev/null; then
        error "${cmd} is not installed. Please install ${cmd} to continue."
        exit 1
    fi
}

# Check if the agents directory exists
check_agents_dir() {
    if [[ ! -d "${AGENTS_DIR}" ]]; then
        error "Agents directory not found: ${AGENTS_DIR}"
        error "Please create an agents directory with agent configurations."
        exit 1
    fi
}

# Extract tmux_session value from an agent's config.yml
# Uses grep/sed for simple YAML parsing (no external dependencies)
extract_tmux_session() {
    local config_file="$1"
    local session_name

    # Check if config file exists
    if [[ ! -f "${config_file}" ]]; then
        return 1
    fi

    # Extract the tmux_session value using grep and sed
    # Handles: tmux_session: "session-name" or tmux_session: session-name
    session_name=$(grep -E '^\s*tmux_session:' "${config_file}" 2>/dev/null | head -n1 | sed -E 's/^[^:]*:\s*["'"'"']?([^"'"'"'\s]+)["'"'"']?.*$/\1/' || true)

    if [[ -z "${session_name}" ]]; then
        return 1
    fi

    echo "${session_name}"
}

# Discover all agents and their tmux sessions
# Returns: agent_name:session_name pairs
discover_agents() {
    local agent_dir config_file session_name

    # Check if agents directory exists (non-fatal for listing)
    if [[ ! -d "${AGENTS_DIR}" ]]; then
        return 0
    fi

    # Scan agent directories
    for agent_dir in "${AGENTS_DIR}"/*/; do
        # Skip if no agent directories found
        [[ -d "${agent_dir}" ]] || continue

        local agent_name
        agent_name=$(basename "${agent_dir}")
        config_file="${agent_dir}config.yml"

        # Extract session name from config
        if session_name=$(extract_tmux_session "${config_file}"); then
            echo "${agent_name}:${session_name}"
        else
            warning "Agent '${agent_name}' has no valid tmux_session in config.yml (skipping)"
        fi
    done
}

# Show help message
show_help() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Tmux Session Attacher with fzf Selection

Options:
    --help, -h      Show this help message
    --list, -l      List all available agent sessions without fzf
    --attach, -a    Attach to a session (default behavior)

Description:
    This script scans the agents/ directory for configured agents and
    presents an interactive fzf menu to attach to any agent's tmux session.

    If fzf is not installed or the session doesn't exist, appropriate
    error messages are displayed.

Examples:
    $(basename "$0")              Attach to an agent session via fzf
    $(basename "$0") --list       List all available agent sessions
    $(basename "$0") --help       Show this help message

Dependencies:
    - bash 4.0+
    - fzf 0.20+
    - tmux 2.0+

Project: Maestro
EOF
}

# =============================================================================
# Main Logic
# =============================================================================

# Parse command line arguments
LIST_MODE=false
ATTACH_MODE=true

while [[ $# -gt 0 ]]; do
    case "$1" in
        --help|-h)
            show_help
            exit 0
            ;;
        --list|-l)
            LIST_MODE=true
            ATTACH_MODE=false
            shift
            ;;
        --attach|-a)
            LIST_MODE=false
            ATTACH_MODE=true
            shift
            ;;
        *)
            error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Check required dependencies
check_dependency "fzf"
check_dependency "tmux"

# Check if agents directory exists
check_agents_dir

# Discover agents and their sessions
mapfile -t agents < <(discover_agents)

if [[ ${#agents[@]} -eq 0 ]]; then
    error "No agents found in ${AGENTS_DIR}"
    error "Please add agent configurations to the agents/ directory."
    exit 1
fi

# If list mode, display all agents and sessions
if [[ "${LIST_MODE}" == "true" ]]; then
    info "Available agent sessions:"
    info "------------------------"
    for agent_entry in "${agents[@]}"; do
        agent_name="${agent_entry%%:*}"
        session_name="${agent_entry##*:}"
        info "  ${agent_name} -> ${session_name}"
    done
    exit 0
fi

# Build fzf menu items: "agent_name:session_name"
IFS=$'\n' fzf_menu=("${agents[@]}")
unset IFS

# Use fzf to select an agent
# --height: 80% of terminal height
# --layout: reverse layout (list on top)
# --info: show search info
# --delimiter: use colon as field separator
# --with-nth: display only agent name
selected=$(printf '%s\n' "${fzf_menu[@]}" | fzf \
    --height=80% \
    --layout=reverse \
    --info=inline \
    --delimiter=':' \
    --with-nth=1 \
    --prompt='Select agent session > ' \
    --preview='echo "Agent: {1}\nSession: {2}"' \
    --preview-window='down:hidden:wrap' 2>/dev/null || true)

# Handle fzf cancellation
if [[ -z "${selected}" ]]; then
    info "No selection made. Exiting."
    exit 130
fi

# Extract session name from selection
session_name="${selected##*:}"

# Verify the session exists before attempting to attach
if ! tmux has-session -t "${session_name}" 2>/dev/null; then
    error "Tmux session '${session_name}' does not exist."
    error "Please start the session first with: tmux new-session -d -s ${session_name}"
    exit 1
fi

# Attach to the selected tmux session
info "Attaching to session: ${session_name}"
tmux attach -t "${session_name}"

exit 0
