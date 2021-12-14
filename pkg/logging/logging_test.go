package logging

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func initTestLogger() *test.Hook {
	l, hook := test.NewNullLogger()
	l.SetLevel(logrus.TraceLevel)

	initLogger(logrusLogger{
		logger: l,
		fields: make(logrus.Fields),
	})

	return hook
}

func TestInitLogger(t *testing.T) {
	err := InitLogger("log.txt")
	assert.NoError(t, err)
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

func TestLogLevels(t *testing.T) {
	hook := initTestLogger()

	messages := []string{
		"Something very low level.",
		"Useful debugging information.",
		"Something noteworthy happened!",
		"You should probably take a look at this.",
		"Something failed but I'm not quitting.",
	}

	levels := []string{"trace", "debug", "info", "warn", "error"}

	log := GetLogger("my-component", "my-activity")
	log.Trace(messages[0])
	log.Debug(messages[1])
	log.Info(messages[2])
	log.Warn(messages[3])
	log.Error(messages[4])

	assert.Equal(t, 5, len(hook.Entries), "Should have 5 log messages")

	for i, entry := range hook.AllEntries() {
		msg, err := entry.String()
		assert.NoError(t, err)
		assert.Contains(t, msg, messages[i])
		assert.Contains(t, msg, fmt.Sprintf("level=%s", levels[i]))
	}
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
