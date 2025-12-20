package logger

import (
	"log"
	"os"
)

// New creates a standard service logger
func New(serviceName string) *log.Logger {
	prefix := "[" + serviceName + "] "
	return log.New(os.Stdout, prefix, log.LstdFlags|log.Lshortfile)
}
