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
	"os"
	"strings"

	en "github.com/adamkpickering/wj/entry"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		workingDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		dirEntries, err := os.ReadDir(workingDir)
		if err != nil {
			return fmt.Errorf("failed to read working directory: %w", err)
		}
		tasks := []en.Task{}
		for _, dirEntry := range dirEntries {
			if !strings.HasSuffix(dirEntry.Name(), ".txt") {
				continue
			}
			entry := &en.Entry{}
			contents, err := os.ReadFile(dirEntry.Name())
			if err != nil {
				return fmt.Errorf("failed to read entry %q: %w", dirEntry.Name(), err)
			}
			if err := entry.UnmarshalText(contents); err != nil {
				return fmt.Errorf("failed to unmarshal entry %q: %w", dirEntry.Name(), err)
			}
			tasks = append(tasks, entry.Tasks...)
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.Encode(tasks)

		return nil
	},
}
