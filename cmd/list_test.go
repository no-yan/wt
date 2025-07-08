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
			want: "main  /repo  (clean)\n",
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
			want: "main          /repo                         (clean)\nfeature-auth  /repo/worktrees/feature-auth  (dirty)\n",
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

func TestFormatWorktreeList_AlignedOutput(t *testing.T) {
	worktrees := []internal.Worktree{
		{
			Path:   "/repo",
			Head:   "abc123",
			Branch: "main",
			Status: internal.StatusClean,
		},
		{
			Path:   "/repo/worktrees/feature-long-name",
			Head:   "def456",
			Branch: "feature/long-name",
			Status: internal.StatusDirty,
		},
		{
			Path:   "/repo/worktrees/fix",
			Head:   "ghi789",
			Branch: "fix/bug",
			Status: internal.StatusStale,
		},
	}

	var buf bytes.Buffer
	formatWorktreeList(worktrees, &buf)

	output := buf.String()
	
	// Test that the output should be properly aligned
	// Expected format: name  path  (status)
	// where name and path columns are padded to align the status column
	expected := "main               /repo                              (clean)\n" +
		"feature-long-name  /repo/worktrees/feature-long-name  (dirty)\n" +
		"fix                /repo/worktrees/fix                (stale)\n"
	
	if output != expected {
		t.Errorf("formatWorktreeList() output not aligned:\nGot:\n%q\nWant:\n%q", output, expected)
	}
}

func TestFormatWorktreeListVerbose_AlignedOutput(t *testing.T) {
	worktrees := []internal.Worktree{
		{
			Path:   "/repo",
			Head:   "abc123",
			Branch: "main",
			Status: internal.StatusClean,
		},
		{
			Path:   "/repo/worktrees/feature-long-name",
			Head:   "def456",
			Branch: "feature/long-name",
			Status: internal.StatusClean, // Changed to clean to avoid nil service issue
		},
		{
			Path:   "/repo/worktrees/fix",
			Head:   "ghi789",
			Branch: "fix/bug",
			Status: internal.StatusStale,
		},
	}

	var buf bytes.Buffer
	// Pass nil for service as we won't trigger detailed status in this test
	formatWorktreeListVerbose(worktrees, &buf, nil)

	output := buf.String()

	// Test that the verbose output should be properly aligned
	// Expected format: name  branch  path  (status)
	// where name, branch, and path columns are padded to align the status column
	expected := "main               main               /repo                              (clean)\n\n" +
		"feature-long-name  feature/long-name  /repo/worktrees/feature-long-name  (clean)\n\n" +
		"fix                fix/bug            /repo/worktrees/fix                (stale)\n\n"

	if output != expected {
		t.Errorf("formatWorktreeListVerbose() output not aligned:\nGot:\n%q\nWant:\n%q", output, expected)
	}
}
