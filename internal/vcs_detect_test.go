package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectVCS(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string) error
		wantType VCSType
		wantErr  bool
	}{
		{
			name: "git repo",
			setup: func(dir string) error {
				return os.Mkdir(filepath.Join(dir, ".git"), 0o755)
			},
			wantType: VCSGit,
			wantErr:  false,
		},
		{
			name: "jj repo",
			setup: func(dir string) error {
				return os.Mkdir(filepath.Join(dir, ".jj"), 0o755)
			},
			wantType: VCSJj,
			wantErr:  false,
		},
		{
			name: "colocated repo (jj + git) - jj takes priority",
			setup: func(dir string) error {
				if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
					return err
				}
				return os.Mkdir(filepath.Join(dir, ".jj"), 0o755)
			},
			wantType: VCSJj,
			wantErr:  false,
		},
		{
			name: "no VCS",
			setup: func(dir string) error {
				return nil
			},
			wantType: "",
			wantErr:  true,
		},
		{
			name: "git worktree (file .git)",
			setup: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: /path/to/main/.git/worktrees/feature"), 0o644)
			},
			wantType: VCSGit,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir, err := os.MkdirTemp("", "vcs-detect-test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Setup the test scenario
			if err := tt.setup(tmpDir); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			// Run detection
			gotType, err := DetectVCS(tmpDir)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectVCS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check VCS type
			if gotType != tt.wantType {
				t.Errorf("DetectVCS() = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}

func TestDetectVCS_WalksUpDirectory(t *testing.T) {
	// Create temp directory structure
	tmpDir, err := os.MkdirTemp("", "vcs-detect-walk-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested structure: tmpDir/.git and tmpDir/subdir/nested
	if err := os.Mkdir(filepath.Join(tmpDir, ".git"), 0o755); err != nil {
		t.Fatalf("failed to create .git: %v", err)
	}
	nestedDir := filepath.Join(tmpDir, "subdir", "nested")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}

	// Detection from nested dir should find git repo
	gotType, err := DetectVCS(nestedDir)
	if err != nil {
		t.Errorf("DetectVCS() unexpected error: %v", err)
	}
	if gotType != VCSGit {
		t.Errorf("DetectVCS() = %v, want %v", gotType, VCSGit)
	}
}

func TestGetRepoRoot(t *testing.T) {
	t.Run("git repo root", func(t *testing.T) {
		mockRunner := &mockCommandRunner{
			outputs: map[string]string{
				"git rev-parse --show-toplevel": "/home/user/project\n",
			},
		}

		got, err := GetRepoRoot(mockRunner, VCSGit)
		if err != nil {
			t.Errorf("GetRepoRoot() unexpected error: %v", err)
		}
		if got != "/home/user/project" {
			t.Errorf("GetRepoRoot() = %q, want %q", got, "/home/user/project")
		}
	})

	t.Run("git path with trailing whitespace", func(t *testing.T) {
		mockRunner := &mockCommandRunner{
			outputs: map[string]string{
				"git rev-parse --show-toplevel": "/path/with/spaces  \n",
			},
		}

		got, err := GetRepoRoot(mockRunner, VCSGit)
		if err != nil {
			t.Errorf("GetRepoRoot() unexpected error: %v", err)
		}
		if got != "/path/with/spaces" {
			t.Errorf("GetRepoRoot() = %q, want %q", got, "/path/with/spaces")
		}
	})
}

func TestGetRepoRoot_Jj(t *testing.T) {
	// jj uses file system traversal to find .jj directory
	tmpDir, err := os.MkdirTemp("", "jj-repo-root-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Resolve symlinks (macOS /tmp -> /private/var/...)
	tmpDir, err = filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("failed to eval symlinks: %v", err)
	}

	// Create .jj directory
	if err := os.Mkdir(filepath.Join(tmpDir, ".jj"), 0o755); err != nil {
		t.Fatalf("failed to create .jj: %v", err)
	}

	// Create nested workspace directory
	workspaceDir := filepath.Join(tmpDir, "worktrees", "feature")
	if err := os.MkdirAll(workspaceDir, 0o755); err != nil {
		t.Fatalf("failed to create workspace dir: %v", err)
	}

	// Change to workspace directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	if err := os.Chdir(workspaceDir); err != nil {
		t.Fatalf("failed to change to workspace dir: %v", err)
	}

	// GetRepoRoot should find the repo root (where .jj is)
	got, err := GetRepoRoot(nil, VCSJj)
	if err != nil {
		t.Errorf("GetRepoRoot() unexpected error: %v", err)
	}
	if got != tmpDir {
		t.Errorf("GetRepoRoot() = %q, want %q", got, tmpDir)
	}
}

func TestCheckVCSAvailable(t *testing.T) {
	tests := []struct {
		name    string
		vcsType VCSType
		cmdErr  bool
		wantErr bool
	}{
		{
			name:    "git available",
			vcsType: VCSGit,
			cmdErr:  false,
			wantErr: false,
		},
		{
			name:    "git not available",
			vcsType: VCSGit,
			cmdErr:  true,
			wantErr: true,
		},
		{
			name:    "jj available",
			vcsType: VCSJj,
			cmdErr:  false,
			wantErr: false,
		},
		{
			name:    "jj not available",
			vcsType: VCSJj,
			cmdErr:  true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &mockCommandRunner{
				shouldError: tt.cmdErr,
			}

			err := CheckVCSAvailable(mockRunner, tt.vcsType)

			if (err != nil) != tt.wantErr {
				t.Errorf("CheckVCSAvailable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockCommandRunner struct {
	outputs     map[string]string
	shouldError bool
}

func (m *mockCommandRunner) Run(command string) (string, error) {
	if m.shouldError {
		return "", os.ErrNotExist
	}
	if output, exists := m.outputs[command]; exists {
		return output, nil
	}
	return "version 1.0", nil // Default for version checks
}
