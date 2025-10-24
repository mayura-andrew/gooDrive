package models

import (
	"encoding/json"
	"os"
)

type File struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	MimeType     string  `json:"mimeType"`
	Size         int64   `json:"size"`
	CreatedTime  string  `json:"createdTime"`
	ModifiedTime string  `json:"modifiedTime"`
	Owners       []Owner `json:"owners"`
	Shared       bool    `json:"shared"`
}

type Owner struct {
	Display      string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Me           bool   `json:"me"`
}

func NewFile(id, name, mimeType string, size int64, createdTime, modifiedTime string, owners []Owner, shared bool) *File {
	return &File{
		ID:           id,
		Name:         name,
		MimeType:     mimeType,
		Size:         size,
		CreatedTime:  createdTime,
		ModifiedTime: modifiedTime,
		Owners:       owners,
		Shared:       shared,
	}
}

// FileMeta represents file metadata for tracking
type FileMeta struct {
	LocalPath string            `json:"local_path"`
	DriveID   string            `json:"drive_id"`
	DriveName string            `json:"drive_name"`
	ModTime   string            `json:"mod_time"`
	IsTracked bool              `json:"is_tracked"`
	Children  map[string]string `json:"children,omitempty"` // name -> driveID
}

const MetaFile = ".drive-cli-meta.json"

// SaveMeta saves metadata for a tracked file
func SaveMeta(localPath, driveID, driveName string, isTracked bool) error {
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

	metaPath := localPath + MetaFile
	file, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(meta)
}

// LoadMeta loads metadata for a tracked file
func LoadMeta(localPath string) (*FileMeta, error) {
	metaPath := localPath + MetaFile
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
