package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/no-yan/wt/internal"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "Add a new worktree",
	Long:  "Add a new git worktree in the worktrees/ subdirectory.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]

		runner := internal.NewExecCommandRunner()
		gitService := internal.NewGitService(runner)
		manager := internal.NewWorktreeManager(gitService, runner)

		repoPath, err := getRepoRoot(runner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding git repository: %v\n", err)
			os.Exit(1)
		}

		worktreePath, err := manager.AddWorktree(repoPath, branch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding worktree: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added worktree: %s\n", worktreePath)
	},
}

func getRepoRoot(runner internal.CommandRunner) (string, error) {
	output, err := runner.Run("git rev-parse --show-toplevel")
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}

	repoPath := strings.TrimSpace(output)
	if repoPath == "" {
		return "", fmt.Errorf("git rev-parse returned empty path")
	}

	return repoPath, nil
}
