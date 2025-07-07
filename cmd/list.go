package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/no-yan/wt/internal"
	"github.com/spf13/cobra"
)

var (
	listDirtyOnly bool
	listVerbose   bool
	listNamesOnly bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all worktrees",
	Long: `List all git worktrees with their status and branch information.

Filtering options:
  --dirty       Show only worktrees with uncommitted changes
  --verbose     Show detailed git status information
  --names-only  Show only worktree names (useful for scripting)`,
	Run: func(cmd *cobra.Command, args []string) {
		runner := internal.NewExecCommandRunner()
		service := internal.NewGitService(runner)

		worktrees, err := service.ListWorktrees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing worktrees: %v\n", err)
			os.Exit(1)
		}

		// Apply filters
		filtered := filterWorktrees(worktrees, listDirtyOnly)

		// Format output based on flags
		if listNamesOnly {
			formatWorktreeNames(filtered, os.Stdout)
		} else if listVerbose {
			formatWorktreeListVerbose(filtered, os.Stdout, service)
		} else {
			formatWorktreeList(filtered, os.Stdout)
		}
	},
}

func init() {
	listCmd.Flags().BoolVar(&listDirtyOnly, "dirty", false, "Show only worktrees with uncommitted changes")
	listCmd.Flags().BoolVar(&listVerbose, "verbose", false, "Show detailed git status information")
	listCmd.Flags().BoolVar(&listNamesOnly, "names-only", false, "Show only worktree names")
}

func filterWorktrees(worktrees []internal.Worktree, dirtyOnly bool) []internal.Worktree {
	if !dirtyOnly {
		return worktrees
	}

	var filtered []internal.Worktree
	for _, wt := range worktrees {
		if wt.Status == internal.StatusDirty {
			filtered = append(filtered, wt)
		}
	}
	return filtered
}

func formatWorktreeList(worktrees []internal.Worktree, w io.Writer) {
	for _, wt := range worktrees {
		status := "clean"
		switch wt.Status {
		case internal.StatusDirty:
			status = "dirty"
		case internal.StatusStale:
			status = "stale"
		}

		if _, err := fmt.Fprintf(w, "%s\t%s\t(%s)\n", wt.Name(), wt.Path, status); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		}
	}
}

func formatWorktreeNames(worktrees []internal.Worktree, w io.Writer) {
	for _, wt := range worktrees {
		if _, err := fmt.Fprintf(w, "%s\n", wt.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		}
	}
}

func formatWorktreeListVerbose(worktrees []internal.Worktree, w io.Writer, service *internal.GitService) {
	for _, wt := range worktrees {
		status := "clean"
		switch wt.Status {
		case internal.StatusDirty:
			status = "dirty"
		case internal.StatusStale:
			status = "stale"
		}

		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", wt.Name(), wt.Branch, wt.Path, status); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		}

		// Show detailed status for dirty worktrees
		if wt.Status == internal.StatusDirty {
			if statusOutput, err := service.GetDetailedStatus(wt.Path); err == nil {
				if _, err := fmt.Fprintf(w, "  Changes:\n"); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
				}
				for _, line := range statusOutput {
					if _, err := fmt.Fprintf(w, "    %s\n", line); err != nil {
						fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
					}
				}
			}
		}
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		}
	}
}
