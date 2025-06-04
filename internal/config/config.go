package config

import (
	"encoding/json"
	"os"
)

// Config holds application configuration loaded from file.
type Config struct {
	DatabaseURL   string `json:"database_url"`
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
}

// Load reads the configuration from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
