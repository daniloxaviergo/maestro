#!/bin/bash
# Wrapper script to run the Python workflow orchestrator

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
python3 "$SCRIPT_DIR/script.py" "$@"
