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

func TestIdentityPoolStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db, err := stormpkg.Open(filepath.Join(dir, "storm.db"), stormpkg.Codec(json.Codec))
	assert.NoError(t, err)

	err = db.Init(&storm.IdentityPool{})
	assert.NoError(t, err)

	state := storm.NewIdentityManager(db)

	err = state.SaveIdentityPool(mock.IdentityPool())
	assert.NoError(t, err)

	m, err := state.GetIdentityPool(mock.StackNameIdentityPool)
	assert.NoError(t, err)
	assert.Equal(t, mock.IdentityPool(), m)

	err = state.RemoveIdentityPool(mock.StackNameIdentityPool)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestIdentityPoolClientStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db, err := stormpkg.Open(filepath.Join(dir, "storm.db"), stormpkg.Codec(json.Codec))
	assert.NoError(t, err)

	err = db.Init(&storm.IdentityPoolClient{})
	assert.NoError(t, err)

	state := storm.NewIdentityManager(db)

	err = state.SaveIdentityPoolClient(mock.IdentityPoolClient())
	assert.NoError(t, err)

	err = state.RemoveIdentityPoolClient(mock.StackNameIdentityPoolClient)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}

func TestIdentityPoolUserStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db, err := stormpkg.Open(filepath.Join(dir, "storm.db"), stormpkg.Codec(json.Codec))
	assert.NoError(t, err)

	err = db.Init(&storm.IdentityPoolUser{})
	assert.NoError(t, err)

	state := storm.NewIdentityManager(db)

	err = state.SaveIdentityPoolUser(mock.IdentityPoolUser())
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}
