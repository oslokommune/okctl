// Package run implements the runnable layer
package run

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
)

type serviceAccountRun struct {
	awsCredentialsPath string
	awsConfigPath      string
	provider           binaries.Provider
	debug              bool
}

func (r *serviceAccountRun) getCli() (*eksctl.Eksctl, error) {
	a, err := r.provider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return nil, err
	}

	k, err := r.provider.Kubectl(kubectl.Version)
	if err != nil {
		return nil, err
	}

	cli, err := r.provider.Eksctl(eksctl.Version)
	if err != nil {
		return nil, err
	}

	cli.Debug(r.debug)
	cli.AddToPath(a.BinaryPath, k.BinaryPath)
	cli.AddToEnv(
		fmt.Sprintf("AWS_CONFIG_FILE=%s", r.awsConfigPath),
		fmt.Sprintf("AWS_SHARED_CREDENTIALS_FILE=%s", r.awsCredentialsPath),
		"AWS_PROFILE=default",
	)

	return cli, nil
}

func (r *serviceAccountRun) DeleteServiceAccount(config *v1alpha1.ClusterConfig) error {
	cli, err := r.getCli()
	if err != nil {
		return err
	}

	_, err = cli.DeleteServiceAccount(config)
	if err != nil {
		return err
	}

	return nil
}

func (r *serviceAccountRun) CreateServiceAccount(config *v1alpha1.ClusterConfig) error {
	cli, err := r.getCli()
	if err != nil {
		return err
	}

	_, err = cli.CreateServiceAccount(config)
	if err != nil {
		return err
	}

	return nil
}

// NewServiceAccountRun returns a runner for creating a service account
func NewServiceAccountRun(debug bool, awsCredentialsPath, awsConfigPath string, provider binaries.Provider) api.ServiceAccountRun {
	return &serviceAccountRun{
		debug:              debug,
		awsCredentialsPath: awsCredentialsPath,
		awsConfigPath:      awsConfigPath,
		provider:           provider,
	}
}
