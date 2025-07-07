# Architecture Documentation

Technical architecture and design decisions for the `wt` project.

## System Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Zsh Layer     │    │   CLI Commands   │    │  Core Logic     │
│   (functions)   │───▶│   (cobra)        │───▶│  (internal/)    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Directory Ops   │    │  User Output     │    │  Worktrees/     │
│ (cd, completion)│    │  (display)       │    │  Organization   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Key Architecture Principles**:
1. **Zsh-only**: Simplified shell integration for single shell
2. **Organized Structure**: All worktrees in `worktrees/` subdirectory  
3. **Exact Matching**: No fuzzy matching complexity
4. **Simple Path Generation**: Deterministic branch → directory mapping

## Module Structure

### `/cmd` - CLI Commands
- **Purpose**: Cobra command definitions and CLI interface
- **Pattern**: One file per command
- **Key Files**:
  - `list.go` - Unified status display with filtering
  - `switch.go` - Path resolution for exact name (no directory change)
  - `shell.go` - Zsh integration code generation
  - `add.go` - Worktree creation in `worktrees/` subdirectory
  - `remove.go`, `clean.go` - Worktree management
- **Responsibilities**:
  - Argument parsing and validation
  - Flag handling
  - User interaction (prompts, confirmations)
  - Output formatting
  - Error handling at CLI level
  - **NOT responsible for**: Actual directory switching (handled by zsh layer)

### `/internal` - Core Logic
- **Purpose**: Business logic and data operations
- **Pattern**: Functionality-based modules
- **Key Files**:
  - `worktree.go` - Core worktree operations and worktrees/ organization
  - `shell.go` - Zsh function generation and integration logic
  - `git.go` - Git command integration and parsing
  - `display.go` - Output formatting and status computation
- **Responsibilities**:
  - Git worktree operations in `worktrees/` subdirectory
  - Data parsing and transformation
  - Simple path generation algorithms
  - Exact name matching logic
  - Status computation
  - **Zsh function generation** (critical component)
  - Auto-setup of worktrees directory and .gitignore

### `/docs` - Documentation
- **Purpose**: Comprehensive documentation
- **Structure**:
  - `api/` - Command reference and API docs
  - `implementation/` - Technical documentation
  - `examples/` - Usage examples and workflows

## Core Components

### 1. Worktree Organization System (`internal/worktree.go`)

**THE FOUNDATION** - predictable worktree structure.

```go
type WorktreeManager struct {
    repoRoot string
}

// Auto-setup and organization
func (wm *WorktreeManager) EnsureWorktreesDir() error
func (wm *WorktreeManager) GenerateWorktreePath(branch string) string
func (wm *WorktreeManager) AddToGitignore() error
```

**Design Decisions**:
- All worktrees in `$REPO_ROOT/worktrees/` subdirectory
- Auto-create directory and add to .gitignore on first use
- Simple path generation: `feature/auth` → `worktrees/feature-auth/`
- Predictable, organized structure eliminates path confusion

### 2. Zsh Integration System (`internal/shell.go`)

**CRITICAL FOR NAVIGATION** - enables directory switching functionality.

```go
type ShellIntegration struct{}

// Generate zsh functions only
func (si *ShellIntegration) GenerateZshInit() string
```

**Design Decisions**:
- Zsh-only implementation eliminates multi-shell complexity
- Generates zsh functions that wrap the `wt` binary
- Intercepts `wt switch` calls to perform actual `cd` operations
- Provides tab completion for all commands
- Clear error message for non-zsh users

### 3. Core Worktree Operations (`internal/worktree.go`)

Central component for all worktree operations.

```go
// Core operations
func (wm *WorktreeManager) List() ([]*Worktree, error)
func (wm *WorktreeManager) Add(branch string) error  // Always creates in worktrees/
func (wm *WorktreeManager) Remove(name string) error
func (wm *WorktreeManager) GetWorktreeByExactName(name string) (*Worktree, error)
func (wm *WorktreeManager) GetWorktreePath(name string) (string, error) // For zsh integration
```

**Design Decisions**:
- Single source of truth for worktree operations
- All worktrees created in `worktrees/` subdirectory
- Exact name matching only - no fuzzy matching complexity
- Error handling with context-rich messages and suggestions
- Git porcelain output parsing for reliability
- **Path resolution method for zsh integration**

### 4. Unified List Command (`cmd/list.go`)

**Replaces separate status functionality** with flexible filtering.

```go
type ListOptions struct {
    Dirty      bool // Show only worktrees with changes
    Verbose    bool // Show detailed git status
    Porcelain  bool // Machine-readable output
    NamesOnly  bool // Names only for completion
}
```

**Design Decisions**:
- Single command handles all status display needs
- Shows worktrees in `worktrees/` subdirectory
- Reduces API surface and user confusion
- Flexible filtering replaces multiple commands
- Supports both human and machine consumption

### 5. Worktree Model (`internal/worktree.go`)

Data structure representing a git worktree in `worktrees/` subdirectory.

```go
type Worktree struct {
    Path     string // Path within worktrees/ subdirectory
    Name     string // Worktree name (directory name)
    Head     string // Current HEAD commit hash
    Branch   string // Current branch name
    Bare     bool   // Is bare repository
    Detached bool   // Is HEAD detached
    Locked   bool   // Is worktree locked
    Prunable bool   // Can be pruned
    Reason   string // Lock reason (if locked)
    Status   WorktreeStatus // Computed status indicators
}
```

