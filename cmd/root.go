package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var dataDirectory string

var rootCmd = &cobra.Command{
	Use:           "wj",
	Short:         "Work with work journals",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	defaultDataDirectory, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("failed to get working directory: %w", err))
	}
	rootCmd.PersistentFlags().StringVarP(&dataDirectory, "data-directory", "d", defaultDataDirectory, "path to the data directory")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
