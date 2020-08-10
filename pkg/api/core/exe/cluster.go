// Package exe implements the exe layer
package exe

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
)

type cluster struct {
	awsCredentialsPath string
	awsConfigPath      string
	provider           binaries.Provider
	debug              bool
}

// CreateCluster invokes a CLI for performing create
func (c *cluster) CreateCluster(kubeConfigPath string, config *api.ClusterConfig) error {
	a, err := c.provider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return err
	}

	k, err := c.provider.Kubectl(kubectl.Version)
	if err != nil {
		return err
	}

	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return err
	}

	cli.Debug(c.debug)
	cli.AddToPath(a.BinaryPath, k.BinaryPath)
	cli.AddToEnv(
		fmt.Sprintf("AWS_CONFIG_FILE=%s", c.awsConfigPath),
		fmt.Sprintf("AWS_SHARED_CREDENTIALS_FILE=%s", c.awsCredentialsPath),
		"AWS_PROFILE=default",
	)

	_, err = cli.CreateCluster(kubeConfigPath, config)

	return err
}

// DeleteClusterConfig invokes a CLI for performing delete
func (c *cluster) DeleteCluster(config *api.ClusterConfig) error {
	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return err
	}

	_, err = cli.DeleteCluster(config)

	return err
}

// NewClusterExe returns a executor for cluster
func NewClusterExe(debug bool, awsCredentialsPath, awsConfigPath string, provider binaries.Provider) api.ClusterExe {
	return &cluster{
		debug:              debug,
		awsCredentialsPath: awsCredentialsPath,
		awsConfigPath:      awsConfigPath,
		provider:           provider,
	}
}
