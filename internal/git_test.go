package internal

import (
	"reflect"
	"testing"
)

func TestParseWorktreeList(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   []Workspace
	}{
		{
			name: "single worktree porcelain format",
			output: `worktree /repo
HEAD abc123
branch refs/heads/main`,
			want: []Workspace{
				{
					Name:     "main",
					Path:     "/repo",
					Revision: "abc123",
					Labels:   []string{"main"},
					Status:   StatusClean,
				},
			},
		},
		{
			name: "multiple worktrees porcelain format",
			output: `worktree /repo
HEAD abc123
branch refs/heads/main

worktree /repo/worktrees/test-feature
HEAD def456
branch refs/heads/test-feature`,
			want: []Workspace{
				{
					Name:     "main",
					Path:     "/repo",
					Revision: "abc123",
					Labels:   []string{"main"},
					Status:   StatusClean,
				},
				{
					Name:     "test-feature",
					Path:     "/repo/worktrees/test-feature",
					Revision: "def456",
					Labels:   []string{"test-feature"},
					Status:   StatusClean,
				},
			},
		},
		{
			name: "detached head porcelain format",
			output: `worktree /repo
HEAD abc123
detached`,
			want: []Workspace{
				{
					Name:     "detached HEAD",
					Path:     "/repo",
					Revision: "abc123",
					Labels:   []string{"detached HEAD"},
					Status:   StatusClean,
				},
			},
		},
		{
			name:   "empty output",
			output: "",
			want:   []Workspace{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseWorktreeList(tt.output)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseWorktreeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWorktreeStatus(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   Status
	}{
		{
			name:   "clean status",
			output: "",
			want:   StatusClean,
		},
		{
			name: "dirty status",
			output: ` M file1.go
?? file2.go`,
			want: StatusDirty,
		},
		{
			name: "staged changes",
			output: `A  file1.go
M  file2.go`,
			want: StatusDirty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseWorktreeStatus(tt.output); got != tt.want {
				t.Errorf("ParseWorktreeStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
