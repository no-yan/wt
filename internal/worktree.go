package internal

import (
	"path/filepath"
	"strings"
)

type Status int

const (
	StatusClean Status = iota
	StatusDirty
	StatusStale
)

type Worktree struct {
	Branch string
	Path   string
	Head   string
	Status Status
}

func (w Worktree) IsClean() bool {
	return w.Status == StatusClean
}

func (w Worktree) Name() string {
	return filepath.Base(w.Path)
}

func BranchToWorktreeName(branch string) string {
	return strings.ReplaceAll(branch, "/", "-")
}
