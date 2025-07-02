# Claude Implementation Context

## 🚨🚨🚨 THE THREE COMMANDMENTS - NEVER VIOLATE 🚨🚨🚨

### COMMANDMENT 1: THOU SHALT NOT use `git checkout` or `git switch` in main directory
### COMMANDMENT 2: THOU SHALT NOT change branches within any worktree
### COMMANDMENT 3: THOU SHALT use wrkt commands for ALL branch/worktree operations

## 🔒 MANDATORY PRE-COMMAND VALIDATION

**COPY AND RUN THIS BEFORE EVERY GIT OPERATION:**

```bash
#!/bin/bash
# DOGFOODING VALIDATION - RUN BEFORE ANY GIT COMMAND
validate_dogfooding() {
    local current_dir=$(pwd)
    local current_branch=$(git branch --show-current 2>/dev/null || echo "unknown")
    
    echo "📍 Current location: $current_dir"
    echo "🌿 Current branch: $current_branch"
    
    # Check if in main directory with wrong branch
    if [[ "$current_dir" == */wrkt ]] && [[ "$current_branch" != "main" ]]; then
        echo "❌ CRITICAL VIOLATION: Main directory on wrong branch!"
        echo "🔧 Required fix: git checkout main"
        echo "⚠️  Future operations: Use worktrees only"
        return 1
    fi
    
    # Check for prohibited commands
    read -p "Enter your intended git command: " git_command
    if [[ "$git_command" =~ "checkout -b" ]] || [[ "$git_command" =~ "switch -c" ]]; then
        echo "❌ PROHIBITED COMMAND: $git_command"
        echo "✅ Use instead: ./wrkt add <branch-name>"
        return 1
    fi
    
    if [[ "$git_command" =~ "checkout" ]] || [[ "$git_command" =~ "switch" ]]; then
        if [[ "$current_dir" == */wrkt ]]; then
            echo "❌ PROHIBITED: No branch switching in main directory"
            echo "✅ Use instead: ./wrkt switch <worktree-name>"
            return 1
        fi
    fi
    
    echo "✅ Validation passed - command allowed"
    return 0
}

# Run validation
validate_dogfooding
```

## 🚨 CRITICAL DEVELOPMENT RULES

### MANDATORY: TodoWrite/TodoRead Usage
**Before starting ANY development work, Claude MUST:**
1. **Read Current Todo List**: Use `TodoRead` tool to check existing tasks
2. **Plan Work**: Use `TodoWrite` tool to create/update task list for the session
3. **Track Progress**: Update todo status throughout development:
   - `"pending"` - Task not yet started
   - `"in_progress"` - Currently working on (limit to ONE at a time)
   - `"completed"` - Task finished successfully
4. **Mark Completion**: IMMEDIATELY mark tasks as completed when finished

### MANDATORY: Dogfooding Principles
**Absolute Rules to Follow:**
- **Always use wrkt commands for worktree operations**
- **Prohibit direct branch operations in main directory**
- **Prohibit branch switching within existing worktrees**

**Correct Operation Procedure:**
1. `./wrkt add <branch>` to create worktree
2. `./wrkt switch <name>` to move (auto-cd after shell integration implementation, manual cd before)
3. Work only within worktrees, no branch changes
4. Always check status with `./wrkt list --verbose` at session start

**Prohibited Operations:**
```bash
# ❌ Prohibited: Feature branch work in main directory
git checkout -b feature/new-feature

# ❌ Prohibited: Branch switching within worktrees
cd /path/to/worktree && git checkout other-branch

# ✅ Correct: Create worktree with wrkt command
./wrkt add feature/new-feature
```

**Design Philosophy:**
- **1 worktree = 1 dedicated branch = 1 feature**
- **Main directory = main branch only**
- **Realize branch switching through worktree creation**

## 📊 TASK MANAGEMENT

### Autonomous Agent Management Protocol

**CRITICAL**: The recurring `git checkout -b` failures occur due to context-switching cognitive load when managing multiple agents simultaneously. The solution is **autonomous agent architecture**.

#### Manager Responsibilities (Instruction Architect)
1. **Create Complete Instructions** - Not commands
2. **Monitor Results** - Not process
3. **Handle Escalations** - Not micromanage
4. **Sequential Task Assignment** - Never concurrent development

#### Agent Instruction Packet Structure

Each agent receives a complete, self-contained instruction document:

