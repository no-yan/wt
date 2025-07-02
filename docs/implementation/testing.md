# Testing Guide

Comprehensive testing strategy and guidelines for the `wrkt` project.

## Testing Philosophy

1. **Test behavior, not implementation** - Focus on what the code does, not how
2. **Test at the right level** - Unit tests for logic, integration tests for workflows
3. **Make tests readable** - Tests serve as documentation
4. **Test edge cases** - Git worktree operations have many edge cases
5. **Fast feedback** - Tests should run quickly during development

## Test Structure

```
wrkt/
├── internal/
│   ├── worktree.go
│   ├── worktree_test.go      # Unit tests
│   ├── git.go
│   └── git_test.go           # Unit tests
├── cmd/
│   ├── list.go
│   ├── list_test.go          # Command tests
│   └── ...
├── test/
│   ├── integration/          # Integration tests
│   ├── fixtures/            # Test data and fixtures
│   └── helpers/             # Test helper functions
└── Makefile                 # Test automation
```

## Unit Tests

### Testing Principles
- Test public interfaces
- Mock external dependencies (git commands)
- Use table-driven tests for multiple scenarios
- Test error conditions explicitly

### Example: Worktree Parser Tests

```go
func TestParseWorktreeList(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected []*Worktree
        wantErr  bool
    }{
        {
            name: "single worktree",
            input: `worktree /path/to/main
HEAD abcd1234abcd1234abcd1234abcd1234abcd1234
branch refs/heads/main

`,
            expected: []*Worktree{
                {
                    Path:   "/path/to/main",
                    Head:   "abcd1234abcd1234abcd1234abcd1234abcd1234",
                    Branch: "main",
                },
            },
        },
        {
            name: "detached HEAD",
            input: `worktree /path/to/detached
HEAD 1234abc1234abc1234abc1234abc1234abc1234a
detached

`,
            expected: []*Worktree{
                {
                    Path:     "/path/to/detached",
                    Head:     "1234abc1234abc1234abc1234abc1234abc1234a",
                    Branch:   "detached HEAD",
                    Detached: true,
                },
            },
        },
        {
            name:    "invalid input",
            input:   "invalid\nformat",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := parseWorktreeList(tt.input)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Testing Guidelines

#### 1. Test Data Setup
```go
func setupTestWorktrees() []*Worktree {
    return []*Worktree{
        {Path: "/path/to/main", Branch: "main", Head: "abc123"},
        {Path: "/path/to/feature-auth", Branch: "feature/auth", Head: "def456"},
        {Path: "/path/to/hotfix", Branch: "hotfix/bug-123", Head: "ghi789"},
    }
}
```

#### 2. Error Testing
```go
func TestWorktreeManager_List_GitError(t *testing.T) {
    // Mock git command failure
    wm := &WorktreeManager{gitCmd: mockFailingGitCmd}
    
    _, err := wm.List()
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to list worktrees")
}
```

#### 3. Fuzzy Matching Tests
```go
func TestFuzzyMatching(t *testing.T) {
    worktrees := setupTestWorktrees()
    wm := &WorktreeManager{worktrees: worktrees}
    
    tests := []struct {
        name     string
        query    string
        expected string
        wantErr  bool
    }{
        {"exact match", "main", "/path/to/main", false},
        {"branch partial", "auth", "/path/to/feature-auth", false},
        {"path partial", "hotfix", "/path/to/hotfix", false},
        {"no match", "nonexistent", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := wm.GetWorktreeByName(tt.query)
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result.Path)
        })
    }
}
```

## Integration Tests

### Test Repository Setup

```go
type TestRepo struct {
    Dir     string
    MainDir string
    Cleanup func()
}

func SetupTestRepo(t *testing.T) *TestRepo {
    tmpDir, err := os.MkdirTemp("", "wrkt-test-*")
    require.NoError(t, err)
    
    mainDir := filepath.Join(tmpDir, "main")
    err = os.MkdirAll(mainDir, 0755)
    require.NoError(t, err)
    
    // Initialize git repo
    cmd := exec.Command("git", "init")
    cmd.Dir = mainDir
    err = cmd.Run()
    require.NoError(t, err)
    
    // Create initial commit
    err = os.WriteFile(filepath.Join(mainDir, "README.md"), []byte("# Test"), 0644)
    require.NoError(t, err)
    
    cmd = exec.Command("git", "add", ".")
    cmd.Dir = mainDir
    err = cmd.Run()
    require.NoError(t, err)
    
    cmd = exec.Command("git", "commit", "-m", "Initial commit")
    cmd.Dir = mainDir
    err = cmd.Run()
    require.NoError(t, err)
    
    return &TestRepo{
        Dir:     tmpDir,
        MainDir: mainDir,
        Cleanup: func() { os.RemoveAll(tmpDir) },
    }
}
```

### Command Integration Tests

```go
func TestListCommand_Integration(t *testing.T) {
    repo := SetupTestRepo(t)
    defer repo.Cleanup()
    
    // Change to repo directory
    oldWd, _ := os.Getwd()
    defer os.Chdir(oldWd)
    os.Chdir(repo.MainDir)
    
    // Create worktrees
    cmd := exec.Command("git", "worktree", "add", "../feature-auth", "-b", "feature/auth")
    err := cmd.Run()
    require.NoError(t, err)
    
    // Test list command
    output, err := execWrktCommand("list")
    require.NoError(t, err)
    
    assert.Contains(t, output, "main")
    assert.Contains(t, output, "feature-auth")
    assert.Contains(t, output, "feature/auth")
}

