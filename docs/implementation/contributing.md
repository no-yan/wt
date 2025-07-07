# Contributing Guide

Welcome to the `wt` project! This guide will help you contribute effectively.

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Git 2.5 or higher (for worktree support)
- Make (optional, for automation)

### Setup Instructions

1. **Clone the repository**
   ```bash
   git clone https://github.com/no-yan/wt.git
   cd wt
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the project**
   ```bash
   go build -o wt
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

5. **Install locally** (optional)
   ```bash
   go install
   ```

### Development Environment

#### Recommended Tools
- **Editor**: VS Code with Go extension, or GoLand
- **Linting**: golangci-lint
- **Testing**: Built-in Go testing tools
- **Debugging**: Delve debugger

#### Project Structure
```
wt/
├── main.go              # Entry point
├── cmd/                 # CLI commands
├── internal/            # Core logic
├── docs/                # Documentation
├── test/                # Test files
├── go.mod               # Go module
└── Makefile            # Build automation
```

## Development Workflow

### 1. Issue-Based Development

1. **Check existing issues** before starting work
2. **Create an issue** for new features or bugs
3. **Assign yourself** to the issue
4. **Reference the issue** in commit messages

### 2. Branch Strategy

- `main` - Stable, production-ready code
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `hotfix/*` - Critical fixes

```bash
# Create feature branch
git checkout -b feature/fuzzy-matching

# Work on changes
git add .
git commit -m "feat: implement fuzzy matching for worktree names"

# Push and create PR
git push origin feature/fuzzy-matching
```

### 3. Commit Message Format

Follow conventional commits:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes
- `refactor` - Code refactoring
- `test` - Test additions/changes
- `chore` - Maintenance tasks

**Examples:**
```
feat(cmd): add fuzzy matching to switch command
fix(internal): handle empty worktree list gracefully
docs(api): update command reference for list command
test(integration): add end-to-end workflow tests
```

## Code Standards

### Go Style Guide

Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and these project-specific guidelines:

#### 1. Formatting
- Use `gofmt` for formatting
- Use `goimports` for import organization
- Maximum line length: 100 characters

#### 2. Naming Conventions
```go
// Good
type WorktreeManager struct { ... }
func (wm *WorktreeManager) List() ([]*Worktree, error) { ... }

// Bad
type worktree_manager struct { ... }
func (w *worktree_manager) get_list() ([]*worktree, error) { ... }
```

#### 3. Error Handling
```go
// Good - wrap errors with context
if err != nil {
    return fmt.Errorf("failed to parse worktree list: %w", err)
}

// Bad - return raw errors
if err != nil {
    return err
}
```

#### 4. Documentation
```go
// WorktreeManager handles all git worktree operations for a repository.
// It maintains the repository root path and provides methods for listing,
// creating, and managing worktrees.
type WorktreeManager struct {
    repoRoot string
}

// List returns all worktrees in the current repository.
// It parses the output of 'git worktree list --porcelain' and returns
// a slice of Worktree structs with complete metadata.
func (wm *WorktreeManager) List() ([]*Worktree, error) {
    // implementation
}
```

### Code Quality

#### 1. Linting
Run linters before submitting:
```bash
golangci-lint run
go vet ./...
```

#### 2. Testing Requirements
- All new code must have tests
- Maintain >80% test coverage
- Include both unit and integration tests

#### 3. Performance
- Commands should execute in <100ms for typical use cases
- Use benchmarks for performance-critical code
- Profile memory usage for large repositories

## Testing Guidelines

### Running Tests