```yaml
task_id: "feature-x-pr"
worktree_path: "/Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/feature-x"
branch: "feature/x"
objective: "Create PR for feature X"

mandatory_constraints:
  - work_only_in: "/Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/feature-x"
  - never_use: ["git checkout", "git switch"]
  - never_leave_worktree: true
  - use_wrkt_commands_only: true

validation_checkpoints:
  - pre_start: "verify_worktree_location_and_branch"
  - pre_commit: "run_all_tests"
  - pre_push: "verify_dogfooding_compliance"
  - pre_pr: "confirm_no_branch_violations"

escalation_triggers:
  - "merge conflicts"
  - "test failures after 3 attempts"
  - "any git command error"
  - "temptation to use git checkout"

autonomous_decisions:
  - fix_lint_errors: "yes"
  - commit_incremental_progress: "yes"
  - revert_failed_changes: "yes"
  - create_pr_when_tests_pass: "yes"
```

#### Example Agent Launch Process

```bash
# 1. Manager creates complete instruction file
cat > /tmp/agent_task_feature_x.md << 'EOF'
# Autonomous Agent Task: Feature X PR Creation

## PRE-FLIGHT VALIDATION (MANDATORY)
```bash
# You MUST run this before starting
cd /Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/feature-x
[[ $(pwd) == */worktrees/feature-x ]] || { echo "ERROR: Wrong directory"; exit 1; }
[[ $(git branch --show-current) == "feature/x" ]] || { echo "ERROR: Wrong branch"; exit 1; }
echo "✅ Validation passed - you are in the correct worktree"
```

## YOUR MISSION
1. Review uncommitted changes with `git status`
2. Run tests with `go test ./...`
3. Fix any failing tests (within this worktree only)
4. Commit changes with descriptive message
5. Push to origin: `git push -u origin feature/x`
6. Create PR with `gh pr create`
7. Notify completion

## DECISION TREE
- Tests failing? → Fix within worktree (NEVER switch branches)
- Need different branch? → STOP and escalate to manager
- Merge conflicts? → ESCALATE immediately
- Lint errors? → Fix within worktree

## ABSOLUTELY PROHIBITED
- `git checkout` (any form)
- `git switch` (any form)
- Leaving worktree directory
- Any branch operations

## SUCCESS CRITERIA
✅ All tests pass
✅ No dogfooding violations
✅ PR created successfully
✅ Sound notification played
EOF

# 2. Launch autonomous agent
claude --autonomous --task-file /tmp/agent_task_feature_x.md
```

#### Cognitive Load Reduction Benefits

1. **Manager Focus**: Creates instructions, not commands
2. **No Context Juggling**: Agents work independently
3. **Constraint Enforcement**: Validation built into instructions
4. **Error Prevention**: Pre-command checks eliminate violations
5. **Clear Boundaries**: Each agent owns exactly one worktree

#### Conflict Resolution Protocol
- **Sequential Development**: Complete one task before starting next
- **Conflict Analysis**: Use `git diff base..branch --name-only` before task assignment
- **Priority-Based Sequencing**: High priority features complete first
- **PR-Based Integration**: All merges via GitHub PRs
- **Test Validation**: All tests must pass before merge

#### GitHub PR Integration Process
- Each autonomous agent creates its own PR
- Manager reviews and merges PRs sequentially
- No parallel merging to prevent conflicts
- Test validation required for each PR

### Sequential Worktree Development Workflow
**Instruction-First Development to Prevent Cognitive Load Issues:**

1. **Check WORKTREE_TRACKING.md** - Review status of all active worktrees
2. **Use TodoWrite** to plan task sequence (one at a time)
3. **Create Complete Instructions** - Write full agent packet before starting
4. **Launch Single Agent** - Never manage multiple agents simultaneously
5. **Wait for Completion** - Agent reports back when done
6. **Verify Results** - Check PR created, tests pass, no violations
7. **Next Task** - Only then move to next worktree

**Mental Model for Managers:**
```
Instruction Creation → Agent Launch → Wait → Verify → Next
```

**Never:**
```
Agent 1 + Agent 2 + Agent 3 → Context switching → git checkout -b errors
```

### Task Categories
- **High Priority**: Core functionality, bug fixes, incomplete features
- **Medium Priority**: Enhancements, new features, optimizations
- **Low Priority**: Documentation, cleanup, nice-to-have features

