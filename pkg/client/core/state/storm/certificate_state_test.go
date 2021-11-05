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

func TestCertificateStateScenario(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	assert.NoError(t, err)

	defer func() {
		err = os.RemoveAll(dir)
		assert.NoError(t, err)
	}()

	db := breeze.New()

	db.SetDatabaseFilePath(filepath.Join(dir, "storm.db"))
	db.SetWritable(true)

	err = db.Init(&storm.Certificate{})
	assert.NoError(t, err)

	state := storm.NewCertificateState(db)

	err = state.SaveCertificate(mock.Certificate())
	assert.NoError(t, err)

	hz, err := state.GetCertificate(mock.DefaultDomain)
	assert.NoError(t, err)
	assert.Equal(t, mock.Certificate(), hz)

	hz.HostedZoneID = "fake2"
	err = state.SaveCertificate(hz)
	assert.NoError(t, err)

	err = state.RemoveCertificate(mock.DefaultDomain)
	assert.NoError(t, err)

	err = state.RemoveCertificate(mock.DefaultDomain)
	assert.NoError(t, err)
}
