package sync

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/mayura-andrew/goodrive/internal/drive"
)

// SyncManager manages the synchronization of files between local storage and Google Drive.
type SyncManager struct {
	localPath string
	drivePath string
	client    *drive.Client
	ctx       context.Context
	watcher   *fsnotify.Watcher
}

// NewSyncManager creates a new SyncManager instance.
func NewSyncManager(client *drive.Client, ctx context.Context, localPath, drivePath string) (*SyncManager, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &SyncManager{
		localPath: localPath,
		drivePath: drivePath,
		client:    client,
		ctx:       ctx,
		watcher:   watcher,
	}, nil
}

// Start begins watching the local directory for changes.
func (s *SyncManager) Start() error {
	err := s.watcher.Add(s.localPath)
	if err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	go s.watch()
	return nil
}

// watch listens for file system events and triggers sync operations.
func (s *SyncManager) watch() {
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}
			s.handleEvent(event)
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("error: %v\n", err)
		}
	}
}

// handleEvent processes file system events and performs sync operations.
func (s *SyncManager) handleEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Write == fsnotify.Write {
		fmt.Printf("modified file: %s\n", event.Name)
		s.syncFile(event.Name)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		fmt.Printf("deleted file: %s\n", event.Name)
		s.deleteFile(event.Name)
	}
}

// syncFile uploads the modified file to Google Drive.
func (s *SyncManager) syncFile(filePath string) {
	relativePath, _ := filepath.Rel(s.localPath, filePath)
	fmt.Printf("Syncing file: %s to Drive path: %s\n", relativePath, s.drivePath)

	_, err := drive.UploadFile(s.client.Service, s.ctx, filePath, s.drivePath)
	if err != nil {
		fmt.Printf("failed to upload file: %s, error: %v\n", filePath, err)
	} else {
		fmt.Printf("successfully synced: %s\n", filePath)
	}
}

// deleteFile removes the file from Google Drive.
func (s *SyncManager) deleteFile(filePath string) {
	relativePath, _ := filepath.Rel(s.localPath, filePath)
	driveFilePath := filepath.Join(s.drivePath, relativePath)

	err := s.client.DeleteFile(driveFilePath)
	if err != nil {
		fmt.Printf("failed to delete file: %s, error: %v\n", driveFilePath, err)
	}
}

// Stop stops the file watcher.
func (s *SyncManager) Stop() {
	s.watcher.Close()
}

// GenerateUniqueID generates a unique identifier for sync operations.
func GenerateUniqueID() string {
	return uuid.New().String()
}

// Synchronize performs a one-time synchronization
func Synchronize() {
	fmt.Println("Starting synchronization...")
	fmt.Println("Note: This is a placeholder. Implement full sync logic as needed.")
	// TODO: Implement full synchronization logic
	// This would typically:
	// 1. Initialize the Drive client
	// 2. Compare local files with remote files
	// 3. Upload/download changes
	// 4. Handle conflicts
}
