package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/no-yan/wrkt/internal"
)

func TestCleanIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directory for test repo
	tempDir, err := ioutil.TempDir("", "wrkt-clean-test-")
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

	// Create test branch
	if err := runGitCommand(tempDir, "git", "checkout", "-b", "feature/stale-test"); err != nil {
		t.Fatalf("Failed to create feature/stale-test branch: %v", err)
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

	t.Run("Create worktree and make it stale", func(t *testing.T) {
		// Add a worktree
		worktreePath, err := manager.AddWorktree(tempDir, "feature/stale-test")
		if err != nil {
			t.Fatalf("Failed to add worktree: %v", err)
		}

		// Verify worktree exists and is clean
		worktrees, err := gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees: %v", err)
		}

		var staleWorktree *internal.Worktree
		for _, wt := range worktrees {
			if wt.Name() == "feature-stale-test" {
				staleWorktree = &wt
				break
			}
		}

		if staleWorktree == nil {
			t.Fatal("feature-stale-test worktree not found")
		}

		if staleWorktree.Status == internal.StatusStale {
			t.Errorf("Worktree should not be stale initially, got %v", staleWorktree.Status)
		}

		// Manually remove the worktree directory to make it stale
		// This simulates what happens when someone deletes a worktree directory
		// without using 'git worktree remove'
		if err := os.RemoveAll(worktreePath); err != nil {
			t.Fatalf("Failed to remove worktree directory: %v", err)
		}
	})

	t.Run("Detect stale worktree", func(t *testing.T) {
		// List worktrees - the removed directory should now be detected as stale
		worktrees, err := gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees: %v", err)
		}

		var staleWorktree *internal.Worktree
		for _, wt := range worktrees {
			if wt.Name() == "feature-stale-test" {
				staleWorktree = &wt
				break
			}
		}

		if staleWorktree == nil {
			t.Fatal("feature-stale-test worktree not found in git worktree list")
		}

		if staleWorktree.Status != internal.StatusStale {
			t.Errorf("Expected worktree to be stale, got %v", staleWorktree.Status)
		}
	})

	t.Run("Clean stale worktree", func(t *testing.T) {
		// Get the stale worktree path
		worktrees, err := gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees: %v", err)
		}

		var staleWorktree *internal.Worktree
		for _, wt := range worktrees {
			if wt.Name() == "feature-stale-test" && wt.Status == internal.StatusStale {
				staleWorktree = &wt
				break
			}
		}

		if staleWorktree == nil {
			t.Fatal("No stale worktree found to clean")
		}

		// Clean the stale worktree
		if err := pruneStaleWorktrees(runner); err != nil {
			t.Fatalf("Failed to clean stale worktree: %v", err)
		}

		// Verify the worktree is no longer listed
		worktrees, err = gitService.ListWorktrees()
		if err != nil {
			t.Fatalf("Failed to list worktrees after cleaning: %v", err)
		}

		for _, wt := range worktrees {
			if wt.Name() == "feature-stale-test" {
				t.Error("Stale worktree should have been removed from git worktree list")
			}
		}
	})
}
