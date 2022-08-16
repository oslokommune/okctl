package state_test

import (
	"bytes"
	"testing"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestBinaryDigests(t *testing.T) {
	// The purpose of this test is to test that the checksums are correct. We don't want to run it
	// as part of the test suite, as we don't want to depend on the files existing.
	t.Skipf("Run this test manually to run it.")

	testCases := []struct {
		name string
	}{
		{
			name: "All digests should be correct",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var stdErr bytes.Buffer
			for _, os := range []string{state.OsDarwin, state.OsLinux} {
				_, err := fetch.New(
					&stdErr,
					logrus.New(),
					true,
					state.Host{
						Os:   os,
						Arch: "amd64",
					},
					state.KnownBinaries(),
					storage.NewEphemeralStorage(),
				)
				if err != nil {
					t.Log(stdErr.String())
				}
				require.NoError(t, err)
			}

			t.Log(stdErr.String())
		})
	}
}
