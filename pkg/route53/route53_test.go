package route53_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/route53"
)

func TestRoute53PublicHostedZones(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			expect: []*route53.HostedZone{
				{
					ID:     "AABBCCDD",
					Domain: "test.oslo.systems",
					FQDN:   "test.oslo.systems.",
					Public: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := route53.New(tc.provider).PublicHostedZones()

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
