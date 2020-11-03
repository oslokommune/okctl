package virtualenv

import (
	"fmt"
	"path"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/storage"
)

// GetVirtualEnvironmentOpts returns data needed to set up a virtual environment.
func GetVirtualEnvironmentOpts(o *okctl.Okctl) (VirtualEnvironmentOpts, error) {
	meta := o.RepoStateWithEnv.GetMetadata()
	cluster := o.RepoStateWithEnv.GetCluster()

	userDataDir, err := o.GetUserDataDir()
	if err != nil {
		return VirtualEnvironmentOpts{}, err
	}

	k, err := o.BinariesProvider.Kubectl(kubectl.Version)
	if err != nil {
		return VirtualEnvironmentOpts{}, err
	}

	a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return VirtualEnvironmentOpts{}, err
	}

	opts := VirtualEnvironmentOpts{
		Region:                 meta.Region,
		AWSAccountID:           cluster.AWSAccountID,
		Environment:            cluster.Environment,
		Repository:             meta.Name,
		ClusterName:            cluster.Name,
		UserDataDir:            userDataDir,
		Debug:                  o.Debug,
		KubectlBinaryDir:       path.Dir(k.BinaryPath),
		AwsIamAuthenticatorDir: path.Dir(a.BinaryPath),
	}

	err = opts.validate()
	if err != nil {
		return VirtualEnvironmentOpts{}, errors.E(err, "failed to validate show credentials options")
	}

	return opts, nil
}

// GetVirtualEnvironmentOptsWithPs1 returns data needed to set up a virtual environment. The ps1Dir must be a path,
// as it will be added to the PATH later.
func GetVirtualEnvironmentOptsWithPs1(o *okctl.Okctl, ps1Dir string) (VirtualEnvironmentOpts, error) {
	opts, err := GetVirtualEnvironmentOpts(o)
	if err != nil {
		return VirtualEnvironmentOpts{}, err
	}

	opts.Ps1Dir = ps1Dir

	return opts, nil
}

// CreatePs1ExecutableIfNotExists creates an executable file that returns "myenv:mynamespace", if it doesn't exist.
// The file will be called in the PS1 environment variable.
//
// This function returns the path to the directory containing the file.
func CreatePs1ExecutableIfNotExists(store storage.Storer) (string, error) {
	ps1Filename := "venv_ps1"
	ps1Dir := "venv"

	ps1FileExists, err := store.Exists(path.Join(ps1Dir, ps1Filename))
	if err != nil {
		return "", fmt.Errorf("couldn't create PS1 helper executable: %w", err)
	}

	if !ps1FileExists {
		ps1File, err := store.Create(ps1Dir, ps1Filename, 0o744)
		if err != nil {
			return "", err
		}

		_, err = ps1File.WriteString(`#!/usr/bin/env bash
ENV=$1
ENV=${ENV:-NOENV}

K8S_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}' 2>/dev/null)"
K8S_NAMESPACE="${K8S_NAMESPACE:-default}"

echo -e "$ENV:$K8S_NAMESPACE"
`)
		if err != nil {
			return "", err
		}

		err = ps1File.Close()
		if err != nil {
			return "", err
		}
	}

	return ps1Dir, nil
}
