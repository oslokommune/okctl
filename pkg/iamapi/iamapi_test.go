package iamapi_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/iamapi"
	"github.com/stretchr/testify/assert"
)

func TestRoleFriendlyName(t *testing.T) {
	testCases := []struct {
		name      string
		roleARN   string
		expect    interface{}
		expectErr bool
	}{
		{
			name:    "Should work",
			roleARN: "arn:aws:iam::123456789012:role/eksctl-okctl-test-cluster-FargatePodExecutionRole-GHUSFFAKE",
			expect:  "eksctl-okctl-test-cluster-FargatePodExecutionRole-GHUSFFAKE",
		},
		{
			name:      "Should fail",
			roleARN:   "notanarn",
			expect:    "getting role friendly name: arn: invalid prefix",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := iamapi.RoleFriendlyName(tc.roleARN)

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

func TestIAMAPIAttachRolePolicy(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		policyARN string
		roleARN   string
		expect    interface{}
		expectErr bool
	}{
		{
			name:      "Should work",
			provider:  mock.NewGoodCloudProvider(),
			policyARN: "arn:iam:::policy/somethingFAHJFAKE",
			roleARN:   mock.DefaultFargateProfilePodExecutionRoleARN,
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			policyARN: "arn:iam:::policy/somethingFAHJFAKE",
			roleARN:   mock.DefaultFargateProfilePodExecutionRoleARN,
			expect:    "attaching policy to role: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := iamapi.New(tc.provider).AttachRolePolicy(tc.policyARN, tc.roleARN)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIAMAPIDetachRolePolicy(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		policyARN string
		roleARN   string
		expect    interface{}
		expectErr bool
	}{
		{
			name:      "Should work",
			provider:  mock.NewGoodCloudProvider(),
			policyARN: "arn:iam:::policy/somethingFAHJFAKE",
			roleARN:   mock.DefaultFargateProfilePodExecutionRoleARN,
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			policyARN: "arn:iam:::policy/somethingFAHJFAKE",
			roleARN:   mock.DefaultFargateProfilePodExecutionRoleARN,
			expect:    "detaching policy from role: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := iamapi.New(tc.provider).DetachRolePolicy(tc.policyARN, tc.roleARN)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
