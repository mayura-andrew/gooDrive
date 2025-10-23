package main

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/spf13/cobra"
)

//go:embed auth_template.html
var authTemplate embed.FS

const (
	tokenFile       = ".drive-cli-token.json"
	credentialsFile = "oauth.json"
	metaFile        = ".drive-cli-meta.json"
	authPort        = 8080 // Fixed port for OAuth callback
)

type FileMeta struct {
	LocalPath string            `json:"local_path"`
	DriveID   string            `json:"drive_id"`
	DriveName string            `json:"drive_name"`
	ModTime   string            `json:"mod_time"`
	IsTracked bool              `json:"is_tracked"`
	Children  map[string]string `json:"children,omitempty"` // name -> driveID
}

type DriveClient struct {
	service *drive.Service
	ctx     context.Context
}

var (
	authCode     string
	authCodeOnce sync.Once
	authServer   *http.Server
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "drive",
		Short: "A CLI tool for Google Drive operations",
		Long:  "Drive-CLI: Access, sync, download, upload files to Google Drive from command line",
	}

	rootCmd.AddCommand(
		cloneCmd(),
		addRemoteCmd(),
		viewFilesCmd(),
		lsCmd(),
		statusCmd(),
		pullCmd(),
		pushCmd(),
		rmCmd(),
		catCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Initialize Drive client with OAuth2
func initDriveClient() (*DriveClient, error) {
	ctx := context.Background()

	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %v", err)
	}

	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive service: %v", err)
	}

	return &DriveClient{service: srv, ctx: ctx}, nil
}

// Get OAuth2 client
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Use fixed port
	port := authPort

	// Update redirect URL
	config.RedirectURL = fmt.Sprintf("http://localhost:%d/callback", port)

	// Setup local HTTP server
	mux := http.NewServeMux()

	// Serve the auth template
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := authTemplate.ReadFile("auth_template.html")
		if err != nil {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(tmpl)
	})

	// Handle OAuth callback
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success": false, "error": "No authorization code received"}`))
			return
		}

		authCodeOnce.Do(func() {
			authCode = code
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))

		// Shutdown server after successful auth
		go func() {
			time.Sleep(2 * time.Second)
			if authServer != nil {
				authServer.Shutdown(context.Background())
			}
		}()
	})

	authServer = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: mux,
	}

	// Start server in background
	go func() {
		if err := authServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Generate auth URL
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║          Google Drive CLI - Authentication Required           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("🌐 Opening browser for authentication...\n\n")
	fmt.Printf("If your browser doesn't open automatically, visit:\n%s\n\n", authURL)
	fmt.Printf("🔗 Local server started at: http://localhost:%d\n", port)
	fmt.Println("\n⏳ Waiting for authentication...")

	// Try to open browser automatically
	openBrowser(authURL)

	// Wait for auth code or timeout
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Fatal("Authentication timeout. Please try again.")
		case <-ticker.C:
			if authCode != "" {
				fmt.Println("\n✅ Authentication successful!")

				tok, err := config.Exchange(context.TODO(), authCode)
				if err != nil {
					log.Fatalf("Unable to retrieve token: %v", err)
				}

				// Cleanup
				if authServer != nil {
					authServer.Shutdown(context.Background())
				}

				return tok
			}
		}
	}
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Clone command - download file/folder from Drive
func cloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone [file-id or sharing-link]",
		Short: "Download file or folder from Google Drive",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}

			fileID := extractFileID(args[0])
			if err := client.cloneFile(fileID); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func extractFileID(input string) string {
	// Extract file ID from sharing link or return as-is
	if strings.Contains(input, "drive.google.com") {
		parts := strings.Split(input, "/")
		for i, part := range parts {
			if part == "d" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}
	return input
}

func (dc *DriveClient) cloneFile(fileID string) error {
	file, err := dc.service.Files.Get(fileID).SupportsAllDrives(true).Fields("id, name, mimeType, size").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve file: %v", err)
	}

	fmt.Printf("Preparing: %s for download\n", file.Name)

	if file.MimeType == "application/vnd.google-apps.folder" {
		return dc.cloneFolder(fileID, file.Name)
	}

	return dc.downloadFile(file)
}

