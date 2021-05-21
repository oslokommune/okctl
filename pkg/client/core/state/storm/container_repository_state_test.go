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

func TestContainerRepositoryStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.ContainerRepository{})
	assert.NoError(t, err)

	state := storm.NewContainerRepositoryState(db)

	err = state.SaveContainerRepository(mock.ContainerRepository())
	assert.NoError(t, err)

	m, err := state.GetContainerRepository(mock.DefaultImage)
	assert.NoError(t, err)
	assert.Equal(t, mock.ContainerRepository(), m)

	m.ImageName = "fake4"
	err = state.SaveContainerRepository(m)
	assert.NoError(t, err)

	err = state.RemoveContainerRepository(mock.DefaultImage)
	assert.NoError(t, err)

	err = state.RemoveContainerRepository(mock.DefaultImage)
	assert.NoError(t, err)
}
