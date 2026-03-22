package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/PC0staS/mesh/internal/monitor"
)

// GetConfigPath devuelve la ruta del archivo config
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "mesh", "config.json"), nil
}
// LoadConfig lee el config.json
func LoadConfig () (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Si no existe el archivo, devuelve config vacío
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{Servers: []monitor.Server{}}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	
	return &config, nil
}

// SaveConfig escribe el config.json
func SaveConfig (config *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	return err
}
func GetSocketPath() string {
	return "/tmp/mesh.sock"
}
