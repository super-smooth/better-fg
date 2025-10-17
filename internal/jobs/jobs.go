package jobs

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Job represents a background job.
type Job struct {
	ID      int
	State   string
	Command string
}

// Parse parses the output of the "jobs" command and returns a slice of Job structs.
func Parse(reader io.Reader) ([]Job, error) {
	var jobs []Job
	scanner := bufio.NewScanner(reader)
	// Regex to capture various job line formats:
	// [1]  - suspended  nvim          (zsh with spaces)
	// [1]+  Stopped     vim file.txt  (bash format)
	// [1]+ Stopped vim               (bash without spaces)
	re := regexp.MustCompile(`^\[(\d+)\]\s*[+-]?\s*(\w+)\s+(.*)$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 {
			id, err := strconv.Atoi(matches[1])
			if err != nil {
				// Skip lines where the ID is not a valid integer
				continue
			}
			state := strings.TrimSpace(matches[2])
			command := strings.TrimSpace(matches[3])
			jobs = append(jobs, Job{ID: id, State: state, Command: command})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}
