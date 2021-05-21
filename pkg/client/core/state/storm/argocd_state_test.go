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

func TestArgoCDStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.ArgoCD{})
	assert.NoError(t, err)

	state := storm.NewArgoCDState(db)

	err = state.SaveArgoCD(mock.ArgoCD())
	assert.NoError(t, err)

	m, err := state.GetArgoCD()
	assert.NoError(t, err)

	a := mock.ArgoCD()

	a.Certificate = nil
	a.ClientSecret = nil
	a.Secret = nil
	a.IdentityClient = nil
	a.PrivateKey = nil
	a.SecretKey = nil
	a.Chart = nil

	assert.Equal(t, a, m)

	a.ArgoURL = "fake1"
	err = state.SaveArgoCD(a)
	assert.NoError(t, err)

	err = state.RemoveArgoCD()
	assert.NoError(t, err)

	err = state.RemoveArgoCD()
	assert.NoError(t, err)
}
