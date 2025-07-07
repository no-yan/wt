package cmd

import (
	"bytes"
	"testing"

	"github.com/no-yan/wt/internal"
)

func TestListCommand(t *testing.T) {
	tests := []struct {
		name      string
		worktrees []internal.Worktree
		want      string
	}{
		{
			name: "single clean worktree",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Head:   "abc123",
					Branch: "main",
					Status: internal.StatusClean,
				},
			},
			want: "main\t/repo\t(clean)\n",
		},
		{
			name: "multiple worktrees mixed status",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Head:   "abc123",
					Branch: "main",
					Status: internal.StatusClean,
				},
				{
					Path:   "/repo/worktrees/feature-auth",
					Head:   "def456",
					Branch: "feature/auth",
					Status: internal.StatusDirty,
				},
			},
			want: "main\t/repo\t(clean)\nfeature-auth\t/repo/worktrees/feature-auth\t(dirty)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			formatWorktreeList(tt.worktrees, &output)

			if got := output.String(); got != tt.want {
				t.Errorf("formatWorktreeList() = %q, want %q", got, tt.want)
			}
		})
	}
}
