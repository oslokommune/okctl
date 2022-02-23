// Package logging API for structured logging
package logging

// The global logger instance
var logger = newNoopLogger() //nolint: gochecknoglobals

// InitLogger initialize the global logger instance
func InitLogger(logFile string) error {
	logger, err := newLogrusLogger(logFile)
	if err == nil {
		initLogger(logger)
	}

	return err
}

func initLogger(l Logger) {
	logger = l
}

// GetLogger returns a logger with the given component and activity as
// context fields on all log entries
func GetLogger(component string, activity string) Logger {
	return logger.WithField("component", component).WithField("activity", activity)
}
