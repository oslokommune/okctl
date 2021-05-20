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

func TestParameterStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.SecretParameter{})
	assert.NoError(t, err)

	state := storm.NewParameterState(db)

	err = state.SaveSecret(mock.SecretParameter("secret"))
	assert.NoError(t, err)

	m, err := state.GetSecret(mock.DefaultSecretParameterName)
	assert.NoError(t, err)
	assert.Equal(t, mock.SecretParameter(""), m)

	m.Path = "fake13"
	err = state.SaveSecret(m)
	assert.NoError(t, err)

	err = state.RemoveSecret(mock.DefaultSecretParameterName)
	assert.NoError(t, err)

	err = state.RemoveSecret(mock.DefaultSecretParameterName)
	assert.NoError(t, err)
}
