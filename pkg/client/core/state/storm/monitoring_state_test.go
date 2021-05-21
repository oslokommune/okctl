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

func TestKubePromStackStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.KubePromStack{})
	assert.NoError(t, err)

	state := storm.NewMonitoringState(db)

	err = state.SaveKubePromStack(mock.KubePromStack())
	assert.NoError(t, err)

	m, err := state.GetKubePromStack()
	assert.NoError(t, err)

	a := mock.KubePromStack()

	a.Certificate = nil
	a.ExternalSecret = nil
	a.IdentityPoolClient = nil
	a.Chart = nil

	assert.Equal(t, a, m)

	a.ClientID = "fake12"
	err = state.SaveKubePromStack(a)
	assert.NoError(t, err)

	err = state.RemoveKubePromStack()
	assert.NoError(t, err)

	err = state.RemoveKubePromStack()
	assert.NoError(t, err)
}
