package infraestructures

import (
	"io"

	"github.com/rs/zerolog"
)

type splitWriter struct {
	io.Writer
	ErrorWriter io.Writer
}

func (lw *splitWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level != zerolog.NoLevel && level >= zerolog.WarnLevel {
		return lw.ErrorWriter.Write(p)
	}
	return lw.Write(p)
}

// NewSplitWriter represents a combination of two writers used according to log level.
// Messages with level higher than or equal to Warn are sent to ErrorWriter.
func NewSplitWriter(infoWriter, errorWriter io.Writer) *splitWriter {
	return &splitWriter{
		Writer:      infoWriter,
		ErrorWriter: errorWriter,
	}
}
