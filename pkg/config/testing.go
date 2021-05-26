package config

import (
	"os"
	"testing"
)

// SkipUnlessIntegration skips a test when not running in a CI environment
func SkipUnlessIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration test due to INTEGRATION_TESTS not being set")
	}
}
