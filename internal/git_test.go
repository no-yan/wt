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
			name:   "single worktree",
			output: `/repo /abc123 [main]`,
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
			name: "multiple worktrees",
			output: `/repo /abc123 [main]
/repo/worktrees/feature-auth /def456 [feature/auth]`,
			want: []Worktree{
				{
					Path:   "/repo",
					Head:   "abc123",
					Branch: "main",
					Status: StatusClean,
				},
				{
					Path:   "/repo/worktrees/feature-auth",
					Head:   "def456",
					Branch: "feature/auth",
					Status: StatusClean,
				},
			},
		},
		{
			name:   "detached head",
			output: `/repo /abc123 (detached HEAD)`,
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
