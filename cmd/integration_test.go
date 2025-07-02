package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/no-yan/wrkt/internal"
)

func TestIntegrationWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip in problematic environments
	if shouldSkipIntegrationTest() {
		t.Skip("Skipping integration test in problematic environment")
	}

	// Create temporary directory for test repo
	tempDir, err := ioutil.TempDir("", "wrkt-integration-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize git repo
	if err := runGitCommand(tempDir, "git", "init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git for testing
	if err := runGitCommand(tempDir, "git", "config", "user.name", "Test User"); err != nil {
		t.Fatalf("Failed to configure git user: %v", err)
	}
	if err := runGitCommand(tempDir, "git", "config", "user.email", "test@example.com"); err != nil {
		t.Fatalf("Failed to configure git email: %v", err)
	}
	if err := runGitCommand(tempDir, "git", "config", "init.defaultBranch", "main"); err != nil {
		t.Fatalf("Failed to configure default branch: %v", err)
	}

	// Create initial commit
	readmeFile := filepath.Join(tempDir, "README.md")
	if err := ioutil.WriteFile(readmeFile, []byte("# Test Repo\n"), 0644); err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}
	if err := runGitCommand(tempDir, "git", "add", "README.md"); err != nil {
		t.Fatalf("Failed to add README: %v", err)
	}
	if err := runGitCommand(tempDir, "git", "commit", "-m", "Initial commit"); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	// Create test branches
	if err := runGitCommand(tempDir, "git", "checkout", "-b", "feature/test1"); err != nil {
		t.Fatalf("Failed to create feature/test1 branch: %v", err)
	}
	if err := runGitCommand(tempDir, "git", "checkout", "main"); err != nil {
		t.Fatalf("Failed to checkout main: %v", err)
	}
	if err := runGitCommand(tempDir, "git", "checkout", "-b", "feature/test2"); err != nil {
		t.Fatalf("Failed to create feature/test2 branch: %v", err)
	}
	if err := runGitCommand(tempDir, "git", "checkout", "main"); err != nil {
		t.Fatalf("Failed to checkout main: %v", err)
	}

	// Change to test repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Test WorktreeManager operations
	runner := internal.NewExecCommandRunner()
	gitService := internal.NewGitService(runner)
	manager := internal.NewWorktreeManager(gitService, runner)

	t.Run("Add worktree", func(t *testing.T) {
		worktreePath, err := manager.AddWorktree(tempDir, "feature/test1")
		if err != nil {
			t.Fatalf("Failed to add worktree: %v", err)
		}

		expectedPath := filepath.Join(tempDir, "worktrees", "feature-test1")
		if worktreePath != expectedPath {
			t.Errorf("Expected worktree path %s, got %s", expectedPath, worktreePath)
		}

		// Verify worktree directory exists
		if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
			t.Errorf("Worktree directory not created: %s", worktreePath)
		}

		// Verify .gitignore entry was added
		gitignorePath := filepath.Join(tempDir, ".gitignore")
		if content, err := ioutil.ReadFile(gitignorePath); err == nil {
			if !strings.Contains(string(content), "worktrees/") {
				t.Error(".gitignore does not contain worktrees/ entry")
			}
		}
	})

	t.Run("List worktrees", func(t *testing.T) {
		worktrees, err := gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees: %v", err)
		}

		if len(worktrees) < 2 { // at least main + feature/test1
			t.Errorf("Expected at least 2 worktrees, got %d", len(worktrees))
		}

		// Check that we have main and feature-test1
		names := make(map[string]bool)
		for _, wt := range worktrees {
			names[wt.Name()] = true
		}

		// The main worktree might have different names depending on repo state
		// Just check that feature-test1 exists
		if !names["feature-test1"] {
			t.Error("feature-test1 worktree not found")
		}
	})

	t.Run("Add second worktree", func(t *testing.T) {
		worktreePath, err := manager.AddWorktree(tempDir, "feature/test2")
		if err != nil {
			t.Fatalf("Failed to add second worktree: %v", err)
		}

		expectedPath := filepath.Join(tempDir, "worktrees", "feature-test2")
		if worktreePath != expectedPath {
			t.Errorf("Expected worktree path %s, got %s", expectedPath, worktreePath)
		}
	})

	t.Run("List all worktrees after adding second", func(t *testing.T) {
		worktrees, err := gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees: %v", err)
		}

		if len(worktrees) < 3 { // at least main + feature/test1 + feature/test2
			t.Errorf("Expected at least 3 worktrees, got %d", len(worktrees))
		}

		// Check that both test worktrees exist
		names := make(map[string]bool)
		for _, wt := range worktrees {
			names[wt.Name()] = true
		}

		if !names["feature-test1"] {
			t.Error("feature-test1 worktree not found")
		}
		if !names["feature-test2"] {
			t.Error("feature-test2 worktree not found")
		}
	})

	t.Run("Remove worktree", func(t *testing.T) {
		err := manager.RemoveWorktree(tempDir, "feature-test1")
		if err != nil {
			t.Fatalf("Failed to remove worktree: %v", err)
		}

		// Verify worktree was removed
		worktrees, err := gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees after removal: %v", err)
		}

		// Check that feature-test1 is gone but feature-test2 still exists
		names := make(map[string]bool)
		for _, wt := range worktrees {
			names[wt.Name()] = true
		}

		if names["feature-test1"] {
			t.Error("feature-test1 worktree should have been removed")
		}
		if !names["feature-test2"] {
			t.Error("feature-test2 worktree should still exist")
		}
	})

	t.Run("Create dirty worktree", func(t *testing.T) {
		// Add a file to the feature-test2 worktree to make it dirty
		test2Path := filepath.Join(tempDir, "worktrees", "feature-test2")
		testFile := filepath.Join(test2Path, "test.txt")
		if err := ioutil.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// List worktrees and check status
		worktrees, err := gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees: %v", err)
		}

		var test2Worktree *internal.Worktree
		for _, wt := range worktrees {
			if wt.Name() == "feature-test2" {
				test2Worktree = &wt
				break
			}
		}

		if test2Worktree == nil {
			t.Fatal("feature-test2 worktree not found")
		}

		if test2Worktree.Status != internal.StatusDirty {
			t.Errorf("Expected feature-test2 to be dirty, got status %v", test2Worktree.Status)
		}
	})

	t.Run("Get detailed status", func(t *testing.T) {
		test2Path := filepath.Join(tempDir, "worktrees", "feature-test2")
		statusLines, err := gitService.GetDetailedStatus(test2Path)
		if err != nil {
			t.Fatalf("Failed to get detailed status: %v", err)
		}

		if len(statusLines) == 0 {
			t.Error("Expected status lines for dirty worktree")
		}

		// Should contain information about the test.txt file
		found := false
		for _, line := range statusLines {
			if strings.Contains(line, "test.txt") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Status lines should mention test.txt file")
		}
	})
}
