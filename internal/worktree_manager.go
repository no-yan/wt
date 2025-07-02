package internal

import (
	"path/filepath"
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
	worktreePath := GenerateWorktreePath(repoPath, branch)
	worktreesDir := filepath.Dir(worktreePath)

	if _, err := wm.runner.Run("mkdir -p " + worktreesDir); err != nil {
		return "", err
	}

	if _, err := wm.runner.Run("git -C " + repoPath + " worktree add " + worktreePath + " " + branch); err != nil {
		return "", err
	}

	return worktreePath, nil
}

func GenerateWorktreePath(repoPath, branch string) string {
	worktreeName := BranchToWorktreeName(branch)
	return filepath.Join(repoPath, "worktrees", worktreeName)
}
