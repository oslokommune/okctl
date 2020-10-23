package virtualenv_test

import (
	"fmt"
	"path"
	"strconv"
	"testing"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/virtualenv"
	"github.com/stretchr/testify/assert"
)

func TestGetVirtualEnvironment(t *testing.T) {
	t.Run("should return expected environment variables", func(t *testing.T) {
		opts := virtualenv.VirtualEnvironmentOpts{
			UserDataDir:            mock.DefaultUserDataDir,
			Debug:                  false,
			Region:                 mock.DefaultRegion,
			AWSAccountID:           mock.DefaultAWSAccountID,
			Environment:            mock.DefaultEnv,
			Repository:             mock.DefaultRepositoryName,
			ClusterName:            mock.DefaultClusterName,
			KubectlBinaryDir:       mock.DefaultKubectlBinaryDir,
			AwsIamAuthenticatorDir: mock.DefaultAwsIamAuthenticatorDir,
		}

		userPath := "/usr/local/go/bin:/usr/local/sbin:/home/johndoe/go/bin"
		osEnvVars := []string{
			"SOME_VAR=A",
			"PATH=" + userPath,
			"LS_COLORS=rs=0:di=01;34:ln=01:*.tar=01;31",
		}

		venv, err := virtualenv.GetVirtualEnvironment(&opts, osEnvVars)
		assert.Nil(t, err)

		expectedVenv := []string{
			"AWS_CONFIG_FILE=" + path.Join(opts.UserDataDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsConfig),
			"AWS_PROFILE=default",
			"AWS_SHARED_CREDENTIALS_FILE=" + path.Join(opts.UserDataDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsCredentials),
			"HELM_CACHE_HOME=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir),
			"HELM_CONFIG_DATA_HOME=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir),
			"HELM_CONFIG_HOME=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir),
			"HELM_DEBUG=" + strconv.FormatBool(opts.Debug),
			"HELM_PLUGINS=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
			"HELM_REGISTRY_CONFIG=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
			"HELM_REPOSITORY_CACHE=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
			"HELM_REPOSITORY_CONFIG=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
			"KUBECONFIG=" + path.Join(opts.UserDataDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterKubeConfig),
			"LS_COLORS=rs=0:di=01;34:ln=01:*.tar=01;31",
			fmt.Sprintf("PATH=%s:%s:%s",
				mock.DefaultKubectlBinaryDir,
				mock.DefaultAwsIamAuthenticatorDir,
				userPath,
			),
			"SOME_VAR=A",
		}

		assert.Equal(t, expectedVenv, venv)
	})
}
