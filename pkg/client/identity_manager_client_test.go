package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

func TestCreateIdentityPoolUserOptsValidate(t *testing.T) {
	testCases := []struct {
		name      string
		opts      client.CreateIdentityPoolUserOpts
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Should work",
			opts: client.CreateIdentityPoolUserOpts{
				ID: api.ID{
					Region:       "eu-west-1",
					AWSAccountID: "123456789012",
					ClusterName:  "okctl-dev",
				},
				Email:      "someone.with-hyphen@origo.oslo.kommune.no",
				UserPoolID: "something",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteIdentityPoolUserOptsValidate(t *testing.T) {
	testCases := []struct {
		name      string
		opts      client.DeleteIdentityPoolUserOpts
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Should work",
			opts: client.DeleteIdentityPoolUserOpts{
				ClusterID: api.ID{
					Region:       "eu-west-1",
					AWSAccountID: "012345678912",
					ClusterName:  "test",
				},
				UserEmail: "test-user",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
