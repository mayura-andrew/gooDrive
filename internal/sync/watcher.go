package sync

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors a directory for changes and triggers sync operations.
type Watcher struct {
	directory string
	watcher   *fsnotify.Watcher
}

// NewWatcher creates a new Watcher for the specified directory.
func NewWatcher(directory string) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		directory: directory,
		watcher:   w,
	}, nil
}

// Start begins watching the directory for changes.
func (w *Watcher) Start() error {
	err := w.watcher.Add(w.directory)
	if err != nil {
		return err
	}

	go w.listen()
	return nil
}

// listen listens for file system events and triggers sync operations.
func (w *Watcher) listen() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}

// handleEvent processes file system events and triggers sync.
func (w *Watcher) handleEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Write == fsnotify.Write {
		fmt.Println("Modified file:", event.Name)
		// Trigger sync operation here
	}
	if event.Op&fsnotify.Create == fsnotify.Create {
		fmt.Println("Created file:", event.Name)
		// Trigger sync operation here
	}
	if event.Op&fsnotify.Remove == fsnotify.Remove {
		fmt.Println("Deleted file:", event.Name)
		// Trigger sync operation here
	}
}

// Stop stops the watcher.
func (w *Watcher) Stop() error {
	return w.watcher.Close()
}
