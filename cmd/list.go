package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/no-yan/wrkt/internal"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long:  "List all git worktrees with their status and branch information.",
	Run: func(cmd *cobra.Command, args []string) {
		runner := internal.NewExecCommandRunner()
		service := internal.NewGitService(runner)

		worktrees, err := service.ListWorktrees()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing worktrees: %v\n", err)
			os.Exit(1)
		}

		formatWorktreeList(worktrees, os.Stdout)
	},
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

		fmt.Fprintf(w, "%s\t%s\t(%s)\n", wt.Name(), wt.Path, status)
	}
}
