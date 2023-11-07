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
	en "github.com/adamkpickering/wj/internal/entry"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

func init() {
	rootCmd.AddCommand(summarizeCmd)
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize a day of work",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		fileName := args[0]
		rawContents, err := os.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", fileName, err)
		}

		entry := &en.Entry{}
		if err := entry.UnmarshalText(rawContents); err != nil {
			return fmt.Errorf("failed to parse entry %q: %w", fileName, err)
		}

		startEndDurationErr := printStartEndDuration(entry)
		fmt.Printf("\n")
		taskTimeTotalsErr := printTaskTimeTotalsTable(entry.Tasks)
		fmt.Printf("\n")
		tasksAsTableErr := printTasksAsTable(entry.Tasks)
		return errors.Join(startEndDurationErr, taskTimeTotalsErr, tasksAsTableErr)
	},
}

func printStartEndDuration(entry *en.Entry) error {
	startTime := entry.Tasks[0].StartTime.Format("15:04")
	endTime := entry.Tasks[len(entry.Tasks)-1].StartTime.Format("15:04")
	var totalTime time.Duration
	for _, task := range entry.Tasks {
		totalTime = totalTime + time.Duration(task.Duration)
	}
	_, err := fmt.Printf("Started %s, ended %s (%s)\n", startTime, endTime, pretty(totalTime))
	return err
}

func printTaskTimeTotalsTable(tasks []en.Task) error {
	tagTimes := map[string]time.Duration{}
	for _, task := range tasks {
		for _, tag := range task.Tags {
			if _, ok := tagTimes[tag]; !ok {
				tagTimes[tag] = time.Duration(task.Duration)
			} else {
				tagTimes[tag] = tagTimes[tag] + time.Duration(task.Duration)
			}
		}
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	if _, err := fmt.Fprintf(writer, "Total Duration\tTag\n"); err != nil {
		return fmt.Errorf("failed to write table header: %w", err)
	}
	for tag, duration := range tagTimes {
		if _, err := fmt.Fprintf(writer, "%s\t%s\n", pretty(duration), tag); err != nil {
			return fmt.Errorf("failed to write table row: %w", err)
		}
	}
	writer.Flush()
	return nil
}

func pretty(duration time.Duration) string {
	hours := duration / time.Hour
	minutes := (duration - hours*time.Hour) / time.Minute
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
