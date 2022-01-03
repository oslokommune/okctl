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

func TestManifestStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New()

	db.SetDatabaseFilePath(filepath.Join(dir, "storm.db"))
	db.SetWritable(true)

	err = db.Init(&storm.KubernetesManifest{})
	assert.NoError(t, err)

	state := storm.NewManifestState(db)

	err = state.SaveKubernetesManifests(mock.KubernetesManifest())
	assert.NoError(t, err)

	m, err := state.GetKubernetesManifests(mock.DefaultManifestName)
	assert.NoError(t, err)
	assert.Equal(t, mock.KubernetesManifest(), m)

	m.Namespace = "fake11"
	err = state.SaveKubernetesManifests(m)
	assert.NoError(t, err)

	err = state.RemoveKubernetesManifests(mock.DefaultManifestName)
	assert.NoError(t, err)

	err = state.RemoveKubernetesManifests(mock.DefaultManifestName)
	assert.NoError(t, err)
}
