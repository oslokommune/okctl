package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func initTestLogger() *test.Hook {
	l, hook := test.NewNullLogger()

	initLogger(logrusLogger{
		logger: l,
		fields: make(logrus.Fields),
	})

	return hook
}

func TestGetLogger(t *testing.T) {
	hook := initTestLogger()
	log := GetLogger("my-component", "my-activity")
	log.Info("Hello!")

	entries := hook.Entries
	assert.Equal(t, 1, len(entries), "Should log 1 message")

	msg, err := entries[0].String()
	assert.NoError(t, err)
	assert.Contains(t, msg, "level=info")
	assert.Contains(t, msg, "component=my-component")
	assert.Contains(t, msg, "activity=my-activity")
	assert.Contains(t, msg, "msg=\"Hello!\"")
}

func TestLogWithField(t *testing.T) {
	hook := initTestLogger()

	log := GetLogger("my-component", "my-activity")
	log.WithField("my-field", "field value").Warn("Hello!")

	msg, err := hook.Entries[0].String()
	assert.NoError(t, err)
	assert.Contains(t, msg, "level=warn")
	assert.Contains(t, msg, "my-field=\"field value\"")
}

func TestLogWithFields(t *testing.T) {
	hook := initTestLogger()

	log := GetLogger("my-component", "my-activity")
	log.WithField("my-field", "field value").
		WithField("second-field", "another value").
		Info("Hello!")

	msg, err := hook.Entries[0].String()
	assert.NoError(t, err)
	assert.Contains(t, msg, "my-field=\"field value\"")
	assert.Contains(t, msg, "second-field=\"another value\"")
}

func TestLogWithDifferentFields(t *testing.T) {
	hook := initTestLogger()

	log := GetLogger("my-component", "my-activity").
		WithField("common-field", "common value")

	log.WithField("first-field", "first value").
		Info("Hello!")

	log.WithField("second-field", "another value").
		Info("Hello again!")

	entries := hook.Entries
	assert.Equal(t, 2, len(entries), "Should log 1 message")

	msg, err := hook.Entries[0].String()
	assert.NoError(t, err)
	assert.Contains(t, msg, "common-field=\"common value\"")
	assert.Contains(t, msg, "first-field=\"first value\"")

	msg, err = hook.Entries[1].String()
	assert.NoError(t, err)
	assert.Contains(t, msg, "common-field=\"common value\"")
	assert.Contains(t, msg, "second-field=\"another value\"")
	assert.NotContains(t, msg, "first-field=\"first value\"")
}
