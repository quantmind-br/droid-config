package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = "config.json"

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".factory", ConfigFileName), nil
}

func Load() (*ConfigData, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ConfigData{CustomModels: []CustomModel{}}, nil
		}
		return nil, err
	}

	var cfg ConfigData
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &ConfigData{CustomModels: []CustomModel{}}, nil
	}

	return &cfg, nil
}

func Save(cfg *ConfigData) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	tmpFile := path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpFile, path)
}
