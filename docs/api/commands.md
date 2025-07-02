# Command Reference

Complete reference for all `wrkt` commands.

## Shell Integration

**IMPORTANT**: For `wrkt switch` to work, you must set up zsh integration:

```bash
# Add to your zsh profile:
eval "$(wrkt shell-init)"
```

**Note**: Only zsh is supported. Other shells will show an informative error message.

## Project Structure

`wrkt` organizes all worktrees in the `worktrees/` subdirectory of your repository:

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

### `wrkt list`

Display worktrees with flexible filtering and status information.

```bash
wrkt list [options]
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
wrkt list                    # Show all worktrees
wrkt list --dirty            # Show only worktrees with changes
wrkt list --verbose          # Show detailed git status
wrkt list --names-only       # Output: main\nfeature-auth\nhotfix-bug-123
```

### `wrkt switch`

Navigate to a worktree directory with exact name matching.

```bash
wrkt switch <exact-name>
```

**Requirements:**
- Zsh shell with integration set up: `eval "$(wrkt shell-init)"`

**Arguments:**
- `<exact-name>` - Exact worktree name (no fuzzy matching)

**Worktree Names:**
- Branch `feature/auth` creates worktree named `feature-auth`
- Branch `hotfix/bug-123` creates worktree named `hotfix-bug-123`
- Branch `main` creates worktree named `main`

**Examples:**
```bash
wrkt switch main             # Switch to main worktree
wrkt switch feature-auth     # Switch to feature-auth worktree
wrkt switch hotfix-bug-123   # Switch to hotfix worktree
```

**Error Handling:**
If a worktree is not found, `wrkt switch` will display available options:
```bash
$ wrkt switch nonexistent
Worktree not found: nonexistent
Available worktrees:
  main
  feature-auth
  hotfix-bug-123
```

**Tab Completion:**
Zsh integration provides tab completion for worktree names:
```bash
wrkt switch <TAB>            # Shows available worktree names
```

### `wrkt add`

Create a new worktree in the `worktrees/` subdirectory.

```bash
wrkt add <branch> [options]
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
On first use in a repository, `wrkt` automatically:
1. Creates `worktrees/` directory
2. Adds `worktrees/` to `.gitignore`

**Examples:**
```bash
wrkt add feature/auth              # Creates worktrees/feature-auth/
wrkt add -b new-feature            # Create new branch and worktree
wrkt add --detach HEAD~1           # Detached worktree at HEAD~1
wrkt add main                      # Create main branch worktree
```

### `wrkt remove`

Safely remove a worktree with confirmation prompts.

```bash
wrkt remove <exact-name> [options]
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
wrkt remove feature-auth     # Remove with confirmation
wrkt remove --force auth     # Force remove without prompt
wrkt remove --keep-branch feature-auth  # Remove worktree, keep branch
```

### `wrkt clean`

Clean up stale and orphaned worktrees.

```bash
wrkt clean [options]
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
wrkt clean               # Interactive cleanup
wrkt clean --dry-run     # Show what would be cleaned
wrkt clean --force       # Clean without prompts
wrkt clean --expire 30d  # Clean worktrees older than 30 days
```

### `wrkt shell-init`

Generate zsh integration code for directory switching.

```bash
wrkt shell-init
```

**Output:**
Generates zsh functions and completions.

**Examples:**
```bash
# Setup for current shell
eval "$(wrkt shell-init)"

# View generated code
wrkt shell-init

# Add to zsh profile
echo 'eval "$(wrkt shell-init)"' >> ~/.zshrc
```

**What It Does:**
- Creates `wrkt` zsh function that wraps the binary
- Intercepts `wrkt switch` calls to perform actual directory changes
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

1. `wrkt shell-init` generates zsh functions
2. These functions intercept `wrkt switch` commands
3. The binary returns the target path
4. The zsh function performs the actual `cd`

### Zsh Integration

```zsh
wrkt() {
    case "$1" in
        switch)
            if [[ -z "$2" ]]; then
                echo "Usage: wrkt switch <worktree-name>" >&2
                return 1
            fi
            
            local target_dir="$(command wrkt path "$2" 2>/dev/null)"
            if [[ $? -eq 0 && -d "$target_dir" ]]; then
                cd "$target_dir"
                echo "→ $2 ($target_dir)"
            else
                echo "Worktree not found: $2" >&2
                echo "Available worktrees:" >&2
                command wrkt list --names-only 2>/dev/null | sed 's/^/  /' >&2
                return 1
            fi
            ;;
        *)
            command wrkt "$@"
            ;;
    esac
}
```

### Tab Completion

Zsh integration includes comprehensive tab completion:

- `wrkt switch <TAB>` - Complete with worktree names
- `wrkt remove <TAB>` - Complete with worktree names  
- `wrkt add <TAB>` - Complete with branch names
- `wrkt <TAB>` - Complete with subcommands

## Migration from git worktree

### Common Commands

| git worktree | wrkt equivalent | Notes |
|--------------|-----------------|-------|
| `git worktree list` | `wrkt list` | Enhanced with status indicators |
| `git worktree add ../feature feature` | `wrkt add feature` | Creates in worktrees/ |
| `git worktree remove ../feature` | `wrkt remove feature-name` | Exact name required |
| `git worktree prune` | `wrkt clean` | Interactive with safety checks |

### Workflow Integration

```bash
# Traditional workflow
git worktree add ../feature-auth feature/auth
cd ../feature-auth
# work...
cd ../main

# wrkt workflow  
wrkt add feature/auth        # Creates worktrees/feature-auth/
wrkt switch feature-auth     # Navigate there
# work...
wrkt switch main            # Return to main worktree
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

**Problem:** `wrkt switch` says "command not found"  
**Solution:** Set up zsh integration: `eval "$(wrkt shell-init)"`

### Directory doesn't change

**Problem:** `wrkt switch` runs but directory doesn't change
**Solution:** Ensure zsh integration is loaded in your current shell session

### Worktree not found

**Problem:** `wrkt switch` can't find worktree
**Solution:** Use `wrkt list` to see exact names, check spelling

### Completion not working

**Problem:** Tab completion doesn't work
**Solution:** Restart zsh after adding shell integration to profile

### Permission denied creating worktree

**Problem:** Cannot create in worktrees/ directory
**Solution:** Check repository permissions, ensure you're in git repository root