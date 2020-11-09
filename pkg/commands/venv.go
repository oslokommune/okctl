package commands

import (
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/virtualenv"
	"os"
	"path"
)

// GetVirtualEnvironmentOpts returns data needed to set up a virtual environment.
func GetVirtualEnvironmentOpts(o *okctl.Okctl) (virtualenv.VirtualEnvironmentOpts, error) {
	meta := o.RepoStateWithEnv.GetMetadata()
	cluster := o.RepoStateWithEnv.GetCluster()

	userDataDir, err := o.GetUserDataDir()
	if err != nil {
		return virtualenv.VirtualEnvironmentOpts{}, err
	}

	k, err := o.BinariesProvider.Kubectl(kubectl.Version)
	if err != nil {
		return virtualenv.VirtualEnvironmentOpts{}, err
	}

	a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return virtualenv.VirtualEnvironmentOpts{}, err
	}

	environ := os.Environ()

	opts := virtualenv.VirtualEnvironmentOpts{
		Region:                 meta.Region,
		AWSAccountID:           cluster.AWSAccountID,
		Environment:            cluster.Environment,
		Repository:             meta.Name,
		ClusterName:            cluster.Name,
		UserDataDir:            userDataDir,
		Debug:                  o.Debug,
		KubectlBinaryDir:       path.Dir(k.BinaryPath),
		AwsIamAuthenticatorDir: path.Dir(a.BinaryPath),
		OsEnvVars:              environ,
	}

	err = opts.Validate()
	if err != nil {
		return virtualenv.VirtualEnvironmentOpts{}, errors.E(err, "failed to validate show credentials options")
	}

	return opts, nil
}

// GetVirtualEnvironmentOptsWithPs1 returns data needed to set up a virtual environment. The ps1Dir must be a path,
// as it will be added to the PATH later.
func GetVirtualEnvironmentOptsWithPs1(o *okctl.Okctl, ps1Dir string) (virtualenv.VirtualEnvironmentOpts, error) {
	opts, err := GetVirtualEnvironmentOpts(o)
	if err != nil {
		return virtualenv.VirtualEnvironmentOpts{}, err
	}

	opts.Ps1Dir = ps1Dir

	return opts, nil
}
