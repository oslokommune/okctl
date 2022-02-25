package logging

type noopLogger struct{}

// Trace does nothing due to this being a NOOP logger
func (l noopLogger) Trace(_ string) { /* this is a NOOP logger, it shouldn't do anything */ }

// Debug does nothing due to this being a NOOP logger
func (l noopLogger) Debug(_ string) { /* this is a NOOP logger, it shouldn't do anything */ }

// Info does nothing due to this being a NOOP logger
func (l noopLogger) Info(_ string) { /* this is a NOOP logger, it shouldn't do anything */ }

// Warn does nothing due to this being a NOOP logger
func (l noopLogger) Warn(_ string) { /* this is a NOOP logger, it shouldn't do anything */ }

// Error does nothing due to this being a NOOP logger
func (l noopLogger) Error(_ string) { /* this is a NOOP logger, it shouldn't do anything */ }

// WithField does nothing due to this being a NOOP logger
func (l noopLogger) WithField(_ string, _ interface{}) Logger {
	return l
}

func newNoopLogger() Logger {
	return &noopLogger{}
}
