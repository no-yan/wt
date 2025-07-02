# Claude Implementation Context

This document contains essential context for Claude to implement and maintain the `wrkt` project.

## Project Overview

**Goal**: Create a CLI tool that makes git worktree operations seamless by organizing worktrees in a predictable structure.

**Problem**: Git worktrees require manual directory navigation (`cd`) and are scattered across filesystem.

**Solution**: Simple CLI that organizes all worktrees in `worktrees/` subdirectory with zsh integration for navigation.

## Technology Stack

- **Language**: Go (1.21+)
- **CLI Framework**: Cobra
- **Shell Support**: Zsh only (no multi-shell complexity)
- **Dependencies**: Minimal (cobra only for MVP)
- **Target**: Single binary, works on macOS/Linux

## Architecture

```
wrkt/
├── main.go                 # Entry point
├── cmd/                    # CLI commands (cobra)
│   ├── root.go            # Root command
│   ├── list.go            # List command (unified status display)
│   ├── switch.go          # Switch command (path resolution only)
│   ├── add.go             # Add command
│   ├── remove.go          # Remove command
│   ├── clean.go           # Clean command
│   └── shell.go           # Shell integration command
├── internal/              # Core logic
│   ├── worktree.go        # Worktree operations
│   ├── git.go             # Git integration
│   ├── shell.go           # Zsh function generation
│   └── display.go         # Output formatting
└── docs/                  # Documentation
```

## Core Components

### 1. Worktree Organization (Foundation)
**All worktrees in `worktrees/` subdirectory**:
- Predictable location: `$REPO_ROOT/worktrees/`
- Auto-setup: Create directory and add to .gitignore
- Simple path mapping: `feature/auth` → `worktrees/feature-auth/`
- No permission issues (same owner as repo)

### 2. Zsh Integration (Critical for Navigation)
**Simple, reliable zsh-only shell integration**:
- Generate zsh functions via `wrkt shell-init`
- Intercept `wrkt switch` calls in zsh functions
- Perform actual `cd` operations in the user's shell
- Tab completion for all commands

### 3. Simplified Command Set
- **list**: Unified status display with filtering options
- **switch**: Exact name matching (no fuzzy matching complexity)
- **add**: Auto-path generation to `worktrees/` subdirectory
- **remove**: Safe removal with confirmations
- **clean**: Automated cleanup of stale worktrees
- **shell-init**: Generate zsh integration code

### 4. Simple Features
- Zsh-only shell integration (eliminates multi-shell complexity)
- Exact name matching (deterministic, predictable behavior)
- Auto-path generation: `feature/auth` → `worktrees/feature-auth/`
- Organized structure in single location
- Clear error messages with available options

## Implementation Guidelines

### Code Style
- Follow Go conventions
- Use structured logging
- Handle errors gracefully
- Write tests for core logic

### Git Integration
- Use `git worktree` commands exclusively
- Parse porcelain output for reliability
- Handle edge cases (bare repos, submodules)
- Validate git repository context

### User Experience
- Consistent command patterns
- Clear error messages
- Interactive prompts for destructive operations
- Helpful suggestions for fuzzy matches

## MVP Scope

**Core Commands**:
- `wrkt list` - Unified status display with filtering options
- `wrkt switch <exact-name>` - Directory switching with zsh integration
- `wrkt add <branch>` - Worktree creation in `worktrees/` subdirectory
- `wrkt remove <exact-name>` - Safe worktree removal
- `wrkt clean` - Cleanup stale worktrees
- `wrkt shell-init` - Zsh integration setup

**Essential Features**:
- **Zsh-only shell integration** (mandatory for switch functionality)
- **Exact name matching** (no fuzzy matching complexity)
- **Worktrees in `worktrees/` subdirectory** (organized structure)
- Auto-path generation with simple rules
- Unified status display (replaces separate status command)
- Tab completion for zsh
- Auto-setup of worktrees directory and .gitignore

**Out of MVP**:
- Multi-shell support (bash, fish, etc.)
- Fuzzy matching
- GitHub PR status integration
- Configuration files
- Claude development status tracking
- Batch operations
- Previous directory tracking (`wrkt switch -`)

## Quality Criteria

### Functionality
- All MVP commands work correctly
- Fuzzy matching finds intended worktrees
- Auto-path generation creates logical paths
- Safe operations with proper error handling

### Robustness
- Handles edge cases (missing worktrees, invalid repos)
- Graceful error messages
- No data loss on operations
- Works across different git repository states

### Usability
- Intuitive command interface
- Clear help documentation
- Consistent behavior patterns
- Fast execution (<100ms for most operations)

