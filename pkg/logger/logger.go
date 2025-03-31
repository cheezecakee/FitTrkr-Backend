package logger

import (
	"log"
	"os"
)

type Logger struct {
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	RequestLog *log.Logger
}

func New() *Logger {
	return &Logger{
		InfoLog:    log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:   log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		RequestLog: log.New(os.Stdout, "REQUEST\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Global logger instance
var Log = New()
