package logging

// Logger Interface for structured logging
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)

	WithField(key string, value interface{}) Logger
}
