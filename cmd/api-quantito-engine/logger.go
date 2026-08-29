package main

import (
	"io"
	"quantito-engine/internal/config"
	"quantito-engine/internal/infraestructures"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger(conf config.LogConfig) {
	writer := createWriter(conf)
	log.Logger = zerolog.New(writer).Level(conf.ToZerologLevel()).With().Timestamp().Caller().Logger()
}

func createWriter(conf config.LogConfig) io.Writer {
	switch conf.Logger {
	case "console":
		return infraestructures.NewSplitWriter(
			zerolog.ConsoleWriter{Out: os.Stdout},
			zerolog.ConsoleWriter{Out: os.Stderr})
	case "split-json":
		return infraestructures.NewSplitWriter(os.Stdout, os.Stderr)
	default:
		return os.Stdout
	}
}