```bash
# All tests
go test ./...

# Unit tests only
go test ./internal/... ./cmd/...

# Integration tests
go test ./test/integration/...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

#### 1. Test Structure
```go
func TestWorktreeManager_List(t *testing.T) {
    tests := []struct {
        name     string
        setup    func() *WorktreeManager
        expected []*Worktree
        wantErr  bool
    }{
        {
            name: "single worktree",
            setup: func() *WorktreeManager {
                // setup test data
            },
            expected: []*Worktree{
                {Path: "/path/to/main", Branch: "main"},
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            wm := tt.setup()
            result, err := wm.List()
            
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

#### 2. Integration Tests
```go
func TestListCommand_Integration(t *testing.T) {
    // Create temporary git repository
    tmpDir := t.TempDir()
    setupGitRepo(t, tmpDir)
    
    // Test the command
    cmd := exec.Command("./wt", "list")
    cmd.Dir = tmpDir
    output, err := cmd.Output()
    
    require.NoError(t, err)
    assert.Contains(t, string(output), "main")
}
```

## Pull Request Process

### 1. Before Submitting

- [ ] All tests pass
- [ ] Code is properly formatted
- [ ] Linting passes
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated (if applicable)

### 2. PR Description Template

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
```

### 3. Review Process

1. **Automated checks** must pass (CI/CD)
2. **Code review** by at least one maintainer
3. **Address feedback** promptly
4. **Squash commits** if requested
5. **Merge** when approved

## Architecture Guidelines

### 1. Module Separation

- **`cmd/`** - CLI interface only, no business logic
- **`internal/`** - Core business logic, git operations
- **`test/`** - Test utilities and integration tests

### 2. Error Handling Strategy

```go
// Wrap errors with context
func (wm *WorktreeManager) Add(branch, path string) error {
    if err := wm.validateBranch(branch); err != nil {
        return fmt.Errorf("invalid branch %q: %w", branch, err)
    }
    
    if err := wm.createWorktree(branch, path); err != nil {
        return fmt.Errorf("failed to create worktree: %w", err)
    }
    
    return nil
}
```

### 3. Configuration Pattern

```go
// Use options pattern for complex configurations
type AddOptions struct {
    Force      bool
    NewBranch  bool
    Detach     bool
    LockReason string
}

func (wm *WorktreeManager) AddWithOptions(branch, path string, opts AddOptions) error {
    // implementation
}
```

## Documentation Standards

### 1. Code Documentation

- Document all public types and functions
- Include examples for complex functions
- Explain the "why" not just the "what"

### 2. API Documentation

- Keep `docs/api/commands.md` updated
- Include examples for all commands
- Document all options and flags

### 3. Implementation Documentation

- Update architecture docs for significant changes
- Document design decisions in `docs/implementation/`
- Include migration guides for breaking changes

## Release Process

### 1. Version Numbering

Follow semantic versioning (SemVer):
- `MAJOR.MINOR.PATCH`
- `MAJOR` - Breaking changes
- `MINOR` - New features, backward compatible
- `PATCH` - Bug fixes, backward compatible

### 2. Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped in relevant files
- [ ] Git tag created
- [ ] Release notes written

### 3. Changelog Format

```markdown
## [1.2.0] - 2024-01-15

### Added
- Fuzzy matching for worktree names
- Auto-path generation from branch names

### Changed
- Improved error messages for git operations

### Fixed
- Handle empty worktree list gracefully

### Deprecated
- Old path generation will be removed in v2.0

### Security
- Validate all user inputs to prevent path traversal
```

## Getting Help

### 1. Documentation
- Check existing documentation in `docs/`
- Read the architecture guide for design decisions
- Look at existing code for patterns

### 2. Communication
- Create an issue for questions or problems
- Use descriptive titles and provide context
- Include relevant code snippets or error messages

### 3. Debugging
- Use the Go debugger (delve) for complex issues
- Add logging for debugging (remove before PR)
- Write tests to reproduce issues

## Common Tasks

### Adding a New Command

1. Create command file in `cmd/`
2. Implement command logic in `internal/`
3. Add tests for both CLI and logic
4. Update documentation
5. Add to root command

### Modifying Git Integration

1. Update `internal/git.go` or `internal/worktree.go`
2. Add comprehensive tests with real git repositories
3. Handle error cases gracefully
4. Test with different git versions if applicable

### Adding Configuration Options

1. Define configuration structure
2. Add parsing logic
3. Update command-line flags
4. Add validation
5. Update documentation

## Code Review Guidelines

### For Authors

- Keep PRs small and focused
- Write clear commit messages
- Test thoroughly before submitting
- Respond to feedback promptly

### For Reviewers

- Be constructive and specific
- Check for edge cases and error handling
- Verify tests cover new functionality
- Consider performance implications
- Ensure documentation is updated

## Performance Considerations

### 1. Git Command Optimization
- Minimize git command invocations
- Use porcelain output for reliable parsing
- Cache results when appropriate

### 2. Memory Usage
- Avoid loading large amounts of data into memory
- Use streaming for large outputs
- Profile memory usage for large repositories

### 3. Response Time
- Target <100ms for common operations
- Use benchmarks to measure performance
- Consider async operations for long-running tasks

Thank you for contributing to `wt`! Your efforts help make git worktree operations smoother for everyone.