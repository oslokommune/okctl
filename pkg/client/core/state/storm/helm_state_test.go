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

func TestHelmStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.Helm{})
	assert.NoError(t, err)

	state := storm.NewHelmState(db)

	err = state.SaveHelmRelease(mock.Helm())
	assert.NoError(t, err)

	m, err := state.GetHelmRelease(mock.DefaultHelmReleaseName)
	assert.NoError(t, err)
	assert.Equal(t, mock.Helm(), m)

	m.Chart.Namespace = "fake8"
	err = state.SaveHelmRelease(m)
	assert.NoError(t, err)

	err = state.RemoveHelmRelease(mock.DefaultHelmReleaseName)
	assert.NoError(t, err)

	err = state.RemoveHelmRelease(mock.DefaultHelmReleaseName)
	assert.NoError(t, err)
}
