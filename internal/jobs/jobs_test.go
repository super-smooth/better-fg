package jobs

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Job
	}{
		{
			name:  "single suspended job",
			input: "[1]  - suspended  nvim",
			expected: []Job{
				{ID: 1, State: "suspended", Command: "nvim"},
			},
		},
		{
			name:  "multiple jobs",
			input: "[1]  - suspended  nvim\n[2]  + suspended  amazon-q chat --resume",
			expected: []Job{
				{ID: 1, State: "suspended", Command: "nvim"},
				{ID: 2, State: "suspended", Command: "amazon-q chat --resume"},
			},
		},
		{
			name:  "running job",
			input: "[1]  + running    sleep 100 &",
			expected: []Job{
				{ID: 1, State: "running", Command: "sleep 100 &"},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []Job{},
		},
		{
			name:     "invalid format",
			input:    "not a job line",
			expected: []Job{},
		},
		// Different shell formats
		{
			name:  "bash format with plus/minus only",
			input: "[1]+  Stopped                 vim file.txt\n[2]-  Running                 sleep 60 &",
			expected: []Job{
				{ID: 1, State: "Stopped", Command: "vim file.txt"},
				{ID: 2, State: "Running", Command: "sleep 60 &"},
			},
		},
		{
			name:  "zsh format with spaces",
			input: "[1]  + suspended (tty output)  vim\n[2]  - running    sleep 100",
			expected: []Job{
				{ID: 1, State: "suspended", Command: "(tty output)  vim"},
				{ID: 2, State: "running", Command: "sleep 100"},
			},
		},
		{
			name:  "format with job state variations",
			input: "[1]  + stopped    vim\n[2]  - done      ls -la",
			expected: []Job{
				{ID: 1, State: "stopped", Command: "vim"},
				{ID: 2, State: "done", Command: "ls -la"},
			},
		},
		{
			name:  "format without spaces after bracket",
			input: "[1]+ Stopped vim\n[2]- Running sleep 60",
			expected: []Job{
				{ID: 1, State: "Stopped", Command: "vim"},
				{ID: 2, State: "Running", Command: "sleep 60"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)

			jobs, err := Parse(reader)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if len(jobs) != len(tt.expected) {
				t.Fatalf("Parse() got %d jobs, want %d", len(jobs), len(tt.expected))
			}

			for i, job := range jobs {
				expected := tt.expected[i]
				if job.ID != expected.ID {
					t.Errorf("Job[%d].ID = %d, want %d", i, job.ID, expected.ID)
				}

				if job.State != expected.State {
					t.Errorf("Job[%d].State = %q, want %q", i, job.State, expected.State)
				}

				if job.Command != expected.Command {
					t.Errorf("Job[%d].Command = %q, want %q", i, job.Command, expected.Command)
				}
			}
		})
	}
}
