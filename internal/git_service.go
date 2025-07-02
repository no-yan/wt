package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

type CommandRunner interface {
	Run(command string) (string, error)
}

type GitService struct {
	runner CommandRunner
}

func NewGitService(runner CommandRunner) *GitService {
	return &GitService{runner: runner}
}

func (g *GitService) ListWorktrees() ([]Worktree, error) {
	output, err := g.runner.Run("git worktree list --porcelain")
	if err != nil {
		return nil, err
	}

	worktrees := ParseWorktreeList(output)

	for i := range worktrees {
		statusOutput, err := g.runner.Run("git -C " + worktrees[i].Path + " status --porcelain")
		if err != nil {
			worktrees[i].Status = StatusStale
		} else {
			worktrees[i].Status = ParseWorktreeStatus(statusOutput)
		}
	}

	return worktrees, nil
}

type ExecCommandRunner struct{}

func NewExecCommandRunner() *ExecCommandRunner {
	return &ExecCommandRunner{}
}

func (e *ExecCommandRunner) Run(command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("command failed: %s: %s", err, exitErr.Stderr)
		}
		return "", fmt.Errorf("command execution failed: %w", err)
	}
	return string(output), nil
}
