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
	"strings"
	"time"

	en "github.com/adamkpickering/wj/internal/entry"
	"github.com/spf13/cobra"
)

const journalFileFormat = "2006-01-02.txt"

var ErrNoLastEntry = errors.New("failed to find last entry")

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a journal file for today",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		fileName := now.Format(journalFileFormat)
		if _, err := os.Stat(fileName); !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("the entry %q already exists", fileName)
		}

		var toDoItems []string
		lastEntry, err := getLastEntry(now)
		if err == nil {
			toDoItems = lastEntry.ToDo
		} else if errors.Is(err, ErrNoLastEntry) {
			toDoItems = []string{}
		} else {
			return fmt.Errorf("failed to get to do items: %w", err)
		}

		lines := []string{"To Do"}
		for _, toDoItem := range toDoItems {
			lines = append(lines, fmt.Sprintf("- %s", toDoItem))
		}
		lines = append(lines, "")
		lines = append(lines, "Done")
		lines = append(lines, "")
		contents := strings.Join(lines, "\n")

		err = os.WriteFile(fileName, []byte(contents), 0o644)
		if err != nil {
			return fmt.Errorf("failed to create file %q: %w", fileName, err)
		}

		return nil
	},
}

func getLastEntry(today time.Time) (en.Entry, error) {
	var (
		contents []byte
		err      error
	)
	for i := 1; i < 15; i++ {
		testDate := today.AddDate(0, 0, -i)
		testFileName := testDate.Format(journalFileFormat)
		contents, err = os.ReadFile(testFileName)
		if err == nil {
			entry, err := parseEntry(contents)
			if err != nil {
				return en.Entry{}, fmt.Errorf("failed to parse entry: %w", err)
			}
			return entry, nil
		} else if errors.Is(err, os.ErrNotExist) {
			continue
		} else {
			return en.Entry{}, fmt.Errorf("failed to open file %s: %w", testFileName, err)
		}
	}
	return en.Entry{}, ErrNoLastEntry
}
