package internal

import (
	"strings"
)

type Status int

const (
	StatusClean Status = iota
	StatusDirty
	StatusStale
)

func BranchToWorktreeName(branch string) string {
	return strings.ReplaceAll(branch, "/", "-")
}
