# Agent Development Handbook

**For Autonomous Task Execution in wrkt Project**

## 🎯 QUICK START GUIDE

### Before Starting Any Task
```bash
# 1. Validate environment
./.claude/validation/pre_task_checklist.sh [TASK_ID]

# 2. Choose task from roadmap
cat .claude/TASK_ROADMAP.md

# 3. Read task specification
cat .claude/tasks/[TASK_ID]*.md

# 4. Create worktree (if needed)
./wrkt add feature/task-name
```

### During Development
```bash
# Frequent checks
go test ./...           # Run tests often
go vet ./...           # Check for issues
gofmt -w .             # Format code

# Test wrkt functionality
./wrkt list --verbose   # Verify wrkt works
```

### Before Completing Task
```bash
# Final validation
./.claude/validation/post_task_validation.sh

# If all passes, commit and mark complete
git add . && git commit -m "feat: implement task XYZ"
```

## 📋 DEVELOPMENT WORKFLOW

### 1. Task Selection & Setup
- **Choose independent task** from TASK_ROADMAP.md
- **Read complete specification** in .claude/tasks/
- **Validate environment** with pre_task_checklist.sh
- **Create worktree** if specified in task

### 2. Implementation Phase
- **Follow task scope strictly** - no scope creep
- **Test frequently** - run `go test ./...` often  
- **Maintain quality** - format code, handle errors
- **Document decisions** in commit messages

### 3. Quality Assurance
- **Run full test suite** - all tests must pass
- **Validate functionality** - test wrkt commands work
- **Check code quality** - linting, formatting
- **Verify task completion** - all acceptance criteria met

### 4. Task Completion
- **Run post-validation** - use post_task_validation.sh
- **Commit with convention** - use conventional commit format
- **Update task status** - mark as completed
- **Coordinate integration** - notify manager if needed

## 🔧 TECHNICAL STANDARDS

### Code Quality Requirements
```bash
# All must pass before task completion
go test ./...                    # Zero test failures
go vet ./...                     # Zero vet warnings  
gofmt -l . | wc -l               # Zero unformatted files
golangci-lint run                # Clean linting (if available)
```

### Testing Standards
- **Unit tests** for all new functions
- **Integration tests** for command-line functionality
- **Regression tests** to ensure no breakage
- **Coverage goal** >80% for new code

### Git Workflow
```bash
# Conventional commit format
git commit -m "feat: add tab completion for remove command"
git commit -m "fix: resolve shell integration test failure"
git commit -m "test: add comprehensive clean command tests"
```

## 🛠️ COMMON PATTERNS

### Working with wrkt Commands
```bash
# Test command functionality
./wrkt add test-branch
./wrkt list --verbose
./wrkt switch test-branch
./wrkt remove test-branch

# Verify in different states
git status  # Check repository state
pwd         # Verify current location
```

### Mock Testing Pattern
```go
// Standard mock setup for tests
mockRunner := &MockCommandRunner{
    outputs: make(map[string]string),
}

// Mock successful commands
mockRunner.outputs["git -C /repo worktree add /path branch"] = ""

// Test with mock
manager := NewWorktreeManager(service, mockRunner)
result, err := manager.AddWorktree("/repo", "branch")
```

### Error Handling Pattern
```go
// Standard error handling
if err := validateInput(input); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

// Command execution with context
cmd := fmt.Sprintf("git -C %s command %s", shellescape(path), shellescape(arg))
if _, err := runner.Run(cmd); err != nil {
    return fmt.Errorf("git command failed: %w", err)
}
```

## 🚨 TROUBLESHOOTING

### Common Issues

#### Test Failures
```bash
# Debug test failures
go test -v ./...                 # Verbose test output
go test -run TestSpecific        # Run specific test
go test -timeout 30s ./...       # Increase timeout
```

#### Build Issues
```bash
# Clean and rebuild
go clean
go mod tidy
go build ./...
```

#### wrkt Command Issues
```bash
# Rebuild and test
go build -o wrkt main.go
./wrkt list --verbose            # Test functionality
PATH=$PATH:$(pwd) ./wrkt add test  # Test with PATH
```

### Environment Problems
```bash
# Check environment
go env                           # Go configuration
git --version                    # Git version
which zsh                        # Shell availability
```

## 📊 QUALITY METRICS

### Success Criteria
- ✅ **All tests pass** (`go test ./...`)
- ✅ **Code builds cleanly** (`go build ./...`)
- ✅ **No linting errors** (if golangci-lint available)
- ✅ **Functionality verified** (manual testing)
- ✅ **Task acceptance criteria met**

### Performance Targets
- **Test execution** <30 seconds for full suite
- **Build time** <10 seconds for full project
- **wrkt command response** <100ms for list operations

### Code Quality Targets
- **Test coverage** >80% for new code
- **Cyclomatic complexity** <10 for functions
- **No code duplication** >10 lines
- **Clear error messages** for all failure cases

## 🔗 REFERENCE LINKS

### Key Files
- `CLAUDE.md` - Core development rules and constraints
- `TASK_ROADMAP.md` - Available tasks and priorities  
- `.claude/tasks/` - Detailed task specifications
- `.claude/validation/` - Quality validation scripts

### Important Commands
```bash
# Project status
./wrkt list --verbose

# Development validation  
./.claude/validation/pre_task_checklist.sh
./.claude/validation/post_task_validation.sh

# Quality checks
go test ./...
go vet ./...
golangci-lint run
```

## 💡 TIPS FOR SUCCESS

### Efficiency Tips
1. **Read task spec completely** before starting
2. **Run validation early** to catch environment issues
3. **Test incrementally** - don't wait until the end
4. **Follow established patterns** in existing code
5. **Keep scope narrow** - resist feature creep

### Quality Tips
1. **Write tests first** or alongside implementation
2. **Test edge cases** - empty inputs, invalid data
3. **Handle errors gracefully** with clear messages
4. **Document non-obvious decisions** in comments
5. **Validate with real usage** - manual testing

### Debugging Tips
1. **Use verbose test output** for failures
2. **Check git status frequently** to avoid confusion
3. **Verify wrkt functionality** after changes
4. **Read error messages carefully** - they're usually clear
5. **Ask for help** if blocked >30 minutes

---

**Remember**: This handbook enables autonomous development. Follow it for consistent, high-quality results that integrate seamlessly with the overall project.