func (dc *DriveClient) cloneFolder(folderID, folderName string) error {
	if err := os.MkdirAll(folderName, 0755); err != nil {
		return err
	}

	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	fileList, err := dc.service.Files.List().Q(query).SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).Fields("files(id, name, mimeType)").Do()
	if err != nil {
		return err
	}

	for _, file := range fileList.Files {
		localPath := filepath.Join(folderName, file.Name)
		if file.MimeType == "application/vnd.google-apps.folder" {
			dc.cloneFolder(file.Id, localPath)
		} else {
			dc.downloadFile(file)
		}
	}

	// Save metadata for tracking
	saveMeta(folderName, folderID, folderName, true)
	return nil
}

func (dc *DriveClient) downloadFile(file *drive.File) error {
	var resp *http.Response
	var err error

	// Handle Google Workspace files (Docs, Sheets, etc.)
	if strings.HasPrefix(file.MimeType, "application/vnd.google-apps") {
		exportType := selectExportType(file.MimeType)
		resp, err = dc.service.Files.Export(file.Id, exportType).Download()
	} else {
		resp, err = dc.service.Files.Get(file.Id).Download()
	}

	if err != nil {
		return fmt.Errorf("unable to download file: %v", err)
	}
	defer resp.Body.Close()

	fileName := file.Name
	if strings.HasPrefix(file.MimeType, "application/vnd.google-apps") {
		fileName += getExtension(file.MimeType)
	}

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("downloading file [####################################] 100%%\n")
	fmt.Printf("completed download of %s (%d bytes)\n", fileName, written)

	// Save metadata
	saveMeta(fileName, file.Id, file.Name, true)
	return nil
}

func selectExportType(mimeType string) string {
	// Default export types for Google Workspace files
	exportMap := map[string]string{
		"application/vnd.google-apps.document":     "application/pdf",
		"application/vnd.google-apps.spreadsheet":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.google-apps.presentation": "application/pdf",
		"application/vnd.google-apps.drawing":      "application/pdf",
	}

	if export, ok := exportMap[mimeType]; ok {
		return export
	}
	return "application/pdf"
}

func getExtension(mimeType string) string {
	extMap := map[string]string{
		"application/vnd.google-apps.document":     ".pdf",
		"application/vnd.google-apps.spreadsheet":  ".xlsx",
		"application/vnd.google-apps.presentation": ".pdf",
	}

	if ext, ok := extMap[mimeType]; ok {
		return ext
	}
	return ".pdf"
}

// Add-remote command - upload local file to Drive
func addRemoteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-remote [local-file-path]",
		Short: "Upload local file to Google Drive and track it",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}

			if err := client.addRemote(args[0]); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func (dc *DriveClient) addRemote(localPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := filepath.Base(localPath)
	driveFile := &drive.File{Name: fileName}

	fmt.Printf("Uploading %s...\n", fileName)

	uploadedFile, err := dc.service.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		return fmt.Errorf("unable to upload file: %v", err)
	}

	fmt.Printf("Successfully uploaded: %s (ID: %s)\n", uploadedFile.Name, uploadedFile.Id)

	// Save metadata
	saveMeta(localPath, uploadedFile.Id, uploadedFile.Name, true)
	return nil
}

// View-files command - list files in Drive
func viewFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view-files",
		Short: "List files in Google Drive",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}

			nameFilter, _ := cmd.Flags().GetString("name")
			typeFilter, _ := cmd.Flags().GetString("type")

			if err := client.viewFiles(nameFilter, typeFilter); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringP("name", "n", "", "Filter by file name")
	cmd.Flags().StringP("type", "t", "", "Filter by file type")
	return cmd
}

