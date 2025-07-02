# Common Workflows

Real-world examples of using `wrkt` for daily development tasks.

## Feature Development Workflow

### Scenario: Working on a new authentication feature

```bash
# Start in main repository
cd ~/projects/myapp
git checkout main
git pull origin main

# Create feature worktree
wrkt add feature/auth-service
# Creates ../feature-auth-service/ with feature/auth-service branch

# Switch to feature worktree
wrkt switch auth
# Now in ~/projects/feature-auth-service/

# Work on the feature
echo "auth code" >> auth.go
git add .
git commit -m "implement basic auth service"

# Switch back to main to check something
wrkt switch -
# Now back in ~/projects/myapp/

# Continue feature work
wrkt switch auth
# Back in feature worktree

# When feature is complete, clean up
wrkt remove auth
# Removes the worktree but keeps the branch for PR
```

## Hotfix Workflow

### Scenario: Critical bug needs immediate fix

```bash
# Currently working on feature branch
pwd  # ~/projects/myapp-feature-xyz/

# Critical bug reported - need hotfix
wrkt add hotfix/security-patch main
# Creates worktree from main branch

wrkt switch hotfix
# Now in ~/projects/hotfix-security-patch/

# Make the fix
vim security.go
git add .
git commit -m "fix: security vulnerability in auth"

# Test the fix
make test

# Push for immediate deployment
git push origin hotfix/security-patch

# Switch back to feature work
wrkt switch -
# Back to feature worktree

# Later, clean up hotfix worktree
wrkt remove hotfix
```

## Multi-Feature Development

### Scenario: Working on multiple features simultaneously

```bash
# Set up multiple feature worktrees
wrkt add feature/user-management
wrkt add feature/api-redesign  
wrkt add feature/performance-improvements

# Check all worktrees
wrkt list
# ✓ main                     /path/main              [main]
# * user-management          /path/user-management   [feature/user-management]
# ✓ api-redesign            /path/api-redesign      [feature/api-redesign]
# ↑ performance-improvements /path/performance       [feature/performance-improvements]

# Work on user management
wrkt switch user
# Work on code...

# Switch to API redesign
wrkt switch api
# Work on different feature...

# Quick check on performance improvements
wrkt switch perf  # Fuzzy match works
# Review performance changes...

# See status across all worktrees
wrkt status
# Shows git status for all dirty worktrees

# Clean up completed features
wrkt remove user  # Feature merged and deployed
wrkt remove api   # Feature completed
```

## Release Preparation

### Scenario: Preparing for version release

```bash
# Create release preparation worktree
wrkt add release/v1.2.0 main

wrkt switch release
# Now in ~/projects/release-v1.2.0/

# Update version files
vim version.go
vim package.json
vim CHANGELOG.md

# Run full test suite
make test-all

# Build release artifacts
make build-release

# Create release commit
git add .
git commit -m "prepare release v1.2.0"

# Tag the release
git tag v1.2.0

# Push release
git push origin release/v1.2.0
git push origin v1.2.0

# Switch back to main work
wrkt switch main

# Clean up release worktree after deployment
wrkt remove release
```

## Code Review Workflow

### Scenario: Reviewing colleague's pull request

```bash
# Create worktree for PR review
wrkt add review/pr-123 origin/feature/new-dashboard

wrkt switch review
# Now in ~/projects/review-pr-123/

# Review the changes
git log --oneline main..HEAD
git diff main..HEAD

# Test the changes
make test
make build

# Try the feature locally
./app --enable-new-dashboard

# Make review notes
echo "Review notes" > REVIEW.md

# Switch back to your work
wrkt switch -

# Clean up after review is complete
wrkt remove review
```

## Bug Investigation

### Scenario: Investigating reported bug

```bash
# Current work in progress
wrkt list
# * main           /path/main         [feature/current-work]

# Create investigation worktree from production tag
wrkt add investigation/bug-report v1.1.5

wrkt switch investigation
# Now in clean state matching production

# Reproduce the bug
./reproduce-bug.sh

# Create minimal test case
echo "test case" > bug-test.go

# Try potential fixes
git checkout -b bugfix/issue-456
# Make changes...

# Verify fix
make test

# Switch back to main work
wrkt switch -

# Later: create proper fix in separate worktree
wrkt add bugfix/issue-456
wrkt switch bugfix
# Implement proper fix...

# Clean up investigation
wrkt remove investigation
```

## Experimental Development

### Scenario: Trying different approaches to a problem

