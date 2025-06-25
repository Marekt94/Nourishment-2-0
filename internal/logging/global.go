package logging

import "fmt"

// Global is the global logger instance for the project.
var Global Logger = &dummyLogger{}

// dummyLogger implements Logger and logs a warning that logger is not initialized
type dummyLogger struct{}

func (d *dummyLogger) Tracef(format string, args ...interface{}) {
	fmt.Println("[dummyLogger] Logger not initialized: TRACE:", fmt.Sprintf(format, args...))
}
func (d *dummyLogger) Debugf(format string, args ...interface{}) {
	fmt.Println("[dummyLogger] Logger not initialized: DEBUG:", fmt.Sprintf(format, args...))
}
func (d *dummyLogger) Infof(format string, args ...interface{}) {
	fmt.Println("[dummyLogger] Logger not initialized: INFO:", fmt.Sprintf(format, args...))
}
func (d *dummyLogger) Warnf(format string, args ...interface{}) {
	fmt.Println("[dummyLogger] Logger not initialized: WARN:", fmt.Sprintf(format, args...))
}
func (d *dummyLogger) Panicf(format string, args ...interface{}) {
	fmt.Println("[dummyLogger] Logger not initialized: PANIC:", fmt.Sprintf(format, args...))
}

// SetGlobalLogger allows to set the global logger instance (for reuse in other projects)
func SetGlobalLogger(l Logger) {
	if l == nil {
		Global = &dummyLogger{}
	} else {
		Global = l
	}
}
