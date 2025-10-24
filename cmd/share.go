package cmd

import (
	"fmt"
	"log"

	"github.com/mayura-andrew/goodrive/internal/drive"
	"github.com/spf13/cobra"
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share [file_id]",
	Short: "Share a file on Google Drive",
	Long:  `This command allows you to share a file on Google Drive by specifying the file ID. You can also manage sharing permissions and generate shareable links.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := drive.InitDriveClient()
		if err != nil {
			log.Fatal(err)
		}

		fileID := args[0]
		email, _ := cmd.Flags().GetString("email")
		role, _ := cmd.Flags().GetString("role")

		// Share with specific email if provided
		if email != "" {
			err = client.ShareFileWithEmail(fileID, email, role)
		} else {
			// Otherwise, create public link
			err = client.ShareFile(fileID)
		}

		if err != nil {
			fmt.Printf("Error sharing file: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().StringP("email", "e", "", "Email address to share with")
	shareCmd.Flags().StringP("role", "r", "reader", "Role for the shared user (reader, writer, commenter)")
}
