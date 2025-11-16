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

type Logger struct {
	writer *log.Logger
}

func New() *Logger {
	return &Logger{
		writer: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) write(level LogLevel, message string) {
	t := time.Now().Format(time.RFC3339)
	l.writer.Printf("[%s] [%s] %s", t, level, message)
}

func (l *Logger) Info(message string) {
	l.write(INFO, message)
}

func (l *Logger) Warn(message string) {
	l.write(WARN, message)
}

func (l *Logger) Error(message string) {
	l.write(ERROR, message)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}
