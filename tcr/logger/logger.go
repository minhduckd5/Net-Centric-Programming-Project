package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

// InitLogger initializes the logging system
func InitLogger(logDir string, logLevel string) error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFile := filepath.Join(logDir, fmt.Sprintf("tcr_%s.log", timestamp))

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Initialize loggers
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Debug logger only if debug level is enabled
	if logLevel == "debug" {
		DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		// Create a no-op debug logger
		DebugLogger = log.New(os.Stdout, "", 0)
		DebugLogger.SetOutput(io.Discard)
	}

	return nil
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	DebugLogger.Printf(format, v...)
}

// Fatal logs a fatal error and exits
func Fatal(format string, v ...interface{}) {
	ErrorLogger.Fatalf(format, v...)
}
