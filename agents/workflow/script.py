#!/usr/bin/env python3
"""
Workflow orchestrator script for Backlog.md task management.
Handles sequential agent assignment with state tracking.
Each agent gets their own entry in the task state, marking completion
when they finish and starting a new entry for the next agent.
"""

import os
import re
import subprocess
import sys
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path


# Paths
SCRIPT_DIR = Path(__file__).parent
CONFIG_PATH = SCRIPT_DIR / "config.yml"
STATE_PATH = SCRIPT_DIR / "tasks.yml"


# ============================================================================
# Logging Functions
# ============================================================================

def log(level: str, message: str) -> None:
    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    print(f"[{timestamp}] {level}: {message}")


def error(msg: str) -> None:
    log("ERROR", msg)


def info(msg: str) -> None:
    log("INFO", msg)


# ============================================================================
# Configuration Parsing
# ============================================================================

def extract_yaml_value(file_path: Path, key: str) -> str | None:
    """Extract a simple key: value from YAML file."""
    try:
        with open(file_path) as f:
            for line in f:
                if line.startswith(f"{key}:"):
                    value = line.split(":", 1)[1].strip()
                    return value.strip('"').strip("'")
    except FileNotFoundError:
        error(f"Config file not found: {file_path}")
        return None
    return None


def get_agents(config_path: Path) -> list[str] | None:
    """Parse comma-separated agents list from config."""
    agents_str = extract_yaml_value(config_path, "agents")
    if agents_str is None:
        return None
    # Handle both "a, b" and "a,b" formats
    return [a.strip() for a in agents_str.split(",") if a.strip()]


def get_backlog_command(config_path: Path) -> str | None:
    """Get backlog command from config."""
    return extract_yaml_value(config_path, "backlog_command")


# ============================================================================
# Task State Management
# ============================================================================

@dataclass
class AgentEntry:
    """Represents one agent's entry in the task state."""
    status: str
    assigned_agent: str
    assigned_at: str
    completed_at: str = ""


def load_task_state(state_file: Path, task_id: str) -> list[AgentEntry]:
    """Load task state from YAML state file as a list of agent entries."""
    entries: list[AgentEntry] = []

    if not state_file.exists():
        return entries

    try:
        with open(state_file) as f:
            content = f.read()
    except OSError as e:
        error(f"Failed to read state file: {e}")
        return entries

    # Find task block using regex - find content under task_id:
    pattern = re.compile(rf"^{re.escape(task_id)}:(.*?)(?=\n\S|\Z)", re.DOTALL | re.MULTILINE)
    match = pattern.search(content)
    if not match:
        return entries

    task_block = match.group(1)

    # Parse each agent entry (lines starting with "  - " are new entries)
    lines = task_block.split("\n")
    current_entry: dict[str, str] = {}

    for line in lines:
        line = line.rstrip()
        if not line:
            continue

        # Check if this is a new entry (starts with "  - " at exactly 2-space indent)
        stripped = line.lstrip()
        indent = len(line) - len(stripped)
        
        if stripped.startswith("- ") and indent == 2:
            # Save previous entry if exists
            if current_entry:
                entries.append(AgentEntry(
                    status=current_entry.get("status", "pending"),
                    assigned_agent=current_entry.get("assigned_agent", ""),
                    assigned_at=current_entry.get("assigned_at", ""),
                    completed_at=current_entry.get("completed_at", ""),
                ))
            # Start new entry - parse "- status: in_progress" format
            rest = stripped[2:].strip()  # Remove "- "
            if ":" in rest:
                key, _, value = rest.partition(":")
                key = key.strip()
                value = value.strip().strip('"').strip("'")
                current_entry = {key: value}
            else:
                current_entry = {"status": "pending"}

        # Extract field values (4-space indented under the list item)
        elif stripped.startswith("status:") or stripped.startswith("assigned_agent:") or \
             stripped.startswith("assigned_at:") or stripped.startswith("completed_at:"):
            key, _, value = stripped.partition(":")
            key = key.strip()
            value = value.strip().strip('"').strip("'")
            current_entry[key] = value

    # Don't forget the last entry
    if current_entry:
        entries.append(AgentEntry(
            status=current_entry.get("status", "pending"),
            assigned_agent=current_entry.get("assigned_agent", ""),
            assigned_at=current_entry.get("assigned_at", ""),
            completed_at=current_entry.get("completed_at", ""),
        ))

    return entries


