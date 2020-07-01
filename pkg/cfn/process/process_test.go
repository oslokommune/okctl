package process_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/process"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestSubnets(t *testing.T) {
	testCases := []struct {
		name        string
		provider    v1alpha1.CloudProvider
		value       string
		expect      interface{}
		expectError bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			value:    "subnet-0bb1c79de3EXAMPLE",
			expect: map[string]v1alpha1.ClusterNetwork{
				"eu-west-1a": {
					ID:   "subnet-0bb1c79de3EXAMPLE",
					CIDR: "192.168.0.0/24",
				},
			},
		},
		{
			name:        "Should fail",
			provider:    mock.NewBadCloudProvider(),
			value:       "subnet-0bb1c79de3EXAMPLE",
			expect:      "failed to describe subnet outputs: something bad",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := map[string]v1alpha1.ClusterNetwork{}

			err := process.Subnets(tc.provider, got)(tc.value)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
