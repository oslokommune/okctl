package logging

type noopLogger struct{}

func (l noopLogger) Trace(_ string) {}
func (l noopLogger) Debug(_ string) {}
func (l noopLogger) Info(_ string)  {}
func (l noopLogger) Warn(_ string)  {}
func (l noopLogger) Error(_ string) {}
func (l noopLogger) WithField(_ string, _ interface{}) Logger {
	return l
}

func newNoopLogger() Logger {
	return &noopLogger{}
}
