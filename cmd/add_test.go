package cmd

import (
	"testing"
)

func TestGetRepoRoot(t *testing.T) {
	tests := []struct {
		name      string
		gitOutput string
		want      string
		wantErr   bool
	}{
		{
			name:      "valid repo",
			gitOutput: "/repo/path\n",
			want:      "/repo/path",
			wantErr:   false,
		},
		{
			name:      "repo with spaces",
			gitOutput: "/repo path/project\n",
			want:      "/repo path/project",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRunner := &testMockCommandRunner{
				outputs: map[string]string{
					"git rev-parse --show-toplevel": tt.gitOutput,
				},
			}

			got, err := getRepoRoot(mockRunner)

			if (err != nil) != tt.wantErr {
				t.Errorf("getRepoRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("getRepoRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testMockCommandRunner struct {
	outputs map[string]string
}

func (m *testMockCommandRunner) Run(command string) (string, error) {
	if output, exists := m.outputs[command]; exists {
		return output, nil
	}
	return "", nil
}