### Example Todo Usage
```
TodoWrite: [
  {"content": "Complete feature-list-filters worktree", "status": "in_progress", "priority": "high", "id": "1"},
  {"content": "Add zsh tab completion", "status": "pending", "priority": "medium", "id": "2"},
  {"content": "Update documentation", "status": "pending", "priority": "low", "id": "3"}
]
```

## 🧠 COGNITIVE LOAD MANAGEMENT

### Mental Model Reconstruction

**WRONG Mental Model (causes failures):**
"Git repository with worktree rules to remember"

**CORRECT Mental Model (prevents failures):**
"Worktree-managed codebase where Git is an implementation detail"

### Pre-Command Decision Framework

**Before ANY command, mandatory thought process:**

```
1. LOCATION CHECK: Where am I?
   → $(pwd) = ? 
   → Main directory or worktree?

2. OPERATION CHECK: What do I want to achieve?
   → New branch? Use `wrkt add`
   → Switch context? Use `wrkt switch`
   → Git operation? Validate against dogfooding

3. CONSTRAINT CHECK: Does this violate the three commandments?
   → Commandment 1: No git checkout in main
   → Commandment 2: No branch switching in worktrees
   → Commandment 3: Use wrkt commands only

4. ALTERNATIVE CHECK: What's the wrkt way?
   → Every git operation has a wrkt equivalent
   → When in doubt, use wrkt commands
```

### Context-Switching Error Prevention

**The Problem:** Managing multiple agents simultaneously causes **cognitive contamination** where standard Git patterns override project-specific constraints.

**The Solution:** Instruction-first autonomous architecture

**Pattern Recognition Training:**

```
❌ DANGER PATTERN: "I need to test this, let me create a branch"
✅ SAFE PATTERN: "I need to test this, am I in the right worktree?"

❌ DANGER PATTERN: "Quick git checkout to..."
✅ SAFE PATTERN: "Wait, let me run dogfooding validation first"

❌ DANGER PATTERN: Managing 3 agents simultaneously
✅ SAFE PATTERN: One complete instruction, one agent, one task
```

### Failure Recovery Protocol

**When dogfooding violation occurs:**

1. **STOP ALL OPERATIONS** immediately
2. **Verify current state** of all worktrees:
   ```bash
   cd /Users/noyan/ghq/github.com/no-yan/wrkt
   ./wrkt list --verbose
   ```
3. **Identify contamination:**
   - Which directory is on wrong branch?
   - Are there uncommitted changes to preserve?
4. **Clean state:**
   ```bash
   # If main directory contaminated
   cd /Users/noyan/ghq/github.com/no-yan/wrkt
   git checkout main
   ```
5. **Redesign approach** using instruction-first method
6. **Update mental model** to prevent recurrence

### Constraint Integration Techniques

**Make dogfooding principles automatic:**

1. **Command Generation Reframe:**
   - Instead of "git checkout -b" → "wrkt add"
   - Instead of "cd worktree && git switch" → "wrkt switch"

2. **Validation Automation:**
   - Run validation script before every git command
   - Embed constraints in instruction templates

3. **Mental Anchoring:**
   - Default assumption: "I cannot change branches"
   - Exception handling: "Different branch = different worktree"

## 🏗️ CORE IMPLEMENTATION NOTES

### Worktree Organization
- All worktrees created in `$REPO_ROOT/worktrees/` subdirectory
- Auto-create worktrees/ directory on first use
- Auto-add "worktrees/" to .gitignore
- Path generation: `feature/auth` → `worktrees/feature-auth/`

### Zsh Integration Implementation
- `wrkt switch` command should only resolve and return the target path
- Actual directory changing is handled by zsh functions
- Zsh functions must be generated by `wrkt shell-init`
- Only support zsh - show clear error for other shells

### Command Design
- `wrkt list` is the primary information command
- No separate `wrkt status` - use `wrkt list --dirty` instead
- `wrkt list --verbose` provides detailed git status
- **Exact name matching only** - no fuzzy matching complexity

## Troubleshooting Common Issues

- **"wrkt switch doesn't change directory"**: User hasn't set up zsh integration
- **"shell not supported"**: User not using zsh
- **"command not found" after switch**: Zsh integration not loaded
- **"worktree not found"**: Check exact name with `wrkt list`
- **Path generation conflicts**: Simple conflict resolution with numbering
- **Tab completion not working**: Zsh integration setup incomplete
