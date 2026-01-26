package cmd

import (
	"fmt"
	"os"

	"github.com/no-yan/wt/internal"
	"github.com/spf13/cobra"
)

var addRevision string

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new workspace",
	Long: `Add a new workspace (git worktree or jj workspace).

For Git repositories: Creates a new worktree with the specified branch.
For jj repositories: Creates a new workspace with an optional base revision.

Examples:
  wt add feature-auth          # Git: create branch + worktree, jj: create workspace
  wt add my-feature -r main    # jj: create workspace based on main`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		runner := internal.NewExecCommandRunner()

		manager, vcsType, err := internal.CreateWorkspaceManager(runner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		repoPath, err := internal.GetRepoRoot(runner, vcsType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding repository: %v\n", err)
			os.Exit(1)
		}

		opts := internal.AddOptions{Revision: addRevision}
		workspacePath, err := manager.AddWorkspace(repoPath, name, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding workspace: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added workspace: %s\n", workspacePath)
	},
}

func init() {
	addCmd.Flags().StringVarP(&addRevision, "revision", "r", "", "Base revision (jj only)")
}