## Testing Strategy

### Unit Tests
- Worktree parser logic
- Path generation algorithms
- Fuzzy matching functions
- Git command integration

### Integration Tests
- End-to-end command execution
- Git repository interactions
- File system operations
- Error scenario handling

### Manual Testing
- Cross-platform compatibility
- Real git repository workflows
- Edge case scenarios
- Performance validation

## Development Commands

```bash
# Build and test
go build -o wrkt
go test ./...

# Run linting
golangci-lint run

# Cross-platform builds
GOOS=linux go build -o wrkt-linux
GOOS=darwin go build -o wrkt-darwin
GOOS=windows go build -o wrkt.exe
```

## Implementation Checklist

- [x] Set up Go module and dependencies
- [x] Implement worktree parser and data structures
- [x] Create cobra command structure
- [x] **Implement worktrees/ organization system** (foundation)
- [x] Auto-setup worktrees directory and .gitignore entry
- [x] Implement path generation: branch → worktree name
- [x] **Implement zsh shell integration system** (critical)
- [x] Implement shell-init command with zsh function generation
- [x] Implement list command with unified status display
- [x] Add filtering options (--dirty, --verbose, --names-only)
- [x] Implement switch command with exact name matching
- [x] Implement add command with worktrees/ path generation
- [x] Implement remove command with safety checks
- [x] Implement clean command for stale worktrees
- [ ] Add zsh tab completion
- [x] Add comprehensive error handling
- [x] Write unit tests for core functions
- [x] Write integration tests for comprehensive workflows
- [ ] Test on macOS and Linux with zsh
- [x] Validate against quality criteria
- [x] Update documentation

## Current Status (as of latest session)

**Completed Features:**
- ✅ Complete worktree management system with `worktrees/` organization
- ✅ Auto-setup of worktrees directory and .gitignore integration
- ✅ Full list command with filtering options (--dirty, --verbose, --names-only)
- ✅ Detailed git status parsing and display for verbose mode
- ✅ Clean command for stale worktree cleanup using `git worktree prune`
- ✅ Comprehensive unit and integration tests
- ✅ Add, remove, and switch commands with safety checks
- ✅ Shell integration system with zsh function generation
- ✅ Robust error handling and input validation
- ✅ Branch name to worktree path conversion
- ✅ Working development workflow using wrkt itself

**Remaining Work:**
- [ ] Zsh tab completion implementation
- [ ] Cross-platform testing (macOS/Linux with zsh)
- [ ] Performance testing and optimization
- [ ] Documentation updates and examples

**Before starting any development work, Claude MUST:**

1. **Read Current Todo List**: Use `TodoRead` tool to check existing tasks
2. **Plan Work**: Use `TodoWrite` tool to create/update task list for the session
3. **Track Progress**: Update todo status throughout development:
   - `"pending"` - Task not yet started
   - `"in_progress"` - Currently working on (limit to ONE at a time)
   - `"completed"` - Task finished successfully
4. **Mark Completion**: IMMEDIATELY mark tasks as completed when finished

### CRITICAL: Parallel Development Manager Responsibilities

**When acting as manager for multiple Claude agents working on different worktrees:**

1. **Conflict Prevention**: Analyze all feature branches before assigning tasks to prevent merge conflicts
2. **Priority-Based Coordination**: Pause lower priority work when conflicts are detected
3. **Sequential Integration**: Merge features one-by-one in dependency order
4. **Test Validation**: Ensure ALL tests pass before any merge attempt
5. **Rollback Strategy**: Abort merges immediately if tests fail
6. **Documentation Updates**: Record all coordination decisions and conflicts in CLAUDE.md

**Conflict Resolution Protocol:**
- Run `git diff f0c4919..branch --name-only` to analyze file changes
- Stop development on conflicting features until higher priority merges complete
- Coordinate timing of commits to avoid simultaneous changes to same files
- Always test merges in clean environment before final integration

**GitHub PR Integration Process:**
- Use GitHub Pull Requests for all feature merges instead of direct git merge
- Create PRs using `gh pr create` with proper titles and descriptions
- Ensure all tests pass in PR before merging
- Use PR reviews to validate changes before integration
- Merge PRs sequentially to prevent conflicts

### Worktree Development Workflow

**When working across multiple worktrees:**

1. **Check WORKTREE_TRACKING.md** - Review status of all active worktrees
2. **Use TodoWrite** to plan which worktrees to work on
3. **Work systematically** - Complete one worktree before moving to next
4. **Update tracking documents** when switching between worktrees
5. **Commit frequently** with descriptive messages

### Task Categories to Track

