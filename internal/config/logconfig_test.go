package config

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestToZerologLevel(t *testing.T) {
	t.Run("trace", func(t *testing.T) {
		config := LogConfig{Level: LevelTrace}
		assert.Equal(t, zerolog.TraceLevel, config.ToZerologLevel())
	})

	t.Run("debug", func(t *testing.T) {
		config := LogConfig{Level: LevelDebug}
		assert.Equal(t, zerolog.DebugLevel, config.ToZerologLevel())
	})

	t.Run("info", func(t *testing.T) {
		config := LogConfig{Level: LevelInfo}
		assert.Equal(t, zerolog.InfoLevel, config.ToZerologLevel())
	})

	t.Run("warn", func(t *testing.T) {
		config := LogConfig{Level: LevelWarn}
		assert.Equal(t, zerolog.WarnLevel, config.ToZerologLevel())
	})

	t.Run("error", func(t *testing.T) {
		config := LogConfig{Level: LevelError}
		assert.Equal(t, zerolog.ErrorLevel, config.ToZerologLevel())
	})

	t.Run("disabled", func(t *testing.T) {
		config := LogConfig{Level: LevelDisabled}
		assert.Equal(t, zerolog.Disabled, config.ToZerologLevel())
	})

	t.Run("default info", func(t *testing.T) {
		config := LogConfig{}
		assert.Equal(t, zerolog.InfoLevel, config.ToZerologLevel())
	})
}
