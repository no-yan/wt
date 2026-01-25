package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

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
	Short:   "List all workspaces",
	Long: `List all workspaces with their status information.

For Git: Shows branch name and worktree path.
For jj: Shows change ID and bookmarks (if any in trunk()..@ range).

Filtering options:
  --dirty       Show only workspaces with uncommitted changes
  --verbose     Show detailed status information
  --names-only  Show only workspace names (useful for scripting)`,
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

		// Apply filters
		filtered := filterWorkspaces(workspaces, listDirtyOnly)

		// Format output based on flags and VCS type
		if listNamesOnly {
			formatWorkspaceNames(filtered, os.Stdout)
		} else if listVerbose {
			formatWorkspaceListVerbose(filtered, os.Stdout, service)
		} else {
			formatWorkspaceList(filtered, os.Stdout, vcsType)
		}
	},
}

func init() {
	listCmd.Flags().BoolVar(&listDirtyOnly, "dirty", false, "Show only workspaces with uncommitted changes")
	listCmd.Flags().BoolVar(&listVerbose, "verbose", false, "Show detailed status information")
	listCmd.Flags().BoolVar(&listNamesOnly, "names-only", false, "Show only workspace names")
}

func filterWorkspaces(workspaces []internal.Workspace, dirtyOnly bool) []internal.Workspace {
	if !dirtyOnly {
		return workspaces
	}

	var filtered []internal.Workspace
	for _, ws := range workspaces {
		if ws.Status == internal.StatusDirty {
			filtered = append(filtered, ws)
		}
	}
	return filtered
}

func formatWorkspaceList(workspaces []internal.Workspace, w io.Writer, vcsType internal.VCSType) {
	cw := calculateWorkspaceColumnWidths(workspaces, vcsType)
	for _, ws := range workspaces {
		status := formatStatus(ws.Status)
		if vcsType == internal.VCSJj {
			// jj: NAME  CHANGE  BOOKMARK  (STATUS)
			bookmark := ""
			if len(ws.Bookmarks) > 0 {
				bookmark = strings.Join(ws.Bookmarks, ", ")
			}
			if _, err := fmt.Fprintf(w, "%-*s  %-*s  %-*s  (%s)\n",
				cw.name, ws.Name,
				cw.revision, ws.Revision,
				cw.bookmark, bookmark,
				status); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			}
		} else {
			// Git: NAME  PATH  (STATUS)
			if _, err := fmt.Fprintf(w, "%-*s  %-*s  (%s)\n",
				cw.name, ws.Name,
				cw.path, ws.Path,
				status); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			}
		}
	}
}

func formatWorkspaceNames(workspaces []internal.Workspace, w io.Writer) {
	for _, ws := range workspaces {
		if _, err := fmt.Fprintf(w, "%s\n", ws.Name); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		}
	}
}

// workspaceColumnWidths holds the maximum width for each column
type workspaceColumnWidths struct {
	name     int
	branch   int
	path     int
	revision int
	bookmark int
}

// calculateWorkspaceColumnWidths calculates the maximum width for each column
func calculateWorkspaceColumnWidths(workspaces []internal.Workspace, vcsType internal.VCSType) workspaceColumnWidths {
	var cw workspaceColumnWidths
	for _, ws := range workspaces {
		cw.name = max(cw.name, utf8.RuneCountInString(ws.Name))
		cw.branch = max(cw.branch, utf8.RuneCountInString(ws.Branch))
		cw.path = max(cw.path, utf8.RuneCountInString(ws.Path))
		cw.revision = max(cw.revision, utf8.RuneCountInString(ws.Revision))
		if len(ws.Bookmarks) > 0 {
			bookmark := strings.Join(ws.Bookmarks, ", ")
			cw.bookmark = max(cw.bookmark, utf8.RuneCountInString(bookmark))
		}
	}
	return cw
}

// formatStatus converts a workspace status to its string representation
func formatStatus(status internal.Status) string {
	switch status {
	case internal.StatusDirty:
		return "dirty"
	case internal.StatusStale:
		return "stale"
	default:
		return "clean"
	}
}

func formatWorkspaceListVerbose(workspaces []internal.Workspace, w io.Writer, service internal.VCSService) {
	cw := calculateWorkspaceColumnWidths(workspaces, service.Type())
	for _, ws := range workspaces {
		status := formatStatus(ws.Status)
		if service.Type() == internal.VCSJj {
			bookmark := ""
			if len(ws.Bookmarks) > 0 {
				bookmark = strings.Join(ws.Bookmarks, ", ")
			}
			if _, err := fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  (%s)\n",
				cw.name, ws.Name,
				cw.revision, ws.Revision,
				cw.bookmark, bookmark,
				cw.path, ws.Path,
				status); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			}
		} else {
			if _, err := fmt.Fprintf(w, "%-*s  %-*s  %-*s  (%s)\n",
				cw.name, ws.Name,
				cw.branch, ws.Branch,
				cw.path, ws.Path,
				status); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			}
		}

		// Show detailed status for dirty workspaces
		if ws.Status == internal.StatusDirty {
			if statusOutput, err := service.GetDetailedStatus(ws.Path); err == nil {
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

// Legacy support for existing tests
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

// columnWidths holds the maximum width for each column (legacy)
type columnWidths struct {
	name   int
	branch int
	path   int
}

// calculateColumnWidths calculates the maximum width for each column (legacy)
func calculateColumnWidths(worktrees []internal.Worktree) columnWidths {
	var ws columnWidths
	for _, wt := range worktrees {
		ws.name = max(ws.name, utf8.RuneCountInString(wt.Name()))
		ws.branch = max(ws.branch, utf8.RuneCountInString(wt.Branch))
		ws.path = max(ws.path, utf8.RuneCountInString(wt.Path))
	}
	return ws
}

// formatWorktreeList formats worktrees for display (legacy for tests)
func formatWorktreeList(worktrees []internal.Worktree, w io.Writer) {
	ws := calculateColumnWidths(worktrees)
	for _, wt := range worktrees {
		status := formatStatus(wt.Status)
		if _, err := fmt.Fprintf(w, "%-*s  %-*s  (%s)\n", ws.name, wt.Name(), ws.path, wt.Path, status); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		}
	}
}

// formatWorktreeListVerbose formats worktrees with verbose output (legacy for tests)
func formatWorktreeListVerbose(worktrees []internal.Worktree, w io.Writer, service *internal.GitService) {
	ws := calculateColumnWidths(worktrees)
	for _, wt := range worktrees {
		status := formatStatus(wt.Status)
		if _, err := fmt.Fprintf(w, "%-*s  %-*s  %-*s  (%s)\n", ws.name, wt.Name(), ws.branch, wt.Branch, ws.path, wt.Path, status); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		}

		// Show detailed status for dirty worktrees
		if wt.Status == internal.StatusDirty && service != nil {
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
