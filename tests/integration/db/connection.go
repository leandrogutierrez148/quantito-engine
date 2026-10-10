package db

import (
	"fmt"
	"os"
	"testing"

	"quantito-engine/internal/config"
	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/interfaces"

	"github.com/BurntSushi/toml"
)

// findConfigFile locates .dev/config.toml relative to the execution working directory
func findConfigFile() string {
	if envPath := os.Getenv("CONFIG_FILE"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	candidates := []string{
		"../configs/config.toml",
		"../../configs/config.toml",
		"tests/integration/configs/config.toml",
		"../../../.dev/config.toml",
		".dev/config.toml",
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "../configs/config.toml"
}

// LoadDBConfig loads database configuration from .dev/config.toml with environment variable fallbacks
func LoadDBConfig() (*config.DBConfig, error) {
	configPath := findConfigFile()

	var fileConfig struct {
		DB config.DBConfig `toml:"database"`
	}

	if _, err := toml.DecodeFile(configPath, &fileConfig); err != nil {
		return nil, fmt.Errorf("failed to decode config file at %s: %w", configPath, err)
	}

	dbConf := fileConfig.DB

	// Allow environment variable overrides or fallback to local 127.0.0.1 if docker host isn't reachable
	if host := os.Getenv("DB_HOST"); host != "" {
		dbConf.Host = host
	} else if dbConf.Host == "mysql:3306" {
		dbConf.Host = "127.0.0.1:3306"
	}

	if user := os.Getenv("DB_USER"); user != "" {
		dbConf.Username = user
	}
	if pass := os.Getenv("DB_PASSWORD"); pass != "" {
		dbConf.Password = pass
	}
	if schema := os.Getenv("DB_SCHEMA"); schema != "" {
		dbConf.Schema = schema
	}

	return &dbConf, nil
}

// NewTestDBHandler creates a real MySQL DB handler using .dev configuration
func NewTestDBHandler(t *testing.T) interfaces.IDBHandler {
	dbConf, err := LoadDBConfig()
	if err != nil {
		t.Fatalf("Failed to load DB config from .dev: %v", err)
	}

	dbHandler, err := infraestructures.NewMysqlDB(dbConf)
	if err != nil {
		t.Fatalf("Failed to connect to test MySQL database at %s: %v. Please run 'make prepare-integration-tests' first.", dbConf.Host, err)
	}

	return dbHandler
}

// CleanSignalsTable truncates or deletes records from the signals and executions tables for test isolation
func CleanSignalsTable(t *testing.T, dbHandler interfaces.IDBHandler) {
	t.Helper()
	var dummy any
	_ = dbHandler.QueryScalar("DELETE FROM `signals`", &dummy)
	_ = dbHandler.QueryScalar("DELETE FROM `execution_variants`", &dummy)
	_ = dbHandler.QueryScalar("DELETE FROM `executions`", &dummy)
}
