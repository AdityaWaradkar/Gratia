package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// Logger wraps standard log.Logger
type Logger struct {
	logger *log.Logger
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
	}
}

// logMessage prints a formatted log message with timestamp and level
func (l *Logger) logMessage(level LogLevel, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	l.logger.Printf("[%s] [%s] %s", timestamp, level, message)
}

// Info logs an info-level message
func (l *Logger) Info(msg string) {
	l.logMessage(INFO, msg)
}

// Warn logs a warning-level message
func (l *Logger) Warn(msg string) {
	l.logMessage(WARN, msg)
}

// Error logs an error-level message
func (l *Logger) Error(msg string) {
	l.logMessage(ERROR, msg)
}

// Errorf logs a formatted error-level message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info-level message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warn-level message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}
