package commands

import (
	"github.com/spf13/cobra"
)

var (
	verbose bool
	rootCmd = &cobra.Command{
		Use:   "drive",
		Short: "A CLI tool for Google Drive operations",
		Long: `gooDrive - A powerful command-line interface for Google Drive
		
Access, sync, download, and upload files to Google Drive directly from your terminal.
Perfect for developers and power users who prefer command-line workflows.`,
		Version: "1.0.0",
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}
package commands
