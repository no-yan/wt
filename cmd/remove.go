package cmd

import (
	"fmt"
	"os"

	"github.com/no-yan/wrkt/internal"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a worktree",
	Long: `Remove a worktree by name.

Safety checks:
- Cannot remove the main worktree
- Cannot remove worktrees with uncommitted changes
- Must be in worktrees/ subdirectory`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		runner := internal.NewExecCommandRunner()
		gitService := internal.NewGitService(runner)
		manager := internal.NewWorktreeManager(gitService, runner)

		repoPath, err := getRepoRoot(runner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding git repository: %v\n", err)
			os.Exit(1)
		}

		if err := manager.RemoveWorktree(repoPath, name); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing worktree: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Removed worktree: %s\n", name)
	},
}
