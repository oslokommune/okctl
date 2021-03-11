package eksapi_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/eksapi"
	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/mock"
)

func TestEKSAPIFargateProfilePodExecutionRoleARN(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
	}{
		{
			name:      "Should work",
			provider:  mock.NewGoodCloudProvider(),
			expect:    mock.DefaultFargateProfilePodExecutionRoleARN,
			expectErr: false,
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			expect:    "getting fargate profile: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := eksapi.New("something", tc.provider).FargateProfilePodExecutionRoleARN("something")

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