def update_task_state(
    state_file: Path,
    task_id: str,
    status: str,
    assigned_agent: str,
    assigned_at: str,
    completed_at: str = "",
) -> None:
    """Update task state by adding a new agent entry at the end of the list."""
    # Initialize state file if needed
    if not state_file.exists():
        _create_state_file(state_file)

    try:
        with open(state_file) as f:
            content = f.read()
    except OSError as e:
        error(f"Failed to read state file: {e}")
        return

    # Check if task exists
    pattern = re.compile(rf"^{re.escape(task_id)}:", re.MULTILINE)
    match = pattern.search(content)

    if not match:
        # Append new task entry
        with open(state_file, "a") as f:
            f.write(f"{task_id}:\n")
            f.write(f"  - status: {status}\n")
            f.write(f"    assigned_agent: {assigned_agent}\n")
            f.write(f"    assigned_at: \"{assigned_at}\"\n")
            if completed_at:
                f.write(f"    completed_at: \"{completed_at}\"\n")
            f.write("\n")
    else:
        # Read the entire file and find the task block
        lines = content.split("\n")
        
        # Find task header position
        task_header_idx = -1
        for i, line in enumerate(lines):
            if line.startswith(f"{task_id}:"):
                task_header_idx = i
                break

        if task_header_idx < 0:
            return

        # Find where the task block ends (next task or end of file)
        task_end_idx = len(lines)
        for i in range(task_header_idx + 1, len(lines)):
            # If line starts at beginning with a word followed by :, it's a new task
            if lines[i].strip() and re.match(r'^\w+-\d+:', lines[i].strip()):
                task_end_idx = i
                break

        # Extract lines before task block (everything before the task header)
        before_task = lines[:task_header_idx]

        # Parse existing list items within task block
        existing_entries: list[dict[str, str]] = []
        current_entry: dict[str, str] = {}

        for i in range(task_header_idx + 1, task_end_idx):
            line = lines[i]
            stripped = line.strip()

            if not stripped:
                continue

            # Check if this is a new list item
            if line.startswith("  - "):
                # Save previous entry if exists
                if current_entry:
                    existing_entries.append(current_entry)
                # Start new entry - parse "- status: in_progress" format
                rest = line[4:].strip()  # Remove "  - " prefix
                if ":" in rest:
                    key, _, value = rest.partition(":")
                    current_entry = {key.strip(): value.strip()}
                else:
                    current_entry = {"status": rest}
            elif line.startswith("    ") and ":" in stripped:
                # This is a field under the list item
                key, _, value = stripped.partition(":")
                current_entry[key.strip()] = value.strip().strip('"').strip("'")

        # Don't forget the last entry
        if current_entry:
            existing_entries.append(current_entry)

        # Build new task block content
        new_task_lines: list[str] = [f"{task_id}:"]
        
        # Add existing entries first
        for entry in existing_entries:
            new_task_lines.append("  - status: " + entry.get("status", "pending"))
            if entry.get("assigned_agent"):
                new_task_lines.append("    assigned_agent: " + entry.get("assigned_agent", ""))
            if entry.get("assigned_at"):
                new_task_lines.append("    assigned_at: \"" + entry.get("assigned_at", "") + "\"")
            if entry.get("completed_at"):
                new_task_lines.append("    completed_at: \"" + entry.get("completed_at", "") + "\"")
            new_task_lines.append("")  # Empty line between entries

        # Add new entry
        new_task_lines.append("  - status: " + status)
        new_task_lines.append("    assigned_agent: " + assigned_agent)
        new_task_lines.append("    assigned_at: \"" + assigned_at + "\"")
        if completed_at:
            new_task_lines.append("    completed_at: \"" + completed_at + "\"")
        new_task_lines.append("")  # Trailing newline

        # Reconstruct file - combine before_task + new_task_lines
        new_lines = before_task + new_task_lines

        with open(state_file, "w") as f:
            f.write("\n".join(new_lines))


