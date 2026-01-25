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

		// Detect VCS type
		vcsType, err := internal.DetectVCS(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting VCS: %v\n", err)
			os.Exit(1)
		}

		// Check if VCS is available
		if err := internal.CheckVCSAvailable(runner, vcsType); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		repoPath, err := internal.GetRepoRoot(runner, vcsType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding repository: %v\n", err)
			os.Exit(1)
		}

		var manager internal.WorkspaceManager
		opts := internal.AddOptions{Revision: addRevision}

		switch vcsType {
		case internal.VCSGit:
			gitService := internal.NewGitService(runner)
			manager = internal.NewWorktreeManager(gitService, runner)
			opts.Branch = name // For Git, name is the branch
		case internal.VCSJj:
			jjService := internal.NewJjService(runner)
			manager = internal.NewJjWorkspaceManager(jjService)
		}

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
