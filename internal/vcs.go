package internal

// VCSType identifies the version control system
type VCSType string

const (
	VCSGit VCSType = "git"
	VCSJj  VCSType = "jj"
)

// Workspace represents a unified view of a Git worktree or jj workspace
type Workspace struct {
	Name      string   // Workspace/worktree name
	Path      string   // Absolute path to workspace directory
	Revision  string   // Git: HEAD commit, jj: change ID
	Status    Status   // Clean/Dirty/Stale
	Branch    string   // Git only: branch name (empty for jj)
	Bookmarks []string // jj only: bookmarks in trunk()..@ range
}

// VCSService defines the interface for version control read operations
type VCSService interface {
	Type() VCSType
	ListWorkspaces() ([]Workspace, error)
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
