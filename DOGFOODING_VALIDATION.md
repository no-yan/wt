# Dogfooding Validation Procedures

**Purpose**: Prevent recurring `git checkout -b` violations through systematic validation and automated compliance checking.

## 🚨 CRITICAL: Pre-Command Validation

**RUN THIS BEFORE EVERY GIT OPERATION**

### Universal Validation Script

```bash
#!/bin/bash
# File: validate_dogfooding.sh
# Usage: source validate_dogfooding.sh && validate_command "git checkout main"

validate_command() {
    local command="$1"
    local current_dir=$(pwd)
    local current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    
    echo "🔍 DOGFOODING VALIDATION"
    echo "📍 Directory: $current_dir"
    echo "🌿 Branch: $current_branch"
    echo "💻 Command: $command"
    echo "---"
    
    # Rule 1: Check for prohibited commands
    if [[ "$command" =~ "checkout -b" ]] || [[ "$command" =~ "checkout -B" ]]; then
        echo "❌ VIOLATION: git checkout -b is prohibited!"
        echo "🔧 Use instead: ./wrkt add <branch-name>"
        echo "📖 Reason: Creates branches in main directory (violates dogfooding)"
        return 1
    fi
    
    if [[ "$command" =~ "switch -c" ]] || [[ "$command" =~ "switch -C" ]]; then
        echo "❌ VIOLATION: git switch -c is prohibited!"
        echo "🔧 Use instead: ./wrkt add <branch-name>"
        echo "📖 Reason: Creates branches in main directory (violates dogfooding)"
        return 1
    fi
    
    # Rule 2: Check for branch switching in main directory
    if [[ "$current_dir" == */wrkt ]] && [[ "$current_dir" != */worktrees/* ]]; then
        if [[ "$command" =~ "checkout" ]] && [[ ! "$command" =~ "checkout main" ]]; then
            echo "❌ VIOLATION: Branch switching in main directory!"
            echo "🔧 Use instead: ./wrkt switch <worktree-name>"
            echo "📖 Reason: Main directory must stay on main branch"
            return 1
        fi
        
        if [[ "$command" =~ "switch" ]] && [[ ! "$command" =~ "switch main" ]]; then
            echo "❌ VIOLATION: Branch switching in main directory!"
            echo "🔧 Use instead: ./wrkt switch <worktree-name>"
            echo "📖 Reason: Main directory must stay on main branch"
            return 1
        fi
    fi
    
    # Rule 3: Check for branch switching within worktrees
    if [[ "$current_dir" == */worktrees/* ]]; then
        if [[ "$command" =~ "checkout" ]] && [[ ! "$command" =~ "checkout -" ]] && [[ ! "$command" =~ "checkout \." ]]; then
            echo "❌ VIOLATION: Branch switching within worktree!"
            echo "🔧 Use instead: Work within current worktree or create new one"
            echo "📖 Reason: 1 worktree = 1 dedicated branch"
            return 1
        fi
        
        if [[ "$command" =~ "switch" ]]; then
            echo "❌ VIOLATION: Branch switching within worktree!"
            echo "🔧 Use instead: Work within current worktree or create new one"
            echo "📖 Reason: 1 worktree = 1 dedicated branch"
            return 1
        fi
    fi
    
    # Rule 4: Check if main directory is on wrong branch
    if [[ "$current_dir" == */wrkt ]] && [[ "$current_dir" != */worktrees/* ]]; then
        if [[ "$current_branch" != "main" ]]; then
            echo "🚨 CRITICAL VIOLATION DETECTED: Main directory on wrong branch!"
            echo "   Current branch: $current_branch"
            echo "   Expected branch: main"
            echo "   This indicates a previous dogfooding violation!"
            echo ""
            echo "🔧 IMMEDIATE FIX REQUIRED:"
            echo "   git checkout main"
            echo ""
            echo "🛡️  FUTURE PREVENTION:"
            echo "   Use ./wrkt add <branch> for new features"
            echo "   Use ./wrkt switch <name> for navigation"
            return 1
        fi
    fi
    
    echo "✅ VALIDATION PASSED: Command is safe to execute"
    return 0
}

# Interactive validation function
interactive_validate() {
    echo "🔍 DOGFOODING COMMAND VALIDATOR"
    echo "Enter your intended git command (or 'quit' to exit):"
    read -p "> " user_command
    
    if [[ "$user_command" == "quit" ]]; then
        echo "Exiting validator"
        return 0
    fi
    
    if validate_command "$user_command"; then
        echo ""
        read -p "🚀 Execute this command? (y/N): " confirm
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            echo "Executing: $user_command"
            eval "$user_command"
        else
            echo "Command cancelled"
        fi
    else
        echo ""
        echo "❌ Command rejected for safety"
    fi
}

# Alias for easy access
alias check-git='validate_command'
alias git-safe='interactive_validate'
```

