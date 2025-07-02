package internal

import (
	"regexp"
	"strings"
)

var worktreeLineRegex = regexp.MustCompile(`^(\S+)\s+/(\S+)\s+\[(.+)\]$|^(\S+)\s+/(\S+)\s+\((.+)\)$`)

func ParseWorktreeList(output string) []Worktree {
	if strings.TrimSpace(output) == "" {
		return []Worktree{}
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	worktrees := make([]Worktree, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := worktreeLineRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		var path, head, branch string
		if matches[1] != "" {
			path, head, branch = matches[1], matches[2], matches[3]
		} else {
			path, head, branch = matches[4], matches[5], matches[6]
		}

		worktree := Worktree{
			Path:   path,
			Head:   head,
			Branch: branch,
			Status: StatusClean,
		}

		worktrees = append(worktrees, worktree)
	}

	return worktrees
}

func ParseWorktreeStatus(output string) Status {
	if strings.TrimSpace(output) == "" {
		return StatusClean
	}
	return StatusDirty
}
