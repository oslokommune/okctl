package storm_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client/mock"

	"github.com/oslokommune/okctl/pkg/client/core/state/storm"
	"github.com/stretchr/testify/assert"
)

func TestGithubRepositoryStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.GithubRepository{})
	assert.NoError(t, err)

	state := storm.NewGithubState(db)

	err = state.SaveGithubRepository(mock.GithubRepository())
	assert.NoError(t, err)

	m, err := state.GetGithubRepository(mock.DefaultGithubFullName)
	assert.NoError(t, err)
	assert.Equal(t, mock.GithubRepository(), m)

	m.Repository = "fake7"
	err = state.SaveGithubRepository(m)
	assert.NoError(t, err)

	err = state.RemoveGithubRepository(mock.DefaultGithubFullName)
	assert.NoError(t, err)

	err = state.RemoveGithubRepository(mock.DefaultGithubFullName)
	assert.NoError(t, err)
}
