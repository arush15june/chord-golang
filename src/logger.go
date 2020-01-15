package main

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

// NewLogger produces a new logger with required name
func NewLogger(name string) *log.Logger {
	logger := log.New(os.Stdout,
		fmt.Sprintf("%s: ", name),
		log.Ldate|log.Ltime|log.LUTC|log.Lshortfile,
	)

	return logger
}

// InitLogger initializes the global logger.
func InitLogger() {
	logger = NewLogger(*HostName)
}
