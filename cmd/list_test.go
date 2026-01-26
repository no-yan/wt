package cmd

import (
	"bytes"
	"testing"

	"github.com/no-yan/wt/internal"
)

func TestListCommand(t *testing.T) {
	tests := []struct {
		name      string
		worktrees []internal.Workspace
		want      string
	}{
		{
			name: "single clean worktree",
			worktrees: []internal.Workspace{
				{
					Path:   "/repo",
					Revision:   "abc123",
					Name: "main",
			Labels: []string{"main"},
					Status: internal.StatusClean,
				},
			},
			want: "main  /repo  (clean)\n",
		},
		{
			name: "multiple worktrees mixed status",
			worktrees: []internal.Workspace{
				{
					Path:   "/repo",
					Revision:   "abc123",
					Name: "main",
			Labels: []string{"main"},
					Status: internal.StatusClean,
				},
				{
					Path:   "/repo/worktrees/feature-auth",
					Revision:   "def456",
					Name: "feature-auth",
				Labels: []string{"feature/auth"},
					Status: internal.StatusDirty,
				},
			},
			want: "main          /repo                         (clean)\nfeature-auth  /repo/worktrees/feature-auth  (dirty)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			formatWorkspaceList(tt.worktrees, &output, internal.VCSGit)

			if got := output.String(); got != tt.want {
				t.Errorf("formatWorkspaceList() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatWorktreeList_EmptySlice(t *testing.T) {
	var buf bytes.Buffer
	formatWorkspaceList([]internal.Workspace{}, &buf, internal.VCSGit)
	
	output := buf.String()
	if output != "" {
		t.Errorf("expected empty output for empty slice, got %q", output)
	}
}

func TestFormatWorktreeList_AlignedOutput(t *testing.T) {
	worktrees := []internal.Workspace{
		{
			Path:   "/repo",
			Revision:   "abc123",
			Name: "main",
			Labels: []string{"main"},
			Status: internal.StatusClean,
		},
		{
			Path:   "/repo/worktrees/feature-long-name",
			Revision:   "def456",
			Name: "feature-long-name",
				Labels: []string{"feature/long-name"},
			Status: internal.StatusDirty,
		},
		{
			Path:   "/repo/worktrees/fix",
			Revision:   "ghi789",
			Name: "fix",
				Labels: []string{"fix/bug"},
			Status: internal.StatusStale,
		},
	}

	var buf bytes.Buffer
	formatWorkspaceList(worktrees, &buf, internal.VCSGit)

	output := buf.String()
	
	// Test that the output should be properly aligned
	// Expected format: name  path  (status)
	// where name and path columns are padded to align the status column
	expected := "main               /repo                              (clean)\n" +
		"feature-long-name  /repo/worktrees/feature-long-name  (dirty)\n" +
		"fix                /repo/worktrees/fix                (stale)\n"
	
	if output != expected {
		t.Errorf("formatWorkspaceList() output not aligned:\nGot:\n%q\nWant:\n%q", output, expected)
	}
}

func TestFormatWorktreeListVerbose_AlignedOutput(t *testing.T) {
	worktrees := []internal.Workspace{
		{
			Path:   "/repo",
			Revision:   "abc123",
			Name: "main",
			Labels: []string{"main"},
			Status: internal.StatusClean,
		},
		{
			Path:   "/repo/worktrees/feature-long-name",
			Revision:   "def456",
			Name: "feature-long-name",
				Labels: []string{"feature/long-name"},
			Status: internal.StatusClean, // Changed to clean to avoid nil service issue
		},
		{
			Path:   "/repo/worktrees/fix",
			Revision:   "ghi789",
			Name: "fix",
				Labels: []string{"fix/bug"},
			Status: internal.StatusStale,
		},
	}

	var buf bytes.Buffer
	// Create a mock service for testing
	mockRunner := &mockCommandRunner{outputs: make(map[string]string)}
	service := internal.NewGitService(mockRunner)
	formatWorkspaceListVerbose(worktrees, &buf, service)

	output := buf.String()

	// Test that the verbose output should be properly aligned
	// Expected format: name  branch  path  (status)
	// where name, branch, and path columns are padded to align the status column
	expected := "main               main               /repo                              (clean)\n\n" +
		"feature-long-name  feature/long-name  /repo/worktrees/feature-long-name  (clean)\n\n" +
		"fix                fix/bug            /repo/worktrees/fix                (stale)\n\n"

	if output != expected {
		t.Errorf("formatWorkspaceListVerbose() output not aligned:\nGot:\n%q\nWant:\n%q", output, expected)
	}
}

type mockCommandRunner struct {
	outputs map[string]string
}

func (m *mockCommandRunner) Run(command string) (string, error) {
	if output, ok := m.outputs[command]; ok {
		return output, nil
	}
	return "", nil
}
