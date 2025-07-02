package internal

import (
	"testing"
)

func TestWorktree_IsClean(t *testing.T) {
	tests := []struct {
		name     string
		worktree Worktree
		want     bool
	}{
		{
			name: "clean worktree",
			worktree: Worktree{
				Branch: "main",
				Path:   "/repo/worktrees/main",
				Head:   "abc123",
				Status: StatusClean,
			},
			want: true,
		},
		{
			name: "dirty worktree",
			worktree: Worktree{
				Branch: "feature/auth",
				Path:   "/repo/worktrees/feature-auth",
				Head:   "def456",
				Status: StatusDirty,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.worktree.IsClean(); got != tt.want {
				t.Errorf("Worktree.IsClean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorktree_Name(t *testing.T) {
	tests := []struct {
		name     string
		worktree Worktree
		want     string
	}{
		{
			name: "main branch",
			worktree: Worktree{
				Branch: "main",
				Path:   "/repo/worktrees/main",
			},
			want: "main",
		},
		{
			name: "feature branch",
			worktree: Worktree{
				Branch: "feature/auth",
				Path:   "/repo/worktrees/feature-auth",
			},
			want: "feature-auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.worktree.Name(); got != tt.want {
				t.Errorf("Worktree.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranchToWorktreeName(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		want   string
	}{
		{
			name:   "simple branch",
			branch: "main",
			want:   "main",
		},
		{
			name:   "feature branch with slash",
			branch: "feature/auth",
			want:   "feature-auth",
		},
		{
			name:   "nested branch path",
			branch: "feature/api/v2",
			want:   "feature-api-v2",
		},
		{
			name:   "bugfix branch",
			branch: "bugfix/user-login",
			want:   "bugfix-user-login",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BranchToWorktreeName(tt.branch); got != tt.want {
				t.Errorf("BranchToWorktreeName() = %v, want %v", got, tt.want)
			}
		})
	}
}
