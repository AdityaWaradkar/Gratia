package logger

import (
    "log/slog"
    "os"
    "strings"
)

// New initializes and sets up a structured JSON logger for the service
func New(serviceName string, logLevel string) *slog.Logger {
    var level slog.Level
    
    // Map the string configuration to the strict slog level types
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

    // Output logs as JSON objects for easy parsing by monitoring tools
    handler := slog.NewJSONHandler(os.Stdout, opts)
    
    // Automatically inject the service name into every single log entry
    logger := slog.New(handler).With(slog.String("service", serviceName))
    
    // Override the global default logger to catch standard library outputs
    slog.SetDefault(logger)

    return logger
}