### Quick Validation One-Liner

```bash
# Add this to your shell profile for quick access
validate_git() {
    [[ $(pwd) == */wrkt ]] && [[ $(git branch --show-current) != "main" ]] && {
        echo "❌ Main directory on wrong branch!"; return 1
    }
    [[ "$1" =~ "checkout -b" ]] && {
        echo "❌ Use ./wrkt add instead!"; return 1
    }
    echo "✅ Safe to proceed"
}

# Usage: validate_git "git checkout main" && git checkout main
```

## 🏥 Recovery Procedures

### Recovery 1: Main Directory on Wrong Branch

**Symptoms:**
- `git branch --show-current` in main directory shows non-main branch
- `./wrkt list --verbose` shows main worktree on wrong branch

**Recovery Steps:**
```bash
# 1. Navigate to main directory
cd /Users/noyan/ghq/github.com/no-yan/wrkt

# 2. Check current state
echo "Current branch: $(git branch --show-current)"
echo "Current directory: $(pwd)"

# 3. Check for uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    echo "⚠️  Uncommitted changes detected!"
    git status
    echo ""
    read -p "🤔 Stash changes before fixing? (y/N): " stash_confirm
    if [[ "$stash_confirm" =~ ^[Yy]$ ]]; then
        git stash push -m "Pre-dogfooding-fix stash $(date)"
        echo "✅ Changes stashed"
    fi
fi

# 4. Return to main branch
git checkout main

# 5. Verify fix
if [[ $(git branch --show-current) == "main" ]]; then
    echo "✅ RECOVERY SUCCESSFUL: Main directory back on main branch"
else
    echo "❌ RECOVERY FAILED: Manual intervention required"
    echo "Current branch: $(git branch --show-current)"
fi

# 6. List current worktree state
./wrkt list --verbose
```

### Recovery 2: Worktree Branch Contamination

**Symptoms:**
- Worktree directory shows wrong branch
- `git branch --show-current` in worktree doesn't match expected branch

**Recovery Steps:**
```bash
# 1. Identify the issue
WORKTREE_PATH=$(pwd)
CURRENT_BRANCH=$(git branch --show-current)
EXPECTED_BRANCH="[determine from worktree name]"

echo "🔍 Worktree contamination detected:"
echo "   Path: $WORKTREE_PATH"
echo "   Current branch: $CURRENT_BRANCH"
echo "   Expected branch: $EXPECTED_BRANCH"

# 2. This should never happen in a proper worktree setup
# If it does, the worktree is corrupted and should be recreated

echo "⚠️  Worktree corruption detected - recreation required"

# 3. Save any important changes
if [[ -n $(git status --porcelain) ]]; then
    git stash push -m "Pre-recreation stash $(date)"
    echo "✅ Changes stashed in corrupted worktree"
fi

# 4. Go back to main and recreate
cd /Users/noyan/ghq/github.com/no-yan/wrkt
WORKTREE_NAME=$(basename "$WORKTREE_PATH")

# 5. Remove corrupted worktree
./wrkt remove "$WORKTREE_NAME"

# 6. Recreate clean worktree
./wrkt add "$EXPECTED_BRANCH"

echo "✅ RECOVERY SUCCESSFUL: Worktree recreated cleanly"
```

### Recovery 3: Multiple Violations

**Symptoms:**
- Multiple worktrees on wrong branches
- Confusion about current state
- Uncommitted changes scattered across worktrees

