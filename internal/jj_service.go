package internal

import (
	"fmt"
	"strings"
)

// jj-specific constants
const (
	// JjDefaultWorkspace is the name of the default workspace in jj
	JjDefaultWorkspace = "default"
	// JjTrunkRange is the revision range for finding bookmarks between trunk and current
	JjTrunkRange = "trunk()..@"
	// jjWorkspaceListTemplate defines the output format for jj workspace list.
	// Format: name<TAB>commit_id<TAB>path<NEWLINE>
	// This template is the single source of truth for the expected format.
	jjWorkspaceListTemplate = `name ++ "\t" ++ working_copy.commit_id().short(8) ++ "\t" ++ working_copy_path ++ "\n"`
)

type JjService struct {
	runner CommandRunner
}

func NewJjService(runner CommandRunner) *JjService {
	return &JjService{runner: runner}
}

// Type implements VCSService interface
func (j *JjService) Type() VCSType {
	return VCSJj
}

// ListWorkspaces implements VCSService interface.
// Use opts to control what additional information to fetch.
// Without opts, only basic workspace info is returned (no N+1 queries).
func (j *JjService) ListWorkspaces(opts ListOptions) ([]Workspace, error) {
	output, err := j.runner.Run("jj workspace list -T '" + jjWorkspaceListTemplate + "'")
	if err != nil {
		return nil, fmt.Errorf("failed to list jj workspaces: %w", err)
	}

	workspaces := ParseJjWorkspaceList(output)

	// Only fetch additional details if requested
	for i := range workspaces {
		if opts.IncludeStatus {
			statusOutput, err := j.runStatus(workspaces[i].Path)
			if err != nil {
				workspaces[i].Status = StatusStale
				workspaces[i].Error = err
			} else {
				workspaces[i].Status = ParseJjStatus(statusOutput)
			}
		}

		if opts.IncludeLabels {
			bookmarks, _ := j.GetBookmarksInRange(workspaces[i].Path, JjTrunkRange)
			if len(bookmarks) > 0 {
				workspaces[i].Labels = bookmarks
			}
		}
	}

	return workspaces, nil
}

// GetDetailedStatus implements VCSService interface
func (j *JjService) GetDetailedStatus(workspacePath string) ([]string, error) {
	output, err := j.runStatus(workspacePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get status for %s: %w", workspacePath, err)
	}

	if strings.TrimSpace(output) == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var statusLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "Working copy") {
			statusLines = append(statusLines, line)
		}
	}

	return statusLines, nil
}

// GetBookmarksInRange returns bookmarks in the specified revision range for a workspace.
// Returns empty slice (not error) if range is empty or invalid, as this is expected for new workspaces.
func (j *JjService) GetBookmarksInRange(workspacePath, revRange string) ([]string, error) {
	cmd := fmt.Sprintf("jj -R %s log -r %s --no-graph -T 'bookmarks ++ \" \"'",
		shellescape(workspacePath),
		shellescape(revRange))

	output, err := j.runner.Run(cmd)
	if err != nil {
		// Range might be empty or invalid - this is expected for new workspaces
		// Return empty slice rather than nil to distinguish from actual errors
		return []string{}, nil
	}

	return ParseJjBookmarks(output), nil
}

// runStatus executes jj status and returns raw output.
// This is the single source for status commands to avoid duplication.
func (j *JjService) runStatus(workspacePath string) (string, error) {
	return j.runner.Run(fmt.Sprintf("jj -R %s status", shellescape(workspacePath)))
}
