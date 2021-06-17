// Package binaries knows how to load CLIs
package binaries

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/binaries/run/okctlupgrade"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/storage"
)

// Provider defines the CLIs that are available
type Provider interface {
	Eksctl(version string) (*eksctl.Eksctl, error)
	Kubectl(version string) (*kubectl.Kubectl, error)
	AwsIamAuthenticator(version string) (*awsiamauthenticator.AwsIamAuthenticator, error)
	OkctlUpgrade(version string) (*okctlupgrade.OkctlUpgrade, error)
	ReloadBinaries() error
}

type provider struct {
	progress io.Writer
	auth     aws.Authenticator
	fetcher  fetch.Provider
	logger   *logrus.Logger

	eksctl              map[string]*eksctl.Eksctl
	kubectl             map[string]*kubectl.Kubectl
	awsIamAuthenticator map[string]*awsiamauthenticator.AwsIamAuthenticator
	okctlUpgrade        map[string]*okctlupgrade.OkctlUpgrade
}

// AwsIamAuthenticator returns an aws-iam-authenticator cli wrapper for running commands
func (p *provider) AwsIamAuthenticator(version string) (*awsiamauthenticator.AwsIamAuthenticator, error) {
	_, ok := p.awsIamAuthenticator[version]

	if !ok {
		binaryPath, err := p.fetcher.Fetch(awsiamauthenticator.Name, version)
		if err != nil {
			return nil, err
		}

		p.awsIamAuthenticator[version] = awsiamauthenticator.New(binaryPath)
	}

	return p.awsIamAuthenticator[version], nil
}

// Kubectl returns a kubectl cli wrapper for running commands
func (p *provider) Kubectl(version string) (*kubectl.Kubectl, error) {
	_, ok := p.kubectl[version]

	if !ok {
		binaryPath, err := p.fetcher.Fetch(kubectl.Name, version)
		if err != nil {
			return nil, err
		}

		p.kubectl[version] = kubectl.New(binaryPath)
	}

	return p.kubectl[version], nil
}

// Eksctl returns an eksctl cli wrapper for running commands
func (p *provider) Eksctl(version string) (*eksctl.Eksctl, error) {
	_, ok := p.eksctl[version]

	if !ok {
		binaryPath, err := p.fetcher.Fetch(eksctl.Name, version)
		if err != nil {
			return nil, err
		}

		store, err := storage.NewTemporaryStorage()
		if err != nil {
			return nil, err
		}

		p.eksctl[version] = eksctl.New(p.logger, store, p.progress, binaryPath, p.auth, run.Cmd())
	}

	return p.eksctl[version], nil
}

// OkctlUpgrade returns an okctl upgrade wrapper. The given version is the version of the okctl upgrade to run.
func (p *provider) OkctlUpgrade(version string) (*okctlupgrade.OkctlUpgrade, error) {
	_, ok := p.okctlUpgrade[version]

	if !ok {
		binaryName := fmt.Sprintf(okctlupgrade.BinaryNameFormat, version)

		binaryPath, err := p.fetcher.Fetch(binaryName, version)
		if err != nil {
			return nil, err
		}

		p.okctlUpgrade[version] = okctlupgrade.New(binaryPath)
	}

	return p.okctlUpgrade[version], nil
}

// ReloadBinaries pre-fetches binaries if possible, so that succeeding runs will not need fetching and thus saves
// time.
func (p *provider) ReloadBinaries() error {
	err := p.fetcher.Load()
	if err != nil {
		return fmt.Errorf("reloading binaries: %w", err)
	}

	return nil
}

// New returns a provider that knows how to fetch binaries and make
// them available for other commands
func New(logger *logrus.Logger, progress io.Writer, auth aws.Authenticator, fetcher fetch.Provider) Provider {
	return &provider{
		progress:            progress,
		auth:                auth,
		fetcher:             fetcher,
		logger:              logger,
		eksctl:              map[string]*eksctl.Eksctl{},
		kubectl:             map[string]*kubectl.Kubectl{},
		awsIamAuthenticator: map[string]*awsiamauthenticator.AwsIamAuthenticator{},
	}
}
