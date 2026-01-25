# Agent Instructions

Instructions for AI agents working on this project.

## Project Overview

`wt` is a Git worktree management CLI tool written in Go. It organizes worktrees in a predictable `worktrees/` subdirectory and provides zsh shell integration.

## Quick Reference

### Build & Test
```bash
go build -o wt .       # Build binary
go test ./...          # Run all tests
go vet ./...           # Static analysis
```

### Test the CLI
```bash
./wt list              # List worktrees
./wt list --verbose    # With git status
./wt add <branch>      # Create worktree
./wt switch <name>     # Switch (requires shell integration)
./wt remove <name>     # Remove worktree
./wt clean             # Remove merged worktrees
```

## Project Structure

```
wt/
├── main.go           # Entry point
├── cmd/              # CLI commands (Cobra)
│   ├── root.go       # Root command
│   ├── list.go       # wt list
│   ├── add.go        # wt add
│   ├── switch.go     # wt switch
│   ├── remove.go     # wt remove
│   └── clean.go      # wt clean
├── internal/         # Core logic
│   ├── worktree.go   # Worktree model
│   ├── worktree_manager.go  # Worktree operations
│   ├── git.go        # Git command parsing
│   └── git_service.go # Git command execution
└── docs/             # Documentation
```

## Code Patterns

### Adding a New Command
1. Create `cmd/<command>.go`
2. Define Cobra command with `Use`, `Short`, `RunE`
3. Register in `cmd/root.go` init()
4. Add tests in `cmd/<command>_test.go`

### Working with Worktrees
```go
// Get git service and manager
service := internal.NewGitService()
manager := internal.NewWorktreeManager(service)

// List worktrees
worktrees, err := manager.List()

// Add worktree
path, err := manager.Add(branch)
```

## Development Rules

1. **No git checkout/switch** in main directory - use `wt` commands
2. **One worktree = one branch = one feature**
3. **Test frequently** - run `go test ./...` often
4. **Follow existing patterns** - check similar code before implementing

## Key Design Decisions

- **Zsh-only** shell integration (no bash/fish support)
- **Exact matching** only (no fuzzy matching)
- **Predictable paths**: `feature/foo` -> `worktrees/feature-foo/`
- **worktrees/ directory** is auto-created and gitignored

## Related Documentation

- `CLAUDE.md` - Claude-specific development rules
- `docs/implementation/architecture.md` - Detailed architecture
- `.claude/WORKFLOWS.md` - Development workflows
