package internal

import (
	"path/filepath"
	"strings"
)

func ParseWorktreeList(output string) []Workspace {
	if strings.TrimSpace(output) == "" {
		return []Workspace{}
	}

	worktreeBlocks := strings.Split(strings.TrimSpace(output), "\n\n")
	workspaces := make([]Workspace, 0, len(worktreeBlocks))

	for _, block := range worktreeBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		workspace := parseWorktreeBlock(block)
		if workspace != nil {
			workspaces = append(workspaces, *workspace)
		}
	}

	return workspaces
}

func parseWorktreeBlock(block string) *Workspace {
	lines := strings.Split(block, "\n")
	if len(lines) < 2 {
		return nil
	}

	var path, head, branch string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "worktree ") {
			path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "HEAD ") {
			head = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			branchRef := strings.TrimPrefix(line, "branch ")
			if strings.HasPrefix(branchRef, "refs/heads/") {
				branch = strings.TrimPrefix(branchRef, "refs/heads/")
			} else {
				branch = branchRef
			}
		} else if line == "detached" {
			branch = "detached HEAD"
		}
	}

	if path == "" || head == "" {
		return nil
	}

	// Calculate workspace name
	name := branch
	if strings.Contains(path, "/worktrees/") {
		name = filepath.Base(path)
	} else {
		name = BranchToWorktreeName(branch)
	}

	return &Workspace{
		Name:     name,
		Path:     path,
		Revision: head,
		Labels:   []string{branch},
		Status:   StatusClean,
	}
}

func ParseWorktreeStatus(output string) Status {
	if strings.TrimSpace(output) == "" {
		return StatusClean
	}
	return StatusDirty
}
