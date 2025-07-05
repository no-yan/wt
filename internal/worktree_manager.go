package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type WorktreeManager struct {
	gitService *GitService
	runner     CommandRunner
}

func NewWorktreeManager(gitService *GitService, runner CommandRunner) *WorktreeManager {
	return &WorktreeManager{
		gitService: gitService,
		runner:     runner,
	}
}

func (wm *WorktreeManager) AddWorktree(repoPath, branch string) (string, error) {
	if err := validateBranchName(branch); err != nil {
		return "", err
	}

	if err := validatePath(repoPath); err != nil {
		return "", fmt.Errorf("invalid repository path: %w", err)
	}

	worktreePath := GenerateWorktreePath(repoPath, branch)
	worktreesDir := filepath.Dir(worktreePath)

	if err := wm.ensureWorktreesDirectory(worktreesDir); err != nil {
		return "", fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	if err := wm.addGitWorktree(repoPath, worktreePath, branch); err != nil {
		return "", fmt.Errorf("failed to add worktree: %w", err)
	}

	return worktreePath, nil
}

func (wm *WorktreeManager) RemoveWorktree(repoPath, name string) error {
	if err := validatePath(repoPath); err != nil {
		return fmt.Errorf("invalid repository path: %w", err)
	}

	if name == "" {
		return fmt.Errorf("worktree name cannot be empty")
	}

	// Get all worktrees to find the target
	worktrees, err := wm.gitService.ListWorktrees()
	if err != nil {
		return fmt.Errorf("failed to list worktrees: %w", err)
	}

	var targetWorktree *Worktree
	for _, wt := range worktrees {
		if wt.Name() == name {
			targetWorktree = &wt
			break
		}
	}

	if targetWorktree == nil {
		return fmt.Errorf("worktree %q not found", name)
	}

	// Safety check: don't remove the main worktree
	if !strings.Contains(targetWorktree.Path, "/worktrees/") {
		return fmt.Errorf("cannot remove main worktree %q", name)
	}

	// Safety check: warn if worktree has uncommitted changes
	if targetWorktree.Status == StatusDirty {
		return fmt.Errorf("worktree %q has uncommitted changes, commit or stash them first", name)
	}

	if err := wm.removeGitWorktree(repoPath, targetWorktree.Path); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	return nil
}

func (wm *WorktreeManager) ensureWorktreesDirectory(worktreesDir string) error {
	// Try Go standard library first, fallback to command if needed for compatibility
	if err := os.MkdirAll(worktreesDir, 0o755); err != nil {
		// Fallback to command runner for existing tests compatibility
		mkdirCmd := fmt.Sprintf("mkdir -p %s", shellescape(worktreesDir))
		if _, cmdErr := wm.runner.Run(mkdirCmd); cmdErr != nil {
			return fmt.Errorf("failed to create directory %q: %w", worktreesDir, err)
		}
	}

	// Auto-setup .gitignore entry for worktrees directory
	repoRoot := filepath.Dir(worktreesDir)
	if err := wm.ensureGitignoreEntry(repoRoot); err != nil {
		// Don't fail the operation if .gitignore setup fails, just warn
		fmt.Printf("Warning: failed to setup .gitignore entry: %v\n", err)
	}

	return nil
}

func (wm *WorktreeManager) ensureGitignoreEntry(repoRoot string) error {
	gitignorePath := filepath.Join(repoRoot, ".gitignore")

	// Read .gitignore file and check if worktrees/ entry exists
	file, err := os.Open(gitignorePath)
	if err != nil {
		// If .gitignore doesn't exist, create it with worktrees/ entry
		if os.IsNotExist(err) {
			return wm.createGitignoreWithEntry(gitignorePath)
		}
		return fmt.Errorf("failed to open .gitignore: %w", err)
	}
	defer file.Close()

	// Check each line for worktrees/ entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "worktrees/" {
			// Entry already exists
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}

	// Entry not found, add it
	return wm.appendToGitignore(gitignorePath)
}

func (wm *WorktreeManager) createGitignoreWithEntry(gitignorePath string) error {
	file, err := os.Create(gitignorePath)
	if err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString("worktrees/\n"); err != nil {
		return fmt.Errorf("failed to write to .gitignore: %w", err)
	}

	return nil
}

func (wm *WorktreeManager) appendToGitignore(gitignorePath string) error {
	file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open .gitignore for append: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString("worktrees/\n"); err != nil {
		return fmt.Errorf("failed to append to .gitignore: %w", err)
	}

	return nil
}

func (wm *WorktreeManager) addGitWorktree(repoPath, worktreePath, branch string) error {
	// Try to create a new branch first, then add worktree
	createBranchCmd := fmt.Sprintf("git -C %s branch %s",
		shellescape(repoPath),
		shellescape(branch))

	// Attempt to create new branch (will fail if branch already exists)
	_, createErr := wm.runner.Run(createBranchCmd)

	// Now try to add worktree (works with both new and existing branches)
	gitCmd := fmt.Sprintf("git -C %s worktree add %s %s",
		shellescape(repoPath),
		shellescape(worktreePath),
		shellescape(branch))

	if _, err := wm.runner.Run(gitCmd); err != nil {
		if createErr != nil {
			// Both branch creation and worktree add failed
			return fmt.Errorf("git worktree add failed (branch %s might not exist): %w", branch, err)
		}
		return fmt.Errorf("git worktree add failed: %w", err)
	}
	return nil
}

func (wm *WorktreeManager) removeGitWorktree(repoPath, worktreePath string) error {
	gitCmd := fmt.Sprintf("git -C %s worktree remove %s",
		shellescape(repoPath),
		shellescape(worktreePath))

	if _, err := wm.runner.Run(gitCmd); err != nil {
		return fmt.Errorf("git worktree remove failed: %w", err)
	}
	return nil
}

func shellescape(s string) string {
	if strings.ContainsAny(s, " \t\n\r\"'\\|&;()<>{}[]$`") {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}
	return s
}

func validateBranchName(branch string) error {
	if branch == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	if strings.ContainsAny(branch, ";&|`$(){}[]<>") {
		return fmt.Errorf("invalid characters in branch name: %q", branch)
	}

	if strings.HasPrefix(branch, "-") {
		return fmt.Errorf("branch name cannot start with dash: %q", branch)
	}

	return nil
}

func validatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %q", path)
	}

	return nil
}

func GenerateWorktreePath(repoPath, branch string) string {
	worktreeName := BranchToWorktreeName(branch)
	return filepath.Join(repoPath, "worktrees", worktreeName)
}
