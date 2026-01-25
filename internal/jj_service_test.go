package internal

import (
	"errors"
	"testing"
)

type jjMockRunner struct {
	outputs map[string]string
	errors  map[string]error
}

func (m *jjMockRunner) Run(command string) (string, error) {
	if err, exists := m.errors[command]; exists {
		return "", err
	}
	if output, exists := m.outputs[command]; exists {
		return output, nil
	}
	return "", nil
}

func TestJjService_Type(t *testing.T) {
	service := NewJjService(&jjMockRunner{})
	if service.Type() != VCSJj {
		t.Errorf("Type() = %v, want %v", service.Type(), VCSJj)
	}
}

func TestJjService_ListWorkspaces(t *testing.T) {
	tests := []struct {
		name       string
		listOutput string
		listErr    error
		wantLen    int
		wantErr    bool
	}{
		{
			name:       "single workspace",
			listOutput: "default\tabc12345\t/repo\n",
			wantLen:    1,
			wantErr:    false,
		},
		{
			name:       "multiple workspaces",
			listOutput: "default\tabc12345\t/repo\nfeature\tdef67890\t/repo/worktrees/feature\n",
			wantLen:    2,
			wantErr:    false,
		},
		{
			name:    "list command fails",
			listErr: errors.New("jj not found"),
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listCmd := "jj workspace list -T 'name ++ \"\\t\" ++ working_copy.commit_id().short(8) ++ \"\\t\" ++ working_copy_path ++ \"\\n\"'"
			mockRunner := &jjMockRunner{
				outputs: map[string]string{
					listCmd: tt.listOutput,
				},
				errors: map[string]error{},
			}
			if tt.listErr != nil {
				mockRunner.errors[listCmd] = tt.listErr
			}

			service := NewJjService(mockRunner)
			opts := ListOptions{} // Basic list, no status/labels
			workspaces, err := service.ListWorkspaces(opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListWorkspaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(workspaces) != tt.wantLen {
				t.Errorf("ListWorkspaces() returned %d workspaces, want %d", len(workspaces), tt.wantLen)
			}
		})
	}
}

func TestJjService_GetDetailedStatus(t *testing.T) {
	tests := []struct {
		name         string
		statusOutput string
		statusErr    error
		wantLines    int
		wantErr      bool
	}{
		{
			name:         "clean workspace",
			statusOutput: "",
			wantLines:    0,
			wantErr:      false,
		},
		{
			name:         "dirty workspace",
			statusOutput: "Working copy changes:\nM file.txt\nA newfile.txt\n",
			wantLines:    2, // Excludes "Working copy changes:" line
			wantErr:      false,
		},
		{
			name:      "status command fails",
			statusErr: errors.New("workspace not found"),
			wantLines: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: /repo doesn't contain special chars, so shellescape doesn't quote it
			statusCmd := "jj -R /repo status"
			mockRunner := &jjMockRunner{
				outputs: map[string]string{
					statusCmd: tt.statusOutput,
				},
				errors: map[string]error{},
			}
			if tt.statusErr != nil {
				mockRunner.errors[statusCmd] = tt.statusErr
			}

			service := NewJjService(mockRunner)
			lines, err := service.GetDetailedStatus("/repo")

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDetailedStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(lines) != tt.wantLines {
				t.Errorf("GetDetailedStatus() returned %d lines, want %d", len(lines), tt.wantLines)
			}
		})
	}
}

func TestJjService_GetBookmarksInRange(t *testing.T) {
	tests := []struct {
		name      string
		logOutput string
		logErr    error
		wantLen   int
	}{
		{
			name:      "has bookmarks",
			logOutput: "main feature-auth",
			wantLen:   2,
		},
		{
			name:      "no bookmarks",
			logOutput: "",
			wantLen:   0,
		},
		{
			name:    "command fails (empty range)",
			logErr:  errors.New("empty revision range"),
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: /repo doesn't contain special chars, so shellescape doesn't quote it
			// But trunk()..@ contains special chars, so it gets quoted
			logCmd := "jj -R /repo log -r 'trunk()..@' --no-graph -T 'bookmarks ++ \" \"'"
			mockRunner := &jjMockRunner{
				outputs: map[string]string{
					logCmd: tt.logOutput,
				},
				errors: map[string]error{},
			}
			if tt.logErr != nil {
				mockRunner.errors[logCmd] = tt.logErr
			}

			service := NewJjService(mockRunner)
			bookmarks, _ := service.GetBookmarksInRange("/repo", JjTrunkRange)

			if len(bookmarks) != tt.wantLen {
				t.Errorf("GetBookmarksInRange() returned %d bookmarks, want %d", len(bookmarks), tt.wantLen)
			}
		})
	}
}

func TestJjService_Run(t *testing.T) {
	mockRunner := &jjMockRunner{
		outputs: map[string]string{
			"test command": "test output",
		},
	}

	service := NewJjService(mockRunner)
	output, err := service.Run("test command")

	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}
	if output != "test output" {
		t.Errorf("Run() = %q, want %q", output, "test output")
	}
}
