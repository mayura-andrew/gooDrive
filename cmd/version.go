package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var version = "v0.1.0-beta.1" // Define the version of the CLI tool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
    Use:     "version",
    Short:   "Display the version of the CLI tool",
    Long:    `This command shows the current version of the Google Drive CLI tool.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Google Drive CLI Tool version: %s\n", version)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd) // Add the version command to the root command
}