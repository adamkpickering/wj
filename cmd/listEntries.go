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
	"path/filepath"
	"strings"
	"text/tabwriter"

	en "github.com/adamkpickering/wj/internal/entry"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listEntriesCmd)
}

var listEntriesCmd = &cobra.Command{
	Use:   "entries",
	Short: "List entries",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		entries, err := readEntries(dataDirectory)
		if err != nil {
			return fmt.Errorf("failed to read entries: %w", err)
		}
		return printEntriesAsTable(entries)
	},
}

func readEntries(dataDir string) ([]*en.Entry, error) {
	dirEntries, err := os.ReadDir(dataDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to read current directory: %w", err)
	}
	entries := make([]*en.Entry, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		fileName := filepath.Join(dataDirectory, dirEntry.Name())
		if !strings.HasSuffix(fileName, ".txt") {
			continue
		}
		contents, err := os.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to read entry %q: %w", fileName, err)
		}
		entry := &en.Entry{}
		if err := entry.UnmarshalText(contents); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entry %q: %w", fileName, err)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func printEntriesAsTable(entries []*en.Entry) error {
	writer := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	if _, err := fmt.Fprintf(writer, "Date\tStart Time\tEnd Time\tDuration\tTask Count\n"); err != nil {
		return fmt.Errorf("failed to write table header: %w", err)
	}
	for _, entry := range entries {
		prettyDate := entry.Date.Format("Mon Jan 02 2006")
		startTime, endTime, duration := getStartEndDuration(entry)
		_, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%d\n",
			prettyDate,
			startTime.Format("15:04"),
			endTime.Format("15:04"),
			pretty(duration),
			len(entry.Tasks),
		)
		if err != nil {
			return fmt.Errorf("failed to write table row: %w", err)
		}
	}
	writer.Flush()
	return nil
}
