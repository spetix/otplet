package logger

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"

	zerolog "github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

func SetupLogger(logLevel string, logFile string, logFormat string) error {
	// Create zerolog logger (destination + format)

	var out io.Writer = os.Stderr
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		out = f
	}

	// --- zerolog base logger ---
	var zl zerolog.Logger
	switch strings.ToLower(logFormat) {
	case "json":
		zl = zerolog.New(out).With().Timestamp().Logger()
	case "text":
		zl = zerolog.New(
			zerolog.ConsoleWriter{Out: out},
		).With().Timestamp().Logger()
	default:
		return errors.New("invalid log-format (use text or json)")
	}

	// --- slog level ---
	level := new(slog.LevelVar)
	switch strings.ToLower(logLevel) {
	case "debug":
		level.Set(slog.LevelDebug)
	case "info":
		level.Set(slog.LevelInfo)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	default:
		return errors.New("invalid log-level (debug, info, warn, error)")
	}

	// --- bridge slog → zerolog ---
	handler := slogzerolog.Option{
		Level:  level,
		Logger: &zl,
	}.NewZerologHandler()

	slog.SetDefault(slog.New(handler))
	return nil

}
