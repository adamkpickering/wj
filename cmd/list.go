package cmd

import (
	"github.com/spf13/cobra"
)

var outputJson bool

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
}
