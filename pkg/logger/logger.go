package logger

import (
	"fmt"

	"github.com/caarlos0/log"
)

// This package wraps the underlying logging library in order to make different log levels more visually distinct
// If any further customization is desired, we'll probably want to change into a more fully featured logging library

// Debug logs at the debug level.
func Debug(msg string) {
	log.Debug(msg)
}

// Info logs at the info level.
func Info(msg string) {
	log.Info(msg)
}

// Warn logs at the warn level.
func Warn(msg string) {
	log.Warn(yellow(msg))
}

// Error logs at the error level.
func Error(msg string) {
	log.Error(red(msg))
}

// Fatal logs at the fatal levl.
func Fatal(msg string) {
	log.Fatal(red(msg))
}

// Debugf logs at the debug level with string formatting
func Debugf(format string, a ...any) {
	Debug(fmt.Sprintf(format, a...))
}

// Infof logs at the debug level with string formatting
func Infof(format string, a ...any) {
	Info(fmt.Sprintf(format, a...))
}

// Warnf logs at the debug level with string formatting
func Warnf(format string, a ...any) {
	Warn(fmt.Sprintf(format, a...))
}

// Errorf logs at the debug level with string formatting
func Errorf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}

// Fatalf logs at the debug level with string formatting
func Fatalf(format string, a ...any) {
	Fatal(fmt.Sprintf(format, a...))
}

// SetLevel sets the log level
func SetLevel(level log.Level) {
	log.SetLevel(level)
}

// DecreasePadding decreases padding
func DecreasePadding() {
	log.DecreasePadding()
}

// IncreasePadding increases padding
func IncreasePadding() {
	log.IncreasePadding()
}

// WithError logs an error
func WithError(err error) *log.Entry {
	return log.WithError(err)
}

func red(msg string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", msg)
}

func yellow(msg string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", msg)
}
