package internal

import (
	"strings"
)

// ParseJjWorkspaceList parses the output of jj workspace list command.
// Expected format (tab-separated):
// workspace-name\tworking-copy-commit-id\tpath
func ParseJjWorkspaceList(output string) []Workspace {
	var workspaces []Workspace

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// Only trim leading/trailing whitespace from the line itself,
		// not from individual fields (to preserve spaces in paths)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}

		// Trim whitespace from name and revision, but preserve path as-is
		// (except trailing whitespace which is likely from output formatting)
		workspaces = append(workspaces, Workspace{
			Name:     strings.TrimSpace(parts[0]),
			Revision: strings.TrimSpace(parts[1]),
			Path:     strings.TrimRight(parts[2], " \t"),
			Status:   StatusClean, // Will be updated by status check
		})
	}

	return workspaces
}

// ParseJjStatus parses jj status output and returns the workspace status.
// jj status shows uncommitted changes.
func ParseJjStatus(output string) Status {
	output = strings.TrimSpace(output)

	// Empty output or only whitespace means clean
	if output == "" {
		return StatusClean
	}

	// Check for "The working copy is clean" message
	if strings.Contains(output, "working copy is clean") {
		return StatusClean
	}

	// Any other output indicates changes
	return StatusDirty
}

// ParseJjBookmarks parses the output of jj log command for bookmarks.
// Expected format: space-separated bookmark names
func ParseJjBookmarks(output string) []string {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}

	var bookmarks []string
	for _, b := range strings.Fields(output) {
		b = strings.TrimSpace(b)
		// Remove trailing markers like @, *, etc.
		b = strings.TrimRight(b, "@*")
		if b != "" {
			bookmarks = append(bookmarks, b)
		}
	}

	return bookmarks
}
