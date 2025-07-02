package internal

import (
	"reflect"
	"testing"
)

func TestParseWorktreeList(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   []Worktree
	}{
		{
			name: "single worktree porcelain format",
			output: `worktree /repo
HEAD abc123
branch refs/heads/main`,
			want: []Worktree{
				{
					Path:   "/repo",
					Head:   "abc123",
					Branch: "main",
					Status: StatusClean,
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
			want: []Worktree{
				{
					Path:   "/repo",
					Head:   "abc123",
					Branch: "main",
					Status: StatusClean,
				},
				{
					Path:   "/repo/worktrees/test-feature",
					Head:   "def456",
					Branch: "test-feature",
					Status: StatusClean,
				},
			},
		},
		{
			name: "detached head porcelain format",
			output: `worktree /repo
HEAD abc123
detached`,
			want: []Worktree{
				{
					Path:   "/repo",
					Head:   "abc123",
					Branch: "detached HEAD",
					Status: StatusClean,
				},
			},
		},
		{
			name:   "empty output",
			output: "",
			want:   []Worktree{},
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
