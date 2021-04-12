package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateValidApplication() OkctlApplication {
	return OkctlApplication{
		Name:            "test-app",
		Namespace:       "test",
		Image:           "testimage",
		Version:         "latest",
		ImagePullSecret: "sometoken",
		SubDomain:       "testapp",
		Port:            80,
		Replicas:        1,
		Environment:     nil,
		Volumes:         nil,
	}
}

// nolint: funlen
func TestOkctlApplicationValidation(t *testing.T) {
	testCases := []struct {
		name string

		withApplication func() OkctlApplication

		expectFail    bool
		expectedError string
	}{
		{
			name: "Sanity check",

			withApplication: generateValidApplication,

			expectFail: false,
		},
		{
			name: "Should allow images from Docker hub",

			withApplication: func() OkctlApplication {
				app := generateValidApplication()

				app.Image = "postgres"

				return app
			},

			expectFail: false,
		},
		{
			name: "Should allow images from GHCR",

			withApplication: func() OkctlApplication {
				app := generateValidApplication()

				app.Image = "ghcr.io/oslokommune/test-app"

				return app
			},

			expectFail: false,
		},
		{
			name: "Should allow images from ECR",

			withApplication: func() OkctlApplication {
				app := generateValidApplication()

				app.Image = "012345678912.dkr.ecr.eu-west-1.amazonaws.com/cluster-test-testapp"

				return app
			},

			expectFail: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.withApplication().Validate()

			if !tc.expectFail {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			}
		})
	}
}
