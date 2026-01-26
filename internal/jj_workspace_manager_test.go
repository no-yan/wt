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
				{Name: "default", Status: StatusClean},
				{Name: "feature", Status: StatusClean},
			},
			statusOutput: map[string]string{
				"default": "",
				"feature": "",
			},
			wantErr: false,
		},
		{
			name:   "cannot remove default",
			wsName: "default",
			workspaces: []Workspace{
				{Name: "default", Status: StatusClean},
			},
			statusOutput: map[string]string{
				"default": "",
			},
			wantErr:   true,
			errSubstr: "cannot remove the default workspace",
		},
		{
			name:   "workspace not found",
			wsName: "nonexistent",
			workspaces: []Workspace{
				{Name: "default", Status: StatusClean},
			},
			statusOutput: map[string]string{
				"default": "",
			},
			wantErr:   true,
			errSubstr: "not found",
		},
		{
			name:   "dirty workspace",
			wsName: "dirty-ws",
			workspaces: []Workspace{
				{Name: "dirty-ws", Status: StatusDirty}, // Path will be set based on tmpDir
			},
			statusOutput: map[string]string{
				"dirty-ws": "M modified.txt\n", // Key is workspace name, path will be computed
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

			// Resolve symlinks (macOS /tmp -> /private/var/...)
			tmpDir, err = filepath.EvalSymlinks(tmpDir)
			if err != nil {
				t.Fatalf("failed to eval symlinks: %v", err)
			}

			// Create .jj directory for repo root detection
			if err := os.Mkdir(filepath.Join(tmpDir, ".jj"), 0o755); err != nil {
				t.Fatalf("failed to create .jj: %v", err)
			}

			// Change to temp directory so findJjRepoRoot can find .jj
			originalWd, _ := os.Getwd()
			defer os.Chdir(originalWd)
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to chdir: %v", err)
			}

			// Create a mock that returns the test workspaces
			listCmd := "jj workspace list -T 'self.name() ++ \"\\t\" ++ self.target().commit_id().short(8) ++ \"\\n\"'"
			listOutput := ""
			for _, ws := range tt.workspaces {
				listOutput += ws.Name + "\tabc123\n"
			}

			mockRunner := &jjMockRunner{
				outputs: map[string]string{
					listCmd: listOutput,
				},
				errors: map[string]error{},
			}

			// Add status outputs for each workspace
			// The key in statusOutput is the workspace name, we need to construct the full path
			for wsName, statusOut := range tt.statusOutput {
				// Construct path based on tmpDir (matching ParseJjWorkspaceList logic)
				var path string
				if wsName == JjDefaultWorkspace {
					path = tmpDir
				} else {
					path = filepath.Join(tmpDir, "worktrees", wsName)
				}
				// shellescape() wraps paths in single quotes only for paths with special chars
				// For most tmpDir paths, they won't have special chars, so add both versions
				statusCmd1 := "jj -R " + path + " status"
				statusCmd2 := "jj -R '" + path + "' status"
				mockRunner.outputs[statusCmd1] = statusOut
				mockRunner.outputs[statusCmd2] = statusOut
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
