package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// shouldSkipIntegrationTest returns true if the integration test should be skipped
// due to environmental conditions that make it unreliable or problematic.
func shouldSkipIntegrationTest() bool {
	// Skip in CI environments
	if os.Getenv("CI") != "" {
		return true
	}

	// Skip in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "" {
		return true
	}

	// Skip if running in a worktree that might conflict with git operations
	if gitCommonDir, err := exec.Command("git", "rev-parse", "--git-common-dir").Output(); err == nil {
		gitTopLevel, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
		if err == nil {
			commonDir := strings.TrimSpace(string(gitCommonDir))
			topLevel := strings.TrimSpace(string(gitTopLevel))

			// Convert relative git-common-dir to absolute path
			if !strings.HasPrefix(commonDir, "/") {
				if cwd, err := os.Getwd(); err == nil {
					commonDir = filepath.Join(cwd, commonDir)
				}
			}

			// If the common dir is not in the current directory, we're in a worktree
			// and creating temp git repos might conflict with the git operations
			if !strings.HasPrefix(commonDir, topLevel) {
				return true
			}
		}
	}

	// Skip if git is not available or not configured
	if err := exec.Command("git", "version").Run(); err != nil {
		return true
	}

	// Check if git user is configured (needed for commits)
	if err := exec.Command("git", "config", "user.name").Run(); err != nil {
		// We'll configure it in the test, but check if global config exists
		if err := exec.Command("git", "config", "--global", "user.name").Run(); err != nil {
			// No global config, but that's ok - we set it in the test
		}
	}

	return false
}

// runGitCommand runs a git command in the specified directory
func runGitCommand(dir string, command ...string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	// Don't output to stdout/stderr to avoid cluttering test output
	return cmd.Run()
}
