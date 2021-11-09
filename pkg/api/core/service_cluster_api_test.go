package core_test

import (
	"context"
	"testing"

	awsmock "github.com/oslokommune/okctl/pkg/mock"

	"github.com/google/go-cmp/cmp"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/stretchr/testify/assert"
)

func TestClusterCreateCluster(t *testing.T) {
	testCases := []struct {
		name        string
		service     api.ClusterService
		opts        api.ClusterCreateOpts
		expect      interface{}
		expectError bool
	}{
		{
			name: "Validate request",
			service: core.NewClusterService(
				mock.NewGoodClusterExe(),
				awsmock.NewGoodCloudProvider(),
			),
			opts:   mock.DefaultClusterCreateOpts(),
			expect: mock.DefaultCluster(),
		},
		{
			name: "Invalid opts",
			service: core.NewClusterService(
				mock.NewGoodClusterExe(),
				awsmock.NewGoodCloudProvider(),
			),
			opts: api.ClusterCreateOpts{},
			//nolint: lll
			expect:      "validating inputs: Cidr: cannot be blank; ID: (AWSAccountID: cannot be blank; ClusterName: cannot be blank; Region: cannot be blank.); Version: cannot be blank; VpcID: cannot be blank; VpcPrivateSubnets: cannot be blank; VpcPublicSubnets: cannot be blank.",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.service.CreateCluster(context.Background(), tc.opts)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Empty(t, cmp.Diff(tc.expect, got))
			}
		})
	}
}

func TestGetClusterSecurityGroupId(t *testing.T) {
	testCases := []struct {
		name        string
		service     api.ClusterService
		opts        api.ClusterSecurityGroupIDGetOpts
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			service: core.NewClusterService(
				mock.NewGoodClusterExe(),
				awsmock.NewGoodCloudProvider(),
			),
			opts:   mock.DefaultClusterSecurityGroupIDGetOpts(),
			expect: mock.DefaultClusterSecurityGroupID(),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			clusterSecurityGroupID, err := tc.service.GetClusterSecurityGroupID(context.Background(), tc.opts)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Empty(t, cmp.Diff(tc.expect, clusterSecurityGroupID))
			}
		})
	}
}
