#!/usr/bin/env bash

set -euo pipefail

PATTERN="/usr/bin/node.*qwen"

# Get matching processes: PID + full command for UX
processes=$(ps -eo pid,command | \
    awk -v pat="$PATTERN" 'NR>1 && $0 ~ pat && !/grep/' | \
    fzf \
        --height=100% \
        --reverse \
        --info=inline \
        --prompt="Select process to kill (Ctrl-C to abort): " \
        --preview='echo {} | cut -f1 | xargs -r ps -p {} -o pid,comm --no-headers 2>/dev/null || echo "PID {} not found"' \
        --preview-window=up:3:hidden \
        --multi \
        --bind 'ctrl-a:select-all' \
        --header='Tab/Shift-Tab to select multiple; Enter to confirm kill')

# Exit on empty selection
if [[ -z "$processes" ]]; then
    echo "No selection made."
    exit 0
fi

# Parse PIDs: trim, split by newline
readarray -t selected_lines <<< "$processes"
pids=()
for line in "${selected_lines[@]}"; do
    # Extract first field (PID)
    pid=$(echo "$line" | awk '{print $1}')
    if [[ "$pid" =~ ^[0-9]+$ ]]; then
        pids+=("$pid")
    else
        echo >&2 "Invalid PID: '$pid' from line: '$line'"
        exit 1
    fi
done

count=${#pids[@]}
echo "Selected $count process(es):"
printf '  PID %s\n' "${pids[@]}"

read -rp "Kill these process(es)? [y/N] " confirm
if [[ "${confirm,,}" != "y" && "${confirm,,}" != "yes" ]]; then
    echo "Cancelled."
    exit 0
fi

# Kill selectively
errors=0
for pid in "${pids[@]}"; do
    if kill -9 "$pid" 2>/dev/null; then
        echo "✓ Sent SIGKILL to $pid"
    else
        echo "⚠ Failed to kill $pid" >&2
        ((errors++)) || true  # ← prevents (()) from returning 1 and triggering set -e
    fi
done

if [[ $errors -gt 0 ]]; then
    echo "⚠ $errors kill(s) failed (process may have exited already)."
fi

echo "Done."
