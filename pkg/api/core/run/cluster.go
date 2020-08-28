// Package run implements the runnable layer
package run

import (
	"fmt"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
)

type clusterRun struct {
	awsCredentialsPath string
	awsConfigPath      string
	provider           binaries.Provider
	debug              bool
}

// CreateCluster invokes a CLI for performing create
func (c *clusterRun) CreateCluster(kubeConfigPath string, config *api.ClusterConfig) error {
	a, err := c.provider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return errors.E(err, "failed to get aws-iam-authenticator cli", errors.Internal)
	}

	k, err := c.provider.Kubectl(kubectl.Version)
	if err != nil {
		return errors.E(err, "failed to get kubectl cli", errors.Internal)
	}

	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return errors.E(err, "failed to get eksctl cli", errors.Internal)
	}

	cli.Debug(c.debug)
	cli.AddToPath(a.BinaryPath, k.BinaryPath)
	cli.AddToEnv(
		fmt.Sprintf("AWS_CONFIG_FILE=%s", c.awsConfigPath),
		fmt.Sprintf("AWS_SHARED_CREDENTIALS_FILE=%s", c.awsCredentialsPath),
		"AWS_PROFILE=default",
	)

	exists, err := cli.HasCluster(config)
	if err != nil {
		return errors.E(err, "unable to determine if cluster exists")
	}

	if !exists {
		_, err = cli.CreateCluster(kubeConfigPath, config)
		if err != nil {
			return errors.E(err, "failed to create cluster", errors.Internal)
		}
	}

	return nil
}

// DeleteCluster invokes a CLI for performing delete
func (c *clusterRun) DeleteCluster(clusterName string) error {
	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return err
	}

	_, err = cli.DeleteCluster(clusterName)

	return err
}

// NewClusterRun returns a executor for clusterRun
func NewClusterRun(debug bool, awsCredentialsPath, awsConfigPath string, provider binaries.Provider) api.ClusterRun {
	return &clusterRun{
		debug:              debug,
		awsCredentialsPath: awsCredentialsPath,
		awsConfigPath:      awsConfigPath,
		provider:           provider,
	}
}
