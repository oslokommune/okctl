package acmapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/acmapi"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

func TestACMPI_InUseBy(t *testing.T) {
	testCases := []struct {
		name      string
		arn       string
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			arn:      mock.DefaultCertificateARN,
			provider: mock.NewGoodCloudProvider(),
			expect: []string{
				mock.DefaultLoadBalancerARN,
			},
			expectErr: false,
		},
		{
			name:      "Should fail",
			arn:       mock.DefaultCertificateARN,
			provider:  mock.NewBadCloudProvider(),
			expect:    "describing certificate: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := acmapi.New(tc.provider).InUseBy(tc.arn)

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
