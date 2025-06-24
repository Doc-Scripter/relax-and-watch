package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Log levels
const (
	LevelError   = "ERROR"
	LevelWarning = "WARNING"
	LevelSuccess = "SUCCESS"
)

// Logger represents a custom logger.
type Logger struct {
	logFile  *os.File
	stdLogger *log.Logger
	mu        sync.Mutex
}

// NewLogger creates a new Logger instance.
func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	stdLogger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		logFile:  file,
		stdLogger: stdLogger,
	}, nil
}

// Close closes the log file.
func (l *Logger) Close() {
	l.logFile.Close()
}

// log writes a message to the log file with the specified level.
func (l *Logger) log(level, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.stdLogger.Printf("[%s], %s", level, message)
}

// Error logs an error message.
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, v...))
}

// Warning logs a warning message.
func (l *Logger) Warning(format string, v ...interface{}) {
	l.log(LevelWarning, fmt.Sprintf(format, v...))
}

// Success logs a success message.
func (l *Logger) Success(format string, v ...interface{}) {
	l.log(LevelSuccess, fmt.Sprintf(format, v...))
}