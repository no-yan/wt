# Command Reference

Complete reference for all `wt` commands.

## Shell Integration

**IMPORTANT**: For `wt switch` to work, you must set up zsh integration:

```bash
# Add to your zsh profile:
eval "$(wt shell-init)"
```

**Note**: Only zsh is supported. Other shells will show an informative error message.

## Project Structure

`wt` organizes all worktrees in the `worktrees/` subdirectory of your repository:

```
~/projects/myapp/
├── .git/
├── src/
├── README.md
├── .gitignore           # Contains "worktrees/"
└── worktrees/           # All worktrees here
    ├── main/           # Main branch worktree
    ├── feature-auth/   # feature/auth branch
    └── hotfix-bug-123/ # hotfix/bug-123 branch
```

## Global Options

```bash
-h, --help     Show help information
-v, --verbose  Enable verbose output
```

## Commands

### `wt list`

Display worktrees with flexible filtering and status information.

```bash
wt list [options]
```

**Options:**
- `--dirty` - Show only worktrees with uncommitted changes
- `--verbose` - Show detailed git status for each worktree (like `git status --short`)
- `--porcelain` - Machine-readable output format
- `--names-only` - Output only worktree names (useful for scripting)

**Default Output Format:**
```
STATUS NAME              PATH                           BRANCH
✓      main              ~/project/worktrees/main       [main]
*      feature-auth      ~/project/worktrees/feature-auth [feature/auth]
↑      hotfix-bug-123    ~/project/worktrees/hotfix-bug-123 [hotfix/bug-123]
```

**Status Indicators:**
- `✓` - Clean worktree (no uncommitted changes)
- `*` - Dirty worktree (uncommitted changes)
- `↑` - Ahead of remote
- `↓` - Behind remote
- `↕` - Diverged from remote
- `L` - Locked worktree
- `P` - Prunable worktree

**Verbose Output (--verbose):**
```
✓ main (~/project/worktrees/main) [main]

* feature-auth (~/project/worktrees/feature-auth) [feature/auth]
   M auth.go
  ?? test.txt

↑ hotfix-bug-123 (~/project/worktrees/hotfix-bug-123) [hotfix/bug-123]
   M security.go
   D old-file.go
```

**Examples:**
```bash
wt list                    # Show all worktrees
wt list --dirty            # Show only worktrees with changes
wt list --verbose          # Show detailed git status
wt list --names-only       # Output: main\nfeature-auth\nhotfix-bug-123
```

### `wt switch`

Navigate to a worktree directory with exact name matching.

```bash
wt switch <exact-name>
```

**Requirements:**
- Zsh shell with integration set up: `eval "$(wt shell-init)"`

**Arguments:**
- `<exact-name>` - Exact worktree name (no fuzzy matching)

**Worktree Names:**
- Branch `feature/auth` creates worktree named `feature-auth`
- Branch `hotfix/bug-123` creates worktree named `hotfix-bug-123`
- Branch `main` creates worktree named `main`

**Examples:**
```bash
wt switch main             # Switch to main worktree
wt switch feature-auth     # Switch to feature-auth worktree
wt switch hotfix-bug-123   # Switch to hotfix worktree
```

**Error Handling:**
If a worktree is not found, `wt switch` will display available options:
```bash
$ wt switch nonexistent
Worktree not found: nonexistent
Available worktrees:
  main
  feature-auth
  hotfix-bug-123
```

**Tab Completion:**
Zsh integration provides tab completion for worktree names:
```bash
wt switch <TAB>            # Shows available worktree names
```

### `wt add`

Create a new worktree in the `worktrees/` subdirectory.

```bash
wt add <branch> [options]
```

**Arguments:**
- `<branch>` - Branch name for the new worktree

**Options:**
- `-b, --new-branch` - Create new branch if it doesn't exist
- `-B, --force-new-branch` - Force create new branch (reset if exists)
- `--detach` - Create detached HEAD worktree
- `--force` - Force creation even if branch is checked out elsewhere

**Path Generation:**
- `feature/auth` → `worktrees/feature-auth/`
- `hotfix/bug-123` → `worktrees/hotfix-bug-123/`
- `docs/api-update` → `worktrees/docs-api-update/`
- `main` → `worktrees/main/`

**Path Generation Rules:**
1. Replace slashes with dashes
2. Replace underscores with dashes  
3. Create in `worktrees/` subdirectory
4. Ensure uniqueness (append number if needed)

**Auto-Setup:**
On first use in a repository, `wt` automatically:
1. Creates `worktrees/` directory
2. Adds `worktrees/` to `.gitignore`

**Examples:**
```bash
wt add feature/auth              # Creates worktrees/feature-auth/
wt add -b new-feature            # Create new branch and worktree
wt add --detach HEAD~1           # Detached worktree at HEAD~1
wt add main                      # Create main branch worktree
```

### `wt remove`

Safely remove a worktree with confirmation prompts.

```bash
wt remove <exact-name> [options]
```

**Arguments:**
- `<exact-name>` - Exact worktree name to remove

**Options:**
- `--force` - Force removal without confirmation
- `--keep-branch` - Remove worktree but keep the branch

**Safety Features:**
- Confirms removal of dirty worktrees
- Prevents removal of current worktree
- Shows what will be removed before confirmation

