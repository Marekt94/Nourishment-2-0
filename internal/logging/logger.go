package logging

// Logger is a generic logging interface for the project.
type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}
