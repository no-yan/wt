# Refactoring Command Execution

This document outlines a plan to refactor the command execution logic within the `wt` tool to improve security and robustness.

## 1. Problem Statement

The current implementation of command execution has two main weaknesses:

1.  **Insecure Command Parsing**: The `internal.ExecCommandRunner.Run` method accepts a single string, which it splits into arguments using `strings.Fields`. This approach is vulnerable to shell injection and does not correctly handle arguments that contain spaces or other special shell characters. For example, a branch name like `"feature/my new feature"` would be incorrectly parsed.

2.  **Brittle Shell Escaping**: The `internal.shellescape` function is a custom implementation for escaping shell characters. While it's a good first step, writing a completely secure shell escaping function is notoriously difficult and prone to errors. It may not cover all edge cases, leaving potential security holes.

These issues make the application less secure and less reliable when dealing with complex or unusually named branches and paths.

## 2. Proposed Solution

The proposed solution is to refactor the `CommandRunner` interface and its implementation to avoid shell interpretation altogether. This is the most secure way to execute external commands.

### 2.1. Modify `CommandRunner` Interface

The `CommandRunner` interface will be changed to pass the command and its arguments as a slice of strings.

**Current Interface:**
```go
type CommandRunner interface {
    Run(command string) (string, error)
}
```

**New Interface:**
```go
type CommandRunner interface {
    Run(name string, args ...string) (string, error)
}
```

### 2.2. Update `ExecCommandRunner`

The `ExecCommandRunner` will be updated to use the new `Run` signature. It will no longer use `strings.Fields` and will instead pass the arguments directly to `exec.Command`.

**Current Implementation:**
```go
func (e *ExecCommandRunner) Run(command string) (string, error) {
    parts := strings.Fields(command)
    // ...
    cmd := exec.Command(parts[0], parts[1:]...)
    // ...
}
```

**New Implementation:**
```go
func (e *ExecCommandRunner) Run(name string, args ...string) (string, error) {
    cmd := exec.Command(name, args...)
    output, err := cmd.Output()
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            return "", fmt.Errorf("command failed: %s: %s", err, exitErr.Stderr)
        }
        return "", fmt.Errorf("command execution failed: %w", err)
    }
    return string(output), nil
}
```

### 2.3. Remove `shellescape`

The `shellescape` function will be removed, as it is no longer needed.

### 2.4. Update Call Sites

All calls to the `runner.Run` method will be updated to the new signature. This primarily affects `internal/git_service.go` and `internal/worktree_manager.go`.

**Example in `WorktreeManager`:**

**Current:**
```go
func (wm *WorktreeManager) addGitWorktree(repoPath, worktreePath, branch string) error {
    gitCmd := fmt.Sprintf("git -C %s worktree add %s %s",
        shellescape(repoPath),
        shellescape(worktreePath),
        shellescape(branch))

    if _, err := wm.runner.Run(gitCmd); err != nil {
        return fmt.Errorf("git worktree add failed: %w", err)
    }
    return nil
}
```

**New:**
```go
func (wm *WorktreeManager) addGitWorktree(repoPath, worktreePath, branch string) error {
    _, err := wm.runner.Run("git", "-C", repoPath, "worktree", "add", worktreePath, branch)
    if err != nil {
        return fmt.Errorf("git worktree add failed: %w", err)
    }
    return nil
}
```

## 3. Benefits

This refactoring will provide several key benefits:

-   **Enhanced Security**: By passing arguments directly to `exec.Command`, we bypass the shell and eliminate the risk of shell injection vulnerabilities. This is the standard, recommended practice for executing commands in Go.
-   **Increased Robustness**: The new implementation will correctly handle branch names, paths, and other arguments that contain spaces, quotes, and other special characters.
-   **Simplified Code**: The removal of the custom `shellescape` function and the simplification of the `Run` method will result in cleaner, more maintainable code.

This change is a critical improvement that will make `wt` a more secure and reliable tool.
