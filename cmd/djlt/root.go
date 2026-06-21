package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	xmlPath string
)

var rootCmd = &cobra.Command{
	Use:   "djlt",
	Short: "DJ Library Tools CLI",
	Long:  `A comprehensive CLI tool for managing DJ libraries and Rekordbox XMLs.`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&xmlPath, "xml", "x", "", "Path to the Rekordbox XML library")
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
