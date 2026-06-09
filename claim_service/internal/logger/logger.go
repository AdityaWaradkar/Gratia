package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level string
}

func New(cfg Config) *slog.Logger {

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: parseLevel(cfg.Level),
		},
	)

	logger := slog.New(handler)

	// Set global logger.
	slog.SetDefault(logger)

	return logger
}

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