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
	Short:   "Remove one or more workspaces",
	Long: `Remove one or more workspaces by name.

Safety checks:
- Cannot remove the main worktree (Git) or default workspace (jj)
- Cannot remove workspaces with uncommitted changes
- When removing multiple workspaces, validates all before removing any (fail-fast)`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runner := internal.NewExecCommandRunner()

		// Detect VCS type
		vcsType, err := internal.DetectVCS(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting VCS: %v\n", err)
			os.Exit(1)
		}

		repoPath, err := internal.GetRepoRoot(runner, vcsType)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding repository: %v\n", err)
			os.Exit(1)
		}

		var manager internal.WorkspaceManager
		switch vcsType {
		case internal.VCSGit:
			gitService := internal.NewGitService(runner)
			manager = internal.NewWorktreeManager(gitService, runner)
		case internal.VCSJj:
			jjService := internal.NewJjService(runner)
			manager = internal.NewJjWorkspaceManager(jjService)
		}

		if len(args) == 1 {
			// Single workspace removal
			if err := manager.RemoveWorkspace(repoPath, args[0]); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing workspace: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Removed workspace: %s\n", args[0])
		} else {
			// Multiple workspace removal with fail-fast validation
			if err := manager.RemoveMultipleWorkspaces(repoPath, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing workspaces: %v\n", err)
				os.Exit(1)
			}
			for _, name := range args {
				fmt.Printf("Removed workspace: %s\n", name)
			}
		}
	},
}
