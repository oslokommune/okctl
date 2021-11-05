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

func TestExternalDNSStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New()

	db.SetDatabaseFilePath(filepath.Join(dir, "storm.db"))
	db.SetWritable(true)

	err = db.Init(&storm.ExternalDNS{})
	assert.NoError(t, err)

	state := storm.NewExternalDNSState(db)

	err = state.SaveExternalDNS(mock.ExternalDNS())
	assert.NoError(t, err)

	m, err := state.GetExternalDNS()
	assert.NoError(t, err)
	assert.Equal(t, mock.ExternalDNS(), m)

	m.Kube.HostedZoneID = "fake6"
	err = state.SaveExternalDNS(m)
	assert.NoError(t, err)

	err = state.RemoveExternalDNS()
	assert.NoError(t, err)

	err = state.RemoveExternalDNS()
	assert.NoError(t, err)
}
