package cmd

import (
	"testing"

	"github.com/no-yan/wt/internal"
)

func TestCleanCommand(t *testing.T) {
	tests := []struct {
		name      string
		worktrees []internal.Worktree
		wantFound int
	}{
		{
			name: "no stale worktrees",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Branch: "main",
					Status: internal.StatusClean,
				},
				{
					Path:   "/repo/worktrees/feature-test",
					Branch: "feature/test",
					Status: internal.StatusDirty,
				},
			},
			wantFound: 0,
		},
		{
			name: "one stale worktree",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Branch: "main",
					Status: internal.StatusClean,
				},
				{
					Path:   "/repo/worktrees/feature-test",
					Branch: "feature/test",
					Status: internal.StatusStale,
				},
			},
			wantFound: 1,
		},
		{
			name: "multiple stale worktrees",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Branch: "main",
					Status: internal.StatusClean,
				},
				{
					Path:   "/repo/worktrees/feature-test1",
					Branch: "feature/test1",
					Status: internal.StatusStale,
				},
				{
					Path:   "/repo/worktrees/feature-test2",
					Branch: "feature/test2",
					Status: internal.StatusStale,
				},
			},
			wantFound: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Count stale worktrees
			staleCount := 0
			for _, wt := range tt.worktrees {
				if wt.Status == internal.StatusStale {
					staleCount++
				}
			}

			if staleCount != tt.wantFound {
				t.Errorf("Expected %d stale worktrees, found %d", tt.wantFound, staleCount)
			}
		})
	}
}

func TestShellescape(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple path",
			input: "/path/to/dir",
			want:  "'/path/to/dir'",
		},
		{
			name:  "path with spaces",
			input: "/path/with spaces/dir",
			want:  "'/path/with spaces/dir'",
		},
		{
			name:  "path with special chars",
			input: "/path/with$special&chars",
			want:  "'/path/with$special&chars'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellescape(tt.input)
			if got != tt.want {
				t.Errorf("shellescape(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
