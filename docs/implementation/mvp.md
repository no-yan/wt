# MVP (Minimum Viable Product) Specification

This document defines the MVP scope for `wt` to ensure focused development and clear success criteria.

## MVP Goals

1. **Eliminate manual directory navigation** between git worktrees
2. **Provide intuitive worktree discovery** with enhanced status display
3. **Enable smart worktree creation** with auto-path generation
4. **Ensure safe worktree management** with proper error handling

## MVP Features

### ✅ Core Commands

#### 1. `wt list` (Unified Status Display)
- Display all worktrees with flexible filtering
- Multiple view modes: default, dirty-only, verbose, names-only
- Show status indicators and detailed information
- Support for bare repositories

**Acceptance Criteria:**
- Lists all worktrees from `git worktree list`
- Shows clean (✓) vs dirty (*) status indicators
- `--dirty` flag shows only worktrees with changes
- `--verbose` shows detailed git status output
- `--names-only` outputs just names for shell completion
- Displays in human-readable format
- Handles empty worktree list gracefully

#### 2. `wt switch <name>` (Shell Integration Required)
- Switch to worktree directory with exact name matching
- Requires shell integration setup via `wt shell-init`
- Support exact name matching only
- Change current working directory through shell functions

**Acceptance Criteria:**
- **Shell integration setup is mandatory**
- Exact name matching finds worktrees by directory name
- Changes to worktree directory on success (via shell function)
- Clear error message for non-existent worktrees with suggestions
- Deterministic behavior with exact matching

#### 3. `wt shell-init` (Critical Infrastructure)
- Generate shell integration code for directory switching
- Support zsh only (eliminates multi-shell complexity)
- Include tab completion for all commands
- Handle shell function wrapper logic

**Acceptance Criteria:**
- Generates working shell functions for zsh
- Provides comprehensive tab completion
- Handles `wt switch` interception correctly
- Easy setup with `eval "$(wt shell-init)"`
- Clear error message for non-zsh shells

#### 4. `wt add <branch> [path]`
- Create new worktree for specified branch
- Auto-generate path if not provided
- Handle existing branches and new branch creation

**Acceptance Criteria:**
- Creates worktree in auto-generated path
- Accepts custom path when provided
- Auto-path follows naming convention
- Validates branch existence and availability
- Handles path conflicts with unique naming

#### 5. `wt remove <name>`
- Remove worktree safely with exact name matching
- Confirm removal of dirty worktrees
- Prevent removal of main worktree

**Acceptance Criteria:**
- Exact name matching finds worktree to remove
- Confirms removal of dirty worktrees
- Prevents main worktree removal
- Cleans up administrative files

### ✅ Smart Features

#### 1. Shell Integration (Mandatory)
- Shell function generation for zsh only
- Tab completion for commands and worktree names
- Seamless integration like zoxide/autojump

#### 2. Unified Status Display
- Single `list` command with multiple view modes
- Filtering options replace separate status command
- Consistent interface reduces cognitive load

#### 3. Exact Name Matching
- Match worktrees by exact directory name
- Deterministic, predictable behavior
- Error suggestions when no match found

**Matching Rules:**
1. Exact match on directory basename only
2. Case-sensitive matching
3. No partial or fuzzy matching complexity

#### 4. Auto-Path Generation
- `feature/auth` → `../feature-auth/`
- `hotfix/bug-123` → `../hotfix-bug-123/`
- `docs/api-update` → `../docs-api-update/`
- `main` → `../main/`

**Rules:**
- Remove common prefixes (feature/, hotfix/, bugfix/, docs/)
- Replace slashes with dashes
- Create in parent directory of repository
- Ensure path uniqueness with numbering if needed

#### 5. Status Indicators
- `✓` - Clean worktree (no uncommitted changes)
- `*` - Dirty worktree (uncommitted changes)
- `↑` - Ahead of remote
- `↓` - Behind remote
- `L` - Locked worktree
- `P` - Prunable worktree

### ✅ Error Handling

#### 1. Git Repository Validation
- Detect if current directory is in git repository
- Provide helpful error message if not

#### 2. Worktree State Validation
- Handle locked worktrees
- Handle missing worktree directories
- Handle invalid worktree states

#### 3. User Input Validation
- Validate command arguments
- Provide helpful error messages
- Suggest corrections for typos

## MVP Constraints

### ❌ Out of Scope

#### 1. GitHub Integration
- Pull request status display
- Review state indicators
- CI/CD status integration

