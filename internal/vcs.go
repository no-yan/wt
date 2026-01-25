package internal

// VCSType identifies the version control system
type VCSType string

const (
	VCSGit VCSType = "git"
	VCSJj  VCSType = "jj"
)

// Workspace represents a unified view of a Git worktree or jj workspace
type Workspace struct {
	Name     string   // Workspace/worktree name
	Path     string   // Absolute path to workspace directory
	Revision string   // Git: HEAD commit, jj: change ID
	Status   Status   // Clean/Dirty/Stale (only populated if ListOptions.IncludeStatus)
	Labels   []string // VCS-agnostic labels: Git branch or jj bookmarks
	Error    error    // Error encountered while fetching details (partial success)
}

// ListOptions controls what information to fetch when listing workspaces.
// This allows callers to avoid unnecessary commands when they don't need all details.
type ListOptions struct {
	IncludeStatus bool // Fetch status for each workspace (adds N commands)
	IncludeLabels bool // Fetch labels (branch/bookmarks) for each workspace
}

// VCSService defines the interface for version control read operations
type VCSService interface {
	Type() VCSType
	ListWorkspaces(opts ListOptions) ([]Workspace, error)
	GetDetailedStatus(workspacePath string) ([]string, error)
}

// WorkspaceManager defines operations that modify workspaces
type WorkspaceManager interface {
	AddWorkspace(repoPath, name string, opts AddOptions) (string, error)
	RemoveWorkspace(repoPath, name string) error
	RemoveMultipleWorkspaces(repoPath string, names []string) error
}

// AddOptions provides VCS-specific options for workspace creation
type AddOptions struct {
	Branch   string // Git: branch name (required)
	Revision string // jj: base revision (optional)
}
