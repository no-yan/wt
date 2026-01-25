package cmd

import (
	"fmt"
	"os"

	"github.com/no-yan/wt/internal"
	"github.com/spf13/cobra"
)

var (
	cleanDryRun bool
	cleanForce  bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up stale workspaces",
	Long: `Clean up stale workspaces that are no longer valid.

For Git: Removes worktree entries that point to non-existent directories.
For jj: Updates stale workspaces.

Use --dry-run to see what would be cleaned without actually removing anything.
Use --force to skip confirmation prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		runner := internal.NewExecCommandRunner()

		// Detect VCS type
		vcsType, err := internal.DetectVCS(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting VCS: %v\n", err)
			os.Exit(1)
		}

		var service internal.VCSService
		switch vcsType {
		case internal.VCSGit:
			service = internal.NewGitService(runner)
		case internal.VCSJj:
			service = internal.NewJjService(runner)
		}

		workspaces, err := service.ListWorkspaces()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing workspaces: %v\n", err)
			os.Exit(1)
		}

		// Find stale workspaces
		var staleWorkspaces []internal.Workspace
		for _, ws := range workspaces {
			if ws.Status == internal.StatusStale {
				staleWorkspaces = append(staleWorkspaces, ws)
			}
		}

		if len(staleWorkspaces) == 0 {
			fmt.Println("No stale workspaces found.")
			return
		}

		// Show what will be cleaned
		fmt.Printf("Found %d stale workspace(s):\n", len(staleWorkspaces))
		for _, ws := range staleWorkspaces {
			fmt.Printf("  %s -> %s\n", ws.Name, ws.Path)
		}

		if cleanDryRun {
			fmt.Println("\nDry run mode - no changes made.")
			return
		}

		// Confirm before proceeding
		if !cleanForce {
			fmt.Print("\nProceed with cleanup? (y/N): ")
			var response string
			_, _ = fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cleanup cancelled.")
				return
			}
		}

		// Clean up stale workspaces
		if err := pruneStaleWorkspaces(runner, vcsType); err != nil {
			fmt.Fprintf(os.Stderr, "Error cleaning stale workspaces: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cleaned %d stale workspace(s).\n", len(staleWorkspaces))
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "Show what would be cleaned without making changes")
	cleanCmd.Flags().BoolVar(&cleanForce, "force", false, "Skip confirmation prompts")
}

func pruneStaleWorkspaces(runner internal.CommandRunner, vcsType internal.VCSType) error {
	switch vcsType {
	case internal.VCSGit:
		_, err := runner.Run("git worktree prune")
		return err
	case internal.VCSJj:
		// jj workspace update-stale updates outdated working copies
		_, err := runner.Run("jj workspace update-stale")
		return err
	default:
		return fmt.Errorf("unsupported VCS type: %s", vcsType)
	}
}

// Legacy function kept for compatibility
func pruneStaleWorktrees(runner internal.CommandRunner) error {
	_, err := runner.Run("git worktree prune")
	return err
}