def _create_state_file(state_file: Path) -> None:
    """Create initial state file with header."""
    state_file.parent.mkdir(parents=True, exist_ok=True)
    header = [
        "# Task state tracking for workflow agent\n",
        "# Format: {task_id}:\n",
        "#   - status: pending|in_progress|finished\n",
        "#     assigned_agent: agent_name\n",
        "#     assigned_at: \"YYYY-MM-DD HH:MM:SS\"\n",
        "#     completed_at: \"YYYY-MM-DD HH:MM:SS\"\n",
        "\n",
    ]
    with open(state_file, "w") as f:
        f.writelines(header)


def get_current_state(entries: list[AgentEntry]) -> AgentEntry:
    """Get the current (last) agent entry from the state."""
    if not entries:
        return AgentEntry(status="pending", assigned_agent="", assigned_at="")
    return entries[-1]


def get_completed_agents_count(entries: list[AgentEntry]) -> int:
    """Count completed agents (entries with completed_at set)."""
    return sum(1 for e in entries if e.completed_at)


def assign_task(backlog_cmd: str, task_id: str, agent: str) -> bool:
    """Assign task to agent via backlog CLI."""
    info(f"Assigning task {task_id} to {agent}")

    command = backlog_cmd.split(" ")
    command.append(task_id)
    command.append("--assignee")
    command.append(agent)

    try:
        result = subprocess.run(
            command,
            capture_output=True,
            text=True,
        )
    except FileNotFoundError:
        error(f"Backlog command not found: {backlog_cmd}")
        return False

    if result.returncode != 0:
        error(f"Failed to assign task {task_id} to {agent}")
        if result.stderr:
            error(f"stderr: {result.stderr}")
        return False

    info(f"Task {task_id} assigned to {agent} successfully")
    return True


def get_backlog_assignee(backlog_cmd: str, task_id: str) -> str | None:
    """Get the current assignee from backlog task."""
    # Use the task ID directly with --plain to get output
    command = ["backlog", "task", task_id, "--plain"]

    try:
        result = subprocess.run(
            command,
            capture_output=True,
            text=True,
        )
    except FileNotFoundError:
        error(f"Backlog command not found")
        return None

    if result.returncode != 0:
        error(f"Failed to get task {task_id} info")
        if result.stderr:
            error(f"stderr: {result.stderr}")
        return None

    # Parse assignee from output (look for "Assignee: @agent" line)
    for line in result.stdout.split("\n"):
        if "Assignee:" in line:
            parts = line.split(":", 1)
            if len(parts) > 1:
                assignee = parts[1].strip()
                # Return agent name without @ prefix
                if assignee.startswith("@"):
                    return assignee[1:]
                if assignee and assignee != "none":
                    return assignee

    return None


def find_project_path(task_file: Path) -> Path:
    """Find project root by walking up from task file."""
    current = task_file.parent
    while current != current.parent:
        tasks_dir = current / "backlog" / "tasks"
        config_file = current / "backlog" / "config.yml"
        if tasks_dir.is_dir() and config_file.is_file():
            return current
        current = current.parent
    return task_file.parent


# ============================================================================
# Main Workflow Logic
# ============================================================================

