package commands

import (
	"github.com/mayura-andrew/gooDrive/internal/drive"
	"github.com/spf13/cobra"
)

var (
	nameFilter string
	typeFilter string
)

func init() {
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(viewFilesCmd)

	viewFilesCmd.Flags().StringVarP(&nameFilter, "name", "n", "", "Filter by file name")
	viewFilesCmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Filter by file type")
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files in current Drive directory",
	Long:  `List all files in the current Google Drive directory context.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}
		return client.ListCurrentDirectory()
	},
}

var viewFilesCmd = &cobra.Command{
	Use:   "view-files",
	Short: "List files in Google Drive",
	Long: `List files in your Google Drive with optional filters.
	
Use flags to filter by name or file type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}

		return client.ViewFiles(nameFilter, typeFilter)
	},
}