func (dc *DriveClient) viewFiles(nameFilter, typeFilter string) error {
	query := "trashed=false"

	if nameFilter != "" {
		query += fmt.Sprintf(" and name contains '%s'", nameFilter)
	}
	if typeFilter != "" {
		query += fmt.Sprintf(" and mimeType contains '%s'", typeFilter)
	}

	fileList, err := dc.service.Files.List().Q(query).
		PageSize(50).Fields("files(id, name, mimeType, size, modifiedTime)").Do()
	if err != nil {
		return err
	}

	fmt.Printf("%-50s %-15s %-15s\n", "Name", "Type", "Size")
	fmt.Println(strings.Repeat("-", 80))

	for _, file := range fileList.Files {
		fileType := filepath.Ext(file.Name)
		if fileType == "" {
			fileType = "folder"
		}
		size := fmt.Sprintf("%d bytes", file.Size)
		fmt.Printf("%-50s %-15s %-15s\n", file.Name, fileType, size)
	}

	return nil
}

// Helper functions for metadata
func saveMeta(localPath, driveID, driveName string, isTracked bool) error {
	info, _ := os.Stat(localPath)
	modTime := ""
	if info != nil {
		modTime = info.ModTime().String()
	}

	meta := FileMeta{
		LocalPath: localPath,
		DriveID:   driveID,
		DriveName: driveName,
		ModTime:   modTime,
		IsTracked: isTracked,
	}

	metaPath := localPath + metaFile
	file, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(meta)
}

