package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorktreeManager_AddWorktree(t *testing.T) {
	tests := []struct {
		name     string
		repoPath string
		branch   string
		wantPath string
		wantErr  bool
	}{
		{
			name:     "simple branch",
			repoPath: "/repo",
			branch:   "feature",
			wantPath: "/repo/worktrees/feature",
			wantErr:  false,
		},
		{
			name:     "branch with slash",
			repoPath: "/repo",
			branch:   "feature/auth",
			wantPath: "/repo/worktrees/feature-auth",
			wantErr:  false,
		},
		{
			name:     "nested branch",
			repoPath: "/repo",
			branch:   "feature/api/v2",
			wantPath: "/repo/worktrees/feature-api-v2",
			wantErr:  false,
		},
		{
			name:     "invalid branch name with semicolon",
			repoPath: "/repo",
			branch:   "feature;rm -rf /",
			wantPath: "",
			wantErr:  true,
		},
		{
			name:     "empty branch name",
			repoPath: "/repo",
			branch:   "",
			wantPath: "",
			wantErr:  true,
		},
		{
			name:     "relative path",
			repoPath: "repo",
			branch:   "feature",
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &MockCommandRunner{
				outputs: make(map[string]string),
			}

			// Setup successful mock responses for auto-setup commands when not expecting error
			if !tt.wantErr {
				// Mock mkdir command
				mockRunner.outputs["mkdir -p "+tt.repoPath+"/worktrees"] = ""
				// Mock gitignore commands
				mockRunner.outputs["echo 'worktrees/' >> "+tt.repoPath+"/.gitignore"] = ""
				// Mock git branch creation command (try to create new branch)
				mockRunner.outputs["git -C "+tt.repoPath+" branch "+tt.branch] = ""
				// Mock git worktree add command
				mockRunner.outputs["git -C "+tt.repoPath+" worktree add "+tt.wantPath+" "+tt.branch] = ""
			}

			service := NewGitService(mockRunner)
			manager := NewWorktreeManager(service, mockRunner)

			gotPath, err := manager.AddWorktree(tt.repoPath, tt.branch)

			if (err != nil) != tt.wantErr {
				t.Errorf("WorktreeManager.AddWorktree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotPath != tt.wantPath {
				t.Errorf("WorktreeManager.AddWorktree() path = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestWorktreeManager_AddWorktree_BranchFallback(t *testing.T) {
	tests := []struct {
		name         string
		repoPath     string
		branch       string
		branchExists bool
		wantPath     string
		wantErr      bool
	}{
		{
			name:         "create new branch successfully",
			repoPath:     "/repo",
			branch:       "new-feature",
			branchExists: false,
			wantPath:     "/repo/worktrees/new-feature",
			wantErr:      false,
		},
		{
			name:         "use existing branch when creation fails",
			repoPath:     "/repo",
			branch:       "existing-feature",
			branchExists: true,
			wantPath:     "/repo/worktrees/existing-feature",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &MockCommandRunner{
				outputs: make(map[string]string),
			}

			// Setup auto-setup commands
			mockRunner.outputs["mkdir -p "+tt.repoPath+"/worktrees"] = ""
			mockRunner.outputs["echo 'worktrees/' >> "+tt.repoPath+"/.gitignore"] = ""

			// Setup branch creation command - add to outputs only if branch doesn't exist
			branchCmd := "git -C " + tt.repoPath + " branch " + tt.branch
			if !tt.branchExists {
				mockRunner.outputs[branchCmd] = ""
			}
			// If branch exists, don't add the command to outputs so it fails

			// Setup worktree add command (always succeeds)
			mockRunner.outputs["git -C "+tt.repoPath+" worktree add "+tt.wantPath+" "+tt.branch] = ""
			service := NewGitService(mockRunner)
			manager := NewWorktreeManager(service, mockRunner)

			gotPath, err := manager.AddWorktree(tt.repoPath, tt.branch)

			if (err != nil) != tt.wantErr {
				t.Errorf("WorktreeManager.AddWorktree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotPath != tt.wantPath {
				t.Errorf("WorktreeManager.AddWorktree() path = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestWorktreeManager_AutoSetup(t *testing.T) {
	// Test only the auto-setup functionality (directory creation + gitignore) without git commands
	tempDir := t.TempDir()

	tests := []struct {
		name               string
		branch             string
		existingGitignore  string
		wantWorktreesDir   bool
		wantGitignoreEntry bool
	}{
		{
			name:               "creates worktrees dir and gitignore entry",
			branch:             "feature",
			existingGitignore:  "",
			wantWorktreesDir:   true,
			wantGitignoreEntry: true,
		},
		{
			name:               "adds gitignore entry to existing file",
			branch:             "feature/auth",
			existingGitignore:  "*.log\n.DS_Store\n",
			wantWorktreesDir:   true,
			wantGitignoreEntry: true,
		},
		{
			name:               "preserves existing gitignore entry",
			branch:             "bugfix",
			existingGitignore:  "*.log\nworktrees/\n.DS_Store\n",
			wantWorktreesDir:   true,
			wantGitignoreEntry: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test repo directory
			testRepo := filepath.Join(tempDir, tt.name)
			if err := os.MkdirAll(testRepo, 0755); err != nil {
				t.Fatal(err)
			}

			// Create existing .gitignore if specified
			gitignorePath := filepath.Join(testRepo, ".gitignore")
			if tt.existingGitignore != "" {
				if err := os.WriteFile(gitignorePath, []byte(tt.existingGitignore), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Use a real WorktreeManager but test only the auto-setup part
			mockRunner := &MockCommandRunner{outputs: make(map[string]string)}
			service := NewGitService(mockRunner)
			manager := NewWorktreeManager(service, mockRunner)

			// Test the auto-setup part by calling ensureWorktreesDirectory directly
			expectedWorktreePath := GenerateWorktreePath(testRepo, tt.branch)
			worktreesDir := filepath.Dir(expectedWorktreePath)

			err := manager.ensureWorktreesDirectory(worktreesDir)

			// Check results - worktrees directory should be created
			if tt.wantWorktreesDir {
				if _, err := os.Stat(worktreesDir); os.IsNotExist(err) {
					t.Errorf("Expected worktrees directory to be created at %s", worktreesDir)
				}
			}

			// Check gitignore entry
			if tt.wantGitignoreEntry {
				content, err := os.ReadFile(gitignorePath)
				if err != nil {
					t.Fatalf("Failed to read .gitignore: %v", err)
				}

				hasWorktreesEntry := false
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) == "worktrees/" {
						hasWorktreesEntry = true
						break
					}
				}

				if !hasWorktreesEntry {
					t.Error("Expected worktrees/ entry in .gitignore, but it was not found")
					t.Logf(".gitignore content: %q", string(content))
				}
			}

			// Verify no unexpected errors in auto-setup
			if err != nil {
				t.Errorf("ensureWorktreesDirectory() failed: %v", err)
			}
		})
	}
}

func TestValidateBranchName(t *testing.T) {
	tests := []struct {
		name    string
		branch  string
		wantErr bool
	}{
		{
			name:    "valid simple branch",
			branch:  "main",
			wantErr: false,
		},
		{
			name:    "valid feature branch",
			branch:  "feature/auth",
			wantErr: false,
		},
		{
			name:    "empty branch name",
			branch:  "",
			wantErr: true,
		},
		{
			name:    "branch with semicolon",
			branch:  "feature;dangerous",
			wantErr: true,
		},
		{
			name:    "branch starting with dash",
			branch:  "-feature",
			wantErr: true,
		},
		{
			name:    "branch with backtick",
			branch:  "feature`dangerous",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBranchName(tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBranchName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateWorktreePath(t *testing.T) {
	tests := []struct {
		name     string
		repoPath string
		branch   string
		want     string
	}{
		{
			name:     "simple branch",
			repoPath: "/repo",
			branch:   "main",
			want:     "/repo/worktrees/main",
		},
		{
			name:     "feature branch",
			repoPath: "/home/user/project",
			branch:   "feature/auth",
			want:     "/home/user/project/worktrees/feature-auth",
		},
		{
			name:     "deep nested branch",
			repoPath: "/repo",
			branch:   "feature/api/v2/auth",
			want:     "/repo/worktrees/feature-api-v2-auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateWorktreePath(tt.repoPath, tt.branch); got != tt.want {
				t.Errorf("GenerateWorktreePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorktreeManager_RemoveWorktree(t *testing.T) {
	tests := []struct {
		name      string
		repoPath  string
		target    string
		worktrees []Worktree
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "remove valid worktree",
			repoPath: "/repo",
			target:   "feature-auth",
			worktrees: []Worktree{
				{
					Path:   "/repo",
					Branch: "main",
					Status: StatusClean,
				},
				{
					Path:   "/repo/worktrees/feature-auth",
					Branch: "feature/auth",
					Status: StatusClean,
				},
			},
			wantErr: false,
		},
		{
			name:     "cannot remove main worktree",
			repoPath: "/repo",
			target:   "main",
			worktrees: []Worktree{
				{
					Path:   "/repo",
					Branch: "main",
					Status: StatusClean,
				},
			},
			wantErr: true,
			errMsg:  "cannot remove main worktree",
		},
		{
			name:     "cannot remove dirty worktree",
			repoPath: "/repo",
			target:   "feature-auth",
			worktrees: []Worktree{
				{
					Path:   "/repo/worktrees/feature-auth",
					Branch: "feature/auth",
					Status: StatusDirty,
				},
			},
			wantErr: true,
			errMsg:  "has uncommitted changes",
		},
		{
			name:      "worktree not found",
			repoPath:  "/repo",
			target:    "nonexistent",
			worktrees: []Worktree{},
			wantErr:   true,
			errMsg:    "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &MockCommandRunner{
				outputs: map[string]string{
					"git worktree list --porcelain": generateMockWorktreeOutput(tt.worktrees),
				},
			}

			// Add status outputs for each worktree
			for _, wt := range tt.worktrees {
				statusKey := "git -C " + wt.Path + " status --porcelain"
				if wt.Status == StatusDirty {
					mockRunner.outputs[statusKey] = " M file.go\n"
				} else {
					mockRunner.outputs[statusKey] = ""
				}
			}

			// Add remove command mock for successful cases
			if !tt.wantErr {
				for _, wt := range tt.worktrees {
					if wt.Name() == tt.target {
						mockRunner.outputs["git -C "+tt.repoPath+" worktree remove "+wt.Path] = ""
						break
					}
				}
			}

			service := NewGitService(mockRunner)
			manager := NewWorktreeManager(service, mockRunner)

			err := manager.RemoveWorktree(tt.repoPath, tt.target)

			if (err != nil) != tt.wantErr {
				t.Errorf("WorktreeManager.RemoveWorktree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("WorktreeManager.RemoveWorktree() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func generateMockWorktreeOutput(worktrees []Worktree) string {
	if len(worktrees) == 0 {
		return ""
	}

	var parts []string
	for _, wt := range worktrees {
		part := fmt.Sprintf("worktree %s\nHEAD abc123\nbranch refs/heads/%s", wt.Path, wt.Branch)
		parts = append(parts, part)
	}
	return strings.Join(parts, "\n\n")
}

func TestWorktreeManager_EnsureGitignoreEntry(t *testing.T) {
	// Create temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name             string
		gitignoreContent string
		expectAppend     bool
		wantErr          bool
	}{
		{
			name:             "add entry to existing gitignore without worktrees entry",
			gitignoreContent: "*.log\n.DS_Store\n",
			expectAppend:     true,
			wantErr:          false,
		},
		{
			name:             "entry already exists",
			gitignoreContent: "*.log\nworktrees/\n.DS_Store\n",
			expectAppend:     false,
			wantErr:          false,
		},
		{
			name:             "gitignore file does not exist",
			gitignoreContent: "",
			expectAppend:     false, // Will create new file
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			testDir := filepath.Join(tempDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatal(err)
			}

			gitignorePath := filepath.Join(testDir, ".gitignore")

			// Create .gitignore file if content is provided
			if tt.gitignoreContent != "" {
				if err := os.WriteFile(gitignorePath, []byte(tt.gitignoreContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			mockRunner := &MockCommandRunner{
				outputs: make(map[string]string),
			}

			service := NewGitService(mockRunner)
			manager := NewWorktreeManager(service, mockRunner)

			err := manager.ensureGitignoreEntry(testDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("ensureGitignoreEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify .gitignore content
			content, err := os.ReadFile(gitignorePath)
			if err != nil {
				t.Fatalf("Failed to read .gitignore: %v", err)
			}

			lines := strings.Split(string(content), "\n")
			hasWorktreesEntry := false
			for _, line := range lines {
				if strings.TrimSpace(line) == "worktrees/" {
					hasWorktreesEntry = true
					break
				}
			}

			if !hasWorktreesEntry {
				t.Error("Expected worktrees/ entry in .gitignore, but it was not found")
			}
		})
	}
}
