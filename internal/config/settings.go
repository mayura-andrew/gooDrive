package config

type Settings struct {
	Theme        string `json:"theme"`         // UI theme (e.g., "dark", "light")
	Language     string `json:"language"`      // Language preference (e.g., "en", "es")
	SyncInterval int    `json:"sync_interval"` // Interval for sync operations in minutes
	DownloadPath string `json:"download_path"` // Default path for downloaded files
	UploadPath   string `json:"upload_path"`   // Default path for uploaded files
}

func LoadSettings(filePath string) (*Settings, error) {
	// Implementation for loading settings from a file
	return &Settings{}, nil
}

func SaveSettings(filePath string, settings *Settings) error {
	// Implementation for saving settings to a file
	return nil
}
