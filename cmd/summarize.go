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
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
	"time"
)

var taskLineRegex regexp.Regexp

func init() {
	taskLineRegex = *regexp.MustCompile(`^([0-9]{2}:[0-9]{2}) ([a-z0-9,]+) (.*)$`)
	rootCmd.AddCommand(summarizeCmd)
}

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Summarize a day of work",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName := args[0]
		rawContents, err := os.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", fileName, err)
		}
		contents := string(rawContents)

		entry := en.Entry{}
		task := en.Task{}
		taskContentLines := make([]string, 0, 100)
		lines := strings.Split(contents, "\n")
		toDoFound := false
		doneFound := false
		doneParsingDone := false
		for _, line := range lines {
			if !toDoFound && !doneFound && !doneParsingDone {
				// Looking for first To Do
				if line == "To Do" {
					toDoFound = true
				}
			} else if toDoFound && !doneFound && !doneParsingDone {
				// Parsing to do statements and skipping empty lines
				if line == "Done" {
					doneFound = true
					continue
				} else if strings.HasPrefix(line, "- ") {
					toDoText := strings.TrimPrefix(line, "- ")
					entry.ToDo = append(entry.ToDo, toDoText)
				}
			} else if toDoFound && doneFound && !doneParsingDone {
				// Parsing done statements and skipping empty lines
				if taskLineRegex.MatchString(line) {
					doneParsingDone = true

					task, err = partialTaskFromTitleLine(line)
					if err != nil {
						return fmt.Errorf("failed to parse first title line %q: %w\n", line, err)
					}
				} else if strings.HasPrefix(line, "- ") {
					doneText := strings.TrimPrefix(line, "- ")
					entry.Done = append(entry.Done, doneText)
				}
			} else if toDoFound && doneFound && doneParsingDone {
				if taskLineRegex.MatchString(line) {
					newTask, err := partialTaskFromTitleLine(line)
					if err != nil {
						return fmt.Errorf("failed to parse title line %q: %w\n", line, err)
					}

					// set .Content of previous task
					task.Content = strings.Join(taskContentLines, "\n")
					taskContentLines = make([]string, 0, 100)

					// calculate duration of previous task
					task.Duration = newTask.StartTime.Sub(task.StartTime)

					entry.Tasks = append(entry.Tasks, task)

					task = newTask
				} else {
					taskContentLines = append(taskContentLines, line)
				}
			}
		}

		// set .Content of last task and add it to list
		task.Content = strings.Join(taskContentLines, "\n")
		taskContentLines = make([]string, 0, 100)
		entry.Tasks = append(entry.Tasks, task)

		printStartEndDuration(entry)
		fmt.Printf("\n")
		printTimeByTaskTag(entry)
		fmt.Printf("\n")
		printSummary(entry)
		return nil
	},
}

// Constructs a Task and reads relevant fields from a Task title line
// into that Task. If the title line does not contain info on
// a given field of the Task, that field is left as its zero value.
func partialTaskFromTitleLine(line string) (en.Task, error) {
	task := en.Task{}

	result := taskLineRegex.FindStringSubmatch(line)
	if len(result) != 4 {
		return en.Task{}, fmt.Errorf("failed to parse line %q", line)
	}

	// parse start time
	rawStartTime := result[1]
	parsedTime, err := time.Parse("15:04", rawStartTime)
	if err != nil {
		return en.Task{}, fmt.Errorf("failed to parse time %q: %w", rawStartTime, err)
	}
	task.StartTime = parsedTime

	// parse tags
	rawTags := result[2]
	parts := strings.Split(rawTags, ",")
	for _, part := range parts {
		task.Tags = append(task.Tags, part)
	}

	task.Title = result[3]

	return task, nil
}

func printStartEndDuration(entry en.Entry) {
	startTime := entry.Tasks[0].StartTime.Format("15:04")
	endTime := entry.Tasks[len(entry.Tasks)-1].StartTime.Format("15:04")
	var totalTime time.Duration
	for _, task := range entry.Tasks {
		totalTime = totalTime + task.Duration
	}
	fmt.Printf("Started %s, ended %s (%s)\n", startTime, endTime, pretty(totalTime))
}

func printTimeByTaskTag(entry en.Entry) {
	tagTimes := map[string]time.Duration{}
	for _, task := range entry.Tasks {
		for _, tag := range task.Tags {
			if _, ok := tagTimes[tag]; !ok {
				tagTimes[tag] = time.Duration(task.Duration)
			} else {
				tagTimes[tag] = tagTimes[tag] + time.Duration(task.Duration)
			}
		}
	}

	for tag, duration := range tagTimes {
		fmt.Printf("%s\t\t%s\n", pretty(duration), tag)
	}
}

func printSummary(entry en.Entry) {
	for _, task := range entry.Tasks {
		tags := strings.Join(task.Tags, ",")
		fmt.Printf("%s\t\t%s\t\t%s\n", pretty(task.Duration), tags, task.Title)
	}
}

func pretty(duration time.Duration) string {
	hours := duration / time.Hour
	minutes := (duration - hours*time.Hour) / time.Minute
	if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