def run_workflow(task_id: str) -> int:
    """Execute the workflow orchestrator logic."""
    # Validate configuration
    if not CONFIG_PATH.exists():
        error(f"Config file not found: {CONFIG_PATH}")
        return 1

    agents = get_agents(CONFIG_PATH)
    if agents is None:
        return 1

    backlog_cmd = get_backlog_command(CONFIG_PATH)
    if backlog_cmd is None:
        return 1

    if not agents:
        error("No agents defined in config")
        return 1

    info(f"Workflow agents: {', '.join(agents)}")
    info(f"Backlog command: {backlog_cmd}")
    info(f"Task ID: {task_id}")

    # Load current task state
    entries = load_task_state(STATE_PATH, task_id)
    current = get_current_state(entries)
    info(f"Current state: status={current.status}, assigned_agent={current.assigned_agent}")

    # Determine next agent based on completed count
    completed_count = get_completed_agents_count(entries)
    if completed_count >= len(agents):
        error(f"No more agents available (completed: {completed_count}, agents: {len(agents)})")
        return 1

    next_agent = agents[completed_count]
    info(f"Next agent: {next_agent} (completed count: {completed_count})")

    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Process based on current status
    if current.status == "pending":
        info(f"Task is pending, assigning to {next_agent}")
        if not assign_task(backlog_cmd, task_id, next_agent):
            return 1
        update_task_state(STATE_PATH, task_id, "in_progress", next_agent, timestamp)
        info(f"State updated: status=in_progress, assigned_agent={next_agent}, assigned_at={timestamp}")

    elif current.status == "in_progress":
        # Check the actual backlog assignee to detect if current agent completed
        backlog_assignee = get_backlog_assignee(backlog_cmd, task_id)
        info(f"Backlog assignee: {backlog_assignee}")

        if backlog_assignee == "workflow":
            # Workflow is assigned, meaning the previous agent completed their work
            info(f"Agent {current.assigned_agent} completed work")
            
            # Update existing entries: mark current as finished, add new entry for next agent
            update_task_state(STATE_PATH, task_id, "finished", current.assigned_agent, current.assigned_at, timestamp)
            info(f"State updated: status=finished, assigned_agent={current.assigned_agent}, completed_at={timestamp}")

            # Next agent is the one after the completed agent
            next_completed_count = completed_count + 1
            if next_completed_count >= len(agents):
                error(f"No more agents available (completed: {next_completed_count}, agents: {len(agents)})")
                return 1
            next_agent = agents[next_completed_count]
            info(f"Assigning to next agent: {next_agent}")
            if not assign_task(backlog_cmd, task_id, next_agent):
                return 1
            # Add new entry for next agent
            update_task_state(STATE_PATH, task_id, "in_progress", next_agent, timestamp)
            info(f"State updated: status=in_progress, assigned_agent={next_agent}, assigned_at={timestamp}")
        elif current.assigned_agent != next_agent:
            # Current agent completed (assigned to someone else), assign next agent
            info(f"Agent {current.assigned_agent} completed work")
            # Mark current as completed
            update_task_state(STATE_PATH, task_id, "finished", current.assigned_agent, current.assigned_at, timestamp)
            info(f"State updated: status=finished, assigned_agent={current.assigned_agent}, completed_at={timestamp}")

            info(f"Assigning to next agent: {next_agent}")
            if not assign_task(backlog_cmd, task_id, next_agent):
                return 1
            update_task_state(STATE_PATH, task_id, "in_progress", next_agent, timestamp)
            info(f"State updated: status=in_progress, assigned_agent={next_agent}, assigned_at={timestamp}")
        else:
            # Same agent, already assigned - nothing to do
            info(f"Agent {next_agent} already assigned")

    elif current.status == "finished":
        info(f"Task {task_id} is already finished (all agents completed)")
        return 0

    else:
        error(f"Unknown status: {current.status}")
        return 1

    info("Workflow orchestration complete")
    return 0


def main() -> int:
    """Main entry point."""
    if len(sys.argv) != 2:
        print(f"Usage: {sys.argv[0]} <task_file>", file=sys.stderr)
        return 1

    task_file = Path(sys.argv[1])
    task_id = task_file.name

    # Find and change to project directory
    project_path = find_project_path(task_file)
    info(f"Project path: {project_path}")
    os.chdir(project_path)

    # Update paths to be relative to project root
    global CONFIG_PATH, STATE_PATH
    CONFIG_PATH = SCRIPT_DIR / "config.yml"
    STATE_PATH = SCRIPT_DIR / "tasks.yml"

    return run_workflow(task_id)


if __name__ == "__main__":
    sys.exit(main())
