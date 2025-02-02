// Package log provides basic logger struct with zerolog's logger
package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is a basic logger struct
type Logger struct {
	*zerolog.Logger
}

// New creates a new Logger
func New() *Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{
		Logger: &logger,
	}
}
