package cmd

import (
	"testing"

	"github.com/no-yan/wt/internal"
)

func TestFindWorktreeByName(t *testing.T) {
	tests := []struct {
		name      string
		worktrees []internal.Worktree
		target    string
		want      string
		wantErr   bool
	}{
		{
			name: "exact match found",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Branch: "main",
				},
				{
					Path:   "/repo/worktrees/feature-auth",
					Branch: "feature/auth",
				},
			},
			target:  "feature-auth",
			want:    "/repo/worktrees/feature-auth",
			wantErr: false,
		},
		{
			name: "main branch match",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Branch: "main",
				},
				{
					Path:   "/repo/worktrees/feature-auth",
					Branch: "feature/auth",
				},
			},
			target:  "main",
			want:    "/repo",
			wantErr: false,
		},
		{
			name: "no match found",
			worktrees: []internal.Worktree{
				{
					Path:   "/repo",
					Branch: "main",
				},
			},
			target:  "nonexistent",
			want:    "",
			wantErr: true,
		},
		{
			name:      "empty worktrees list",
			worktrees: []internal.Worktree{},
			target:    "anything",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findWorktreeByName(tt.worktrees, tt.target)

			if (err != nil) != tt.wantErr {
				t.Errorf("findWorktreeByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("findWorktreeByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
