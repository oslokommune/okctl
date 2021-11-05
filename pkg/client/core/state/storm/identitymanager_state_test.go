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

func TestIdentityPoolStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New()

	db.SetDatabaseFilePath(filepath.Join(dir, "storm.db"))
	db.SetWritable(true)

	err = db.Init(&storm.IdentityPool{})
	assert.NoError(t, err)

	state := storm.NewIdentityManager(db)

	err = state.SaveIdentityPool(mock.IdentityPool())
	assert.NoError(t, err)

	m, err := state.GetIdentityPool(mock.StackNameIdentityPool)
	assert.NoError(t, err)
	assert.Equal(t, mock.IdentityPool(), m)

	m.UserPoolID = "fake16"
	err = state.SaveIdentityPool(m)
	assert.NoError(t, err)

	err = state.RemoveIdentityPool(mock.StackNameIdentityPool)
	assert.NoError(t, err)

	err = state.RemoveIdentityPool(mock.StackNameIdentityPool)
	assert.NoError(t, err)
}

func TestIdentityPoolClientStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New()

	db.SetDatabaseFilePath(filepath.Join(dir, "storm.db"))
	db.SetWritable(true)

	err = db.Init(&storm.IdentityPoolClient{})
	assert.NoError(t, err)

	state := storm.NewIdentityManager(db)

	err = state.SaveIdentityPoolClient(mock.IdentityPoolClient())
	assert.NoError(t, err)

	m, err := state.GetIdentityPoolClient(mock.StackNameIdentityPoolClient)
	assert.NoError(t, err)

	m.ClientID = "fake15"
	err = state.SaveIdentityPoolClient(m)
	assert.NoError(t, err)

	err = state.RemoveIdentityPoolClient(mock.StackNameIdentityPoolClient)
	assert.NoError(t, err)

	err = state.RemoveIdentityPoolClient(mock.StackNameIdentityPoolClient)
	assert.NoError(t, err)
}

func TestIdentityPoolUserStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New()

	db.SetDatabaseFilePath(filepath.Join(dir, "storm.db"))
	db.SetWritable(true)

	err = db.Init(&storm.IdentityPoolUser{})
	assert.NoError(t, err)

	state := storm.NewIdentityManager(db)

	err = state.SaveIdentityPoolUser(mock.IdentityPoolUser())
	assert.NoError(t, err)

	m, err := state.GetIdentityPoolUser(mock.StackNameIdentityPoolUser)
	assert.NoError(t, err)
	assert.Equal(t, mock.IdentityPoolUser(), m)

	m.UserPoolID = "fake9"
	err = state.SaveIdentityPoolUser(m)
	assert.NoError(t, err)

	err = state.RemoveIdentityPoolUser(mock.StackNameIdentityPoolUser)
	assert.NoError(t, err)

	err = state.RemoveIdentityPoolUser(mock.StackNameIdentityPoolUser)
	assert.NoError(t, err)
}
