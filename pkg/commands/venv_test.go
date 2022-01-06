package commands

import (
	"fmt"
	"testing"

	"github.com/oslokommune/okctl/pkg/context"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/stretchr/testify/assert"
)

func TestCleanOsEnvVars(t *testing.T) {
	testCases := []struct {
		name            string
		withVariables   []string
		expectVariables []string
	}{
		{
			name: "Should do nothing with no blacklisted variables",

			withVariables:   []string{"PATH=/bin", "EDITOR=vim"},
			expectVariables: []string{"PATH=/bin", "EDITOR=vim"},
		},
		{
			name: "Should remove blacklisted variables",

			withVariables: []string{
				"PATH=/bin",
				fmt.Sprintf("%s=somevalue", constant.EnvClusterDeclaration),
				"EDITOR=vim",
			},
			expectVariables: []string{"EDITOR=vim", "PATH=/bin"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result := CleanOsEnvVars(tc.withVariables)

			assert.Equal(t, len(tc.expectVariables), len(result))

			for _, expectedItem := range tc.expectVariables {
				assert.True(t, contains(result, expectedItem))
			}
		})
	}
}

func TestGetOkctlEnvVarsForAccessKey(t *testing.T) {
	okctlEnv := OkctlEnvironment{
		AWSCredentialsType: context.AWSCredentialsTypeAccessKey,
		UserDataDir:        "/home/user/.okctl/",
		ClusterName:        "testcluster",
	}

	envVars, err := GetVenvEnvVars(okctlEnv)

	assert.Nil(t, err, "should not return error")
	assert.Equal(t, "default", envVars["AWS_PROFILE"])
	assert.Equal(t, "/home/user/.okctl/credentials/testcluster/aws-config", envVars["AWS_CONFIG_FILE"])
	assert.Equal(t, "/home/user/.okctl/credentials/testcluster/aws-credentials", envVars["AWS_SHARED_CREDENTIALS_FILE"])
}

func TestGetOkctlEnvVarsForAwsProfile(t *testing.T) {
	okctlEnv := OkctlEnvironment{
		AWSCredentialsType: context.AWSCredentialsTypeAwsProfile,
		AwsProfile:         "testprofile",
		UserHomeDir:        "/home/user",
	}

	envVars, err := GetOkctlEnvVars(okctlEnv)

	assert.Nil(t, err, "should not return error")
	assert.Equal(t, "testprofile", envVars["AWS_PROFILE"])
	assert.Equal(t, "/home/user/.aws/config", envVars["AWS_CONFIG_FILE"])
	assert.Equal(t, "/home/user/.aws/credentials", envVars["AWS_SHARED_CREDENTIALS_FILE"])
	assert.Empty(t, envVars["AWS_ACCESS_KEY_ID"], "access key id should not be set")
	assert.Empty(t, envVars["AWS_SECRET_ACCESS_KEY"], "secret key should not be set")
}

func TestGetOkctlEnvVarMissingAwsProfile(t *testing.T) {
	okctlEnv := OkctlEnvironment{
		AWSCredentialsType: context.AWSCredentialsTypeAwsProfile,
	}

	envVars, err := GetOkctlEnvVars(okctlEnv)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "environment variable AWS_PROFILE not set")
	assert.Nil(t, envVars)
}

func contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}
