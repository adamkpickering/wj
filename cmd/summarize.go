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
	"encoding/json"
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
	taskLineRegex = *regexp.MustCompile(`^[0-9]{2}:[0-9]{2} .*$`)
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
					task.Duration = en.JSONStringDuration(newTask.StartTime.Sub(task.StartTime))

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

		printSummary(entry)
		return nil
	},
}

// Constructs a Task and reads relevant fields from a Task title line
// into that Task. If the title line does not contain info on
// a given field of the Task, that field is left as its zero value.
func partialTaskFromTitleLine(line string) (en.Task, error) {
	task := en.Task{}

	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return en.Task{}, fmt.Errorf("failed to split line %q", line)
	}

	rawStartTime := parts[0]
	parsedTime, err := time.Parse("15:04", rawStartTime)
	if err != nil {
		return en.Task{}, fmt.Errorf("failed to parse time %q: %w", rawStartTime, err)
	}
	task.StartTime = parsedTime

	task.Description = parts[1]

	return task, nil
}

func printEntry(entry en.Entry) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(entry)
	if err != nil {
		panic(err)
	}
}

func printSummary(entry en.Entry) {
	for _, task := range entry.Tasks {
		fmt.Printf("%s\t\t%s\n", task.Duration.Pretty(), task.Description)
	}
}
