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
	"os"
	"strings"

	en "github.com/adamkpickering/wj/internal/entry"
	"github.com/spf13/cobra"
	"slices"
)

var tag string

func init() {
	listCmd.AddCommand(listTasksCmd)
	listTasksCmd.Flags().StringVarP(&tag, "tag", "t", "", "filter by tags")
}

var listTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true

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
		if len(tag) == 0 {
			printTasksAsTable(tasks)
			return nil
		}

		// Filter the tasks
		filteredTasks := make([]en.Task, 0, len(tasks))
		for _, task := range tasks {
			if slices.Contains(task.Tags, tag) {
				filteredTasks = append(filteredTasks, task)
			}
		}
		printTasksAsTable(filteredTasks)
		return nil
	},
}
