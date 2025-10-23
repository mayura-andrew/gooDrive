package commands

import (
	"github.com/mayura-andrew/gooDrive/internal/drive"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone [file-id or sharing-link]",
	Short: "Download file or folder from Google Drive",
	Long: `Clone downloads a file or folder from Google Drive to your local machine.
	
You can provide either a file ID or a sharing link. The command will automatically
detect the type and download accordingly.

Examples:
  drive clone 1abc123def456
  drive clone https://drive.google.com/file/d/1abc123def456/view`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}

		fileID := drive.ExtractFileID(args[0])
		return client.Clone(fileID)
	},
}
