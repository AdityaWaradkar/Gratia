package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Config holds the logger configuration
type Config struct {
	Level string
}

// New creates a new structured JSON logger instance
func New(cfg Config) *slog.Logger {

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: parseLevel(cfg.Level),
		},
	)

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}

// parseLevel converts a string log level to slog.Level
func parseLevel(level string) slog.Level {

	switch strings.ToLower(level) {

	case "debug":
		return slog.LevelDebug

	case "warn", "warning":
		return slog.LevelWarn

	case "error":
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}