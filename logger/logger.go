package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Level represents the severity level of a log message
type Level int

const (
	// DEBUG level for detailed debugging information
	DEBUG Level = iota
	// INFO level for general operational information
	INFO
	// WARN level for warning messages
	WARN
	// ERROR level for error messages
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// Logger is the interface that defines the logging methods
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	SetLevel(level Level)
	GetLevel() Level
}

// DefaultLogger is the default implementation of the Logger interface
type DefaultLogger struct {
	mu    sync.Mutex
	out   io.Writer
	level Level
}

// NewLogger creates a new DefaultLogger with the specified output writer and log level
func NewLogger(out io.Writer, level Level) *DefaultLogger {
	return &DefaultLogger{
		out:   out,
		level: level,
	}
}

// Default is the default logger instance
var Default = NewLogger(os.Stdout, INFO)

// SetLevel sets the minimum log level
func (l *DefaultLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current log level
func (l *DefaultLogger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// log logs a message at the specified level
func (l *DefaultLogger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelName := levelNames[level]
	message := fmt.Sprintf(format, args...)

	fmt.Fprintf(l.out, "[%s] [%s] %s\n", timestamp, levelName, message)
}

// Debug logs a message at DEBUG level
func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs a message at INFO level
func (l *DefaultLogger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a message at WARN level
func (l *DefaultLogger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs a message at ERROR level
func (l *DefaultLogger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Global convenience functions that use the default logger

// Debug logs a message at DEBUG level using the default logger
func Debug(format string, args ...interface{}) {
	Default.Debug(format, args...)
}

// Info logs a message at INFO level using the default logger
func Info(format string, args ...interface{}) {
	Default.Info(format, args...)
}

// Warn logs a message at WARN level using the default logger
func Warn(format string, args ...interface{}) {
	Default.Warn(format, args...)
}

// Error logs a message at ERROR level using the default logger
func Error(format string, args ...interface{}) {
	Default.Error(format, args...)
}

// SetLevel sets the minimum log level for the default logger
func SetLevel(level Level) {
	Default.SetLevel(level)
}

// GetLevel returns the current log level for the default logger
func GetLevel() Level {
	return Default.GetLevel()
}
