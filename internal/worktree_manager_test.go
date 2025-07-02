package internal

import (
	"testing"
)

func TestWorktreeManager_AddWorktree(t *testing.T) {
	tests := []struct {
		name         string
		repoPath     string
		branch       string
		wantPath     string
		wantCommands []string
		wantErr      bool
	}{
		{
			name:     "simple branch",
			repoPath: "/repo",
			branch:   "feature",
			wantPath: "/repo/worktrees/feature",
			wantCommands: []string{
				"mkdir -p /repo/worktrees",
				"git -C /repo worktree add /repo/worktrees/feature feature",
			},
			wantErr: false,
		},
		{
			name:     "branch with slash",
			repoPath: "/repo",
			branch:   "feature/auth",
			wantPath: "/repo/worktrees/feature-auth",
			wantCommands: []string{
				"mkdir -p /repo/worktrees",
				"git -C /repo worktree add /repo/worktrees/feature-auth feature/auth",
			},
			wantErr: false,
		},
		{
			name:     "nested branch",
			repoPath: "/repo",
			branch:   "feature/api/v2",
			wantPath: "/repo/worktrees/feature-api-v2",
			wantCommands: []string{
				"mkdir -p /repo/worktrees",
				"git -C /repo worktree add /repo/worktrees/feature-api-v2 feature/api/v2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &MockCommandRunner{
				outputs: make(map[string]string),
			}

			service := NewGitService(mockRunner)
			manager := NewWorktreeManager(service, mockRunner)

			gotPath, err := manager.AddWorktree(tt.repoPath, tt.branch)

			if (err != nil) != tt.wantErr {
				t.Errorf("WorktreeManager.AddWorktree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotPath != tt.wantPath {
				t.Errorf("WorktreeManager.AddWorktree() path = %v, want %v", gotPath, tt.wantPath)
			}

			if len(mockRunner.commands) != len(tt.wantCommands) {
				t.Errorf("WorktreeManager.AddWorktree() commands = %v, want %v", mockRunner.commands, tt.wantCommands)
				return
			}

			for i, cmd := range tt.wantCommands {
				if i >= len(mockRunner.commands) || mockRunner.commands[i] != cmd {
					t.Errorf("WorktreeManager.AddWorktree() command[%d] = %v, want %v", i, mockRunner.commands[i], cmd)
				}
			}
		})
	}
}

func TestGenerateWorktreePath(t *testing.T) {
	tests := []struct {
		name     string
		repoPath string
		branch   string
		want     string
	}{
		{
			name:     "simple branch",
			repoPath: "/repo",
			branch:   "main",
			want:     "/repo/worktrees/main",
		},
		{
			name:     "feature branch",
			repoPath: "/home/user/project",
			branch:   "feature/auth",
			want:     "/home/user/project/worktrees/feature-auth",
		},
		{
			name:     "deep nested branch",
			repoPath: "/repo",
			branch:   "feature/api/v2/auth",
			want:     "/repo/worktrees/feature-api-v2-auth",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateWorktreePath(tt.repoPath, tt.branch); got != tt.want {
				t.Errorf("GenerateWorktreePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