#### 2. Configuration System
- Configuration files (.wt.yaml)
- Custom path templates
- User preference storage

#### 3. Batch Operations
- Multiple worktree operations at once
- Bulk cleanup commands
- Scripting automation features

#### 4. Advanced Git Features
- Submodule support
- Sparse checkout integration
- Git hooks integration
- Remote synchronization

#### 5. Claude Development Tracking
- Development session management
- Progress tracking
- Status persistence

#### 6. Advanced Display Features
- Commit information display
- Detailed remote tracking
- Performance metrics

## Success Criteria

### Functional Requirements

1. **All core commands work correctly**
   - `wt list` shows all worktrees with status
   - `wt switch` navigates to correct worktree
   - `wt add` creates worktree in correct location
   - `wt remove` safely removes worktree

2. **Exact matching is predictable**
   - Deterministic behavior with exact names
   - Error messages are helpful with suggestions

3. **Auto-path generation is logical**
   - Generated paths are predictable
   - No path conflicts occur
   - Paths are human-readable

4. **Error handling is robust**
   - No crashes on invalid input
   - Clear error messages with suggestions
   - Graceful handling of git repository issues

### Non-Functional Requirements

1. **Performance**
   - Commands execute in <100ms for typical repositories
   - List command handles 20+ worktrees efficiently
   - Exact matching is instant

2. **Usability**
   - Commands are intuitive for git users
   - Help text is clear and complete
   - Error messages are actionable

3. **Reliability**
   - No data loss on any operation
   - Consistent behavior across platforms
   - Handles edge cases gracefully

## Quality Assurance

### Testing Requirements

1. **Unit Tests**
   - Worktree parsing logic
   - Exact name matching functions
   - Path generation functions
   - Error handling scenarios

2. **Integration Tests**
   - Git command integration
   - File system operations
   - End-to-end command flows

3. **Manual Testing**
   - Real git repository scenarios
   - Cross-platform compatibility
   - Edge cases and error conditions

### Code Quality

1. **Go Best Practices**
   - Proper error handling
   - Clear function signatures
   - Comprehensive documentation
   - Consistent code style

2. **Git Integration**
   - Use porcelain commands for parsing
   - Handle all git worktree states
   - Validate git command output

3. **User Experience**
   - Consistent command patterns
   - Clear help documentation
   - Predictable behavior

## MVP Validation

## Ordinary Development Workflow

### Daily Feature Development

**Scenario**: Starting work on a new authentication feature

```bash
# 1. Setup (one-time - zsh required)
eval "$(wt shell-init)"  # Add to ~/.zshrc for permanent setup

# 2. Start from main repository
cd ~/projects/myapp
git checkout main
git pull origin main

# 3. Create feature worktree
wt add feature/auth-service
# → Creates worktrees/feature-auth-service/ with feature/auth-service branch

# 4. Switch to feature worktree
wt switch feature-auth-service
# → Now in ~/projects/myapp/worktrees/feature-auth-service/
pwd  # /Users/you/projects/myapp/worktrees/feature-auth-service

# 5. Develop the feature
echo "auth code" >> auth.go
git add .
git commit -m "implement basic auth service"

# 6. Check overall status
wt list
# ✓ main                    ~/projects/myapp/worktrees/main              [main]
# * feature-auth-service    ~/projects/myapp/worktrees/feature-auth-service  [feature/auth-service]

# 7. Switch back to main to check something
wt switch main
# → Back in ~/projects/myapp/worktrees/main

# 8. Return to feature work
wt switch feature-auth-service
# → Back in feature worktree

# 9. Push feature when ready
git push origin feature/auth-service

# 10. Clean up after merge
wt switch main
wt remove feature-auth-service
```

### Multi-Feature Development

**Scenario**: Working on multiple features simultaneously

```bash
# Set up multiple features (exact names required)
wt add feature/user-management
wt add feature/api-redesign
wt add hotfix/security-patch

# Check all worktrees
wt list
# ✓ main                      ~/projects/myapp/worktrees/main                 [main]
# * feature-user-management   ~/projects/myapp/worktrees/feature-user-mgmt    [feature/user-management]
# ✓ feature-api-redesign     ~/projects/myapp/worktrees/feature-api-redesign [feature/api-redesign]
# ↑ hotfix-security-patch    ~/projects/myapp/worktrees/hotfix-security      [hotfix/security-patch]

# Work on user management (exact name required)
wt switch feature-user-management
# Develop feature...

# Quick switch to API redesign
wt switch feature-api-redesign
# Work on different feature...

# See what needs attention
wt list --dirty
# * feature-user-management   ~/projects/myapp/worktrees/feature-user-mgmt   [feature/user-management]

# Get detailed status
wt list --verbose
# * feature-user-management (~/projects/myapp/worktrees/feature-user-mgmt) [feature/user-management]
#    M user.go
#   ?? test.go
# 
# ✓ feature-api-redesign (~/projects/myapp/worktrees/feature-api-redesign) [feature/api-redesign]
```

