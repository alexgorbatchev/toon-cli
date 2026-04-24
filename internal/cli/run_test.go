package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		args        []string
		stdin       string
		stdinIsTTY  bool
		version     string
		wantStdout  string
		wantStderr  string
		wantErrText string
	}{
		{
			name:       "converts piped json",
			stdin:      `{"name":"Ada","id":1}`,
			version:    "dev",
			wantStdout: "name: Ada\nid: 1\n",
		},
		{
			name:       "prints help",
			args:       []string{"--help"},
			version:    "dev",
			wantStdout: "Usage:\n  toon < input.json > output.toon\n\nSupported input formats: JSON, JSONC, NDJSON.\n\nFlags:\n  --help     Show this help message\n  --version  Print the CLI version\n",
		},
		{
			name:       "prints version",
			args:       []string{"--version"},
			version:    "v1.2.3",
			wantStdout: "toon v1.2.3\n",
		},
		{
			name:        "rejects unknown flag",
			args:        []string{"--wat"},
			version:     "dev",
			wantStderr:  "unknown flag: --wat\n",
			wantErrText: "unknown flag",
		},
		{
			name:        "requires piped input without flags",
			stdinIsTTY:  true,
			version:     "dev",
			wantStderr:  "no input detected on stdin\n",
			wantErrText: "no input detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			err := Run(tt.args, Environment{
				Stdin:      strings.NewReader(tt.stdin),
				Stdout:     stdout,
				Stderr:     stderr,
				StdinIsTTY: tt.stdinIsTTY,
				Version:    tt.version,
			})

			if tt.wantErrText == "" && err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if tt.wantErrText != "" {
				if err == nil {
					t.Fatal("Run() error = nil, want error")
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("Run() error = %q, want substring %q", err.Error(), tt.wantErrText)
				}
			}

			if got := stdout.String(); got != tt.wantStdout {
				t.Fatalf("stdout = %q, want %q", got, tt.wantStdout)
			}

			if got := stderr.String(); got != tt.wantStderr {
				t.Fatalf("stderr = %q, want %q", got, tt.wantStderr)
			}
		})
	}
}
