package s3api_test

import (
	"bytes"
	"testing"

	"github.com/oslokommune/okctl/pkg/s3api"
	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/mock"
)

func TestS3APIPutObject(t *testing.T) {
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
			expect:    "something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := s3api.New(tc.provider).PutObject("bucket", "key", bytes.NewReader([]byte{}))

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestS3APIDeleteObject(t *testing.T) {
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
			expect:    "calling delete object API: something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := s3api.New(tc.provider).DeleteObject("bucket", "key")

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestS3API_GetObject(t *testing.T) {
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
			expect:    "something bad",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			_, err := s3api.New(tc.provider).GetObject("bucket", "key")

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
