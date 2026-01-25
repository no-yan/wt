package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type JjWorkspaceManager struct {
	jjService *JjService
}

func NewJjWorkspaceManager(jjService *JjService) *JjWorkspaceManager {
	return &JjWorkspaceManager{
		jjService: jjService,
	}
}

// AddWorkspace implements WorkspaceManager interface for jj
func (m *JjWorkspaceManager) AddWorkspace(repoPath, name string, opts AddOptions) (string, error) {
	if err := validateWorkspaceName(name); err != nil {
		return "", err
	}

	if err := validatePath(repoPath); err != nil {
		return "", fmt.Errorf("invalid repository path: %w", err)
	}

	workspacePath := filepath.Join(repoPath, "worktrees", name)

	// Ensure worktrees directory exists
	worktreesDir := filepath.Dir(workspacePath)
	if err := os.MkdirAll(worktreesDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	// Build jj workspace add command
	cmd := fmt.Sprintf("jj -R %s workspace add %s --name %s",
		shellescape(repoPath),
		shellescape(workspacePath),
		shellescape(name))

	// Add revision option if specified
	if opts.Revision != "" {
		cmd += fmt.Sprintf(" --revision %s", shellescape(opts.Revision))
	}

	if _, err := m.jjService.Run(cmd); err != nil {
		return "", fmt.Errorf("jj workspace add failed: %w", err)
	}

	return workspacePath, nil
}

// RemoveWorkspace implements WorkspaceManager interface for jj
func (m *JjWorkspaceManager) RemoveWorkspace(repoPath, name string) error {
	if err := validatePath(repoPath); err != nil {
		return fmt.Errorf("invalid repository path: %w", err)
	}

	if name == "" {
		return fmt.Errorf("workspace name cannot be empty")
	}

	// Get all workspaces to find the target
	workspaces, err := m.jjService.ListWorkspaces()
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	var targetWorkspace *Workspace
	for _, ws := range workspaces {
		if ws.Name == name {
			targetWorkspace = &ws
			break
		}
	}

	if targetWorkspace == nil {
		return fmt.Errorf("workspace %q not found", name)
	}

	// Safety check: don't remove the default workspace
	if name == JjDefaultWorkspace {
		return fmt.Errorf("cannot remove the default workspace")
	}

	// Safety check: warn if workspace has uncommitted changes
	if targetWorkspace.Status == StatusDirty {
		return fmt.Errorf("workspace %q has uncommitted changes, commit or abandon them first", name)
	}

	// jj workspace forget removes the workspace tracking
	forgetCmd := fmt.Sprintf("jj -R %s workspace forget %s",
		shellescape(repoPath),
		shellescape(name))

	if _, err := m.jjService.Run(forgetCmd); err != nil {
		return fmt.Errorf("jj workspace forget failed: %w", err)
	}

	// Remove the workspace directory
	if err := os.RemoveAll(targetWorkspace.Path); err != nil {
		return fmt.Errorf("failed to remove workspace directory: %w", err)
	}

	return nil
}

// RemoveMultipleWorkspaces implements WorkspaceManager interface for jj
func (m *JjWorkspaceManager) RemoveMultipleWorkspaces(repoPath string, names []string) error {
	if err := validatePath(repoPath); err != nil {
		return fmt.Errorf("invalid repository path: %w", err)
	}

	if len(names) == 0 {
		return fmt.Errorf("at least one workspace name is required")
	}

	// Get all workspaces once for validation
	workspaces, err := m.jjService.ListWorkspaces()
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	// Create lookup map
	workspaceMap := make(map[string]*Workspace)
	for i := range workspaces {
		workspaceMap[workspaces[i].Name] = &workspaces[i]
	}

	// Phase 1: Validate all targets
	var targetsToRemove []*Workspace
	for _, name := range names {
		if name == "" {
			return fmt.Errorf("workspace name cannot be empty")
		}

		targetWorkspace, exists := workspaceMap[name]
		if !exists {
			return fmt.Errorf("workspace %q not found", name)
		}

		if name == JjDefaultWorkspace {
			return fmt.Errorf("cannot remove the default workspace")
		}

		if targetWorkspace.Status == StatusDirty {
			return fmt.Errorf("workspace %q has uncommitted changes, commit or abandon them first", name)
		}

		targetsToRemove = append(targetsToRemove, targetWorkspace)
	}

	// Phase 2: Execute removals
	for _, target := range targetsToRemove {
		forgetCmd := fmt.Sprintf("jj -R %s workspace forget %s",
			shellescape(repoPath),
			shellescape(target.Name))

		if _, err := m.jjService.Run(forgetCmd); err != nil {
			return fmt.Errorf("failed to forget workspace %q: %w", target.Name, err)
		}

		if err := os.RemoveAll(target.Path); err != nil {
			return fmt.Errorf("failed to remove workspace directory %q: %w", target.Name, err)
		}
	}

	return nil
}

func validateWorkspaceName(name string) error {
	if name == "" {
		return fmt.Errorf("workspace name cannot be empty")
	}

	if strings.ContainsAny(name, ";&|`$(){}[]<>") {
		return fmt.Errorf("invalid characters in workspace name: %q", name)
	}

	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("workspace name cannot start with dash: %q", name)
	}

	return nil
}