- **High Priority**: Core functionality, bug fixes, incomplete features
- **Medium Priority**: Enhancements, new features, optimizations  
- **Low Priority**: Documentation, cleanup, nice-to-have features

### Example Todo Usage

```
TodoWrite: [
  {"content": "Complete feature-list-filters worktree", "status": "in_progress", "priority": "high", "id": "1"},
  {"content": "Add zsh tab completion", "status": "pending", "priority": "medium", "id": "2"},
  {"content": "Update documentation", "status": "pending", "priority": "low", "id": "3"}
]
```

## Notes for Future Claude Sessions

1. **Zsh-only integration** - Don't implement multi-shell support
2. **Worktrees go in `worktrees/` subdirectory** - Always use this location
3. **No fuzzy matching** - Exact name matching only
4. **Always check existing code** before implementing new features
5. **Follow the established patterns** in cmd/ and internal/ directories
6. **Test zsh integration** thoroughly
7. **Update tests** when adding new functionality
8. **Validate MVP scope** - don't add features beyond core requirements
9. **Test with real git repositories** to ensure robustness
10. **Update documentation** when changing behavior
11. **Clean command uses `git worktree prune`** - This is the correct approach for stale worktrees
12. **Integration tests are comprehensive** - Cover full worktree lifecycle and edge cases

## Critical Implementation Notes

### Worktree Organization
- All worktrees created in `$REPO_ROOT/worktrees/` subdirectory
- Auto-create worktrees/ directory on first use
- Auto-add "worktrees/" to .gitignore
- Path generation: `feature/auth` → `worktrees/feature-auth/`

### Zsh Integration Implementation
- `wrkt switch` command should only resolve and return the target path
- Actual directory changing is handled by zsh functions
- Zsh functions must be generated by `wrkt shell-init`
- Only support zsh - show clear error for other shells

### Command Design
- `wrkt list` is the primary information command
- No separate `wrkt status` - use `wrkt list --dirty` instead
- `wrkt list --verbose` provides detailed git status
- Exact name matching only - no fuzzy matching complexity

## Implementation Details (Current State)

### Completed Commands
1. **`wrkt list`** - Unified status display with comprehensive filtering
   - `--dirty`: Filter to show only worktrees with uncommitted changes
   - `--verbose`: Show detailed git status information with parsed change descriptions
   - `--names-only`: Output only worktree names for scripting
   - Supports combining filters (e.g., `--dirty --names-only`)

2. **`wrkt add <branch>`** - Create worktrees in organized structure
   - Auto-creates `worktrees/` directory on first use
   - Auto-adds `worktrees/` to `.gitignore`
   - Path generation: `feature/auth` → `worktrees/feature-auth/`
   - Input validation and security checks

3. **`wrkt remove <name>`** - Safe worktree removal
   - Exact name matching for safety
   - Prevents removal of main worktree
   - Blocks removal of dirty worktrees
   - Comprehensive error handling

4. **`wrkt clean`** - Stale worktree cleanup
   - `--dry-run`: Preview what would be cleaned
   - `--force`: Skip confirmation prompts
   - Uses `git worktree prune` for safe cleanup
   - Detects and handles missing worktree directories

5. **`wrkt switch <name>`** - Path resolution for zsh integration
   - Exact name matching (no fuzzy matching)
   - Returns target path for shell functions

6. **`wrkt shell-init`** - Generate zsh integration code
   - Creates shell functions for directory changing
   - Enables seamless worktree navigation

### Test Coverage
- **Unit Tests**: Core logic, parsers, validation functions
- **Integration Tests**: Full worktree lifecycle with real git operations
- **Command Tests**: All CLI commands with mock and real scenarios
- **Edge Case Coverage**: Stale worktrees, dirty states, error conditions

### Status Detection System
- **StatusClean**: No uncommitted changes
- **StatusDirty**: Has uncommitted changes (detected via `git status --porcelain`)
- **StatusStale**: Worktree directory missing or invalid

### Security Features
- Branch name validation (prevents command injection)
- Path validation (absolute paths required)
- Shell escaping for all command execution
- Input sanitization throughout

## Troubleshooting Common Issues

- **"wrkt switch doesn't change directory"**: User hasn't set up zsh integration
- **"shell not supported"**: User not using zsh
- **"command not found" after switch**: Zsh integration not loaded
- **"worktree not found"**: Check exact name with `wrkt list`
- **Path generation conflicts**: Simple conflict resolution with numbering
- **Tab completion not working**: Zsh integration setup incomplete
- **Stale worktrees not cleaning**: Use `wrkt clean` which uses `git worktree prune`