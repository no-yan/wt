package internal

import (
	"strings"
)

func ParseWorktreeList(output string) []Worktree {
	if strings.TrimSpace(output) == "" {
		return []Worktree{}
	}

	worktreeBlocks := strings.Split(strings.TrimSpace(output), "\n\n")
	worktrees := make([]Worktree, 0, len(worktreeBlocks))

	for _, block := range worktreeBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		worktree := parseWorktreeBlock(block)
		if worktree != nil {
			worktrees = append(worktrees, *worktree)
		}
	}

	return worktrees
}

func parseWorktreeBlock(block string) *Worktree {
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

	return &Worktree{
		Path:   path,
		Head:   head,
		Branch: branch,
		Status: StatusClean,
	}
}

func ParseWorktreeStatus(output string) Status {
	if strings.TrimSpace(output) == "" {
		return StatusClean
	}
	return StatusDirty
}