func execWrktCommand(args ...string) (string, error) {
    cmd := exec.Command("go", append([]string{"run", "main.go"}, args...)...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

### End-to-End Workflow Tests

```go
func TestWorktreeWorkflow_E2E(t *testing.T) {
    repo := SetupTestRepo(t)
    defer repo.Cleanup()
    
    oldWd, _ := os.Getwd()
    defer os.Chdir(oldWd)
    os.Chdir(repo.MainDir)
    
    // 1. Add worktree
    output, err := execWrktCommand("add", "feature/auth")
    require.NoError(t, err)
    assert.Contains(t, output, "Created worktree")
    
    // 2. List worktrees
    output, err = execWrktCommand("list")
    require.NoError(t, err)
    assert.Contains(t, output, "feature-auth")
    
    // 3. Switch to worktree
    output, err = execWrktCommand("switch", "auth")
    require.NoError(t, err)
    
    // 4. Remove worktree
    output, err = execWrktCommand("remove", "auth")
    require.NoError(t, err)
    
    // 5. Verify removal
    output, err = execWrktCommand("list")
    require.NoError(t, err)
    assert.NotContains(t, output, "feature-auth")
}
```

## Test Helpers

### Git Command Mocking

```go
type MockGitCmd struct {
    responses map[string]string
    errors    map[string]error
}

func (m *MockGitCmd) Output(args ...string) ([]byte, error) {
    key := strings.Join(args, " ")
    if err, exists := m.errors[key]; exists {
        return nil, err
    }
    if response, exists := m.responses[key]; exists {
        return []byte(response), nil
    }
    return nil, fmt.Errorf("unexpected git command: %s", key)
}

func NewMockGitCmd() *MockGitCmd {
    return &MockGitCmd{
        responses: make(map[string]string),
        errors:    make(map[string]error),
    }
}
```

### Test Assertions

```go
func AssertWorktreeExists(t *testing.T, path string) {
    _, err := os.Stat(path)
    assert.NoError(t, err, "worktree should exist at %s", path)
}

func AssertWorktreeNotExists(t *testing.T, path string) {
    _, err := os.Stat(path)
    assert.True(t, os.IsNotExist(err), "worktree should not exist at %s", path)
}

func AssertGitWorktreeExists(t *testing.T, name string) {
    cmd := exec.Command("git", "worktree", "list", "--porcelain")
    output, err := cmd.Output()
    require.NoError(t, err)
    assert.Contains(t, string(output), name)
}
```

## Test Categories

### 1. Fast Unit Tests
- Pure functions (parsing, matching, path generation)
- Data structure operations
- Error handling logic
- No external dependencies

**Run with**: `go test ./internal/...`

### 2. Integration Tests  
- Git command integration
- File system operations
- Command execution
- Real git repository interactions

**Run with**: `go test ./test/integration/...`

### 3. End-to-End Tests
- Complete user workflows
- Cross-platform compatibility
- Performance validation
- Real-world scenarios

**Run with**: `go test ./test/e2e/...`

## Test Automation

### Makefile

```makefile
.PHONY: test test-unit test-integration test-e2e test-coverage

# Run all tests
test: test-unit test-integration test-e2e

# Fast unit tests only
test-unit:
	go test -v ./internal/... ./cmd/...

# Integration tests with git
test-integration:
	go test -v ./test/integration/...

# End-to-end tests
test-e2e:
	go test -v ./test/e2e/...

# Coverage report
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Test with race detection
test-race:
	go test -race ./...

# Benchmark tests
test-bench:
	go test -bench=. ./...
```

### GitHub Actions

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]
    
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: make test-unit
    
    - name: Run integration tests
      run: make test-integration
    
    - name: Run e2e tests
      run: make test-e2e
    
    - name: Generate coverage
      run: make test-coverage
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
```

## Testing Best Practices

### 1. Test Naming
- Use descriptive test names: `TestWorktreeManager_List_ReturnsAllWorktrees`
- Include the scenario: `TestFuzzyMatching_ExactMatch_ReturnsWorktree`
- Indicate expected outcome: `TestAdd_InvalidBranch_ReturnsError`

### 2. Test Organization
- Group related tests in the same file
- Use subtests for variations: `t.Run("exact match", func(t *testing.T) {...})`
- Keep test functions focused and small

### 3. Test Data
- Use realistic test data
- Cover edge cases and boundary conditions
- Test with different git repository states

### 4. Test Isolation
- Each test should be independent
- Clean up resources in defer statements
- Use temporary directories for integration tests

### 5. Error Testing
- Test both success and failure paths
- Verify error messages are helpful
- Test error recovery scenarios

## Continuous Testing

### Pre-commit Hooks
```bash
#!/bin/sh
# .git/hooks/pre-commit
make test-unit || exit 1
```

### IDE Integration
- Configure VS Code/GoLand to run tests on save
- Set up test coverage visualization
- Enable test debugging

### Performance Testing
```go
func BenchmarkFuzzyMatching(b *testing.B) {
    worktrees := generateLargeWorktreeList(1000)
    wm := &WorktreeManager{worktrees: worktrees}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = wm.GetWorktreeByName("feature-auth-service")
    }
}
```

## Quality Gates

### Coverage Targets
- Unit tests: >90% coverage
- Integration tests: All critical paths covered
- E2E tests: All user workflows covered

### Performance Targets
- Unit tests: <10ms per test
- Integration tests: <1s per test
- E2E tests: <10s per workflow

### Quality Checks
- All tests pass on multiple Go versions (1.21+)
- Tests pass on multiple platforms (Linux, macOS, Windows)
- No race conditions detected
- Memory usage within acceptable limits