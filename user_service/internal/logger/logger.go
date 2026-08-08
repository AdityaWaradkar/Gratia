package logger

import (
    "log/slog"
    "os"
    "strings"
)

// New initializes and sets up a structured JSON logger for the service
func New(serviceName string, logLevel string) *slog.Logger {
    var level slog.Level
    
    switch strings.ToLower(logLevel) {
    case "debug":
        level = slog.LevelDebug
    case "warn":
        level = slog.LevelWarn
    case "error":
        level = slog.LevelError
    default:
        level = slog.LevelInfo
    }

    opts := &slog.HandlerOptions{
        Level: level,
    }

    handler := slog.NewJSONHandler(os.Stdout, opts)
    logger := slog.New(handler).With(slog.String("service", serviceName))
    
    slog.SetDefault(logger)

    return logger
}