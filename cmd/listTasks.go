/*
Copyright © 2021 ADAM PICKERING

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"slices"

	en "github.com/adamkpickering/wj/internal/entry"
	"github.com/spf13/cobra"
)

var tag string
var last string

var dateDurationRegexp *regexp.Regexp

func init() {
	listCmd.AddCommand(listTasksCmd)
	listTasksCmd.Flags().StringVarP(&tag, "tag", "t", "", "filter by tags")
	listTasksCmd.Flags().StringVarP(&last, "last", "l", "", "only list tags from the last ex. 7d")

	dateDurationRegexp = regexp.MustCompile("^([0-9]+)([dw])$")
}

var listTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true

		var cutoffTime time.Time
		if last != "" {
			parsedDateDuration, err := parseDateDuration(last)
			if err != nil {
				return fmt.Errorf("failed to parse --last value %q as time.Duration: %w", last, err)
			}
			cutoffTime = time.Now().Add(-parsedDateDuration)
		}

		// Read all entries
		dirEntries, err := os.ReadDir(".")
		if err != nil {
			return fmt.Errorf("failed to read current directory: %w", err)
		}
		entries := make([]en.Entry, 0, len(dirEntries))
		for _, dirEntry := range dirEntries {
			fileName := dirEntry.Name()
			if !strings.HasSuffix(fileName, ".txt") {
				continue
			}
			contents, err := os.ReadFile(fileName)
			if err != nil {
				return fmt.Errorf("failed to read entry %q: %w", fileName, err)
			}
			entry := en.Entry{}
			if err := entry.UnmarshalText(contents); err != nil {
				return fmt.Errorf("failed to unmarshal entry %q: %w", fileName, err)
			}
			entries = append(entries, entry)
		}

		// Compile a list of tasks
		taskCount := 0
		for _, entry := range entries {
			taskCount += len(entry.Tasks)
		}
		tasks := make([]en.Task, 0, taskCount)
		for _, entry := range entries {
			tasks = append(tasks, entry.Tasks...)
		}
		if tag == "" && last == "" {
			printTasksAsTable(tasks)
			return nil
		}

		// Filter the tasks
		filteredTasks := make([]en.Task, 0, len(tasks))
		for _, task := range tasks {
			if tag != "" && !slices.Contains(task.Tags, tag) {
				continue
			} else if !cutoffTime.IsZero() && task.StartTime.Before(cutoffTime) {
				continue
			}
			filteredTasks = append(filteredTasks, task)
		}
		printTasksAsTable(filteredTasks)
		return nil
	},
}

// time.ParseDuration only deals with hours and below. We need to
// deal with days and weeks.
func parseDateDuration(rawDuration string) (time.Duration, error) {
	submatches := dateDurationRegexp.FindStringSubmatch(rawDuration)
	if len(submatches) == 0 {
		return 0, errors.New("no match for date duration regex")
	}
	count, err := strconv.ParseInt(submatches[1], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse count: %w", err)
	}
	var multiplier time.Duration
	switch submatches[2] {
	case "d":
		multiplier = 24 * time.Hour
	case "w":
		multiplier = 24 * 7 * time.Hour
	default:
		return 0, fmt.Errorf("unknown unit %q", submatches[2])
	}
	return time.Duration(count) * multiplier, nil
}
