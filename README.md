# wrkt - Git Worktree Made Easy

A CLI tool that makes working with multiple git worktrees seamless by organizing them in a predictable structure and providing simple navigation.

## Problem

Git worktrees are powerful but cumbersome:
- Manual `cd` between worktree directories
- Hard to discover existing worktrees and their status  
- Need to remember paths and branch associations
- Worktrees scattered across filesystem

## Solution

`wrkt` organizes all worktrees in `worktrees/` subdirectory with simple commands:

```bash
# Discovery
wrkt list                    # Show all worktrees with status
wrkt list --dirty            # Show only worktrees with changes
wrkt list --verbose          # Detailed status information

# Navigation (zsh only)
wrkt switch feature-auth     # Switch to worktree (exact name)
wrkt switch main             # Switch to main worktree

# Management
wrkt add feature/auth        # Create at worktrees/feature-auth/
wrkt remove feature-auth     # Safe removal with cleanup
wrkt clean                   # Clean up stale worktrees
```

## Installation

```bash
# Install the binary
go install github.com/no-yan/wrkt@latest

# REQUIRED: Add zsh integration for directory switching
eval "$(wrkt shell-init)"

# Add to your zsh profile for permanent setup:
echo 'eval "$(wrkt shell-init)"' >> ~/.zshrc

# Note: Only zsh is supported for shell integration
```

## Quick Start

```bash
# Set up zsh integration (one-time)
eval "$(wrkt shell-init)"

# Create worktrees for parallel development
wrkt add feature/auth        # Creates worktrees/feature-auth/
wrkt add hotfix/bug-123      # Creates worktrees/hotfix-bug-123/

# Switch between worktrees (requires zsh)
wrkt switch feature-auth     # Navigate to worktrees/feature-auth/
# work on feature...
wrkt switch main             # Switch to main worktree

# See all worktrees with status
wrkt list
# Output:
# ✓ main              ~/project/worktrees/main              [main]
# * feature-auth      ~/project/worktrees/feature-auth      [feature/auth] 
# ↑ hotfix-bug-123    ~/project/worktrees/hotfix-bug-123    [hotfix/bug-123]

# Show only worktrees with uncommitted changes
wrkt list --dirty

# Detailed status (like git status) for all worktrees
wrkt list --verbose

# Clean up when done
wrkt remove feature-auth
wrkt clean
```

## Features

- **Organized Structure**: All worktrees in predictable `worktrees/` subdirectory
- **Zsh Integration**: True directory switching via zsh functions
- **Simple Navigation**: Exact name matching with tab completion
- **Auto-Path Generation**: Intelligent path creation from branch names
- **Unified Status Display**: All worktree information in one command with flexible filtering
- **Safe Management**: Prevent accidental deletion with confirmation prompts

## Documentation

- [Command Reference](docs/api/commands.md) - Complete command documentation
- [Implementation Guide](docs/implementation/architecture.md) - Technical details for contributors
- [Examples](docs/examples/) - Common workflows and use cases

## Development Status

**Current**: MVP in development
**Next**: Shell integration, configuration, advanced features

See [CLAUDE.md](CLAUDE.md) for implementation context and [docs/implementation/mvp.md](docs/implementation/mvp.md) for MVP scope.

## Contributing

1. Read [docs/implementation/contributing.md](docs/implementation/contributing.md)
2. Check [docs/implementation/testing.md](docs/implementation/testing.md) for testing guidelines
3. Submit PRs with tests and documentation updates

## License

MIT License - see [LICENSE](LICENSE) file for details.