package logging

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogrusLogging(t *testing.T) {
	l, hook := test.NewNullLogger()

	initLogger(logrusLogger{
		logger: l,
		fields: make(logrus.Fields),
	})

	log := GetLogger("my-component", "my-activity")
	log.WithField("my-field", "field value").Info("Hello!")

	for _, entry := range hook.AllEntries() {
		msg, err := entry.String()
		assert.NoError(t, err)
		assert.Contains(t, msg, "level=info")
		assert.Contains(t, msg, "msg=\"Hello!\"")
		assert.Contains(t, msg, "component=my-component")
		assert.Contains(t, msg, "activity=my-activity")
		assert.Contains(t, msg, "my-field=\"field value\"")
	}
}
