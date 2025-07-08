# Claude Development Handoff Guide

This document enables Claude to quickly orient and continue development across sessions.

## Quick Start

### Build & Test
```bash
# Build the project
go build -o wt

# Run tests
go test ./...

# Install locally (optional)
go install
```

### Current Commands
- `wt list` - Show all worktrees (with `--dirty`, `--verbose` flags)
- `wt add <branch>` - Create new worktree in `worktrees/` directory
- `wt switch <name>` - Return path for zsh integration (no actual cd)
- `wt remove <name>` - Remove worktree
- `wt clean` - Clean up stale worktrees
- `wt shell-init` - Generate zsh integration functions

### Dependencies
- Go 1.24+ (current: Go 1.24.0)
- Git 2.5+ (for worktree support)
- Cobra CLI framework (`github.com/spf13/cobra v1.9.1`)

### Key Architecture Facts
- **Zsh-only integration**: Generates shell functions for directory switching
- **Exact matching**: No fuzzy matching - names must match exactly
- **Organized structure**: All worktrees in `worktrees/` subdirectory
- **Separation**: CLI commands (`cmd/`) vs business logic (`internal/`)

## Project Context

### Critical Design Decisions
1. **Zsh Integration**: `wt switch` returns path, zsh function handles `cd`
2. **Worktrees Organization**: All worktrees in `$REPO_ROOT/worktrees/` subdirectory
3. **Exact Matching**: No fuzzy matching complexity - exact names only
4. **Auto-setup**: Creates `worktrees/` directory and adds to `.gitignore`

### Current Command Structure
```
wt/
├── main.go              # Entry point: cmd.Execute()
├── cmd/                 # CLI commands using Cobra
│   ├── root.go         # Root command with aliases (sw, rm, ls)
│   ├── add.go          # Create worktree in worktrees/
│   ├── list.go         # List with filtering options
│   ├── switch.go       # Path resolution for zsh
│   ├── remove.go       # Remove worktree
│   ├── clean.go        # Clean stale worktrees
│   └── shell_init.go   # Generate zsh functions
├── internal/           # Business logic
│   ├── worktree_manager.go  # Core worktree operations
│   ├── git_service.go       # Git command integration
│   └── worktree.go         # Worktree data model
└── test/               # Integration tests
```

### Worktree Organization System
- All worktrees created in `worktrees/` subdirectory
- Path generation: `feature/auth` → `worktrees/feature-auth/`
- Auto-creates directory and adds to `.gitignore` on first use
- Safety checks prevent removing main worktree

## Actual Code Patterns

### Command Pattern (based on `cmd/add.go`)
```go
var addCmd = &cobra.Command{
    Use:   "add <branch>",
    Short: "Add a new worktree",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        // 1. Create services
        runner := internal.NewExecCommandRunner()
        gitService := internal.NewGitService(runner)
        manager := internal.NewWorktreeManager(gitService, runner)
        
        // 2. Get repo root
        repoPath, err := getRepoRoot(runner)
        
        // 3. Execute business logic
        worktreePath, err := manager.AddWorktree(repoPath, branch)
        
        // 4. Handle output/errors
        fmt.Printf("Added worktree: %s\n", worktreePath)
    },
}
```

### Error Handling Pattern (from `internal/worktree_manager.go`)
```go
func (wm *WorktreeManager) RemoveWorktree(repoPath, name string) error {
    if err := validatePath(repoPath); err != nil {
        return fmt.Errorf("invalid repository path: %w", err)
    }
    
    // Safety check: don't remove main worktree
    if !strings.Contains(targetWorktree.Path, "/worktrees/") {
        return fmt.Errorf("cannot remove main worktree %q", name)
    }
    
    // Safety check: warn if worktree has uncommitted changes
    if targetWorktree.Status == StatusDirty {
        return fmt.Errorf("worktree %q has uncommitted changes, commit or stash them first", name)
    }
    
    return nil
}
```

### Atomic Operations Pattern
```go
// Phase 1: Validate all targets before removing any (fail-fast strategy)
for _, name := range names {
    if targetWorktree.Status == StatusDirty {
        return fmt.Errorf("worktree %q has uncommitted changes", name)
    }
}

// Phase 2: All validations passed, execute removals
for _, target := range targetsToRemove {
    if err := wm.removeGitWorktree(repoPath, target.Path); err != nil {
        return fmt.Errorf("failed to remove worktree %q: %w", target.Name(), err)
    }
}
```

### Testing Pattern (from integration tests)
```go
func TestWorktreeManager_Add(t *testing.T) {
    tmpDir := t.TempDir()
    setupGitRepo(t, tmpDir)
    
    runner := NewExecCommandRunner()
    gitService := NewGitService(runner)
    manager := NewWorktreeManager(gitService, runner)
    
    worktreePath, err := manager.AddWorktree(tmpDir, "feature-branch")
    require.NoError(t, err)
    
    // Verify worktree exists
    assert.Contains(t, worktreePath, "worktrees/feature-branch")
}
```

## Development Workflow

### Adding New Commands
1. Create command file in `cmd/` (follow `add.go` pattern)
2. Implement business logic in `internal/`
3. Add integration tests with real git repos
4. Register command in `cmd/root.go`

### Testing Changes
```bash
# Run all tests
go test ./...

# Run specific test
go test ./cmd -run TestAdd

# Run integration tests
go test ./cmd -run Integration
```

### Validating Functionality
- Test with real git repository
- Verify worktree creation in `worktrees/` directory
- Check error handling for edge cases
- Test zsh integration if shell-related changes

### Common Development Tasks
- **Modify git operations**: Update `internal/git_service.go`
- **Add command flags**: Update cobra command definitions
- **Change worktree behavior**: Update `internal/worktree_manager.go`
- **Test shell integration**: Use `eval "$(./wt shell-init)"` in zsh

## Documentation References

For comprehensive information, see:
- **Architecture details**: `docs/implementation/architecture.md`
- **Command reference**: `docs/api/commands.md`
- **Dogfooding instructions**: `CLAUDE.md`
- **Usage examples**: `docs/examples/workflows.md`

This document focuses on quick Claude session startup. Refer to existing documentation for detailed design decisions and comprehensive API reference.