**Recovery Steps:**
```bash
#!/bin/bash
# comprehensive_recovery.sh

echo "🚨 COMPREHENSIVE DOGFOODING RECOVERY"
echo "===================================="

# 1. Survey the damage
echo "📊 CURRENT STATE SURVEY:"
cd /Users/noyan/ghq/github.com/no-yan/wrkt
./wrkt list --verbose

echo ""
echo "📋 MAIN DIRECTORY STATUS:"
echo "   Branch: $(git branch --show-current)"
echo "   Uncommitted: $(git status --porcelain | wc -l) files"

# 2. Fix main directory first
if [[ $(git branch --show-current) != "main" ]]; then
    echo ""
    echo "🔧 FIXING MAIN DIRECTORY:"
    if [[ -n $(git status --porcelain) ]]; then
        git stash push -m "Recovery stash - main directory $(date)"
        echo "   ✅ Stashed changes in main directory"
    fi
    git checkout main
    echo "   ✅ Main directory returned to main branch"
fi

# 3. Check each worktree
echo ""
echo "🔍 WORKTREE VALIDATION:"
for worktree_path in /Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/*/; do
    if [[ -d "$worktree_path" ]]; then
        cd "$worktree_path"
        worktree_name=$(basename "$worktree_path")
        current_branch=$(git branch --show-current)
        echo "   $worktree_name: $current_branch"
        
        # Basic validation - each worktree should be on its dedicated branch
        # More complex validation would require knowing expected branch names
        cd /Users/noyan/ghq/github.com/no-yan/wrkt
    fi
done

echo ""
echo "✅ RECOVERY COMPLETE"
echo "🔮 NEXT STEPS:"
echo "   1. Review ./wrkt list --verbose output"
echo "   2. Verify each worktree is on correct branch"
echo "   3. Recreate any corrupted worktrees"
echo "   4. Use only wrkt commands going forward"
```

## 🛡️ Prevention Measures

### Pre-Command Checklist

Before running ANY git command:

```
□ Run validation script
□ Verify current directory
□ Check current branch
□ Consider wrkt alternative
□ Confirm command necessity
```

### Shell Integration

Add to your `.zshrc` or `.bashrc`:

```bash
# Dogfooding protection
source /Users/noyan/ghq/github.com/no-yan/wrkt/validate_dogfooding.sh

# Override git with validation
git() {
    local cmd="$1"
    shift
    local full_command="git $cmd $*"
    
    # Skip validation for safe read-only commands
    case "$cmd" in
        status|log|show|diff|blame|ls-files)
            command git "$cmd" "$@"
            return $?
            ;;
    esac
    
    # Validate dangerous commands
    if validate_command "$full_command"; then
        command git "$cmd" "$@"
    else
        echo "🚫 Command blocked by dogfooding protection"
        return 1
    fi
}
```

### IDE Integration

For VS Code, add to settings.json:

```json
{
    "git.confirmSync": true,
    "git.confirmEmptyCommits": true,
    "terminal.integrated.shellArgs.osx": [
        "-c", 
        "source ~/.zshrc && source /Users/noyan/ghq/github.com/no-yan/wrkt/validate_dogfooding.sh && zsh"
    ]
}
```

## 🧪 Testing Validation Scripts

### Test Suite for Validation

```bash
#!/bin/bash
# test_validation.sh

echo "🧪 TESTING DOGFOODING VALIDATION"

# Source the validation functions
source validate_dogfooding.sh

# Test cases
test_commands=(
    "git status"                          # Should pass
    "git checkout -b new-feature"         # Should fail
    "git switch -c another-feature"       # Should fail
    "git checkout main"                   # Should pass (if in main dir)
    "git log --oneline"                   # Should pass
    "git checkout feature-branch"         # Should fail (if in main dir)
)

echo "Running validation tests..."

for cmd in "${test_commands[@]}"; do
    echo ""
    echo "Testing: $cmd"
    if validate_command "$cmd"; then
        echo "  Result: ✅ PASSED"
    else
        echo "  Result: ❌ BLOCKED"
    fi
done

echo ""
echo "🏁 Validation testing complete"
```

## 📊 Monitoring and Metrics

### Violation Tracking

```bash
# Log violations for analysis
log_violation() {
    local command="$1"
    local location="$2"
    local timestamp=$(date)
    
    echo "$timestamp,$command,$location,VIOLATION" >> ~/.wrkt_violations.log
}

# Weekly violation report
weekly_report() {
    echo "📊 DOGFOODING VIOLATIONS REPORT"
    echo "Week of: $(date)"
    echo "================================"
    
    if [[ -f ~/.wrkt_violations.log ]]; then
        tail -50 ~/.wrkt_violations.log | grep "$(date +%Y-%m-%d)"
    else
        echo "No violations logged ✅"
    fi
}
```

This validation system creates multiple layers of protection against dogfooding violations:

1. **Pre-command validation** prevents violations before they happen
2. **Recovery procedures** fix violations when they occur
3. **Shell integration** makes validation automatic
4. **Testing** ensures validation works correctly
5. **Monitoring** tracks violations for pattern analysis

By implementing these procedures, the recurring `git checkout -b` failures will be eliminated through systematic prevention rather than relying on memory.