package logger

import (
	"log"
)

var Logger *log.Logger

func Init(level string) {
	Logger = log.Default()
	Logger.Println("Logger initialized with level:", level)
}
