package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJjWorkspaceManager_AddWorkspace(t *testing.T) {
	tests := []struct {
		name      string
		wsName    string
		opts      AddOptions
		wantErr   bool
		errSubstr string
	}{
		{
			name:    "valid workspace name",
			wsName:  "feature-auth",
			opts:    AddOptions{},
			wantErr: false,
		},
		{
			name:    "with revision option",
			wsName:  "my-feature",
			opts:    AddOptions{Revision: "main"},
			wantErr: false,
		},
		{
			name:      "empty name",
			wsName:    "",
			opts:      AddOptions{},
			wantErr:   true,
			errSubstr: "cannot be empty",
		},
		{
			name:      "invalid characters",
			wsName:    "test;rm -rf",
			opts:      AddOptions{},
			wantErr:   true,
			errSubstr: "invalid characters",
		},
		{
			name:      "starts with dash",
			wsName:    "-invalid",
			opts:      AddOptions{},
			wantErr:   true,
			errSubstr: "cannot start with dash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir, err := os.MkdirTemp("", "jj-wm-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			mockRunner := &jjMockRunner{
				outputs: map[string]string{},
				errors:  map[string]error{},
			}
			service := NewJjService(mockRunner)
			manager := NewJjWorkspaceManager(service, mockRunner)

			path, err := manager.AddWorkspace(tmpDir, tt.wsName, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errSubstr != "" {
				if err == nil || !contains(err.Error(), tt.errSubstr) {
					t.Errorf("AddWorkspace() error = %v, want error containing %q", err, tt.errSubstr)
				}
				return
			}

			if !tt.wantErr {
				expectedPath := filepath.Join(tmpDir, "worktrees", tt.wsName)
				if path != expectedPath {
					t.Errorf("AddWorkspace() path = %q, want %q", path, expectedPath)
				}
			}
		})
	}
}

func TestJjWorkspaceManager_RemoveWorkspace(t *testing.T) {
	tests := []struct {
		name         string
		wsName       string
		workspaces   []Workspace
		statusOutput map[string]string // path -> status output
		wantErr      bool
		errSubstr    string
	}{
		{
			name:   "remove existing workspace",
			wsName: "feature",
			workspaces: []Workspace{
				{Name: "default", Path: "/repo", Status: StatusClean},
				{Name: "feature", Path: "/repo/worktrees/feature", Status: StatusClean},
			},
			statusOutput: map[string]string{
				"/repo":                     "",
				"/repo/worktrees/feature": "",
			},
			wantErr: false,
		},
		{
			name:   "cannot remove default",
			wsName: "default",
			workspaces: []Workspace{
				{Name: "default", Path: "/repo", Status: StatusClean},
			},
			statusOutput: map[string]string{
				"/repo": "",
			},
			wantErr:   true,
			errSubstr: "cannot remove the default workspace",
		},
		{
			name:   "workspace not found",
			wsName: "nonexistent",
			workspaces: []Workspace{
				{Name: "default", Path: "/repo", Status: StatusClean},
			},
			statusOutput: map[string]string{
				"/repo": "",
			},
			wantErr:   true,
			errSubstr: "not found",
		},
		{
			name:   "dirty workspace",
			wsName: "dirty-ws",
			workspaces: []Workspace{
				{Name: "dirty-ws", Path: "/repo/worktrees/dirty-ws", Status: StatusDirty},
			},
			statusOutput: map[string]string{
				"/repo/worktrees/dirty-ws": "M modified.txt\n", // Non-empty means dirty
			},
			wantErr:   true,
			errSubstr: "uncommitted changes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory and workspace directory for test
			tmpDir, err := os.MkdirTemp("", "jj-wm-rm-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create a mock that returns the test workspaces
			listCmd := "jj workspace list -T 'name ++ \"\\t\" ++ working_copy.commit_id().short(8) ++ \"\\t\" ++ working_copy_path ++ \"\\n\"'"
			listOutput := ""
			for _, ws := range tt.workspaces {
				listOutput += ws.Name + "\tabc123\t" + ws.Path + "\n"
			}

			mockRunner := &jjMockRunner{
				outputs: map[string]string{
					listCmd: listOutput,
				},
				errors: map[string]error{},
			}

			// Add status outputs for each workspace path
			// Note: shellescape() wraps paths in single quotes
			for path, statusOut := range tt.statusOutput {
				statusCmd := "jj -R " + path + " status"
				mockRunner.outputs[statusCmd] = statusOut
				// Also add quoted version for paths without special chars
				statusCmdQuoted := "jj -R '" + path + "' status"
				mockRunner.outputs[statusCmdQuoted] = statusOut
			}

			service := NewJjService(mockRunner)
			manager := NewJjWorkspaceManager(service, mockRunner)

			err = manager.RemoveWorkspace(tmpDir, tt.wsName)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errSubstr != "" {
				if err == nil || !contains(err.Error(), tt.errSubstr) {
					t.Errorf("RemoveWorkspace() error = %v, want error containing %q", err, tt.errSubstr)
				}
			}
		})
	}
}

func TestValidateWorkspaceName(t *testing.T) {
	tests := []struct {
		name    string
		wsName  string
		wantErr bool
	}{
		{"valid simple name", "feature", false},
		{"valid with dash", "feature-auth", false},
		{"valid with underscore", "feature_auth", false},
		{"empty name", "", true},
		{"semicolon", "test;cmd", true},
		{"pipe", "test|cmd", true},
		{"backtick", "test`cmd", true},
		{"dollar", "test$var", true},
		{"starts with dash", "-invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.wsName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName(%q) error = %v, wantErr %v", tt.wsName, err, tt.wantErr)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
