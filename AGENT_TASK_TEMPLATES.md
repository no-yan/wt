# Agent Task Templates

**Purpose**: Standardized, self-contained instruction templates for autonomous Claude agents to prevent dogfooding violations and reduce cognitive load.

## Template Usage Guidelines

1. **Copy template** for specific task
2. **Fill in variables** marked with `{VARIABLE}`
3. **Verify all paths** are correct
4. **Launch single agent** with complete instructions
5. **Wait for completion** before next task

## Template: Feature PR Creation

### Agent Instruction Packet: Feature PR

```markdown
# Autonomous Agent Task: Create PR for {FEATURE_NAME}

## 🔒 MANDATORY PRE-FLIGHT VALIDATION

**YOU MUST RUN THIS FIRST - NO EXCEPTIONS:**

```bash
#!/bin/bash
# Pre-flight validation script
echo "🚀 Starting pre-flight validation..."

# Validate location
EXPECTED_PATH="{WORKTREE_PATH}"
CURRENT_PATH=$(pwd)
if [[ "$CURRENT_PATH" != "$EXPECTED_PATH" ]]; then
    echo "❌ CRITICAL ERROR: Wrong directory!"
    echo "   Expected: $EXPECTED_PATH"
    echo "   Current:  $CURRENT_PATH"
    echo "   Fix: cd $EXPECTED_PATH"
    exit 1
fi

# Validate branch
EXPECTED_BRANCH="{BRANCH_NAME}"
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]]; then
    echo "❌ CRITICAL ERROR: Wrong branch!"
    echo "   Expected: $EXPECTED_BRANCH"
    echo "   Current:  $CURRENT_BRANCH"
    echo "   This indicates a dogfooding violation!"
    exit 1
fi

# Validate git status
if ! git status &>/dev/null; then
    echo "❌ CRITICAL ERROR: Not a git repository!"
    exit 1
fi

echo "✅ Pre-flight validation PASSED"
echo "📍 Location: $CURRENT_PATH"
echo "🌿 Branch: $CURRENT_BRANCH"
echo "🎯 Ready to proceed with task"
```

## 🎯 YOUR MISSION

**Objective**: Create a pull request for {FEATURE_NAME}

**Your Workspace** (NEVER LEAVE THIS DIRECTORY):
- Directory: {WORKTREE_PATH}
- Branch: {BRANCH_NAME}
- Repository: {REPO_NAME}

## 📋 EXECUTION STEPS

### Step 1: Check Current State
```bash
git status
git log --oneline -5
```

### Step 2: Validate Changes
- Review all uncommitted changes
- Ensure changes are related to {FEATURE_NAME}
- If unrelated changes exist, create separate commits

### Step 3: Run All Tests
```bash
go test ./...
```
- If tests fail, fix within this worktree
- NEVER switch branches to fix tests
- Commit fixes incrementally

### Step 4: Commit Changes
```bash
git add .
git commit -m "{COMMIT_MESSAGE}"
```

### Step 5: Push to Remote
```bash
git push -u origin {BRANCH_NAME}
```

### Step 6: Create Pull Request
```bash
gh pr create --title "{PR_TITLE}" --body "$(cat <<'EOF'
## Summary
{PR_SUMMARY}

## Test Plan
- [ ] All unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
)"
```

### Step 7: Notify Completion
```bash
afplay /System/Library/Sounds/Hero.aiff
echo "✅ PR created successfully!"
```

## 🚫 ABSOLUTELY PROHIBITED OPERATIONS

**These commands will cause immediate failure:**
- `git checkout` (any variation)
- `git switch` (any variation)
- `cd` outside of {WORKTREE_PATH}
- Any branch creation/switching operations

## 🤖 AUTONOMOUS DECISION TREE

**Scenario: Tests are failing**
- ✅ Fix tests within current worktree
- ❌ Switch to different branch to fix

**Scenario: Merge conflicts during push**
- ✅ Resolve conflicts within current worktree
- ❌ Switch branches to avoid conflicts

**Scenario: Lint errors**
- ✅ Fix lint errors and commit
- ❌ Skip linting to speed up process

**Scenario: Need to check other code**
- ✅ Use `grep`, `find`, or `git log` within worktree
- ❌ Switch to main branch to check

## 🚨 ESCALATION TRIGGERS

**STOP and escalate to manager if:**
- Git operations fail unexpectedly
- Tests fail after 3 fix attempts
- Merge conflicts cannot be resolved
- Any urge to use `git checkout` or `git switch`
- Directory/branch validation fails

## ✅ SUCCESS CRITERIA

- [ ] All tests pass (`go test ./...`)
- [ ] Code compiles without errors
- [ ] Linting passes (if applicable)
- [ ] PR created successfully
- [ ] Sound notification played
- [ ] No dogfooding violations occurred

## 📞 COMPLETION REPORT

When finished, report:
1. PR URL
2. Test results summary
3. Any issues encountered
4. Confirmation of no dogfooding violations
```

