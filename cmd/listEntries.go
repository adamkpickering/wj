package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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
		if filepath.Ext(fileName) != ".wj" {
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
