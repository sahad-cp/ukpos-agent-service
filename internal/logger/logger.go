package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(logFile string) zerolog.Logger {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	zerolog.TimeFieldFormat = time.RFC3339

	return zerolog.New(file).
		With().
		Timestamp().
		Logger()
}