package smapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/smapi"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/mock"
)

func TestRotateSecret(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			expect:    "creating secret rotation: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := smapi.New(tc.provider).RotateSecret("arn:::lambda/function", "arn::::secret/something")

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCancelRotateSecret(t *testing.T) {
	testCases := []struct {
		name      string
		provider  v1alpha1.CloudProvider
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Should work",
			provider: mock.NewGoodCloudProvider(),
		},
		{
			name:      "Should fail",
			provider:  mock.NewBadCloudProvider(),
			expect:    "canceling secret rotation: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := smapi.New(tc.provider).CancelRotateSecret("arn:::lambda/function")

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
