package virtualenv_test

import (
	"fmt"
	"path"
	"sort"
	"testing"

	"github.com/oslokommune/okctl/pkg/virtualenv/shellgetter"

	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

type testHelper struct {
	t               *testing.T
	ps1Dir          string
	currentUsername string
	tmpBasedir      string
}

func newTestHelper(t *testing.T) *testHelper {
	const currentUsername = "mickeymouse"

	return &testHelper{
		t:               t,
		ps1Dir:          fmt.Sprintf("/home/%s/.okctl/%s", currentUsername, commandlineprompter.Ps1Dir),
		currentUsername: currentUsername,
		tmpBasedir:      "/tmp/okctl34796",
	}
}

func (h *testHelper) createEtcStorage(username string, shellCmd string) (storage.Storer, error) {
	s := storage.NewEphemeralStorage()

	passwdFile, err := s.Create("/", "passwd", 0o644)
	assert.Nil(h.t, err)

	_, err = passwdFile.WriteString(fmt.Sprintf(`root:x:0:0:root:/root:/bin/bash
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
bin:x:2:2:bin:/bin:/usr/sbin/nologin
johndoe:x:1000:1000:John Doe,,,:/home/%s:%s
tcpdump:x:108:116::/nonexistent:/usr/sbin/nologin
`, username, shellCmd))
	if err != nil {
		return nil, fmt.Errorf("could not write text to file: %w", err)
	}

	return s, nil
}

type userDirStorage struct {
	storage  storage.Storer
	basePath string
}

func (h *testHelper) createUserDirStorage(basepath string) *userDirStorage {
	s := storage.NewEphemeralStorage()
	s.BasePath = basepath

	return &userDirStorage{
		storage:  s,
		basePath: basepath,
	}
}

func (h *testHelper) createUserHomeDirStorage(createZshrcFile bool) (*storage.EphemeralStorage, error) {
	s := storage.NewEphemeralStorage()

	if createZshrcFile {
		file, err := s.Create(".", ".zshrc", 0o644)
		if err != nil {
			return nil, fmt.Errorf("couldn't create .zshrc file: %w", err)
		}

		err = file.Close()
		if err != nil {
			return nil, fmt.Errorf("couldn't close .zshrc file: %w", err)
		}
	}

	return s, nil
}

func (h *testHelper) toSlice(m map[string]string) []string {
	s := make([]string, 0)

	for k, v := range m {
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(s)

	return s
}

func (h *testHelper) assertGoldenVenvPs1(t *testing.T, opts commandlineprompter.CommandLinePromptOpts) {
	// Make sure executable venv_ps1 is a file that exists on the PATH
	content, err := opts.UserDirStorage.ReadAll(path.Join(commandlineprompter.Ps1Dir, commandlineprompter.Ps1Filename))
	assert.Nil(t, err)

	g := goldie.New(t)
	g.Assert(t, "venv_ps1", content)
}

// NewMacOsCmdGetter returns a MacOsUserShellGetter
func (h *testHelper) NewTestMacOsLoginShellGetter() shellgetter.MacOsUserShellCmdGetter {
	return &TestMacOsShellGetter{}
}

// DefaultMacOsShellGetter is the default implementation for getting the user's login shell on macOS
type TestMacOsShellGetter struct{}

const MacOsUserLoginShell string = "/bin/someMacShell"

// FindUserLoginShell returns
func (m *TestMacOsShellGetter) RunDsclUserShell() (string, error) {
	return MacOsUserLoginShell, nil
}
