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

func TestManagedPolicyStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.ManagedPolicy{})
	assert.NoError(t, err)

	state := storm.NewManagedPolicyState(db)

	err = state.SavePolicy(mock.ManagedPolicy())
	assert.NoError(t, err)

	m, err := state.GetPolicy(mock.StackNameManagedPolicy)
	assert.NoError(t, err)
	assert.Equal(t, mock.ManagedPolicy(), m)

	m.PolicyARN = "fake10"
	err = state.SavePolicy(m)
	assert.NoError(t, err)

	err = state.RemovePolicy(mock.StackNameManagedPolicy)
	assert.NoError(t, err)

	err = state.RemovePolicy(mock.StackNameManagedPolicy)
	assert.NoError(t, err)
}
