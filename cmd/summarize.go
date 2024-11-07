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
	"path/filepath"
	"text/tabwriter"
	"time"
)

func init() {
	rootCmd.AddCommand(summarizeCmd)
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize [<date>]",
	Short: "Summarize an entry",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true

		fileName := ""
		switch len(args) {
		case 0:
			fileName = time.Now().Format(journalFileFormat)
		case 1:
			fileName = args[0] + ".wj"
		}
		filePath := filepath.Join(dataDirectory, fileName)
		rawContents, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", fileName, err)
		}

		entry := &en.Entry{}
		if err := entry.UnmarshalText(rawContents); err != nil {
			return fmt.Errorf("failed to parse entry %q: %w", fileName, err)
		}

		fmt.Printf("Summary of %s\n\n", entry.Date.Format("Mon Jan 2, 2006"))
		startEndDurationErr := printStartEndDuration(entry)
		fmt.Printf("\n")
		taskTimeTotalsErr := printTaskTimeTotalsTable(entry.Tasks)
		fmt.Printf("\n")
		tasksAsTableErr := printTasksAsTable(entry.Tasks)
		return errors.Join(startEndDurationErr, taskTimeTotalsErr, tasksAsTableErr)
	},
}

func getStartEndDuration(entry *en.Entry) (time.Time, time.Time, time.Duration) {
	startTime := entry.Tasks[0].StartTime
	endTime := entry.Tasks[len(entry.Tasks)-1].StartTime
	duration := endTime.Sub(startTime)
	return startTime, endTime, duration
}

func printStartEndDuration(entry *en.Entry) error {
	startTime, endTime, duration := getStartEndDuration(entry)
	_, err := fmt.Printf("Started %s, ended %s (%s)\n",
		startTime.Format("15:04"),
		endTime.Format("15:04"),
		pretty(duration),
	)
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
	if _, err := fmt.Fprintf(writer, "Tag\t Time\n"); err != nil {
		return fmt.Errorf("failed to write table header: %w", err)
	}
	for tag, duration := range tagTimes {
		if _, err := fmt.Fprintf(writer, "%s\t%s\n", tag, pretty(duration)); err != nil {
			return fmt.Errorf("failed to write table row: %w", err)
		}
	}
	writer.Flush()
	return nil
}

// pretty converts a duration into the format "XXhYYm". The output is
// always 6 characters wide. Panics if the duration is long enough that
// this cannot be the case.
func pretty(duration time.Duration) string {
	hours := duration / time.Hour
	minutes := (duration - hours*time.Hour) / time.Minute
	if hours > 99 {
		panic(fmt.Sprintf(`hour count "%d" has too many decimal places`, hours))
	}
	switch {
	case hours == 0 && minutes < 10:
		return fmt.Sprintf("    %dm", minutes)
	case hours == 0 && minutes >= 10:
		return fmt.Sprintf("   %dm", minutes)
	case hours < 10 && minutes < 10:
		return fmt.Sprintf(" %dh %dm", hours, minutes)
	case hours < 10 && minutes >= 10:
		return fmt.Sprintf(" %dh%dm", hours, minutes)
	case hours >= 10 && minutes < 10:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	case hours >= 10 && minutes >= 10:
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	panic(fmt.Sprintf("unhandled case for %v", duration))
}
