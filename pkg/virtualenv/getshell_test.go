package virtualenv

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func returnEnvOkctlShell(key string) (string, bool) {
	if key == "OKCTL_SHELL" {
		return "/bin/fish", true
	}

	return "", false
}

func returnFalse(key string) (string, bool) {
	return "", false
}

func TestGetShell(t *testing.T) {
	t.Run("Should return correct shell", func(t *testing.T) {
		shellCmd, err := GetShellCmd(returnEnvOkctlShell, nil)

		assert.Nil(t, err)
		assert.Equal(t, "/bin/fish", shellCmd)
	})

	t.Run("Should return shell command specified in /etc/passwd", func(t *testing.T) {
		store := storage.NewEphemeralStorage()
		passwdFile, err := store.Create("/", "passwd", 0o644)
		assert.Nil(t, err)

		_, err = passwdFile.WriteString(`root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
bin:x:2:2:bin:/bin:/usr/sbin/nologin
johndoe:x:1000:1000:John Doe,,,:/home/yngvar:/bin/zsh
tcpdump:x:108:116::/nonexistent:/usr/sbin/nologin
`)
		assert.Nil(t, err)

		err = passwdFile.Close()
		assert.Nil(t, err)

		shellCmd, err := GetShellCmd(returnFalse, store)
		assert.Nil(t, err)

		assert.Equal(t, "/bin/zsh", shellCmd)
	})
}
