package core

import (
	"errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/spf13/afero"
	"gotest.tools/assert"
	"testing"
)

// nolint: funlen
func TestBinaryService(t *testing.T) {
	t.Run("Should contain known binaries", func(t *testing.T) {
		// Given
		testHelper := newBinaryServiceTestHelper(t)

		// Then
		expectedBinaries := append(state.KnownBinaries())
		testHelper.assertBinaryServiceContains(expectedBinaries)
		testHelper.assertPersistedConfigContains(expectedBinaries)
	})

	t.Run("Should add two binaries", func(t *testing.T) {
		// Given
		var err error
		testHelper := newBinaryServiceTestHelper(t)

		b1 := state.Binary{Name: "my-binary", Version: "1.0.0"}
		b2 := state.Binary{Name: "my-other-binary", Version: "2.0.5"}

		// When
		err = testHelper.binaryService.Add(b1)
		assert.NilError(t, err)

		err = testHelper.binaryService.Add(b2)
		assert.NilError(t, err)

		// Then
		expectedBinaries := append(state.KnownBinaries(), b1, b2)
		testHelper.assertBinaryServiceContains(expectedBinaries)
		testHelper.assertPersistedConfigContains(expectedBinaries)
	})

	t.Run("Should add and remove binaries", func(t *testing.T) {
		// Given
		var err error
		testHelper := newBinaryServiceTestHelper(t)

		b1 := state.Binary{Name: "my-binary", Version: "1.0.0"}
		b2 := state.Binary{Name: "my-other-binary", Version: "2.0.5"}

		// When
		err = testHelper.binaryService.Add(b1)
		assert.NilError(t, err)

		err = testHelper.binaryService.Add(b2)
		assert.NilError(t, err)

		err = testHelper.binaryService.Remove(b2)
		assert.NilError(t, err)

		// Then
		expectedBinaries := append(state.KnownBinaries(), b1)
		testHelper.assertBinaryServiceContains(expectedBinaries)
		testHelper.assertPersistedConfigContains(expectedBinaries)
	})

	t.Run("Should return error if trying to remove a binary that isn't stored", func(t *testing.T) {
		// Given
		var err error
		testHelper := newBinaryServiceTestHelper(t)

		b1 := state.Binary{Name: "my-binary", Version: "1.0.0"}

		// When
		err = testHelper.binaryService.Remove(b1)

		// Then
		assert.ErrorContains(t, err, "binary my-binary-1.0.0 not found")
	})
}

type binaryServiceTestHelper struct {
	t             *testing.T
	binaryService client.BinaryService
	cfg           *config.Config
}

func newBinaryServiceTestHelper(t *testing.T) binaryServiceTestHelper {
	h := binaryServiceTestHelper{
		t: t,
	}
	h.init()

	return h
}

func (h *binaryServiceTestHelper) init() {
	h.cfg = config.New()
	h.cfg.Context.FileSystem = &afero.Afero{Fs: afero.NewMemMapFs()}

	h.cfg.UserDataLoader = load.BuildUserDataLoader(
		load.LoadDefaultUserData,
	)

	err := h.cfg.LoadUserData()
	assert.NilError(h.t, err)

	err = h.cfg.WriteCurrentUserData()
	assert.NilError(h.t, err)

	h.binaryService = NewBinaryService(h.cfg)
}

func (h *binaryServiceTestHelper) loadConfig() *config.Config {
	notFoundFn := func(c *config.Config) error {
		return errors.New("should not happen")
	}

	cfg := config.New()
	cfg.Context.FileSystem = h.cfg.FileSystem
	cfg.UserDataLoader = load.BuildUserDataLoader(
		load.LoadDefaultUserData,
		load.LoadStoredUserData(notFoundFn),
	)

	err := cfg.LoadUserData()
	assert.NilError(h.t, err)

	return cfg
}

// assertBinaryServiceContains asserts that binary service contains the expected binaries
func (h *binaryServiceTestHelper) assertBinaryServiceContains(expectedBinaries []state.Binary) {
	binaryList := h.binaryService.List()
	assert.Equal(h.t, len(expectedBinaries), len(binaryList))

	for i, b := range binaryList {
		assert.Equal(h.t, expectedBinaries[i].Id(), b.Id())
	}
}

// assertPersistedConfigContains asserts that the persisted configuration contains the expected binaries
func (h *binaryServiceTestHelper) assertPersistedConfigContains(expectedBinaries []state.Binary) {
	cfg := h.loadConfig()

	assert.Equal(h.t, len(expectedBinaries), len(cfg.UserState.Binaries))

	for i, b := range cfg.UserState.Binaries {
		assert.Equal(h.t, expectedBinaries[i].Id(), b.Id())
	}
}
