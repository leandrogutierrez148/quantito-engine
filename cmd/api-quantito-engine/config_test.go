package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigPath(t *testing.T) {
	t.Run("default when no env set", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")
		os.Unsetenv("CONFIG_PATH")
		assert.Equal(t, defaultConfigFile, getConfigPath())
	})

	t.Run("CONFIG_FILE takes precedence", func(t *testing.T) {
		os.Setenv("CONFIG_FILE", "/custom/path/config.toml")
		defer os.Unsetenv("CONFIG_FILE")
		assert.Equal(t, "/custom/path/config.toml", getConfigPath())
	})

	t.Run("CONFIG_PATH used if CONFIG_FILE not set", func(t *testing.T) {
		os.Unsetenv("CONFIG_FILE")
		os.Setenv("CONFIG_PATH", "/another/path/config.toml")
		defer os.Unsetenv("CONFIG_PATH")
		assert.Equal(t, "/another/path/config.toml", getConfigPath())
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "config.toml")
		content := `
[api]
listen = ":9090"

[database]
host = "localhost:3306"
schema = "quantito-test"
username = "testuser"
password = "testpassword"
maxOpenConn = 20

[log]
level = "debug"
logger = "console"
`
		err := os.WriteFile(configFile, []byte(content), 0644)
		assert.NoError(t, err)

		cfg, err := loadConfig(configFile)
		assert.NoError(t, err)
		assert.Equal(t, ":9090", cfg.API.Listen)
		assert.Equal(t, "localhost:3306", cfg.DB.Host)
		assert.Equal(t, "quantito-test", cfg.DB.Schema)
		assert.Equal(t, "testuser", cfg.DB.Username)
		assert.Equal(t, "testpassword", cfg.DB.Password)
		assert.Equal(t, 20, cfg.DB.MaxOpenConn)
		assert.Equal(t, 10, cfg.DB.MaxIdleConn) // default value preserved
		assert.Equal(t, "debug", cfg.Log.Level)
		assert.Equal(t, "console", cfg.Log.Logger)
	})

	t.Run("non-existent config file returns error", func(t *testing.T) {
		_, err := loadConfig("/path/does/not/exist.toml")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode configuration file")
	})

	t.Run("invalid toml syntax returns error", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "invalid.toml")
		err := os.WriteFile(configFile, []byte("invalid toml content = ="), 0644)
		assert.NoError(t, err)

		_, err = loadConfig(configFile)
		assert.Error(t, err)
	})
}
