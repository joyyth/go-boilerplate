package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Options struct {
	Level       string // "debug", "info", "warn", "error"
	Pretty      bool   // json in prod
	ServiceName string
}

func NewLogger(opts Options) zerolog.Logger {
	level, err := zerolog.ParseLevel(opts.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var logger zerolog.Logger

	if opts.Pretty {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(output)
	} else {
		logger = zerolog.New(os.Stdout)
	}

	return logger.With().Timestamp().Str("service", opts.ServiceName).Logger()
}
