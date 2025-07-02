package cmd

import (
	"fmt"
	"os"

	"github.com/no-yan/wrkt/internal"
	"github.com/spf13/cobra"
)

var (
	cleanDryRun bool
	cleanForce  bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up stale worktrees",
	Long: `Clean up stale worktrees that are no longer valid.

This command removes worktree entries that point to non-existent directories
or directories that are no longer valid git worktrees.

Use --dry-run to see what would be cleaned without actually removing anything.
Use --force to skip confirmation prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		runner := internal.NewExecCommandRunner()
		service := internal.NewGitService(runner)

		worktrees, err := service.ListWorktrees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing worktrees: %v\n", err)
			os.Exit(1)
		}

		// Find stale worktrees
		var staleWorktrees []internal.Worktree
		for _, wt := range worktrees {
			if wt.Status == internal.StatusStale {
				staleWorktrees = append(staleWorktrees, wt)
			}
		}

		if len(staleWorktrees) == 0 {
			fmt.Println("No stale worktrees found.")
			return
		}

		// Show what will be cleaned
		fmt.Printf("Found %d stale worktree(s):\n", len(staleWorktrees))
		for _, wt := range staleWorktrees {
			fmt.Printf("  %s -> %s\n", wt.Name(), wt.Path)
		}

		if cleanDryRun {
			fmt.Println("\nDry run mode - no changes made.")
			return
		}

		// Confirm before proceeding
		if !cleanForce {
			fmt.Print("\nProceed with cleanup? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cleanup cancelled.")
				return
			}
		}

		// Clean up stale worktrees using git worktree prune
		if err := pruneStaleWorktrees(runner); err != nil {
			fmt.Fprintf(os.Stderr, "Error cleaning stale worktrees: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cleaned %d stale worktree(s).\n", len(staleWorktrees))
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Show what would be cleaned without making changes")
	cleanCmd.Flags().BoolVar(&cleanForce, "force", false, "Skip confirmation prompts")
}

func pruneStaleWorktrees(runner internal.CommandRunner) error {
	// Use git worktree prune to remove stale worktree entries
	_, err := runner.Run("git worktree prune")
	return err
}

func shellescape(s string) string {
	// Simple shell escaping for paths
	return "'" + s + "'"
}
