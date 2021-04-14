package hostedzone_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/janitor/hostedzone"
	"github.com/oslokommune/okctl/pkg/mock"
)

func TestUndelegatedZonesInHostedZones(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		hzID      string
		fn        hostedzone.NameServersFunc
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should return nil",
			provider: mock.NewGoodCloudProvider(),
			hzID:     "something",
			fn: func(_ string) ([]string, error) {
				return []string{"ns1.something.com"}, nil
			},
			expect: []*hostedzone.UndelegatedHostedZone(nil),
		},
		{
			name:     "Should return undelegated",
			provider: mock.NewGoodCloudProvider(),
			hzID:     "something",
			fn: func(_ string) ([]string, error) {
				return []string{}, nil
			},
			expect: []*hostedzone.UndelegatedHostedZone{
				{
					Name: "mine.oslo.systems",
					NameServers: []string{
						"ns1.something.com",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := hostedzone.New(tc.provider).UndelegatedZonesInHostedZones(tc.hzID, tc.fn)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
