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

		service, err := internal.CreateVCSService(runner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Only fetch status when needed (for filtering or display)
		opts := internal.ListOptions{
			IncludeStatus: listDirtyOnly || listVerbose || !listNamesOnly,
			IncludeLabels: !listNamesOnly, // Labels for display
		}

		workspaces, err := service.ListWorkspaces(opts)
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
			formatWorkspaceList(filtered, os.Stdout, service.Type())
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
			// jj: NAME  CHANGE  LABELS  (STATUS)
			labels := ""
			if len(ws.Labels) > 0 {
				labels = strings.Join(ws.Labels, ", ")
			}
			if _, err := fmt.Fprintf(w, "%-*s  %-*s  %-*s  (%s)\n",
				cw.name, ws.Name,
				cw.revision, ws.Revision,
				cw.labels, labels,
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
	path     int
	revision int
	labels   int
}

// calculateWorkspaceColumnWidths calculates the maximum width for each column
func calculateWorkspaceColumnWidths(workspaces []internal.Workspace, vcsType internal.VCSType) workspaceColumnWidths {
	var cw workspaceColumnWidths
	for _, ws := range workspaces {
		cw.name = max(cw.name, utf8.RuneCountInString(ws.Name))
		cw.path = max(cw.path, utf8.RuneCountInString(ws.Path))
		cw.revision = max(cw.revision, utf8.RuneCountInString(ws.Revision))
		if len(ws.Labels) > 0 {
			labels := strings.Join(ws.Labels, ", ")
			cw.labels = max(cw.labels, utf8.RuneCountInString(labels))
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
		labels := ""
		if len(ws.Labels) > 0 {
			labels = strings.Join(ws.Labels, ", ")
		}

		if service.Type() == internal.VCSJj {
			if _, err := fmt.Fprintf(w, "%-*s  %-*s  %-*s  %-*s  (%s)\n",
				cw.name, ws.Name,
				cw.revision, ws.Revision,
				cw.labels, labels,
				cw.path, ws.Path,
				status); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			}
		} else {
			if _, err := fmt.Fprintf(w, "%-*s  %-*s  %-*s  (%s)\n",
				cw.name, ws.Name,
				cw.labels, labels,
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
