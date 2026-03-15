---
id: GOT-003
title: 'Task 3: Go Development Tools Installation'
status: To Do
assignee: []
created_date: '2026-03-15 00:12'
updated_date: '2026-03-15 18:58'
labels:
  - tools
  - linters
  - debugger
  - gopls
dependencies: []
references:
  - backlog/docs/doc-001 - PRD-Go-Development-Environment-Setup.md
priority: high
ordinal: 12750
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Install essential Go development tools (linters, formatters, debuggers)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 gopls (Go language server) is installed and configured
- [ ] #2 goimports is installed for import management
- [ ] #3 A linter (golint or golangci-lint) is installed
- [ ] #4 Delve (dlv) debugger is installed
- [ ] #5 All tools work with the installed Go version
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task installs essential Go development tools using `go install` with specific versions pinned for stability. The approach:

- Use `go install` with `@version` syntax to install tools at pinned versions
- Install tools to `$HOME/go/bin` which should already be in PATH from Task 1
- Use versions compatible with Go 1.21+ (current standard)
- Install: gopls, goimports, golangci-lint, and delve (dlv)
- Verify each tool installs and runs correctly
- Record installed versions for reproducibility

**Version choices:**
- gopls v0.15.3 - stable, compatible with Go 1.21+
- goimports v0.25.0 - actively maintained import manager
- golangci-lint v1.63.4 - popular linter with many checks
- dlv v1.23.1 - debugger with good Go 1.21+ support

**Alternative considered:** Using the older golint package, but golangci-lint provides more comprehensive linting and golint is deprecated.

### 2. Files to Modify

**New Files to Create:**
- `scripts/install-go-tools.sh` - Main script to install all Go development tools

**Existing Files to Reference:**
- `backlog/tasks/got-001` - For Go installation location and PATH setup
- `backlog/tasks/got-007` - For gopls version compatibility and vim integration notes

**Files to Create for Record:**
- `docs/go-tools.md` - Documentation for installed tools and their usage
- `scripts/verify-go-tools.sh` - Helper script to verify tool installation

### 3. Dependencies

**Prerequisites:**
- Go 1.21+ must be installed (Task 1 completion required)
- `$HOME/go/bin` must be in user's PATH
- Network access for downloading Go modules

**Blocking Issues:**
- If Task 1 is not completed first, tools may not install to correct location
- If PATH is not configured, tools will not be accessible after installation

**External Requirements:**
- Internet access to fetch Go modules
- No root/sudo required (all tools install to user home directory)

### 4. Code Patterns

**Shell Script Conventions (follow Task 1 pattern):**
```bash
#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_PREFIX="[GO-TOOLS-INSTALL]"

log_info() { echo "$LOG_PREFIX [INFO] $*"; }
log_error() { echo "$LOG_PREFIX [ERROR] $*" >&2; }

# Verify Go is available before proceeding
if ! command -v go &> /dev/null; then
    log_error "Go is not installed or not in PATH"
    exit 1
fi

# Verify go version meets minimum
GO_VERSION=$(go version | sed -E 's/.*go version go([0-9]+\.[0-9]+).*/\1/')
if [[ "$GO_VERSION" != "1.21" && "$GO_VERSION" != "1.22" && "$GO_VERSION" != "1.23" ]]; then
    log_warn "Go version $GO_VERSION may not be fully compatible"
fi
```

**Tool Installation Pattern:**
```bash
# Install with pinned version
go install golang.org/x/tools/gopls@v0.15.3
```

**Error Handling:**
- Each tool installation should be wrapped in try/like error handling
- Failures should not prevent other tools from installing
- Provide clear error messages indicating which tool failed

### 5. Testing Strategy

**Verification Steps:**
1. Run `go install` for each tool
2. Verify binary exists in `$HOME/go/bin`
3. Run each tool with `--version` or `-version` flag
4. For gopls, run `gopls version` to verify
5. For dlv, run `dlv version` to verify
6. For golangci-lint, run `golangci-lint version` to verify

**Edge Cases to Cover:**
- Tools already installed at different version (should upgrade)
- Network failure during installation (should handle gracefully)
- Partial installation (some tools succeed, some fail)
- Go modules disabled (should fail with clear error)

**Verification Script:**
- Create `scripts/verify-go-tools.sh` to check all tools
- Report installed versions
- Indicate which tools are missing or broken

### 6. Risks and Considerations

**Known Risks:**
- `golangci-lint` may have build issues on some systems (requires CGO for some linters)
- `dlv` may require debugging symbols on some distributions
- Tool versions may become outdated quickly; document pinned versions clearly

**Trade-offs:**
- Using pinned versions ensures reproducibility but requires manual updates
- golangci-lint vs golint: golangci-lint is more comprehensive but heavier

**Deployment Considerations:**
- Should be part of main setup script (Task 6) or run as standalone
- User should be informed to reload PATH after installation
- Consider adding tool version check to verification script

**Follow-ups:**
- Consider creating a tool version management file (like go.mod for tools)
- Document recommended tool versions in `docs/go-tools.md`
- Consider adding tool update command to setup script
<!-- SECTION:PLAN:END -->
