package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateValidApplication() Application {
	app := NewApplication(Cluster{})

	app.Metadata.Name = "test-app"
	app.Metadata.Namespace = "test"

	app.Image.URI = "testimage:latest"
	app.ImagePullSecret = "sometoken"

	app.SubDomain = "testapp"
	app.Port = 80

	app.Replicas = 1
	app.Environment = nil
	app.Volumes = nil

	return app
}

// nolint: funlen
func TestApplicationValidation(t *testing.T) {
	testCases := []struct {
		name string

		withApplication func() Application

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

			withApplication: func() Application {
				app := generateValidApplication()

				app.Image.URI = "postgres"

				return app
			},

			expectFail: false,
		},
		{
			name: "Should allow images from GHCR",

			withApplication: func() Application {
				app := generateValidApplication()

				app.Image.URI = "ghcr.io/oslokommune/test-app"

				return app
			},

			expectFail: false,
		},
		{
			name: "Should allow images from ECR",

			withApplication: func() Application {
				app := generateValidApplication()

				app.Image.URI = "012345678912.dkr.ecr.eu-west-1.amazonaws.com/cluster-test-testapp"

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
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			}
		})
	}
}
