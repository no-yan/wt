package internal

import (
	"reflect"
	"testing"
)

func TestGitService_ListWorktrees(t *testing.T) {
	tests := []struct {
		name          string
		gitOutput     string
		statusOutputs map[string]string
		want          []Worktree
		wantErr       bool
	}{
		{
			name:      "single worktree",
			gitOutput: `/repo /abc123 [main]`,
			statusOutputs: map[string]string{
				"/repo": "",
			},
			want: []Worktree{
				{
					Path:   "/repo",
					Head:   "abc123",
					Branch: "main",
					Status: StatusClean,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple worktrees with mixed status",
			gitOutput: `/repo /abc123 [main]
/repo/worktrees/feature-auth /def456 [feature/auth]`,
			statusOutputs: map[string]string{
				"/repo":                        "",
				"/repo/worktrees/feature-auth": " M file.go\n?? new.go",
			},
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
					Status: StatusDirty,
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
			got, err := service.ListWorktrees()

			if (err != nil) != tt.wantErr {
				t.Errorf("GitService.ListWorktrees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitService.ListWorktrees() = %v, want %v", got, tt.want)
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
	return "", nil
}
