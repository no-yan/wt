package internal

import (
	"testing"
)

func TestParseJjWorkspaceList(t *testing.T) {
	repoRoot := "/repo"
	tests := []struct {
		name   string
		output string
		want   []Workspace
	}{
		{
			name:   "empty output",
			output: "",
			want:   nil,
		},
		{
			name:   "single workspace (default)",
			output: "default\tabc12345\n",
			want: []Workspace{
				{Name: "default", Revision: "abc12345", Path: "/repo", Status: StatusClean},
			},
		},
		{
			name:   "multiple workspaces",
			output: "default\tabc12345\nfeature\tdef67890\n",
			want: []Workspace{
				{Name: "default", Revision: "abc12345", Path: "/repo", Status: StatusClean},
				{Name: "feature", Revision: "def67890", Path: "/repo/worktrees/feature", Status: StatusClean},
			},
		},
		{
			name:   "output with extra whitespace",
			output: "  default\tabc12345  \n",
			want: []Workspace{
				{Name: "default", Revision: "abc12345", Path: "/repo", Status: StatusClean},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseJjWorkspaceList(tt.output, repoRoot)

			if len(got) != len(tt.want) {
				t.Errorf("ParseJjWorkspaceList() returned %d workspaces, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].Name != tt.want[i].Name {
					t.Errorf("workspace[%d].Name = %q, want %q", i, got[i].Name, tt.want[i].Name)
				}
				if got[i].Revision != tt.want[i].Revision {
					t.Errorf("workspace[%d].Revision = %q, want %q", i, got[i].Revision, tt.want[i].Revision)
				}
				if got[i].Path != tt.want[i].Path {
					t.Errorf("workspace[%d].Path = %q, want %q", i, got[i].Path, tt.want[i].Path)
				}
			}
		})
	}
}

func TestParseJjStatus(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   Status
	}{
		{
			name:   "empty output means clean",
			output: "",
			want:   StatusClean,
		},
		{
			name:   "whitespace only means clean",
			output: "   \n  ",
			want:   StatusClean,
		},
		{
			name:   "working copy is clean message",
			output: "The working copy is clean",
			want:   StatusClean,
		},
		{
			name:   "has changes means dirty",
			output: "Working copy changes:\nM file.txt\n",
			want:   StatusDirty,
		},
		{
			name:   "added file means dirty",
			output: "A newfile.txt\n",
			want:   StatusDirty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseJjStatus(tt.output)
			if got != tt.want {
				t.Errorf("ParseJjStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJjBookmarks(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   []string
	}{
		{
			name:   "empty output",
			output: "",
			want:   nil,
		},
		{
			name:   "single bookmark",
			output: "main",
			want:   []string{"main"},
		},
		{
			name:   "multiple bookmarks",
			output: "main feature-auth bugfix",
			want:   []string{"main", "feature-auth", "bugfix"},
		},
		{
			name:   "bookmark with @ marker",
			output: "main@ feature",
			want:   []string{"main", "feature"},
		},
		{
			name:   "bookmark with * marker",
			output: "main* feature",
			want:   []string{"main", "feature"},
		},
		{
			name:   "output with extra whitespace",
			output: "  main   feature  ",
			want:   []string{"main", "feature"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseJjBookmarks(tt.output)

			if len(got) != len(tt.want) {
				t.Errorf("ParseJjBookmarks() returned %d bookmarks, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("bookmark[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
