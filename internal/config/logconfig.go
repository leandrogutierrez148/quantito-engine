package config

import (
	"strings"

	"github.com/rs/zerolog"
)

// LogConfig represents Logging configuration model
type LogConfig struct {
	Level  string `toml:"level"`
	Logger string `toml:"logger"`
}

const (
	LevelTrace    = "trace"
	LevelDebug    = "debug"
	LevelInfo     = "info"
	LevelWarn     = "warn"
	LevelError    = "error"
	LevelDisabled = "disabled"
)

const defaultLevel = LevelInfo

func (l *LogConfig) ToZerologLevel() zerolog.Level {
	return mapZerologLevel(l.Level)
}

func mapZerologLevel(logLevel string) zerolog.Level {
	switch strings.ToLower(logLevel) {
	case LevelTrace:
		return zerolog.TraceLevel
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelDisabled:
		return zerolog.Disabled
	default:
		return mapZerologLevel(defaultLevel)
	}
}
