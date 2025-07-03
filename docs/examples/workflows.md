# Common Workflows

Real-world examples of using `wrkt` for daily development tasks.

## Quick Start

This comprehensive example demonstrates the core `wrkt` commands in a realistic development scenario:

```bash
# Start in main repository
cd ~/projects/myapp
git checkout main
git pull origin main

# 1. CREATE: Set up multiple parallel worktrees
wrkt add feature/auth-service           # New feature from main
wrkt add hotfix/security-patch main     # Critical fix from main  
wrkt add review/pr-123 origin/feature/dashboard  # PR review from remote branch

# 2. LIST: Check all worktrees
wrkt list
#   main            /path/main           [main]
#   auth-service     /path/auth-service   [feature/auth-service]
#   security-patch   /path/security-patch [hotfix/security-patch]  
# * pr-123          /path/pr-123         [feature/dashboard]

# 3. SWITCH: Work on different tasks
wrkt switch auth
# Now in ~/projects/feature-auth-service/

# Work on feature
echo "auth code" >> auth.go
git add .
git commit -m "implement basic auth service"

# 4. SWITCH: Handle urgent hotfix
wrkt switch hotfix  # Fuzzy matching finds security-patch
# Now in ~/projects/hotfix-security-patch/

# Make critical fix
vim security.go
git add .
git commit -m "fix: security vulnerability"
git push origin hotfix/security-patch

# 5. SWITCH: Review colleague's PR
wrkt switch review  # Fuzzy matching finds pr-123
# Now in ~/projects/review-pr-123/

# Test the PR
make test
git log --oneline main..HEAD

# Switch back to feature work
wrkt switch -  # Returns to previous (hotfix)
wrkt switch auth  # Back to feature work

# 6. REMOVE: Clean up completed work
wrkt remove hotfix    # Hotfix deployed
wrkt remove review    # PR review complete
# Keep feature worktree for continued development

# Final state
wrkt list
#   main            /path/main           [main]
# * auth-service     /path/auth-service   [feature/auth-service]
```

## Advanced Scenarios

### Experimental Development
When you need to try multiple approaches to solve a problem:

```bash
# Create multiple experimental worktrees from same base
wrkt add experiment/approach-a feature/base
wrkt add experiment/approach-b feature/base  
wrkt add experiment/approach-c feature/base

# Test each approach independently
wrkt switch approach-a && make benchmark
wrkt switch approach-b && make benchmark  
wrkt switch approach-c && make benchmark

# Compare results and keep the best
wrkt remove approach-a approach-c  # Remove failed experiments
# Continue with approach-b as main feature
```

### Long-Running Feature Development
For features that take weeks/months with frequent main branch updates:

```bash
# Create long-running feature worktree
wrkt add feature/new-architecture main

# Periodically sync with main to avoid conflicts
wrkt switch new-architecture
git fetch origin
git rebase origin/main  # Keep feature current with main

# Work continues in isolated worktree
# Main development continues unaffected in other worktrees
```

## Tips and Best Practices

### Naming Conventions
```bash
# Use consistent prefixes for easy identification
wrkt add feature/auth-service
wrkt add bugfix/login-issue  
wrkt add hotfix/security-patch
wrkt add experiment/new-algorithm
wrkt add review/pr-123
```

### Efficient Switching
```bash
# Use fuzzy matching for speed
wrkt switch auth     # Matches feature/auth-service
wrkt switch -        # Quick return to previous worktree

# Create shell aliases for common operations
alias wtmain='wrkt switch main'
alias wtlist='wrkt list'
```

### Status Monitoring
```bash
# Check all worktrees for uncommitted changes
wrkt list

# Find specific worktree types
wrkt list | grep "feature/"    # All feature worktrees
wrkt list | grep "*"           # Current worktree
wrkt list | grep "!"           # Dirty worktrees
```

These workflows demonstrate how `wrkt` integrates into real development scenarios, making parallel development efficient and organized.
