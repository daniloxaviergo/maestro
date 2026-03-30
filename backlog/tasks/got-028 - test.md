---
id: GOT-028
title: test
status: To Do
assignee:
  - workflow
created_date: '2026-03-16 11:14'
updated_date: '2026-03-30 18:20'
labels: []
dependencies: []
ordinal: 6250
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a new version of readme and save in ./new_readme.md
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The comprehensive README update will be structured to serve both new users and developers with clear information hierarchy and practical examples.

**Approach:**
- **Restructure existing content**: Consolidate information from current README.md, QWEN.md, AGENTS.md, and docs/ directory
- **Add new sections**: Project goals, quick-start guide, detailed configuration, and architecture overview
- **Improve organization**: Logical flow from installation to advanced usage
- **Maintain consistency**: Align with existing code patterns and project conventions

**Key Sections:**
1. Overview and features (expand current list)
2. Quick-start guide for new users
3. Detailed installation instructions
4. Usage examples (basic, tmux notifications, testing)
5. Configuration (agent setup, environment variables)
6. Architecture overview (component diagram description)
7. Development workflow (building, testing, code style)
8. Backlog.md integration (task management workflow)
9. Troubleshooting section
10. API reference (key packages)

**Trade-offs:**
- Keep README under ~500 lines for readability
- Move detailed architecture to QWEN.md (keep README high-level)
- Use code blocks for examples, tables for reference
- Link to docs/ for deep dives on specific topics

### 2. Files to Modify

| File | Action | Purpose |
|------|--------|--|
| `README.md` | **Rewrite** | Comprehensive project documentation |
| `QWEN.md` | **Update** | Remove redundancy, link to README for installation/usage |
| `docs/` | **Review** | Verify docs are referenced correctly and consistent |

**Files to Reference (read-only for planning):**
- `cmd/monitor/main.go` - CLI entry point with signal handling and component wiring
- `pkg/watcher/*.go` - fsnotify wrapper for file monitoring
- `pkg/parser/*.go` - YAML frontmatter extraction
- `pkg/cache/*.go` - File state caching with debouncing
- `pkg/change_detect/*.go` - Assignee change detection and logging
- `pkg/notifier/*.go` - Tmux notification and script execution
- `pkg/matcher/*.go` - Agent-assignee matching logic
- `pkg/agent/*.go` - Agent identity and configuration management
- `pkg/config/*.go` - Configuration loading
- `pkg/logs/*.go` - JSON logging
- `Makefile` - Build commands (build, run, tmux-*)
- `docs/agent-configuration.md` - Agent config details
- `docs/agent-orchestration-quickstart.md` - Agent setup guide
- `docs/setup-monitor.md` - Monitor setup details

### 3. Dependencies

**Prerequisites:**
- All features from tasks GOT-008 through GOT-029 must be implemented (all marked as Done in backlog)
- tmux installed for notifications and script execution (verified via `tmux --version`)
- Go 1.25.7+ for building (verified via `go version`)
- `./backlog/tasks` directory exists

**Build Verification:**
- `make build` succeeds without errors
- `go vet ./...` passes with no warnings

**No blocking tasks** - All core functionality is implemented and documented.

### 4. Code Patterns

**Conventions to Follow:**
1. **YAML code blocks**: Use triple backticks with `yaml` language identifier
2. **Bash commands**: Use triple backticks with `bash` language identifier
3. **Go code blocks**: Use triple backticks with `go` language identifier
4. **Tables**: Use Markdown tables for configuration options and reference
5. **Error handling**: Show `if err != nil { return err }` pattern in code examples
6. **Paths**: Use relative paths from project root (`./backlog/tasks`, not `/absolute/path`)
7. **Package names**: lowercase, short (e.g., `cache`, `watcher`, `parser`)
8. **Function names**: CamelCase (e.g., `NewWatcher`, `ProcessFile`)
9. **Variables**: camelCase (e.g., `fileWatcher`, `eventQueue`)

