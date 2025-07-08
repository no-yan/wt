package cmd

import (
	"fmt"
	"os"

	"github.com/no-yan/wt/internal"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <name>...",
	Aliases: []string{"rm"},
	Short:   "Remove one or more worktrees",
	Long: `Remove one or more worktrees by name.

Safety checks:
- Cannot remove the main worktree
- Cannot remove worktrees with uncommitted changes
- Must be in worktrees/ subdirectory
- When removing multiple worktrees, validates all before removing any (fail-fast)`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runner := internal.NewExecCommandRunner()
		gitService := internal.NewGitService(runner)
		manager := internal.NewWorktreeManager(gitService, runner)

		repoPath, err := getRepoRoot(runner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding git repository: %v\n", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			// Single worktree removal
			if err := manager.RemoveWorktree(repoPath, args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing worktree: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Removed worktree: %s\n", args[0])
		} else {
			// Multiple worktree removal with fail-fast validation
			if err := manager.RemoveMultipleWorktrees(repoPath, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing worktrees: %v\n", err)
				os.Exit(1)
			}
			for _, name := range args {
				fmt.Printf("Removed worktree: %s\n", name)
			}
		}
	},
}
