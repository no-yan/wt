# Architecture

## System Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Zsh Layer     │    │   CLI Commands   │    │  Core Logic     │
│   (shell funcs) │───>│   (cmd/)         │───>│  (internal/)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
        │                      │                       │
        v                      v                       v
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Directory Ops   │    │  User Output     │    │  Git Commands   │
│ (cd via shell)  │    │  (display)       │    │  (git worktree) │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Module Responsibilities

### `/cmd` - CLI Layer
- Cobra command definitions
- Argument parsing and validation
- User interaction (prompts, confirmations)
- Output formatting
- **Does NOT** change directories (shell layer does this)

### `/internal` - Core Logic
- `worktree.go` - Worktree model and helpers
- `worktree_manager.go` - Worktree CRUD operations
- `git.go` - Git output parsing
- `git_service.go` - Git command execution

## Key Data Structures

### Worktree
```go
type Worktree struct {
    Branch string   // Current branch name
    Path   string   // Filesystem path
    Head   string   // HEAD commit hash
    Status Status   // Clean/Dirty/Stale
}
```

### WorktreeManager
Provides all worktree operations:
- `List()` - Get all worktrees
- `Add(branch)` - Create new worktree
- `Remove(name)` - Delete worktree
- `FindByName(name)` - Exact name lookup

## Data Flow

### Switch Command
```
User: wt switch auth
  │
  v
Shell function intercepts "switch"
  │
  v
Binary: wt switch auth (resolves path)
  │
  v
Shell function: cd "$resolved_path"
  │
  v
Directory changed
```

### Add Command
```
User: wt add feature/login
  │
  v
Generate path: worktrees/feature-login
  │
  v
git worktree add worktrees/feature-login feature/login
  │
  v
Success message with path
```

## Design Principles

1. **Predictable Structure**
   - All worktrees in `worktrees/` subdirectory
   - Branch `feature/foo` -> directory `worktrees/feature-foo`

2. **Shell Integration**
   - `wt` binary handles logic
   - Shell functions handle directory changes
   - Generated via `wt shell-init`

3. **Exact Matching**
   - No fuzzy matching to avoid ambiguity
   - Names must match exactly

4. **Minimal Surface**
   - Few commands: list, add, switch, remove, clean
   - Each command does one thing well

## Error Handling

- Wrap errors with context: `fmt.Errorf("...: %w", err)`
- Provide actionable messages
- Fail fast for unrecoverable errors

## Testing Strategy

- **Unit tests** for core logic (`*_test.go`)
- **Integration tests** for Git operations
- **Mock interfaces** for testability (`CommandRunner`)

## See Also

- `docs/implementation/architecture.md` - Detailed architecture documentation
- `docs/implementation/testing.md` - Testing guidelines
