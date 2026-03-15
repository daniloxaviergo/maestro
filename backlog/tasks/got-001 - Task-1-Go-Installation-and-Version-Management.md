---
id: GOT-001
title: 'Task 1: Go Installation and Version Management'
status: To Do
assignee: []
created_date: '2026-03-15 00:12'
updated_date: '2026-03-15 00:19'
labels:
  - go
  - installation
  - automation
dependencies: []
references:
  - backlog/docs/doc-001 - PRD-Go-Development-Environment-Setup.md
priority: high
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement Go installation using gvm with version switching support
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Script installs Go without requiring manual intervention
- [ ] #2 go version command works after installation
- [ ] #3 Multiple Go versions can be installed and switched between
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task implements Go installation using **gvm** (Go Version Manager), a tool for managing multiple Go versions. The script will:

- **gvm Installation**: Download and install gvm from its GitHub repository using the official installation script or direct git clone
- **Go Installation**: Use `gvm install` to install the specified Go version (default: latest stable)
- **Version Switching**: Configure gvm to allow users to switch between installed versions using `gvm use`
- **PATH Management**: Update shell profile files (`.bashrc`, `.zshrc`) to source gvm on shell startup

**Why gvm**: 
- It's the established Go version manager (similar to rbenv/rvm for Ruby)
- Supports installing multiple versions and switching between them
- Handles binary path management automatically
- Well-maintained and widely used in the Go community

**Trade-offs**:
- Requires shell profile modification (sources gvm script on startup)
- Users need to start a new shell or run `source ~/.gvm/scripts/gvm` for changes to take effect
- Alternative: `goenv` or manual tarball installation, but gvm is the most established for Go

### 2. Files to Modify

**New Files to Create:**
- `scripts/install-go.sh` - Main script for Go installation and version management
- `scripts/gvm-setup.sh` - Helper script to source gvm (may be embedded in main script)
- `docs/setup-go.md` - Documentation for Go installation and version management

**Shell Profile Modifications (runtime):**
- `~/.bashrc` or `~/.zshrc` - gvm sourcing line appended (only if shell is bash/zsh)

**No Files to Delete or Modify:**
- Existing project files remain unchanged
- No system-wide modifications (gvm installs to `~/.gvm`)

### 3. Dependencies

**Prerequisites for Script Execution:**
- Linux operating system (Ubuntu/Debian or CentOS/RHEL-based)
- `curl` or `wget` for downloading files (HTTPS)
- `git` for gvm installation (or fallback to tarball download)
- Package manager (`apt`, `yum`, or `dnf`) for system dependencies
- User with no root access required (installs to home directory)

**Existing Tasks/Dependencies:**
- Task 6 (Setup Script and Documentation): Task 1 should be implemented as part of the main setup script
- No other tasks blocking this implementation

**System Dependencies to Install:**
- Build tools: `build-essential` (Debian/Ubuntu) or `gcc`, `make`, `autoconf` (RHEL/CentOS)
- OpenSSL development libraries for Go's crypto packages
- `bzip2` for compression utilities

### 4. Code Patterns

**Shell Script Conventions:**
```bash
#!/bin/bash
set -euo pipefail  # Strict error handling

# Use absolute paths for all file operations
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Log functions for visibility
log_info() { echo "[INFO] $*"; }
log_error() { echo "[ERROR] $* >&2; }

# Check for prerequisites before proceeding
if ! command -v curl &> /dev/null; then
    log_error "curl is required but not installed"
    exit 1
fi
```

**gvm Integration Patterns:**
```bash
# Source gvm in the current shell
source ~/.gvm/scripts/gvm

# Install a specific Go version
gvm install go1.21.0

# Switch to a specific version
gvm use go1.21.0 --default

# List installed versions
gvm list
```

**Error Handling:**
- All commands use `set -euo pipefail` for strict error handling
- Each major step includes success/failure logging
- Cleanup on failure (remove partially downloaded files)

### 5. Testing Strategy

**Unit Tests (using `bats` or shell assertions):**
- Verify gvm is installed successfully (check `~/.gvm/bin/gvm` exists)
- Verify Go installation works (`go version` returns expected version)
- Verify multiple versions can be installed
- Verify version switching works (`gvm use <version>`)

**Manual Testing Checklist:**
```bash
# Test 1: Fresh installation
./scripts/install-go.sh
source ~/.bashrc
go version  # Should show installed version

# Test 2: Multiple versions
./scripts/install-go.sh --version go1.20.14
gvm list  # Should show both versions
gvm use go1.20.14
go version  # Should show go1.20.14

# Test 3: PATH verification
which go  # Should point to ~/.gvm/pkgsets/...
```

**Edge Cases to Cover:**
- gvm already installed (graceful skip or upgrade)
- Go version already installed (skip or force reinstall)
- Network failure during download (retry logic)
- Insufficient disk space (check before downloading)
- Shell not bash/zsh (warning + manual configuration needed)

### 6. Risks and Considerations

**Blocking Issues:**
- None identified. This is a well-understood process.

**Potential Pitfalls:**
- **Shell configuration**: Users may have non-standard shells (fish, zsh with custom config) that require manual gvm sourcing
- **Existing Go installation**: Script may conflict with system-installed Go; should detect and warn
- **Permissions**: While gvm installs to home directory, some system dependencies may require sudo (document clearly)

**Trade-offs:**
- **Approach**: Using official gvm installation method (git clone + script) vs. alternative tools like `goenv` or manual tarball extraction
- **Default version**: Should the script install the latest stable, or allow user-specified version? (Support both via flag)
- **Shell modification**: Modifying `.bashrc`/`.zshrc` may conflict with user's custom configs; should include backup option

**Deployment Considerations:**
- Script is idempotent (can be run multiple times safely)
- No restart required (new shell needed to pick up PATH changes)
- No system-wide changes (does not modify `/usr/local`)
- Cleanup path: Users can remove `~/.gvm` to uninstall
<!-- SECTION:PLAN:END -->
