package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/super-smooth/better-fg/internal/jobs"
	"github.com/super-smooth/better-fg/internal/tui"
)

const shellFunc = `
bfg() {
  # Get the output of the 'jobs' command from the CURRENT shell
  jobs_output=$(jobs)

  # If there are no jobs, do nothing
  if [ -z "$jobs_output" ]; then
    echo "No background jobs." >&2
    return
  fi

  # Run the Go program and get the selected job
  # Make sure the 'better-fg' binary is in your PATH
  selected_job=$(echo "$jobs_output" | command better-fg)

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

		jobs, err := jobs.Parse(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to parse jobs: %w", err)
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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
