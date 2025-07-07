# Common Workflows

Real-world examples of using `wt` for daily development tasks.

## Quick Start

This comprehensive example demonstrates the core `wt` commands in a realistic development scenario:

```bash
# Start in main repository
cd ~/projects/myapp
git checkout main
git pull origin main

# 1. CREATE: Set up multiple parallel worktrees
wt add feature/auth-service           # New feature from main
wt add hotfix/security-patch main     # Critical fix from main  
wt add review/pr-123 origin/feature/dashboard  # PR review from remote branch

# 2. LIST: Check all worktrees
wt list
#   main            /path/main           [main]
#   auth-service     /path/auth-service   [feature/auth-service]
#   security-patch   /path/security-patch [hotfix/security-patch]  
# * pr-123          /path/pr-123         [feature/dashboard]

# 3. SWITCH: Work on different tasks
wt switch auth
# Now in ~/projects/feature-auth-service/

# Work on feature
echo "auth code" >> auth.go
git add .
git commit -m "implement basic auth service"

# 4. SWITCH: Handle urgent hotfix
wt switch hotfix  # Fuzzy matching finds security-patch
# Now in ~/projects/hotfix-security-patch/

# Make critical fix
vim security.go
git add .
git commit -m "fix: security vulnerability"
git push origin hotfix/security-patch

# 5. SWITCH: Review colleague's PR
wt switch review  # Fuzzy matching finds pr-123
# Now in ~/projects/review-pr-123/

# Test the PR
make test
git log --oneline main..HEAD

# Switch back to feature work
wt switch -  # Returns to previous (hotfix)
wt switch auth  # Back to feature work

# 6. REMOVE: Clean up completed work
wt remove hotfix    # Hotfix deployed
wt remove review    # PR review complete
# Keep feature worktree for continued development

# Final state
wt list
#   main            /path/main           [main]
# * auth-service     /path/auth-service   [feature/auth-service]
```

## Advanced Scenarios

### Experimental Development
When you need to try multiple approaches to solve a problem:

```bash
# Create multiple experimental worktrees from same base
wt add experiment/approach-a feature/base
wt add experiment/approach-b feature/base  
wt add experiment/approach-c feature/base

# Test each approach independently
wt switch approach-a && make benchmark
wt switch approach-b && make benchmark  
wt switch approach-c && make benchmark

# Compare results and keep the best
wt remove approach-a approach-c  # Remove failed experiments
# Continue with approach-b as main feature
```

### Long-Running Feature Development
For features that take weeks/months with frequent main branch updates:

```bash
# Create long-running feature worktree
wt add feature/new-architecture main

# Periodically sync with main to avoid conflicts
wt switch new-architecture
git fetch origin
git rebase origin/main  # Keep feature current with main

# Work continues in isolated worktree
# Main development continues unaffected in other worktrees
```

## Tips and Best Practices

### Naming Conventions
```bash
# Use consistent prefixes for easy identification
wt add feature/auth-service
wt add bugfix/login-issue  
wt add hotfix/security-patch
wt add experiment/new-algorithm
wt add review/pr-123
```

### Efficient Switching
```bash
# Use fuzzy matching for speed
wt switch auth     # Matches feature/auth-service
wt switch -        # Quick return to previous worktree

# Create shell aliases for common operations
alias wtmain='wt switch main'
alias wtlist='wt list'
```

### Status Monitoring
```bash
# Check all worktrees for uncommitted changes
wt list

# Find specific worktree types
wt list | grep "feature/"    # All feature worktrees
wt list | grep "*"           # Current worktree
wt list | grep "!"           # Dirty worktrees
```

These workflows demonstrate how `wt` integrates into real development scenarios, making parallel development efficient and organized.
