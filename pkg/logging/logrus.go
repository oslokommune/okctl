package logging

import (
	"os"

	"github.com/oslokommune/okctl/pkg/context"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

func (l logrusLogger) Debug(msg string) {
	l.logger.WithFields(l.fields).Debug(msg)
}

func (l logrusLogger) Info(msg string) {
	l.logger.WithFields(l.fields).Info(msg)
}

func (l logrusLogger) Warn(msg string) {
	l.logger.WithFields(l.fields).Warn(msg)
}

func (l logrusLogger) Error(msg string) {
	l.logger.WithFields(l.fields).Error(msg)
}

func (l logrusLogger) WithField(key string, value interface{}) Logger {
	newFields := make(logrus.Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}

	newFields[key] = value

	return logrusLogger{
		logger: l.logger,
		fields: newFields,
	}
}

func newLogrusLogger() Logger {
	_, debug := os.LookupEnv(context.DefaultDebugEnv)

	logger := logrus.New()

	logger.Out = os.Stderr
	logger.Formatter = &logrus.TextFormatter{}
	logger.Level = logrus.InfoLevel

	if debug {
		logger.Level = logrus.TraceLevel
	}

	return logrusLogger{
		logger: logger,
		fields: make(logrus.Fields),
	}
}
