---
id: GOT-028
title: test
status: To Do
assignee:
  - workflow
created_date: '2026-03-16 11:14'
updated_date: '2026-03-30 18:27'
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

The comprehensive README update will consolidate and enhance documentation across multiple files while maintaining consistency with the codebase. The approach focuses on creating a well-structured, maintainable README that serves both users and developers.

**Key objectives:**
- **Restructure and expand**: Transform the current README into a comprehensive guide covering installation, usage, configuration, and architecture
- **Consolidate from QWEN.md**: Integrate key architectural information while keeping QWEN.md as the detailed technical reference
- **Incorporate docs/ content**: Reference and integrate information from the docs directory (agent configuration, orchestration quickstart, setup guide)
- **Add practical examples**: Include working code examples for all major features
- **Maintain readability**: Keep the README under ~500 lines with clear section hierarchy

**Approach:**
- **Content audit**: Review all existing documentation files to identify overlaps and gaps
- **Section-by-section rewrite**: Build new README with logical flow from overview to advanced usage
- **Code block verification**: Ensure all bash and YAML examples are syntactically correct
- **Cross-reference strategy**: Use relative links to QWEN.md for deep architecture, docs/ for specific topics

**Trade-offs:**
- Keep README high-level; move detailed architecture to QWEN.md
- Link to docs/ for deep dives rather than duplicating content
- Use tables for configuration options, code blocks for examples
- Include "Last updated" or version info to track documentation drift

### 2. Files to Modify

| File | Action | Purpose |
|------|--------|---------|
| `README.md` | **Rewrite** | Comprehensive project documentation with all key sections |
| `QWEN.md` | **Update** | Add link to README for installation/usage; keep architecture/development details |
| `docs/agent-configuration.md` | **Review** | Verify consistency with README agent section |
| `docs/agent-orchestration-quickstart.md` | **Review** | Verify consistency with README usage examples |
| `docs/setup-monitor.md` | **Review** | Verify consistency with README installation section |

**Files to Reference (read-only):**
- `cmd/monitor/main.go` - CLI entry point and component wiring
- `pkg/watcher/*.go` - File monitoring implementation
- `pkg/parser/*.go` - YAML frontmatter extraction
- `pkg/cache/*.go` - State caching and debouncing
- `pkg/change_detect/*.go` - Assignee change detection
- `pkg/notifier/*.go` - Tmux notifications
- `pkg/matcher/*.go` - Agent matching logic
- `pkg/agent/*.go` - Agent configuration management
- `pkg/config/*.go` - Configuration loading
- `pkg/logs/*.go` - JSON logging
- `Makefile` - Build and run commands
- `agents/workflow/*` - Workflow agent implementation
- `scripts/agent_status.sh` - Agent status checker
- `maestro.yml` - Example configuration

### 3. Dependencies

**Prerequisites:**
- All core functionality from tasks GOT-008 through GOT-029 implemented (all marked Done)
- Go 1.25.7+ available for building (`go version` verifies)
- tmux installed for notifications (`tmux --version` verifies)
- `./backlog/tasks` directory exists
- Project builds successfully (`make build` passes)

**Build Verification:**
```bash
go vet ./...      # Static analysis (must pass with no warnings)
make build        # Build binary (must succeed)
```

**No blocking tasks** - All core features are implemented and ready for documentation update.

### 4. Code Patterns

**Conventions to Follow:**

