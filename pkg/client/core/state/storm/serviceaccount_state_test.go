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

func TestServiceAccountStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db, err := stormpkg.Open(filepath.Join(dir, "storm.db"), stormpkg.Codec(json.Codec))
	assert.NoError(t, err)

	err = db.Init(&storm.ServiceAccount{})
	assert.NoError(t, err)

	state := storm.NewServiceAccountState(db)

	err = state.SaveServiceAccount(mock.ServiceAccount())
	assert.NoError(t, err)

	sa, err := state.GetServiceAccount(mock.DefaultServiceAccountName)
	assert.NoError(t, err)
	assert.Equal(t, mock.ServiceAccount(), sa)

	sa.PolicyArn = "arn:aws:iam::123456789012:policy/new-policy"
	err = state.UpdateServiceAccount(sa)
	assert.NoError(t, err)

	err = state.RemoveServiceAccount(mock.DefaultServiceAccountName)
	assert.NoError(t, err)

	err = db.Close()
	assert.NoError(t, err)
}
