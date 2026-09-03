package main

import (
	"fmt"
	"os"

	configpkg "quantito-engine/internal/config"

	"github.com/BurntSushi/toml"
)

const defaultConfigFile = "./config.toml"

// Config represents the application configuration model
type Config struct {
	API configpkg.APIConfig `toml:"api"`
	DB  configpkg.DBConfig  `toml:"database"`
	Log configpkg.LogConfig `toml:"log"`
}

// getConfigPath determines the configuration file path from environment variables or default
func getConfigPath() string {
	if path := os.Getenv("CONFIG_FILE"); path != "" {
		return path
	}
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return defaultConfigFile
}

// loadConfig loads configuration from a specific path
func loadConfig(path string) (*Config, error) {
	config := buildDefaultConfig()
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration file %q: %w", path, err)
	}
	return config, nil
}

// getConfig returns loaded configuration structure and the path it was loaded from
func getConfig() (*Config, string) {
	path := getConfigPath()
	config, err := loadConfig(path)
	if err != nil {
		panic(err)
	}
	return config, path
}

// buildDefaultConfig creates a configuration object with default values
func buildDefaultConfig() *Config {
	return &Config{
		Log: configpkg.LogConfig{
			Level:  configpkg.LevelInfo,
			Logger: "split-json",
		},
		DB: configpkg.DBConfig{
			MaxOpenConn:     10,
			MaxIdleConn:     10,
			MaxTTLConn:      60 * 3,
			MaxIdleTimeConn: 60,
		},
	}
}

