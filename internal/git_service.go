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

func (g *GitService) GetDetailedStatus(worktreePath string) ([]string, error) {
	output, err := g.runner.Run("git -C " + worktreePath + " status --porcelain")
	if err != nil {
		return nil, fmt.Errorf("failed to get status for %s: %w", worktreePath, err)
	}

	if strings.TrimSpace(output) == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var statusLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse git status porcelain format
		if len(line) >= 3 {
			indexStatus := line[0]
			workTreeStatus := line[1]
			filename := line[3:]

			var status string
			switch {
			case indexStatus == 'M' && workTreeStatus == ' ':
				status = "modified (staged)"
			case indexStatus == ' ' && workTreeStatus == 'M':
				status = "modified"
			case indexStatus == 'A' && workTreeStatus == ' ':
				status = "added (staged)"
			case indexStatus == ' ' && workTreeStatus == 'A':
				status = "added"
			case indexStatus == 'D' && workTreeStatus == ' ':
				status = "deleted (staged)"
			case indexStatus == ' ' && workTreeStatus == 'D':
				status = "deleted"
			case indexStatus == '?' && workTreeStatus == '?':
				status = "untracked"
			case indexStatus == 'R':
				status = "renamed"
			case indexStatus == 'C':
				status = "copied"
			default:
				status = "modified"
			}

			statusLines = append(statusLines, fmt.Sprintf("%s: %s", status, filename))
		}
	}

	return statusLines, nil
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
