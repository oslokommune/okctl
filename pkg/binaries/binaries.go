// Package binaries knows how to load CLIs
package binaries

import (
	"io"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/storage"
)

// Provider defines the CLIs that are available
type Provider interface {
	Eksctl(version string) (*eksctl.Eksctl, error)
	Kubectl(version string) (*kubectl.Kubectl, error)
}

type provider struct {
	progress io.Writer
	auth     aws.Authenticator
	fetcher  fetch.Provider

	eksctl  map[string]*eksctl.Eksctl
	kubectl map[string]*kubectl.Kubectl
}

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

		p.eksctl[version] = eksctl.New(store, p.progress, binaryPath, p.auth, run.Cmd())
	}

	return p.eksctl[version], nil
}

// New returns a provider that knows how to fetch binaries and make
// them available for other commands
func New(progress io.Writer, auth aws.Authenticator, fetcher fetch.Provider) Provider {
	return &provider{
		progress: progress,
		auth:     auth,
		fetcher:  fetcher,
		eksctl:   map[string]*eksctl.Eksctl{},
		kubectl:  map[string]*kubectl.Kubectl{},
	}
}