```bash
# Create multiple experimental worktrees
wrkt add experiment/approach-a feature/base
wrkt add experiment/approach-b feature/base  
wrkt add experiment/approach-c feature/base

# Try first approach
wrkt switch approach-a
# Implement solution A...
git commit -m "experiment: implement approach A"

# Try second approach
wrkt switch approach-b
# Implement solution B...
git commit -m "experiment: implement approach B"

# Try third approach  
wrkt switch approach-c
# Implement solution C...
git commit -m "experiment: implement approach C"

# Compare approaches
wrkt list
# See which experiments are dirty/clean

# Benchmark different solutions
wrkt switch approach-a && make benchmark
wrkt switch approach-b && make benchmark
wrkt switch approach-c && make benchmark

# Choose best approach and continue development
wrkt switch approach-b  # Best performance
# Continue with chosen approach...

# Clean up failed experiments
wrkt remove approach-a
wrkt remove approach-c
# Keep approach-b as the main feature branch
```

## Conference/Demo Preparation

### Scenario: Preparing demo for conference presentation

```bash
# Create demo-specific worktree
wrkt add demo/conference-2024 main

wrkt switch demo
# Now in ~/projects/demo-conference-2024/

# Set up demo data
./scripts/setup-demo-data.sh

# Create demo-specific features
echo "demo features" >> demo.go
git add .
git commit -m "add demo-specific features"

# Test demo flow
./run-demo.sh

# Package demo for conference
make package-demo

# After conference, clean up
wrkt switch main
wrkt remove demo
```

## Maintenance Workflows

### Regular Cleanup

```bash
# Weekly cleanup routine
wrkt list
# Review all worktrees

# Clean up stale worktrees
wrkt clean
# Removes orphaned/prunable worktrees

# Remove completed feature worktrees
wrkt remove old-feature-1
wrkt remove old-feature-2

# Keep only active worktrees
wrkt list
# Should show only current work
```

### Status Monitoring

```bash
# Check status across all worktrees
wrkt status
# Shows uncommitted changes in all worktrees

# Full status check
wrkt list --verbose
# Shows detailed information including lock status

# Find worktrees that need attention
wrkt list | grep "*"  # Dirty worktrees
wrkt list | grep "↑"  # Ahead of remote
```

## Integration with Other Tools

### With Git Hooks

```bash
# Pre-commit hook that runs across all worktrees
#!/bin/sh
# .git/hooks/pre-commit

# Check all worktrees for style issues
wrkt status --short | while read line; do
  path=$(echo $line | cut -d' ' -f2)
  (cd $path && make lint) || exit 1
done
```

### With IDE/Editor

```bash
# VS Code workspace for all worktrees
wrkt list --json | jq -r '.[] | .path' | while read path; do
  echo "Adding $path to workspace"
  code --add "$path"
done
```

### With CI/CD

```bash
# Deploy script that handles multiple worktrees
#!/bin/bash

# Deploy main branch
wrkt switch main
git pull origin main
make deploy-production

# Deploy staging branches
for worktree in $(wrkt list --format="{{.Name}}" | grep "staging/"); do
  wrkt switch "$worktree"
  make deploy-staging
done
```

## Tips and Best Practices

### Naming Conventions

```bash
# Use consistent prefixes
wrkt add feature/auth-service
wrkt add bugfix/login-issue
wrkt add hotfix/security-patch
wrkt add experiment/new-algorithm
wrkt add review/pr-123
wrkt add demo/conference-2024
```

### Path Organization

```bash
# Default auto-paths work well:
# feature/auth-service    → ../feature-auth-service/
# bugfix/login-issue      → ../bugfix-login-issue/  
# hotfix/security-patch   → ../hotfix-security-patch/

# But you can customize when needed:
wrkt add feature/auth /tmp/auth-work  # Temporary location
wrkt add demo/v2 ./demo-v2           # Relative to current dir
```

### Switching Efficiency

```bash
# Use fuzzy matching for speed
wrkt switch auth     # Matches feature/auth-service
wrkt switch bug      # Matches bugfix/login-issue
wrkt switch -        # Quick return to previous

# Create aliases for frequently used worktrees
alias wtmain='wrkt switch main'
alias wtfeat='wrkt switch feature'
```

### Status Monitoring

```bash
# Add to your shell prompt
export PS1='$(wrkt list --current --format="{{.Name}}") $ '

# Or use in scripts
current_worktree=$(wrkt list --current --format="{{.Name}}")
echo "Working in: $current_worktree"
```

These workflows demonstrate how `wrkt` integrates into real development scenarios, making multi-worktree development natural and efficient.