**Example format:**
```markdown
```yaml
# Configuration example
script_path: "./agents/my-agent/script.sh"
tmux_session: "my-agent"
enabled: true
```
``

### 5. Testing Strategy

**Verification approach:**
1. **Build test**: Run `make build` and verify binary exists at `bin/monitor`
2. **Run test**: Execute `make run` with sample task file, verify file events are output
3. **Documentation test**: Verify all code examples are syntactically correct
4. **Link check**: Ensure all internal links (docs/, backlog/) are valid
5. **Example test**: Run manual testing examples from README to verify they work

**No unit tests required** - README changes are documentation-only.

### 6. Risks and Considerations

**Known risks:**
1. **Information overload**: README may become too long - need to balance comprehensiveness with readability
2. **Documentation drift**: Changes to code may not update README - need maintenance process
3. **Link rot**: Internal links may break if file structure changes
4. **Example accuracy**: Code examples must match actual implementation
5. **Version mismatch**: README may reference outdated features or versions

**Mitigation strategies:**
- Use modular section structure for easy updates
- Link to detailed docs instead of duplicating information
- Include "Last updated" date or version information
- Create maintenance checklist for documentation updates
- Review against current codebase before finalizing plan

**Rollout considerations:**
- No code changes required - pure documentation update
- No deployment required - static file update
- No breaking changes for existing users
- No configuration changes required
- Backward compatible - existing users unaffected

### 7. Detailed README Structure

**Proposed section order:**
1. **Project Header**: Name, description, status badges (if any)
2. **Overview**: High-level description of what Maestro does
3. **Features**: Bullet list of key capabilities
4. **Technology Stack**: Language, core libraries, tools
5. **Quick Start**: 3-step setup for immediate use
6. **Installation**: Detailed installation instructions
7. **Usage**: Basic monitoring, tmux notifications, testing
8. **Configuration**: Agent setup, environment variables
9. **Architecture**: Component overview with data flow
10. **Development**: Building, testing, code style
11. **Backlog.md Integration**: Task management workflow
12. **Troubleshooting**: Common issues and solutions
13. **Contributing**: How to contribute (if applicable)
14. **License**: MIT license notice

**Sections to reference:**
- Link to `QWEN.md` for detailed architecture (component interactions, data structures)
- Link to `docs/` for deep dives on specific topics
- Link to `backlog/tasks/` for implementation details

### 8. Implementation Steps

**Phase 1: Content Audit**
1. Read current README.md, QWEN.md, AGENTS.md
2. Review all files in `docs/` directory
3. Identify overlapping content and gaps
4. Create content outline

**Phase 2: README Rewrite**
1. Rewrite Overview and Features sections
2. Add Quick Start section with working example
3. Expand Installation with prerequisites
4. Add Usage examples for different scenarios
5. Document Configuration options
6. Include Architecture overview
7. Add Development workflow
8. Document Backlog.md integration

**Phase 3: QWEN.md Update**
1. Add link to README for installation/usage
2. Remove redundant installation sections
3. Keep architecture and development details
4. Update links to docs/ directory

**Phase 4: Verification**
1. Run `make build` with README changes
2. Run `make run` and verify examples work
3. Check all links are valid
4. Verify code examples compile

**Phase 5: Documentation Review**
1. Review docs/ files for consistency
2. Update any outdated information
3. Ensure all features are documented
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 Code follows existing project conventions package structure naming error handling
- [ ] #2 go vet passes with no warnings
- [ ] #3 go build succeeds without errors
- [ ] #4 Unit tests added or updated for new or changed functionality
- [ ] #5 go test ... passes with no failures
- [ ] #6 Code comments added for non-obvious logic
- [ ] #7 README or docs updated if public behavior changes
- [ ] #8 make build succeeds
- [ ] #9 make run works as expected
- [ ] #10 Errors are logged not silently ignored
- [ ] #11 Graceful degradation monitor continues if individual file processing fails
- [ ] #12 No resource leaks channels closed files closed goroutines stopped
- [ ] #13 #1 UPDATED_README.md created with comprehensive project coverage
- [ ] #14 #2 All implemented features documented (file watcher, parser, change detection, notifier, agent orchestration)
- [ ] #15 #3 Installation and usage examples provided for both quick-start and advanced scenarios
- [ ] #16 #4 Configuration options documented with YAML examples
- [ ] #17 #5 Architecture overview explaining component interactions
- [ ] #18 #6 Development workflow and testing guidance included
- [ ] #19 #7 References to existing docs/ directory and backlog task files
- [ ] #20 #8 go vet passes with no warnings on project code
- [ ] #21 #9 go build succeeds without errors
- [ ] #22 #10 make build succeeds
- [ ] #23 #11 make run works as expected (verify with sample task file)
<!-- DOD:END -->
