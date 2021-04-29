package arn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/arn"
)

func TestIs(t *testing.T) {
	testCases := []struct {
		name      string
		arn       string
		service   arn.ServiceType
		resource  arn.ResourceType
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should return true",
			arn:      mock.DefaultLoadBalancerARN,
			service:  arn.ServiceElasticLoadBalancing,
			resource: arn.ResourceLoadBalancer,
			expect:   true,
		},
		{
			name:     "Should return false",
			arn:      mock.DefaultListenerARN,
			service:  arn.ServiceElasticLoadBalancing,
			resource: arn.ResourceLoadBalancer,
			expect:   false,
		},
		{
			name:      "Should fail",
			arn:       "arn::ekfoj",
			service:   arn.ServiceElasticLoadBalancing,
			resource:  arn.ResourceLoadBalancer,
			expect:    "not a valid arn: arn::ekfoj",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := arn.Is(tc.service, tc.resource, tc.arn)

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
