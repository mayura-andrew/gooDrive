package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	DrivePath string `json:"drive_path"`
	Sync      bool   `json:"sync"`
	LogLevel  string `json:"log_level"`
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig(filePath string) (*Config, error) {
	var err error
	once.Do(func() {
		file, err := os.Open(filePath)
		if err != nil {
			return
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		instance = &Config{}
		err = decoder.Decode(instance)
	})

	return instance, err
}

func SaveConfig(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(instance)
}
