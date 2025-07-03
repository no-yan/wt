# Development Workflows & Procedures

- [Development Workflows \& Procedures](#development-workflows--procedures)
  - [🎯 AGENT TASK TEMPLATES](#-agent-task-templates)
    - [Template 1: PR Creation](#template-1-pr-creation)
    - [Template 2: Feature Development](#template-2-feature-development)
    - [Template 3: Manual Execution (System Constraints)](#template-3-manual-execution-system-constraints)
  - [📋 EXECUTION PROCEDURES](#-execution-procedures)
    - [Sequential Task Management](#sequential-task-management)
    - [Error Recovery Procedures](#error-recovery-procedures)
    - [System Constraint Handling](#system-constraint-handling)
  - [🔧 SPECIALIZED PROCEDURES](#-specialized-procedures)
    - [Creating New Worktrees](#creating-new-worktrees)
    - [Worktree Cleanup](#worktree-cleanup)
    - [Branch Management](#branch-management)
  - [🧪 TESTING PROCEDURES](#-testing-procedures)
    - [Standard Test Workflow](#standard-test-workflow)
    - [Integration Testing](#integration-testing)
  - [📊 TROUBLESHOOTING GUIDE](#-troubleshooting-guide)
    - [Common Issues \& Solutions](#common-issues--solutions)
    - [Performance Optimization](#performance-optimization)
  - [🔄 WORKFLOW PATTERNS](#-workflow-patterns)
    - [Daily Development Cycle](#daily-development-cycle)
    - [Feature Development Lifecycle](#feature-development-lifecycle)
    - [Bug Fix Workflow](#bug-fix-workflow)

## 🎯 AGENT TASK TEMPLATES

### Template 1: PR Creation

**Use when**: Creating pull requests for completed features

**Pre-flight Validation** (MANDATORY):
```bash
#!/bin/bash
echo "🚀 Starting pre-flight validation..."

# Validate location
EXPECTED_PATH="/Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/[WORKTREE_NAME]"
CURRENT_PATH=$(pwd)
if [[ "$CURRENT_PATH" != "$EXPECTED_PATH" ]]; then
    echo "❌ CRITICAL ERROR: Wrong directory!"
    echo "   Expected: $EXPECTED_PATH"
    echo "   Current:  $CURRENT_PATH"
    exit 1
fi

# Validate branch
EXPECTED_BRANCH="[BRANCH_NAME]"
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "$EXPECTED_BRANCH" ]]; then
    echo "❌ CRITICAL ERROR: Wrong branch!"
    echo "   Expected: $EXPECTED_BRANCH"
    echo "   Current:  $CURRENT_BRANCH"
    exit 1
fi

echo "✅ Pre-flight validation PASSED"
```

**Execution Steps**:
1. Run pre-flight validation
2. Check git status and add files: `git add .`
3. Run tests: `go test ./...`
4. Commit changes with descriptive message
5. Push to remote: `git push -u origin [BRANCH_NAME]`
6. Create PR: `gh pr create --title "[TITLE]" --body "[DESCRIPTION]"`

**Success Criteria**:
- [ ] All tests pass
- [ ] Code compiles without errors
- [ ] All changes committed and pushed
- [ ] PR created with proper title and description
- [ ] No dogfooding violations occurred

### Template 2: Feature Development

**Use when**: Implementing new features within existing worktrees

**Setup Process**:
1. Validate current worktree and branch
2. Update `.claude/SESSION.md` with current task
3. Plan implementation steps using TodoWrite

**Development Cycle**:
1. Implement feature incrementally
2. Run tests after each significant change: `go test ./...`
3. Commit working states frequently
4. Never switch branches or worktrees during development
5. Document progress in session notes

**Completion Process**:
1. Ensure all tests pass
2. Update documentation if needed
3. Prepare for PR creation using Template 1

### Template 3: Manual Execution (System Constraints)

**Use when**: Bash tool is non-functional or system constraints prevent automation

**Pre-execution Validation**:
```bash
# Confirm location and branch
pwd
git branch --show-current

# Verify no dogfooding violations
if [[ $(pwd) == */wrkt ]] && [[ $(pwd) != */worktrees/* ]]; then
    if [[ $(git branch --show-current) != "main" ]]; then
        echo "🚨 DOGFOODING VIOLATION: Main directory on wrong branch"
        echo "Fix required: git checkout main"
        exit 1
    fi
fi

# List worktrees for context
./wrkt list --verbose
```

**Manual PR Creation Steps**:
1. Navigate to worktree: `cd /Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/[NAME]`
2. Validate location and branch
3. Check status and add files: `git status && git add .`
4. Commit with proper message format:
   ```bash
   git commit -m "[Description of changes]

   - [Bullet point summary]
   - [Key features implemented]
   - [Technical details]

   🤖 Generated with [Claude Code](https://claude.ai/code)

   Co-Authored-By: Claude <noreply@anthropic.com>"
   ```
5. Push to remote: `git push -u origin [BRANCH]`
6. Create PR: `gh pr create --title "[TITLE]" --body "[DESCRIPTION]"`
7. Update `.claude/SESSION.md` with PR URL and completion status

## 📋 EXECUTION PROCEDURES

### Sequential Task Management

**Core Principle**: Complete one task entirely before starting the next to prevent cognitive contamination.

**Workflow**:
1. **Update Session State**: Record current task in `.claude/SESSION.md`
2. **Validate Environment**: Run session start checklist from `CLAUDE.md`
3. **Focus Single Task**: Work on one worktree/feature at a time
4. **Complete Fully**: Finish, test, and create PR before moving on
5. **Document Progress**: Update session state with completion

**Task Prioritization**:
- **High Priority**: Core functionality, bug fixes, incomplete features
- **Medium Priority**: Enhancements, new features, optimizations
- **Low Priority**: Documentation, cleanup, nice-to-have features

### Error Recovery Procedures

**If Git Operations Fail**:
```bash
# 1. Verify repository state
git status
git log --oneline -5

# 2. Check for conflicts or issues
git diff
git diff --staged

# 3. If needed, reset to clean state
git reset --soft HEAD~1  # Undo last commit but keep changes
# or
git reset --hard HEAD    # Discard all changes (use carefully)

# 4. Re-attempt operation with fresh state
```

**If PR Creation Fails**:
1. **Authentication Error**:
   ```bash
   gh auth status
   gh auth login  # If needed
   ```

2. **Remote Branch Doesn't Exist**:
   ```bash
   git push -u origin <branch-name>
   # Then retry PR creation
   ```

3. **Title/Body Too Long**:
   - Shorten title to < 72 characters
   - Move detailed info to PR description

**If Manual Updates Corrupt State**:
1. **Check repository state**:
   ```bash
   git status
   ./wrkt list --verbose
   ```

2. **Restore from known good state**:
   ```bash
   # If in main directory with wrong branch:
   git checkout main

   # If worktree corrupted:
   ./wrkt remove <corrupted-worktree>
   ./wrkt add <branch-name>
   ```

3. **Validate restoration**:
   ```bash
   # Run full validation
   pwd
   git branch --show-current
   ./wrkt list --verbose
   ```

### System Constraint Handling

**When Bash Tool is Non-Functional**:
1. **New Session**: Most reliable solution - restart Claude session
2. **Manual Execution**: Use Template 3 procedures above
3. **Focus on Documentation**: Update guides and session notes
4. **Prepare for Recovery**: Ensure clear next steps documented

**When Context is Lost**:
1. **Read Session State**: Check `.claude/SESSION.md` for current context
2. **Validate Repository State**: Check worktrees and current locations
3. **Resume from Last Known State**: Follow next actions from session notes

**When Violations are Detected**:
1. **Follow Emergency Procedures**: Use procedures from `CLAUDE.md`
2. **Document the Violation**: Add to session learnings
3. **Update Prevention**: Enhance validation if new pattern discovered

## 🔧 SPECIALIZED PROCEDURES

### Creating New Worktrees

**Standard Process**:
```bash
# 1. From main directory only
cd /Users/noyan/ghq/github.com/no-yan/wrkt

# 2. Create new worktree
./wrkt add feature/new-feature-name

# 3. Switch to new worktree
./wrkt switch feature-new-feature-name

# 4. Verify setup
pwd  # Should show worktree path
git branch --show-current  # Should show feature branch
```

**With New Branch Creation**:
```bash
# For completely new branches
./wrkt add -b feature/new-branch origin/main
```

### Worktree Cleanup

**Remove Completed Worktrees**:
```bash
# 1. Ensure you're not in the worktree to be removed
cd /Users/noyan/ghq/github.com/no-yan/wrkt

# 2. Remove worktree
./wrkt remove worktree-name

# 3. Verify removal
./wrkt list
```

### Branch Management

**Push New Branch**:
```bash
# From within worktree
git push -u origin <branch-name>
```

**Delete Remote Branch** (after PR merge):
```bash
git push origin --delete <branch-name>
```

## 🧪 TESTING PROCEDURES

### Standard Test Workflow

**Before Any Commit**:
```bash
# 1. Run all tests
go test ./...

# 2. Check for compilation errors
go build -o wrkt

# 3. Test basic functionality
./wrkt list
./wrkt list --verbose
```

**Before PR Creation**:
```bash
# 1. Full test suite
go test ./... -v

# 2. Test with coverage (if available)
go test ./... -cover

# 3. Lint check (if configured)
go vet ./...
```

### Integration Testing

**Shell Integration Tests**:
```bash
# If shell integration tests are available
go test ./cmd -v -run TestShellIntegration
```

**Worktree Manager Tests**:
```bash
# Test worktree operations
go test ./internal -v -run TestWorktreeManager
```

## 📊 TROUBLESHOOTING GUIDE

### Common Issues & Solutions

**"wrkt switch doesn't change directory"**:
- Setup zsh integration: `./wrkt shell-init`
- Source the integration in shell profile

**"shell not supported"**:
- Ensure using zsh (other shells not supported)
- Check shell with: `echo $SHELL`

**"command not found" after switch**:
- Zsh integration not loaded
- Run: `source ~/.zshrc` or restart terminal

**"worktree not found"**:
- Check exact names: `./wrkt list`
- Use exact name matching (no fuzzy matching)

**Path generation conflicts**:
- Use numbered suffixes for conflicts
- Check existing paths: `./wrkt list --verbose`

**Tab completion not working**:
- Ensure zsh integration setup complete
- Check compdef function is loaded

### Performance Optimization

**Large Repository Handling**:
- Use `./wrkt list --dirty` for quick status
- Avoid `--verbose` flag for large numbers of worktrees

**Memory Usage**:
- Clean up unused worktrees regularly
- Remove merged branches to reduce clutter

## 🔄 WORKFLOW PATTERNS

### Daily Development Cycle

1. **Session Start**: Run checklist from `CLAUDE.md`
2. **Load Context**: Read `.claude/SESSION.md`
3. **Plan Tasks**: Update TodoWrite with today's priorities
4. **Focus Work**: One worktree at a time
5. **Test Frequently**: After each significant change
6. **Document Progress**: Update session state regularly
7. **Session End**: Record handoff notes and next actions

### Feature Development Lifecycle

1. **Planning**: Create worktree with `./wrkt add feature/name`
2. **Development**: Implement incrementally with frequent commits
3. **Testing**: Continuous testing throughout development
4. **Review**: Self-review before PR creation
5. **PR Creation**: Use Template 1 procedures
6. **Cleanup**: Remove worktree after merge

### Bug Fix Workflow

1. **Reproduce**: Create worktree for bug investigation
2. **Diagnose**: Identify root cause within worktree
3. **Fix**: Implement minimal necessary changes
4. **Test**: Verify fix and ensure no regressions
5. **Document**: Clear commit message explaining the fix
6. **Fast-track**: Priority PR creation and review

## 🔄 PARALLEL TASK MANAGEMENT

### Multi-Task Session Setup

**Task Creation Process**:
1. **Check Capacity**: Verify <3 concurrent tasks in SESSION.md
2. **Create Worktree**: `./wrkt add feature/task-name`
3. **Create Task File**: Copy `.claude/tasks/TASK_TEMPLATE.md` → `TASK_XXX_name.md`
4. **Update SESSION.md**: Add task to Active Tasks Matrix
5. **Assign Priority**: Set task priority and dependencies

### Task Monitoring & Coordination

**Master Claude Responsibilities**:
- Monitor all active tasks via SESSION.md matrix
- Resolve inter-task dependencies and conflicts
- Coordinate PR integration sequence
- Handle escalations from task-specific Claude instances

**Task-Specific Claude Responsibilities**:
- Focus exclusively on assigned task and worktree
- Update task file every 30 minutes with progress
- Escalate blockers immediately to master
- Never switch to different worktrees or tasks

### Task Integration Workflow

**PR Integration Sequence**:
1. **Dependency Check**: Verify prerequisite tasks completed
2. **Conflict Analysis**: Check for file overlap conflicts
3. **Sequential Merge**: Integrate PRs one at a time
4. **Integration Testing**: Run full test suite after each merge
5. **Task Cleanup**: Remove completed task files and worktrees

### Parallel Development Safety

**Isolation Requirements**:
- Each task must use dedicated worktree
- No shared file editing between tasks
- All task communication via SESSION.md updates
- Master Claude coordinates all cross-task activities

**Conflict Prevention**:
- Check file overlap before task assignment
- Prioritize tasks with potential conflicts
- Use feature flags for incompatible changes
- Maintain clear task scope boundaries

## 📐 DOCUMENTATION MAINTENANCE

### File-Specific Update Guidelines

**CLAUDE.md Updates**:
- **Target**: Ultra-high frequency information only
- **Line Limit**: Keep under 250 lines total
- **Content**: Session checklist, rules, core patterns
- **Avoid**: Detailed procedures, examples, changing information

**WORKFLOWS.md Updates**:
- **Target**: Operational procedures and templates
- **Structure**: Templates → Procedures → Troubleshooting
- **Cross-Reference**: Link to CLAUDE.md for rules
- **Practical Focus**: Real workflow steps and examples

**SESSION.md Updates**:
- **Frequency**: Continuous during session
- **Content**: Current state, progress, insights, next actions
- **Archival**: Summarize completed work, detail current focus
- **Handoff**: Clear next session starting points

**Task Files (.claude/tasks/TASK_XXX.md)**:
- **Size Limit**: Maximum 15 lines
- **Content**: Task-specific context, progress, next actions
- **References**: Link to WORKFLOWS.md for procedures
- **Lifecycle**: Create → Update → Archive/Delete on completion

### Information Architecture Best Practices

**Redundancy Management**:
```markdown
❌ Bad: Same information in multiple files
✅ Good: Single source with references
✅ Best: "Procedure: WORKFLOWS.md#section-name"
```

**Access Optimization**:
- **Level 1**: Immediate needs (CLAUDE.md)
- **Level 2**: Working procedures (WORKFLOWS.md)  
- **Level 3**: Current state (SESSION.md)
- **Level 4**: Task specifics (tasks/*.md)

**Visual Hierarchy**:
```markdown
# Main Purpose (1 per file)
## Functional Sections (grouped by purpose)
### Detailed Categories (when needed)
**Emphasis**: For critical points
- Lists: For parallel information
```

### Quality Assurance Checklist

**Before Publishing Updates**:
- [ ] Information placed in correct file by frequency
- [ ] No duplicate content across files
- [ ] References properly formatted and working
- [ ] Visual hierarchy maintained
- [ ] Line limits respected per file type

**After Updates**:
- [ ] Information findable within 2 minutes
- [ ] Consistent terminology used
- [ ] Navigation paths remain clear
- [ ] File purpose boundaries maintained
