package main

import (
	configpkg "quantito-engine/internal/config"

	"github.com/BurntSushi/toml"
)

const configFile = "./config.toml"

// Config represents the application configuration model
type Config struct {
	API configpkg.APIConfig `toml:"api"`
	DB  configpkg.DBConfig  `toml:"database"`
	Log configpkg.LogConfig `toml:"log"`
}

// getConfig returns loaded configuration structure
func getConfig() *Config {
	config := buildDefaultConfig()
	if _, err := toml.DecodeFile(configFile, config); err != nil {
		panic(err)
	}
	return config
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
