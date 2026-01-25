package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DetectVCS determines the VCS type for a given directory.
// Priority: .jj/ (jj) > .git/ (git)
// Reason: jj repos often have .git/ for colocated mode
func DetectVCS(path string) (VCSType, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	return detectVCSRecursive(absPath)
}

func detectVCSRecursive(path string) (VCSType, error) {
	// Check for jj first (handles colocated repos)
	jjPath := filepath.Join(path, ".jj")
	if info, err := os.Stat(jjPath); err == nil && info.IsDir() {
		return VCSJj, nil
	}

	// Check for git
	gitPath := filepath.Join(path, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		// .git can be a file (worktree) or directory
		return VCSGit, nil
	}

	// Walk up to find repo root
	parent := filepath.Dir(path)
	if parent == path {
		return "", fmt.Errorf("no VCS detected in %s or its parents", path)
	}

	return detectVCSRecursive(parent)
}

// GetRepoRoot returns the repository root for the detected VCS
func GetRepoRoot(runner CommandRunner, vcsType VCSType) (string, error) {
	switch vcsType {
	case VCSGit:
		output, err := runner.Run("git rev-parse --show-toplevel")
		if err != nil {
			return "", fmt.Errorf("not in a git repository: %w", err)
		}
		return strings.TrimSpace(output), nil
	case VCSJj:
		output, err := runner.Run("jj workspace root")
		if err != nil {
			return "", fmt.Errorf("not in a jj repository: %w", err)
		}
		return strings.TrimSpace(output), nil
	default:
		return "", fmt.Errorf("unsupported VCS type: %s", vcsType)
	}
}

// CheckVCSAvailable verifies that the VCS binary is installed and accessible
func CheckVCSAvailable(runner CommandRunner, vcsType VCSType) error {
	switch vcsType {
	case VCSGit:
		_, err := runner.Run("git --version")
		if err != nil {
			return fmt.Errorf("git is not installed or not in PATH: %w", err)
		}
		return nil
	case VCSJj:
		_, err := runner.Run("jj --version")
		if err != nil {
			return fmt.Errorf("jj is not installed or not in PATH: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported VCS type: %s", vcsType)
	}
}

// CreateVCSService detects the VCS type and returns the appropriate service.
// This is the preferred way to obtain a VCSService instance.
func CreateVCSService(runner CommandRunner) (VCSService, error) {
	vcsType, err := DetectVCS(".")
	if err != nil {
		return nil, err
	}

	if err := CheckVCSAvailable(runner, vcsType); err != nil {
		return nil, err
	}

	switch vcsType {
	case VCSGit:
		return NewGitService(runner), nil
	case VCSJj:
		return NewJjService(runner), nil
	default:
		return nil, fmt.Errorf("unsupported VCS type: %s", vcsType)
	}
}

// CreateWorkspaceManager detects the VCS type and returns the appropriate manager.
func CreateWorkspaceManager(runner CommandRunner) (WorkspaceManager, VCSType, error) {
	vcsType, err := DetectVCS(".")
	if err != nil {
		return nil, "", err
	}

	if err := CheckVCSAvailable(runner, vcsType); err != nil {
		return nil, "", err
	}

	switch vcsType {
	case VCSGit:
		gitService := NewGitService(runner)
		return NewWorktreeManager(gitService, runner), vcsType, nil
	case VCSJj:
		jjService := NewJjService(runner)
		return NewJjWorkspaceManager(jjService), vcsType, nil
	default:
		return nil, "", fmt.Errorf("unsupported VCS type: %s", vcsType)
	}
}
