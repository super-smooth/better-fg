package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/super-smooth/better-fg/internal/jobs"
	"github.com/super-smooth/better-fg/internal/tui"
)

const shellFunc = `
bfg() {
  # Use a temporary file to capture all job output reliably
  local temp_file=$(mktemp)
  jobs > "$temp_file" 2>&1

  # If there are no jobs, clean up and return
  if [ ! -s "$temp_file" ]; then
    rm "$temp_file"
    echo "No background jobs." >&2
    return
  fi

  # Run the Go program and get the selected job (pass through any arguments like --verbose)
  selected_job=$(cat "$temp_file" | command better-fg "$@")
  rm "$temp_file"

  # If a job was selected, bring it to the foreground
  if [ -n "$selected_job" ]; then
    fg "$selected_job"
  fi
}
`

var rootCmd = &cobra.Command{
	Use:   "better-fg",
	Short: "A more userfriendly and interactive version of `fg`.",
	Long: `better-fg is a CLI tool that provides an interactive TUI for selecting from multiple background jobs.
It also supports fuzzy searching for jobs and falls back to normal fg behavior when there is only one background job.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Check if stdin has data
		stat, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("failed to check stdin: %w", err)
		}

		if (stat.Mode() & os.ModeCharDevice) != 0 {
			// No piped input, stdin is a terminal
			fmt.Fprintln(os.Stderr, "No job data provided. Use 'jobs | better-fg' or the 'bfg' shell function.")
			return nil
		}

		// Read all input for debugging
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "DEBUG: Received input: %q\n", string(input))
		}

		jobs, err := jobs.Parse(strings.NewReader(string(input)))
		if err != nil {
			return fmt.Errorf("failed to parse jobs: %w", err)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "DEBUG: Parsed %d jobs\n", len(jobs))
			for i, job := range jobs {
				fmt.Fprintf(os.Stderr, "DEBUG: Job %d: ID=%d, State=%s, Command=%s\n", i, job.ID, job.State, job.Command)
			}
		}

		if len(jobs) == 0 {
			fmt.Fprintln(os.Stderr, "No background jobs.")
			return nil
		}

		if len(jobs) == 1 {
			fmt.Printf("%%%d\n", jobs[0].ID)
			return nil
		}

		choice, err := tui.Run(jobs)
		if err != nil {
			return fmt.Errorf("TUI failed: %w", err)
		}

		if choice != nil {
			fmt.Printf("%%%d\n", choice.ID)
		}

		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Prints the shell function to be evaluated in your shell config.",
	Long:  `Prints the shell function to be evaluated in your shell config.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(shellFunc)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.Flags().BoolP("verbose", "v", false, "Enable verbose debug output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