## Template: Worktree Cleanup Task

### Agent Instruction Packet: Cleanup

```markdown
# Autonomous Agent Task: Cleanup {WORKTREE_NAME}

## 🔒 MANDATORY PRE-FLIGHT VALIDATION

```bash
# Validate we're in the correct worktree
EXPECTED_PATH="{WORKTREE_PATH}"
CURRENT_PATH=$(pwd)
[[ "$CURRENT_PATH" == "$EXPECTED_PATH" ]] || { echo "❌ Wrong directory!"; exit 1; }

# Validate branch
EXPECTED_BRANCH="{BRANCH_NAME}"
CURRENT_BRANCH=$(git branch --show-current)
[[ "$CURRENT_BRANCH" == "$EXPECTED_BRANCH" ]] || { echo "❌ Wrong branch!"; exit 1; }

echo "✅ Cleanup validation passed"
```

## 🎯 YOUR MISSION

**Objective**: Clean up and finalize {WORKTREE_NAME} worktree

## 📋 EXECUTION STEPS

### Step 1: Status Check
```bash
git status
git log --oneline -10
```

### Step 2: Commit Any Pending Changes
```bash
# If there are uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    git add .
    git commit -m "Final cleanup of {WORKTREE_NAME}"
fi
```

### Step 3: Final Testing
```bash
go test ./...
go build -o wrkt
```

### Step 4: Push Final State
```bash
git push origin {BRANCH_NAME}
```

## 🚫 PROHIBITED OPERATIONS
- `git checkout` or `git switch`
- Leaving {WORKTREE_PATH} directory
- Any branch operations
```

## Template: Investigation Task

### Agent Instruction Packet: Investigation

```markdown
# Autonomous Agent Task: Investigate {INVESTIGATION_TOPIC}

## 🔒 MANDATORY PRE-FLIGHT VALIDATION

```bash
# Validate location
EXPECTED_PATH="{WORKTREE_PATH}"
CURRENT_PATH=$(pwd)
[[ "$CURRENT_PATH" == "$EXPECTED_PATH" ]] || { 
    echo "❌ Wrong directory! Expected: $EXPECTED_PATH, Got: $CURRENT_PATH"
    exit 1
}
echo "✅ Investigation validation passed"
```

## 🎯 YOUR MISSION

**Objective**: Investigate {INVESTIGATION_TOPIC} within {WORKTREE_NAME}

## 📋 INVESTIGATION TOOLS (ALLOWED)

- `grep -r "pattern" .` - Search code
- `find . -name "*.go" -exec grep "pattern" {} \;` - Find files
- `git log --grep="pattern"` - Search commit history
- `git blame filename` - Check code authorship
- `git show commit-hash` - View specific commit

## 🚫 PROHIBITED INVESTIGATION METHODS

- `git checkout` to check other branches
- `cd` outside worktree to check main branch
- `git switch` to compare with other branches

## 📊 INVESTIGATION REPORT FORMAT

```markdown
# Investigation Report: {INVESTIGATION_TOPIC}

## Findings
- [ ] Finding 1
- [ ] Finding 2

## Evidence
- File: path/to/file.go:line_number
- Commit: commit_hash - "commit message"

## Recommendations
- [ ] Recommendation 1
- [ ] Recommendation 2
```
```

## Template Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| `{FEATURE_NAME}` | Name of the feature | "list-filtering" |
| `{WORKTREE_PATH}` | Full path to worktree | "/Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/feature-list-filters" |
| `{BRANCH_NAME}` | Git branch name | "feature/list-filters-fixed" |
| `{REPO_NAME}` | Repository name | "wrkt" |
| `{COMMIT_MESSAGE}` | Commit message template | "Add list filtering functionality" |
| `{PR_TITLE}` | Pull request title | "Add list filtering options (--dirty, --verbose, --names-only)" |
| `{PR_SUMMARY}` | PR description | "Implements MVP requirements for list command filtering" |

## Template Customization Guidelines

1. **Always include pre-flight validation** - This prevents 90% of dogfooding violations
2. **Specify exact paths** - No relative paths or assumptions
3. **Include decision trees** - Help agents make autonomous choices
4. **Define escalation triggers** - When to stop and ask for help
5. **Add success criteria** - Clear definition of completion
6. **Prohibit dangerous operations** - Explicit blacklist of commands

## Manager Instructions for Template Usage

1. **Copy appropriate template**
2. **Fill in ALL variables** (use find/replace)
3. **Verify paths exist** and are correct
4. **Save as dedicated instruction file**
5. **Launch ONE agent** with instruction file
6. **Wait for completion** before next task
7. **Never manage multiple agents simultaneously**

This template system eliminates context-switching cognitive load by providing complete, self-contained instructions that prevent dogfooding violations through systematic validation and constraint enforcement.