package internal

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
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

// Type implements VCSService interface
func (g *GitService) Type() VCSType {
	return VCSGit
}

// ListWorkspaces implements VCSService interface
func (g *GitService) ListWorkspaces(opts ListOptions) ([]Workspace, error) {
	output, err := g.runner.Run("git worktree list --porcelain")
	if err != nil {
		return nil, err
	}

	workspaces := ParseWorktreeList(output)

	// Only fetch status if requested
	if opts.IncludeStatus {
		for i := range workspaces {
			statusOutput, err := g.runner.Run("git -C " + workspaces[i].Path + " status --porcelain")
			if err != nil {
				workspaces[i].Status = StatusStale
				workspaces[i].Error = err
			} else {
				workspaces[i].Status = ParseWorktreeStatus(statusOutput)
			}
		}
	}

	return workspaces, nil
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

// TODO: Consider refactoring CommandRunner interface to accept arguments separately
// instead of a single command string. This would eliminate the need for shell parsing
// and be more type-safe. Example:
//
//	type CommandRunner interface {
//	    Run(cmd string, args ...string) (string, error)
//	}
//
// This would require updates to all callers but would be cleaner long-term.
func (e *ExecCommandRunner) Run(command string) (string, error) {
	parts, err := shellquote.Split(command)
	if err != nil {
		return "", fmt.Errorf("failed to parse command: %w", err)
	}
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
