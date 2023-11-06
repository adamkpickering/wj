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
	"fmt"
	en "github.com/adamkpickering/wj/internal/entry"
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"os"
	"strings"
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

		printStartEndDuration(entry)
		fmt.Printf("\n")
		printTaskTimeTotalsTable(entry.Tasks)
		fmt.Printf("\n")
		printTasksAsTable(entry.Tasks)
		return nil
	},
}

func printStartEndDuration(entry *en.Entry) {
	startTime := entry.Tasks[0].StartTime.Format("15:04")
	endTime := entry.Tasks[len(entry.Tasks)-1].StartTime.Format("15:04")
	var totalTime time.Duration
	for _, task := range entry.Tasks {
		totalTime = totalTime + time.Duration(task.Duration)
	}
	fmt.Printf("Started %s, ended %s (%s)\n", startTime, endTime, pretty(totalTime))
}

func printTaskTimeTotalsTable(tasks []en.Task) {
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

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactClassic)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Total Duration"},
			{Text: "Tag"},
		},
	}
	for tag, duration := range tagTimes {
		row := []*simpletable.Cell{
			{Text: pretty(duration)},
			{Text: tag},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}
	fmt.Println(table.String())
}

func printTasksAsTable(tasks []en.Task) {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactClassic)
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Duration"},
			{Text: "Tags"},
			{Text: "Title"},
		},
	}
	for _, task := range tasks {
		tags := strings.Join(task.Tags, ",")
		row := []*simpletable.Cell{
			{Text: pretty(time.Duration(task.Duration))},
			{Text: tags},
			{Text: task.Title},
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}
	fmt.Println(table.String())
}

func pretty(duration time.Duration) string {
	hours := duration / time.Hour
	minutes := (duration - hours*time.Hour) / time.Minute
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