**Examples:**
```bash
wt remove feature-auth     # Remove with confirmation
wt remove --force auth     # Force remove without prompt
wt remove --keep-branch feature-auth  # Remove worktree, keep branch
```

### `wt clean`

Clean up stale and orphaned worktrees.

```bash
wt clean [options]
```

**Options:**
- `--dry-run` - Show what would be cleaned without removing
- `--force` - Remove without confirmation prompts
- `--expire <time>` - Only remove worktrees older than specified time

**What Gets Cleaned:**
- Prunable worktrees (administrative files without directories)
- Worktrees with missing directories
- Unlocked worktrees that are no longer valid

**Examples:**
```bash
wt clean               # Interactive cleanup
wt clean --dry-run     # Show what would be cleaned
wt clean --force       # Clean without prompts
wt clean --expire 30d  # Clean worktrees older than 30 days
```

### `wt shell-init`

Generate zsh integration code for directory switching.

```bash
wt shell-init
```

**Output:**
Generates zsh functions and completions.

**Examples:**
```bash
# Setup for current shell
eval "$(wt shell-init)"

# View generated code
wt shell-init

# Add to zsh profile
echo 'eval "$(wt shell-init)"' >> ~/.zshrc
```

**What It Does:**
- Creates `wt` zsh function that wraps the binary
- Intercepts `wt switch` calls to perform actual directory changes
- Adds tab completion for worktree names and commands
- Preserves all other commands to pass through to the binary

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Command syntax error
- `3` - Git repository error
- `4` - Worktree not found
- `5` - Operation cancelled by user
- `6` - Shell not supported (non-zsh)

## Shell Integration Details

### How It Works

The zsh integration works similarly to tools like `zoxide`:

1. `wt shell-init` generates zsh functions
2. These functions intercept `wt switch` commands
3. The binary returns the target path
4. The zsh function performs the actual `cd`

### Zsh Integration

```zsh
wt() {
    case "$1" in
        switch)
            if [[ -z "$2" ]]; then
                echo "Usage: wt switch <worktree-name>" >&2
                return 1
            fi
            
            local target_dir="$(command wt path "$2" 2>/dev/null)"
            if [[ $? -eq 0 && -d "$target_dir" ]]; then
                cd "$target_dir"
                echo "→ $2 ($target_dir)"
            else
                echo "Worktree not found: $2" >&2
                echo "Available worktrees:" >&2
                command wt list --names-only 2>/dev/null | sed 's/^/  /' >&2
                return 1
            fi
            ;;
        *)
            command wt "$@"
            ;;
    esac
}
```

### Tab Completion

Zsh integration includes comprehensive tab completion:

- `wt switch <TAB>` - Complete with worktree names
- `wt remove <TAB>` - Complete with worktree names  
- `wt add <TAB>` - Complete with branch names
- `wt <TAB>` - Complete with subcommands

## Migration from git worktree

### Common Commands

| git worktree | wt equivalent | Notes |
|--------------|-----------------|-------|
| `git worktree list` | `wt list` | Enhanced with status indicators |
| `git worktree add ../feature feature` | `wt add feature` | Creates in worktrees/ |
| `git worktree remove ../feature` | `wt remove feature-name` | Exact name required |
| `git worktree prune` | `wt clean` | Interactive with safety checks |

### Workflow Integration

```bash
# Traditional workflow
git worktree add ../feature-auth feature/auth
cd ../feature-auth
# work...
cd ../main

# wt workflow  
wt add feature/auth        # Creates worktrees/feature-auth/
wt switch feature-auth     # Navigate there
# work...
wt switch main            # Return to main worktree
```

## Path Structure

### Repository Layout
```
~/projects/myapp/           # Repository root
├── .git/                   # Git directory
├── src/                    # Source code
├── .gitignore              # Contains "worktrees/"
└── worktrees/              # All worktrees
    ├── main/              # Main branch
    ├── feature-auth/      # feature/auth branch  
    ├── hotfix-bug-123/    # hotfix/bug-123 branch
    └── docs-update/       # docs/update branch
```

### Branch to Directory Mapping
- `main` → `worktrees/main/`
- `feature/auth` → `worktrees/feature-auth/`
- `hotfix/bug-123` → `worktrees/hotfix-bug-123/`
- `docs/api_update` → `worktrees/docs-api-update/`

## Troubleshooting

### "Shell not supported"

**Problem:** Using bash, fish, or other non-zsh shell
**Solution:** Switch to zsh: `zsh` then set up integration

### "Command not found" after switch

**Problem:** `wt switch` says "command not found"  
**Solution:** Set up zsh integration: `eval "$(wt shell-init)"`

### Directory doesn't change

**Problem:** `wt switch` runs but directory doesn't change
**Solution:** Ensure zsh integration is loaded in your current shell session

### Worktree not found

**Problem:** `wt switch` can't find worktree
**Solution:** Use `wt list` to see exact names, check spelling

### Completion not working

**Problem:** Tab completion doesn't work
**Solution:** Restart zsh after adding shell integration to profile

### Permission denied creating worktree

**Problem:** Cannot create in worktrees/ directory
**Solution:** Check repository permissions, ensure you're in git repository root