**Design Decisions**:
- All worktrees are in `worktrees/` subdirectory
- Name field represents exact worktree directory name
- Maps directly to git worktree list output
- Supports all git worktree features
- **Added computed status for display**

### 5. Command Pattern (`cmd/*.go`)

Each command follows a consistent pattern with special handling for shell integration:

```go
var switchCmd = &cobra.Command{
    Use:   "switch <name>",
    Short: "Switch to a worktree directory",
    RunE:  runSwitch,
}

func runSwitch(cmd *cobra.Command, args []string) error {
    // For switch command:
    // 1. Resolve worktree name to path
    // 2. Return path only (shell function handles cd)
    // 3. Exit with error code if not found
    
    // For other commands:
    // 1. Validate arguments
    // 2. Create WorktreeManager  
    // 3. Execute business logic
    // 4. Format and display output
    // 5. Handle errors
}
```

**Design Decisions**:
- **Switch command is special**: only resolves paths
- Separation of CLI concerns from business logic
- Consistent error handling patterns
- Structured output formatting
- User-friendly error messages with suggestions

## Data Flow

### Shell Integration Flow (Critical)
```
User: wt switch auth
  ↓
Shell Function: wt() { case switch → }
  ↓  
Binary: wt path auth
  ↓
WorktreeManager: GetWorktreePath("auth")
  ↓
Shell Function: cd "$target_path"
  ↓
User Directory Changed
```

### List Command Flow
```
User Input → Argument Parsing → WorktreeManager.List() → 
Git Command → Parse Output → Apply Filters → Format Display → User Output
```

### Add Command Flow
```
User Input → Path Generation → WorktreeManager.Add() →
Git Worktree Add → Validation → Success/Error Response
```

### Shell Integration Setup Flow
```
User: eval "$(wt shell-init)"
  ↓
Binary: ShellIntegration.GenerateBashInit()
  ↓
Shell: Loads wt() function and completion
  ↓
User: Can now use wt switch functionality
```

## Key Algorithms

### 1. Fuzzy Matching Algorithm

```go
func (wm *WorktreeManager) GetWorktreeByName(name string) (*Worktree, error) {
    // Priority order:
    // 1. Exact match on basename
    // 2. Substring match on branch name
    // 3. Substring match on full path
    // 4. Fuzzy match with edit distance
}
```

**Rationale**: Provides intuitive matching while maintaining predictable behavior.

### 2. Path Generation Algorithm

```go
func generateWorktreePath(repoRoot, branch string) string {
    // 1. Remove common prefixes (feature/, hotfix/, etc.)
    // 2. Replace slashes with dashes
    // 3. Create path relative to repository parent
    // 4. Ensure uniqueness
}
```

**Rationale**: Creates logical, conflict-free paths while maintaining readability.

### 3. Status Computation

```go
func (w *Worktree) ComputeStatus() WorktreeStatus {
    // 1. Check git status for uncommitted changes
    // 2. Compare with remote branches
    // 3. Determine status indicators
    // 4. Cache results for performance
}
```

**Rationale**: Provides rich status information while minimizing git command overhead.

## Error Handling Strategy

### 1. Error Types
- **Git Errors**: Repository not found, invalid worktree state
- **File System Errors**: Permission denied, path conflicts
- **User Errors**: Invalid arguments, missing worktrees
- **System Errors**: Command execution failures

### 2. Error Handling Patterns
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to list worktrees: %w", err)
}

// Provide actionable error messages
if !isGitRepo {
    return errors.New("not in a git repository (run 'git init' or 'cd' to a git repository)")
}
```

### 3. Error Recovery
- Graceful degradation when possible
- Clear instructions for resolution
- Fail-fast for unrecoverable errors

## Performance Considerations

### 1. Git Command Optimization
- Use porcelain output for parsing reliability
- Minimize git command invocations
- Cache worktree list for multiple operations

### 2. Fuzzy Matching Performance
- Early termination on exact matches
- Efficient string matching algorithms
- Limit search scope for large repositories

### 3. Output Formatting
- Stream output for large worktree lists
- Lazy evaluation of status information
- Parallel status computation where safe

## Testing Strategy

### 1. Unit Tests
- Worktree parsing logic
- Path generation algorithms
- Fuzzy matching functions
- Error handling paths

### 2. Integration Tests
- Git command integration
- File system operations
- End-to-end command flows
- Cross-platform compatibility

### 3. Property-Based Tests
- Path generation uniqueness
- Fuzzy matching consistency
- Status computation accuracy

## Security Considerations

### 1. Input Validation
- Sanitize all user inputs
- Validate git repository state
- Prevent path traversal attacks

### 2. File System Safety
- Check permissions before operations
- Validate paths before creation
- Prevent accidental deletion of important files

### 3. Git Integration Security
- Use safe git commands only
- Validate git output parsing
- Handle malicious repository states

## Extension Points

### 1. Configuration System
- YAML/JSON configuration files
- Environment variable support
- Per-repository settings

### 2. Plugin Architecture
- Custom status providers
- Additional output formats
- Integration with other tools

### 3. Shell Integration
- Advanced shell functions
- Completion scripts
- Directory stack management

## Migration and Compatibility

### 1. Git Version Compatibility
- Support git 2.5+ (when worktree was introduced)
- Graceful handling of older git versions
- Feature detection for newer git capabilities

### 2. Cross-Platform Support
- Windows path handling
- macOS and Linux compatibility
- Shell integration across platforms

### 3. Backward Compatibility
- Stable command interface
- Deprecation warnings for changes
- Migration guides for breaking changes