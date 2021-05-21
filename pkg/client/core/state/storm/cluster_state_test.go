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

func TestClusterStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.Cluster{})
	assert.NoError(t, err)

	state := storm.NewClusterState(db)

	err = state.SaveCluster(mock.Cluster())
	assert.NoError(t, err)

	m, err := state.GetCluster(mock.DefaultClusterName)
	assert.NoError(t, err)
	assert.Equal(t, mock.Cluster(), m)

	m.Config = nil
	err = state.SaveCluster(m)
	assert.NoError(t, err)

	err = state.RemoveCluster(mock.DefaultClusterName)
	assert.NoError(t, err)

	err = state.RemoveCluster(mock.DefaultClusterName)
	assert.NoError(t, err)
}
