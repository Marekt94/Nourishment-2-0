package logging

import (
	"io"
	"os"
)

// getLogWriter returns the centralized log output writer
func getLogWriter() io.Writer {
	return os.Stdout
}

// Logger is a generic logging interface for the project.
type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Writer() io.Writer
}
