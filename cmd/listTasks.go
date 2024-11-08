package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
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
	listTasksCmd.Flags().BoolVarP(&outputJson, "json", "j", false, "output JSON")

	dateDurationRegexp = regexp.MustCompile("^([0-9]+)([dw])$")
}

var listTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cutoffTime time.Time
		if last != "" {
			parsedDateDuration, err := parseDateDuration(last)
			if err != nil {
				return fmt.Errorf("failed to parse --last value %q as time.Duration: %w", last, err)
			}
			cutoffTime = time.Now().Add(-parsedDateDuration)
		}

		// Compile a list of tasks
		entries, err := readEntries(dataDirectory)
		if err != nil {
			return fmt.Errorf("failed to read entries: %w", err)
		}
		taskCount := 0
		for _, entry := range entries {
			taskCount += len(entry.Tasks)
		}
		tasks := make([]en.Task, 0, taskCount)
		for _, entry := range entries {
			tasks = append(tasks, entry.Tasks...)
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

		if outputJson {
			return outputTasksAsJson(filteredTasks)
		}
		return printTasksAsTable(filteredTasks)
	},
}

func outputTasksAsJson(tasks []en.Task) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tasks); err != nil {
		return fmt.Errorf("failed to encode tasks: %w", err)
	}
	return nil
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

func printTasksAsTable(tasks []en.Task) error {
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	if _, err := fmt.Fprintf(writer, "Date\tStart Time\tDuration\tTags\tTitle\n"); err != nil {
		return fmt.Errorf("failed to write table header: %w", err)
	}
	for _, task := range tasks {
		date := task.StartTime.Format("Mon Jan 02 2006")
		startTime := task.StartTime.Format("15:04")
		duration := time.Duration(task.Duration)
		tags := strings.Join(task.Tags, ",")
		if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", date, startTime, pretty(duration), tags, task.Title); err != nil {
			return fmt.Errorf("failed to write table row: %w", err)
		}
	}
	writer.Flush()
	return nil
}
