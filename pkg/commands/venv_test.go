package commands

import (
	"testing"

	"github.com/bmizerany/assert"
)

func returnEnvOkctlShell(key string) (string, bool) {
	if key == "OKCTL_SHELL" {
		return "/bin/fish", true
	}

	return "", false
}

func returnEnvShell(key string) (string, bool) {
	if key == "SHELL" {
		return "/bin/zsh", true
	}

	return "", false
}

func returnFalse(key string) (string, bool) {
	return "", false
}

func TestGetShell(t *testing.T) {
	t.Run("Should return correct shell", func(t *testing.T) {
		assert.Equal(t, "/bin/fish", GetShell(returnEnvOkctlShell))
		assert.Equal(t, "/bin/zsh", GetShell(returnEnvShell))
		assert.Equal(t, "/bin/sh", GetShell(returnFalse))
	})
}