### Hotfix Workflow

**Scenario**: Critical bug needs immediate attention

```bash
# Currently working on feature
wt switch feature-xyz

# Critical bug reported - create hotfix from main
wt add hotfix/security-vulnerability main
wt switch hotfix

# Make the fix
vim security.go
git add .
git commit -m "fix: security vulnerability"

# Test and push
make test
git push origin hotfix/security-vulnerability

# Return to feature work
wt switch -  # Back to feature-xyz

# Clean up hotfix after deployment
wt remove hotfix
```

### Code Review Workflow

**Scenario**: Reviewing colleague's pull request

```bash
# Create review worktree from PR branch
wt add review/pr-456 origin/feature/new-dashboard

# Switch to review
wt switch review

# Review changes
git log --oneline main..HEAD
git diff main..HEAD
make test

# Switch back to work
wt switch -

# Clean up after review
wt remove review
```

### End of Day Cleanup

**Scenario**: Regular maintenance

```bash
# Check all worktrees
wt list

# See what has uncommitted changes
wt list --dirty

# Clean up completed features
wt remove old-feature
wt remove completed-task

# Clean up stale worktrees
wt clean
```

### Acceptance Test Scenarios

1. **Basic Workflow**
   ```bash
   # Setup
   cd /path/to/git/repo
   eval "$(wt shell-init)"
   
   # Create worktrees
   wt add feature/auth
   wt add hotfix/bug-123
   
   # List worktrees
   wt list
   # Should show main, feature-auth, hotfix-bug-123
   
   # Switch between worktrees
   wt switch auth
   pwd # Should be in feature-auth directory
   
   # Remove worktree
   wt remove auth
   wt list
   # Should not show feature-auth
   ```

2. **Zsh Integration**
   ```bash
   # Without zsh integration
   wt switch feature-auth  # Should show setup instructions
   
   # With zsh integration
   eval "$(wt shell-init)"
   wt switch feature-auth  # Should actually change directory
   pwd  # Should be in feature-auth worktree directory
   ```

3. **Exact Name Matching**
   ```bash
   wt add feature/authentication-service
   wt switch feature-authentication-service  # Exact name required
   wt switch auth  # Should show "worktree not found" with suggestions
   ```

4. **Error Handling**
   ```bash
   cd /tmp  # Not a git repository
   wt list  # Should show helpful error
   
   cd /path/to/git/repo
   wt switch nonexistent  # Should show error with suggestions
   wt remove main  # Should prevent removal
   ```

### Performance Benchmarks

- `wt list` with 10 worktrees: <50ms
- `wt switch` with exact matching: <10ms
- `wt add` with auto-path: <200ms
- `wt remove` with confirmation: <100ms

## MVP Timeline

### Phase 1: Core Implementation (Week 1)
- Set up Go project structure
- Implement worktree parsing
- Create basic command framework

### Phase 2: Command Implementation (Week 2)
- Implement list command
- Implement switch command
- Implement add command
- Implement remove command

### Phase 3: Smart Features (Week 3)
- Implement exact name matching
- Implement auto-path generation
- Add status indicators
- Enhance error handling

### Phase 4: Testing & Polish (Week 4)
- Write comprehensive tests
- Manual testing across platforms
- Documentation updates
- Bug fixes and improvements

## Success Metrics

### Quantitative
- All acceptance tests pass
- Code coverage >80%
- Performance benchmarks met
- Zero critical bugs

### Qualitative
- Commands feel intuitive to git users
- Error messages are helpful
- Documentation is clear and complete
- Code is maintainable and extensible

## Post-MVP Roadmap

### v1.1 - Shell Integration
- Shell functions for seamless switching
- Completion scripts
- Environment variable support

### v1.2 - Enhanced Status
- Ahead/behind indicators
- Detailed status aggregation
- Performance optimizations

### v1.3 - Configuration
- Configuration file support
- Custom path templates
- User preferences

### v2.0 - Advanced Features
- Batch operations
- Submodule support
- Plugin architecture