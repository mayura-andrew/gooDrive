package cmd

import (
	"github.com/spf13/cobra"
    "github.com/mayura-andrew/goodrive/internal/sync"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize files between local storage and Google Drive",
	Long: `The sync command allows you to synchronize files between your local storage and Google Drive.
It supports bidirectional sync, conflict resolution, and selective sync options.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the sync function from the internal sync package
		sync.Synchronize()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Define flags for the sync command
	syncCmd.Flags().BoolP("watch", "w", false, "Enable watch mode for real-time synchronization")
	syncCmd.Flags().StringP("folder", "f", "", "Specify the folder to synchronize")
	syncCmd.Flags().BoolP("selective", "s", false, "Enable selective sync for specific files/folders")
	syncCmd.Flags().BoolP("force", "F", false, "Force synchronization even if there are conflicts")
}
