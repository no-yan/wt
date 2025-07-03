# T101: Fix Shell Integration Test Failure

**Priority**: High (Ship-blocking)  
**Complexity**: Medium (2-3h)  
**Status**: Ready for assignment  
**Worktree**: `test-shell-integration` (existing)

## 🎯 PROBLEM STATEMENT

Shell integration tests are failing due to `git worktree add -b` command issues when creating new branches that don't exist yet.

**Current Error**:
```
Error adding worktree: failed to add worktree: git worktree add failed: command failed: exit status 128: fatal: invalid reference: test-branch
```

## 📋 TASK SCOPE

### Primary Objective
Fix the branch creation logic in shell integration testing to handle both new and existing branches correctly.

### Specific Changes Required
1. **Fix `addGitWorktree` method** in `internal/worktree_manager.go`
2. **Update test expectations** in shell integration tests
3. **Ensure backward compatibility** with existing worktree operations

## 🔧 TECHNICAL REQUIREMENTS

### Implementation Details

#### Option A: Adopt proven fallback logic
```go
// Try to create new branch first, then add worktree
createBranchCmd := fmt.Sprintf("git -C %s branch %s", repoPath, branch)
_, createErr := wm.runner.Run(createBranchCmd)

// Then add worktree (works with both new and existing branches)
gitCmd := fmt.Sprintf("git -C %s worktree add %s %s", repoPath, worktreePath, branch)
```

#### Option B: Fix current -b approach
```go
// Try creating new branch first
gitCmd := fmt.Sprintf("git -C %s worktree add -b %s %s", repoPath, branch, worktreePath)
if _, err := wm.runner.Run(gitCmd); err != nil {
    // Fallback to existing branch
    gitCmd = fmt.Sprintf("git -C %s worktree add %s %s", repoPath, worktreePath, branch)
    // ... handle error
}
```

### Files to Modify
- `internal/worktree_manager.go` (line ~96-114)
- `internal/worktree_manager_test.go` (update mocks)
- `test/test_zsh.sh` (ensure test compatibility)

## ✅ ACCEPTANCE CRITERIA

### Functional Requirements
- [ ] `wrkt add test-branch` creates new branch successfully
- [ ] `wrkt add existing-branch` uses existing branch successfully  
- [ ] Shell integration tests pass completely
- [ ] All existing tests continue to pass

### Quality Requirements
- [ ] No regression in existing functionality
- [ ] Error messages are clear and helpful
- [ ] Test coverage maintained or improved

### Integration Requirements
- [ ] Compatible with auto-setup feature (T001 completed)
- [ ] Works with branch fallback logic pattern

## 🧪 TESTING STRATEGY

### Unit Tests
```bash
# Ensure core logic works
go test ./internal/... -v

# Focus on specific methods
go test ./internal/... -run TestWorktreeManager_AddWorktree -v
```

### Integration Tests
```bash
# Shell integration test (main target)
export PATH=$PATH:$(pwd) && go test ./cmd/shell_init_integration_test.go -v

# Manual verification
./wrkt add test-new-branch
./wrkt list --verbose
```

### Regression Testing
```bash
# Ensure no breakage
go test ./... -v
./wrkt add feature/existing-test  # Should work
./wrkt add feature/new-test       # Should work
```

## 📂 DEVELOPMENT WORKFLOW

### Setup
```bash
# Navigate to existing worktree
cd /Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/test-shell-integration

# Check current status
git status
./wrkt list --verbose
```

### Implementation Process
1. **Analyze current implementation** - understand the existing -b approach
2. **Choose implementation strategy** - Option A (proven) or Option B (fix current)
3. **Update worktree_manager.go** - implement chosen approach
4. **Update tests** - ensure mocks and expectations align
5. **Run targeted tests** - focus on shell integration
6. **Run full test suite** - ensure no regression

### Completion Checklist
- [ ] Implementation complete and tested
- [ ] Shell integration test passes
- [ ] All unit tests pass
- [ ] No linting errors
- [ ] Changes committed with clear message
- [ ] Task marked complete in SESSION.md

## 🚨 POTENTIAL BLOCKERS

### Known Issues
- **Git behavior variations**: Different Git versions may behave differently
- **Test environment**: Shell integration tests require zsh and specific PATH setup

### Escalation Criteria
- If both Option A and B fail, escalate for architectural review
- If test failures persist after implementation, escalate for debugging support
- If merge conflicts arise, coordinate with manager

## 🔗 RELATED WORK

### Completed Dependencies
- ✅ T001: Auto-setup feature (provides context)
- ✅ T002: Branch fallback logic (provides implementation pattern)

### Future Dependencies
- T301: Merge to main (depends on this completion)

### Reference Implementation
See `/Users/noyan/ghq/github.com/no-yan/wrkt/worktrees/feature-branch-fallback-logic/internal/worktree_manager.go` for proven fallback pattern.

---

**Estimated Completion**: 2-3 hours focused work  
**Success Metric**: Shell integration tests pass completely  
**Next Action**: Assign to available developer