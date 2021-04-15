package storm_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/oslokommune/okctl/pkg/client/mock"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/json"
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

	db, err := stormpkg.Open(filepath.Join(dir, "storm.db"), stormpkg.Codec(json.Codec))
	assert.NoError(t, err)

	err = db.Init(&storm.ContainerRepository{})
	assert.NoError(t, err)

	state := storm.NewContainerRepositoryState(db)

	err = state.SaveContainerRepository(mock.ContainerRepository())
	assert.NoError(t, err)

	m, err := state.GetContainerRepository(mock.StackNameContainerRepository)
	assert.NoError(t, err)
	assert.Equal(t, mock.ContainerRepository(), m)

	err = state.RemoveContainerRepository(mock.StackNameContainerRepository)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}