func loadMeta(localPath string) (*FileMeta, error) {
	metaPath := localPath + metaFile
	file, err := os.Open(metaPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var meta FileMeta
	if err := json.NewDecoder(file).Decode(&meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// Placeholder commands (implement as needed)
func lsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List files in current Drive directory",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}
			if err := client.listCurrentDirectory(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show status of tracked files",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}
			if err := client.checkStatus(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func pullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pull latest changes from Drive",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}
			if err := client.pullChanges(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func pushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Push local changes to Drive",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}
			if err := client.pushChanges(); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func rmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm [file-id]",
		Short: "Remove file from Drive",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}
			if err := client.removeFile(args[0]); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func catCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cat [file-id]",
		Short: "View file contents without downloading",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := initDriveClient()
			if err != nil {
				log.Fatal(err)
			}
			if err := client.viewFileContents(args[0]); err != nil {
				log.Fatal(err)
			}
		},
	}
}

// LS Command - List files in current Drive directory
func (dc *DriveClient) listCurrentDirectory() error {
	cwd, _ := os.Getwd()
	meta, err := loadMeta(cwd)

	var parentID string
	if err != nil || meta == nil {
		// List root directory
		parentID = "root"
	} else {
		parentID = meta.DriveID
	}

	query := fmt.Sprintf("'%s' in parents and trashed=false", parentID)
	fileList, err := dc.service.Files.List().Q(query).
		PageSize(100).
		Fields("files(id, name, mimeType, size, modifiedTime, owners)").
		OrderBy("folder,name").Do()
	if err != nil {
		return fmt.Errorf("unable to list files: %v", err)
	}

	fmt.Printf("\nFiles in current directory:\n")
	fmt.Printf("%-50s %-12s %-15s %-20s\n", "Name", "Type", "Size", "Modified")
	fmt.Println(strings.Repeat("-", 100))

	for _, file := range fileList.Files {
		fileType := "file"
		size := fmt.Sprintf("%d B", file.Size)

		if file.MimeType == "application/vnd.google-apps.folder" {
			fileType = "folder"
			size = "-"
		} else if strings.HasPrefix(file.MimeType, "application/vnd.google-apps") {
			fileType = strings.TrimPrefix(file.MimeType, "application/vnd.google-apps.")
			size = "-"
		}

		modTime := ""
		if file.ModifiedTime != "" {
			t, _ := time.Parse(time.RFC3339, file.ModifiedTime)
			modTime = t.Format("2006-01-02 15:04")
		}

		fmt.Printf("%-50s %-12s %-15s %-20s\n",
			truncate(file.Name, 50), fileType, size, modTime)
	}

	fmt.Printf("\nTotal: %d files\n", len(fileList.Files))
	return nil
}

// STATUS Command - Show status of tracked files
func (dc *DriveClient) checkStatus() error {
	trackedFiles := make(map[string]*FileMeta)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// Check if this file has metadata
		meta, err := loadMeta(path)
		if err == nil && meta.IsTracked {
			trackedFiles[path] = meta
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(trackedFiles) == 0 {
		fmt.Println("No tracked files found in current directory")
		return nil
	}

	fmt.Println("\nTracked files status:")
	fmt.Printf("%-50s %-15s\n", "File", "Status")
	fmt.Println(strings.Repeat("-", 70))

	for path, meta := range trackedFiles {
		status := "up-to-date"

		// Check if local file was modified
		info, err := os.Stat(path)
		if err != nil {
			status = "deleted locally"
		} else {
			localModTime := info.ModTime()
			savedModTime, _ := time.Parse(time.RFC3339, meta.ModTime)

			if localModTime.After(savedModTime) {
				status = "modified locally"
			}
		}

		// Check if file still exists on Drive
		_, err = dc.service.Files.Get(meta.DriveID).Fields("id").Do()
		if err != nil {
			status = "deleted on Drive"
		}

		fmt.Printf("%-50s %-15s\n", truncate(path, 50), status)
	}

	return nil
}

// PULL Command - Download latest changes from Drive
func (dc *DriveClient) pullChanges() error {
	// Find all tracked files
	trackedFiles := make(map[string]*FileMeta)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		meta, err := loadMeta(path)
		if err == nil && meta.IsTracked {
			trackedFiles[path] = meta
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(trackedFiles) == 0 {
		fmt.Println("No tracked files to pull")
		return nil
	}

	fmt.Printf("Pulling changes for %d tracked files...\n", len(trackedFiles))

	for path, meta := range trackedFiles {
		// Get file info from Drive
		driveFile, err := dc.service.Files.Get(meta.DriveID).
			Fields("id, name, modifiedTime, mimeType, size").Do()
		if err != nil {
			fmt.Printf("⚠ Skipping %s: file not found on Drive\n", path)
			continue
		}

		// Check if Drive file is newer
		driveModTime, _ := time.Parse(time.RFC3339, driveFile.ModifiedTime)
		localInfo, _ := os.Stat(path)

		if localInfo != nil {
			localModTime := localInfo.ModTime()
			if !driveModTime.After(localModTime) {
				fmt.Printf("✓ %s: up-to-date\n", path)
				continue
			}
		}

		// Download the file
		fmt.Printf("⬇ Downloading %s...\n", path)
		if err := dc.downloadFileToPath(driveFile, path); err != nil {
			fmt.Printf("✗ Error downloading %s: %v\n", path, err)
		} else {
			fmt.Printf("✓ Updated %s\n", path)
			// Update metadata
			saveMeta(path, driveFile.Id, driveFile.Name, true)
		}
	}

	return nil
}

func (dc *DriveClient) downloadFileToPath(file *drive.File, localPath string) error {
	var resp *http.Response
	var err error

	if strings.HasPrefix(file.MimeType, "application/vnd.google-apps") {
		exportType := selectExportType(file.MimeType)
		resp, err = dc.service.Files.Export(file.Id, exportType).Download()
	} else {
		resp, err = dc.service.Files.Get(file.Id).Download()
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// PUSH Command - Upload local changes to Drive
func (dc *DriveClient) pushChanges() error {
	trackedFiles := make(map[string]*FileMeta)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		meta, err := loadMeta(path)
		if err == nil && meta.IsTracked {
			trackedFiles[path] = meta
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(trackedFiles) == 0 {
		fmt.Println("No tracked files to push")
		return nil
	}

	fmt.Printf("Pushing changes for %d tracked files...\n", len(trackedFiles))

	for path, meta := range trackedFiles {
		// Check if local file was modified
		localInfo, err := os.Stat(path)
		if err != nil {
			fmt.Printf("⚠ Skipping %s: file not found locally\n", path)
			continue
		}

		localModTime := localInfo.ModTime()
		savedModTime, _ := time.Parse(time.RFC3339, meta.ModTime)

		if !localModTime.After(savedModTime) {
			fmt.Printf("✓ %s: no changes\n", path)
			continue
		}

		// Upload the file
		fmt.Printf("⬆ Uploading %s...\n", path)
		if err := dc.updateDriveFile(meta.DriveID, path); err != nil {
			fmt.Printf("✗ Error uploading %s: %v\n", path, err)
		} else {
			fmt.Printf("✓ Updated %s on Drive\n", path)
			// Update metadata with new mod time
			saveMeta(path, meta.DriveID, meta.DriveName, true)
		}
	}

	return nil
}

func (dc *DriveClient) updateDriveFile(fileID, localPath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = dc.service.Files.Update(fileID, &drive.File{}).Media(file).Do()
	return err
}

// RM Command - Remove file from Drive
func (dc *DriveClient) removeFile(fileID string) error {
	// Get file info first
	file, err := dc.service.Files.Get(fileID).Fields("id, name").Do()
	if err != nil {
		return fmt.Errorf("file not found: %v", err)
	}

	fmt.Printf("Are you sure you want to delete '%s'? (y/N): ", file.Name)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	// Delete the file
	err = dc.service.Files.Delete(fileID).Do()
	if err != nil {
		return fmt.Errorf("unable to delete file: %v", err)
	}

	fmt.Printf("✓ Deleted '%s' from Drive\n", file.Name)
	return nil
}

// CAT Command - View file contents without downloading
func (dc *DriveClient) viewFileContents(fileID string) error {
	// Get file metadata
	file, err := dc.service.Files.Get(fileID).Fields("id, name, mimeType, size").Do()
	if err != nil {
		return fmt.Errorf("file not found: %v", err)
	}

	fmt.Printf("=== %s ===\n\n", file.Name)

	// Check if it's a text-based file
	if !isTextFile(file.MimeType) {
		return fmt.Errorf("cannot display non-text file (type: %s)", file.MimeType)
	}

	var resp *http.Response

	if strings.HasPrefix(file.MimeType, "application/vnd.google-apps") {
		// Export Google Workspace files as plain text
		resp, err = dc.service.Files.Export(fileID, "text/plain").Download()
	} else {
		resp, err = dc.service.Files.Get(fileID).Download()
	}

	if err != nil {
		return fmt.Errorf("unable to download file: %v", err)
	}
	defer resp.Body.Close()

	// Stream contents to stdout
	scanner := bufio.NewScanner(resp.Body)
	lineNum := 1
	for scanner.Scan() {
		fmt.Printf("%4d | %s\n", lineNum, scanner.Text())
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("\n=== End of %s ===\n", file.Name)
	return nil
}

func isTextFile(mimeType string) bool {
	textTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/javascript",
		"application/vnd.google-apps.document",
		"application/vnd.google-apps.spreadsheet",
	}

	for _, t := range textTypes {
		if strings.Contains(mimeType, t) {
			return true
		}
	}
	return false
}

// openBrowser tries to open the URL in a browser
func openBrowser(url string) {
	var err error
	switch {
	case fileExists("/usr/bin/xdg-open"):
		err = exec.Command("xdg-open", url).Start()
	case fileExists("/usr/bin/open"):
		err = exec.Command("open", url).Start()
	default:
		// Windows or other
		err = exec.Command("cmd", "/c", "start", url).Start()
	}

	if err != nil {
		fmt.Printf("⚠ Could not open browser automatically: %v\n", err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