1. **Markdown formatting:**
   - Use triple backticks with language identifier: ` ```yaml `, ` ```bash `, ` ```go `
   - Tables for configuration options and reference
   - Horizontal rules (`---`) between major sections

2. **Bash code examples:**
   - Use `make` commands where applicable
   - Show full commands including quotes and heredocs
   - Include `chmod +x` for scripts
   - Show file creation with `cat > file <<EOF`

3. **YAML code examples:**
   - Use `yaml` language identifier
   - Show minimal working configuration first, then extended
   - Document default values in text, not just in comments

4. **Go code examples:**
   - Use `go` language identifier
   - Show error handling patterns: `if err != nil { return err }`
   - Include necessary imports

5. **Path conventions:**
   - Use relative paths from project root (`./backlog/tasks`, `agents/agent-name/`)
   - Avoid absolute paths in examples

6. **Naming conventions:**
   - Package names: lowercase (`cache`, `watcher`, `parser`)
   - Function names: CamelCase (`NewWatcher`, `ProcessFile`)
   - Variable names: camelCase (`fileWatcher`, `eventQueue`)
   - Error variables: `Err` prefix (`ErrWatcherStopped`)

### 5. Testing Strategy

**Documentation verification (no unit tests required for README):**

1. **Build verification:**
   ```bash
   make build         # Verify binary builds
   go vet ./...       # Verify static analysis passes
   ```

2. **Example verification:**
   - Run all bash examples in isolated shell to verify syntax
   - Verify YAML examples are valid (can be parsed by Go yaml.v3)
   - Check that file creation examples work (`cat > file` patterns)

3. **Link verification:**
   - Verify all relative links to docs/ files are valid
   - Verify QWEN.md link is correct
   - Verify backlog/ task file references are valid

4. **Manual smoke test:**
   ```bash
   # Test quick-start example
   mkdir -p backlog/tasks
   touch backlog/tasks/test.md
   make run  # Verify monitor starts and shows output
   ```

**No code changes required** - README update is documentation-only.

### 6. Risks and Considerations

**Known risks:**

1. **Information overload:** README may become too long
   - *Mitigation:* Keep sections concise; use tables for reference; link to QWEN.md for depth

2. **Documentation drift:** Changes to code may not update README
   - *Mitigation:* Add maintenance checklist; include version/date marker

3. **Link rot:** Internal links may break if file structure changes
   - *Mitigation:* Use relative paths; review links after any file reorganization

4. **Example inaccuracy:** Code examples must match implementation
   - *Mitigation:* Test all examples; reference actual config files

5. **Version mismatch:** README may reference outdated features
   - *Mitigation:* Update README when features are added/changed

**Rollout considerations:**
- Pure documentation update - no code changes
- No deployment required
- No breaking changes for existing users
- No configuration changes required
- Backward compatible

### 7. Detailed README Structure

**Proposed section order:**

1. **Project Header:** Name, description, status (To Do/In Progress/Done badges)
2. **Overview:** High-level description of what Maestro does
3. **Features:** Bullet list of key capabilities (expand current list)
4. **Technology Stack:** Language, core libraries, tools
5. **Quick Start:** 3-step setup for immediate use (minimal example)
6. **Installation:** Detailed installation instructions with prerequisites
7. **Usage:** Basic monitoring, tmux notifications, agent orchestration
8. **Configuration:** maestro.yml, agent configuration, environment variables
9. **Architecture:** Component overview with data flow description
10. **Development:** Building, testing, code style, contribution guidelines
11. **Backlog.md Integration:** Task management workflow and task IDs
12. **Agent Orchestration:** Setting up and using agents (detailed)
13. **Troubleshooting:** Common issues and solutions
14. **Documentation:** Links to QWEN.md, docs/, and backlog/
15. **Contributing:** How to contribute (if applicable)
16. **License:** MIT license notice

**Sections to reference:**
- Link to `QWEN.md` for detailed architecture (component interactions, data structures)
- Link to `docs/` for deep dives on specific topics (agent config, orchestration setup)
- Link to `backlog/tasks/` for implementation details and task history
- Link to `scripts/` for utility scripts (agent_status.sh, attach.sh)

### 8. Implementation Steps

**Phase 1: Content Audit and Planning (Current Phase)**
1. Read all existing documentation files (README, QWEN, AGENTS, docs/)
2. Identify overlapping content and gaps
3. Create content outline with section priorities
4. Draft test examples to verify accuracy
5. **Decision point:** Present plan to user for approval

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

**Phase 4: Documentation Review**
1. Review docs/ files for consistency
2. Update any outdated information
3. Ensure all features are documented
4. Verify example commands work

**Phase 5: Verification**
1. Run `make build` with README changes
2. Run `go vet ./...` with README changes
3. Check all links are valid
4. Verify code examples compile
5. Test manual examples from README
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
