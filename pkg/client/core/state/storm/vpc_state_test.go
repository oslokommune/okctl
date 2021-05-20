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

func TestVpcStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New(filepath.Join(dir, "storm.db"))

	err = db.Init(&storm.Vpc{})
	assert.NoError(t, err)

	state := storm.NewVpcState(db)

	err = state.SaveVpc(mock.Vpc())
	assert.NoError(t, err)

	m, err := state.GetVpc(mock.StackNameVpc)
	assert.NoError(t, err)
	assert.Equal(t, mock.Vpc(), m)

	m.VpcID = "fake"
	err = state.SaveVpc(m)
	assert.NoError(t, err)

	err = state.RemoveVpc(mock.StackNameVpc)
	assert.NoError(t, err)

	err = state.RemoveVpc(mock.StackNameVpc)
	assert.NoError(t, err)
}
