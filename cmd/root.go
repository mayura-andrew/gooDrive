package cmd

import (
	"github.com/mayura-andrew/goodrive/internal/config"
	"github.com/mayura-andrew/goodrive/internal/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gooDrive",
	Short: "A CLI tool for managing Google Drive",
	Long:  `gooDrive is a command-line interface for interacting with Google Drive, allowing users to upload, download, list, sync, search, and share files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommands are provided
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.LogError(err)
		return
	}
}

func init() {
	// Initialize configuration
	// Load config if exists, otherwise use defaults
	_, err := config.LoadConfig("config.json")
	if err != nil {
		// Config file not found, continue with defaults
	}

	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(shareCmd)
	rootCmd.AddCommand(versionCmd)
}
