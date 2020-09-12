package state_test

import (
	"log"
	"testing"

	"github.com/sanity-io/litter"

	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSaver(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "Should work",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			s := state.NewRepositoryStateWithEnv("test", state.NewRepository(), state.DefaultFileSystemSaver("something", "", &afero.Afero{
				Fs: afero.NewMemMapFs(),
			}))

			zone := state.HostedZone{
				IsDelegated: false,
				Primary:     true,
				Domain:      "test.oslo.systems",
				FQDN:        "huh",
				NameServers: []string{
					"umm",
					"ok",
				},
			}

			_, err := s.SaveHostedZone(zone.Domain, zone)
			assert.NoError(t, err)

			got := s.GetHostedZone(zone.Domain)
			assert.Equal(t, zone, got)

			log.Println(litter.Sdump(s.GetCluster()))
		})
	}
}
