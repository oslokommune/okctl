package virtualenv_test

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type testHelper struct {
	t               *testing.T
	ps1Dir          string
	currentUsername string
}

func NewTestHelper(t *testing.T) *testHelper {
	currentUsername := "mickeymouse"

	return &testHelper{
		t:               t,
		ps1Dir:          fmt.Sprintf("/home/%s/.okctl/%s", currentUsername, commandlineprompter.Ps1Dir),
		currentUsername: currentUsername,
	}
}

func (h *testHelper) CreateEtcStorage(username string, shellCmd string) storage.Storer {
	s := storage.NewEphemeralStorage()

	passwdFile, err := s.Create("/", "passwd", 0o644)
	assert.Nil(h.t, err)

	_, err = passwdFile.WriteString(fmt.Sprintf(`root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
bin:x:2:2:bin:/bin:/usr/sbin/nologin
johndoe:x:1000:1000:John Doe,,,:/home/%s:%s
tcpdump:x:108:116::/nonexistent:/usr/sbin/nologin
`, username, shellCmd))

	return s
}

type userDirStorage struct {
	storage  storage.Storer
	basePath string
}

func (h *testHelper) CreateUserDirStorage(basepath string) *userDirStorage {
	s := storage.NewEphemeralStorage()
	s.BasePath = basepath

	return &userDirStorage{
		storage:  s,
		basePath: basepath,
	}
}

func (h *testHelper) toSlice(m map[string]string) []string {
	s := make([]string, 0)

	for k, v := range m {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	return s
}

func (h *testHelper) toMap(slice []string) map[string]string {
	m := make(map[string]string)

	for _, env := range slice {
		split := strings.SplitN(env, "=", 2)
		key := split[0]
		val := split[1]
		m[key] = val
	}

	return m
}
