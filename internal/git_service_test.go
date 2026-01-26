package internal

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGitService_ListWorkspaces(t *testing.T) {
	tests := []struct {
		name          string
		gitOutput     string
		statusOutputs map[string]string
		want          []Workspace
		wantErr       bool
	}{
		{
			name: "single worktree",
			gitOutput: `worktree /repo
HEAD abc123
branch refs/heads/main`,
			statusOutputs: map[string]string{
				"/repo": "",
			},
			want: []Workspace{
				{
					Name:     "main",
					Path:     "/repo",
					Revision: "abc123",
					Labels:   []string{"main"},
					Status:   StatusClean,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple worktrees with mixed status",
			gitOutput: `worktree /repo
HEAD abc123
branch refs/heads/main

worktree /repo/worktrees/feature-auth
HEAD def456
branch refs/heads/feature/auth`,
			statusOutputs: map[string]string{
				"/repo":                        "",
				"/repo/worktrees/feature-auth": " M file.go\n?? new.go",
			},
			want: []Workspace{
				{
					Name:     "main",
					Path:     "/repo",
					Revision: "abc123",
					Labels:   []string{"main"},
					Status:   StatusClean,
				},
				{
					Name:     "feature-auth",
					Path:     "/repo/worktrees/feature-auth",
					Revision: "def456",
					Labels:   []string{"feature/auth"},
					Status:   StatusDirty,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &MockCommandRunner{
				outputs: map[string]string{
					"git worktree list --porcelain": tt.gitOutput,
				},
			}

			for path, status := range tt.statusOutputs {
				mockRunner.outputs["git -C "+path+" status --porcelain"] = status
			}

			service := NewGitService(mockRunner)
			opts := ListOptions{IncludeStatus: true}
			got, err := service.ListWorkspaces(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("GitService.ListWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitService.ListWorkspaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

type MockCommandRunner struct {
	outputs  map[string]string
	commands []string
}

func (m *MockCommandRunner) Run(command string) (string, error) {
	m.commands = append(m.commands, command)
	if output, exists := m.outputs[command]; exists {
		return output, nil
	}
	// Simulate command failure for commands not in outputs
	return "", fmt.Errorf("command failed: %s", command)
}

func (m *MockCommandRunner) GetCommands() []string {
	return m.commands
}
