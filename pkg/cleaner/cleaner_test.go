package cleaner_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cleaner"
	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

func TestCleaner_RemoveThingsThatAreUsingCertificate(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		arn       string
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			arn:      mock.DefaultCertificateARN,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := cleaner.New(tc.provider).RemoveThingsThatAreUsingCertificate(tc.arn)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCleaner_RemoveThingsUsingCertForDomain(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		domain    string
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			domain:   mock.DefaultDomain,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := cleaner.New(tc.provider).RemoveThingsUsingCertForDomain(tc.domain)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCleaner_DeleteDanglingTargetGroups(t *testing.T) {
	testCases := []struct {
		name        string
		provider    v1alpha1.CloudProvider
		clusterName string
		expect      interface{}
		expectErr   bool
	}{
		{
			name:        "Should work",
			provider:    mock.NewGoodCloudProvider(),
			clusterName: mock.DefaultClusterName,
		},
		{
			name:        "Should fail",
			provider:    mock.NewBadCloudProvider(),
			clusterName: mock.DefaultClusterName,
			expect:      "listing target groups: something bad",
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := cleaner.New(tc.provider).DeleteDanglingTargetGroups(tc.clusterName)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCleaner_DeleteDanglingSecurityGroups(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		vpcID     string
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			vpcID:    mock.DefaultVpcID,
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			vpcID:     mock.DefaultVpcID,
			expect:    "listing security groups for vpc: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := cleaner.New(tc.provider).DeleteDanglingSecurityGroups(tc.vpcID)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCleaner_DeleteDanglingALBs(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		vpcID     string
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
			vpcID:    mock.DefaultVpcID,
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			vpcID:     mock.DefaultVpcID,
			expect:    "listing load balancers: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := cleaner.New(tc.provider).DeleteDanglingALBs(tc.vpcID)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
