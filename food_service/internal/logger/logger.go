package logger

import (
    "log"
    "os"
)

// Init initializes application-wide logger
func Init() {
    log.SetOutput(os.Stdout)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    log.Println("logger initialized")
}