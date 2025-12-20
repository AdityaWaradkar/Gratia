package logger

import (
	"log"
	"os"
)

// Logger is the shared service logger
var Logger *log.Logger

// Init initializes the logger
func Init() {
	Logger = log.New(
		os.Stdout,
		"[USER_SERVICE] ",
		log.LstdFlags|log.Lshortfile,
	)
}
