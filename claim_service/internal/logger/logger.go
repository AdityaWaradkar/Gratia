package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Level string
}

func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)

	// Set as global default logger
	slog.SetDefault(logger)

	return logger
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
