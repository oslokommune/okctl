package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/rotatefilehook"

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

func newLogrusLogger(logFile string) (Logger, error) {
	_, debug := os.LookupEnv(context.DefaultDebugEnv)

	logger := logrus.New()

	logger.Out = os.Stderr
	logger.Formatter = &logrus.TextFormatter{}
	logger.Level = logrus.InfoLevel

	if debug {
		logger.Level = logrus.TraceLevel
	}

	err := AddLogFileHook(logger, logFile)
	if err != nil {
		return nil, err
	}

	return logrusLogger{
		logger: logger,
		fields: make(logrus.Fields),
	}, nil
}

// AddLogFileHook Add a Hook to write logs to rotating log files
func AddLogFileHook(logger *logrus.Logger, logFile string) error {
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logFile,
		MaxSize:    constant.DefaultLogSizeInMb,
		MaxBackups: constant.DefaultLogBackups,
		MaxAge:     constant.DefaultLogDays,
		Levels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		},
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})
	if err != nil {
		return fmt.Errorf("initialising the file rotate hook: %v", err)
	}

	logger.AddHook(rotateFileHook)

	return nil
}
