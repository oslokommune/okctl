// Package binaries knows how to load CLIs
package binaries

import (
	"io"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/credentials"
)

// Provider defines the CLIs that are available
type Provider interface {
	Eksctl(version string) (*eksctl.Eksctl, error)
}

type provider struct {
	progress    io.Writer
	credentials credentials.Provider
	fetcher     fetch.Provider
	eksctl      map[string]*eksctl.Eksctl
}

// Eksctl returns an eksctl cli wrapper for running commands
func (p *provider) Eksctl(version string) (*eksctl.Eksctl, error) {
	_, ok := p.eksctl[version]

	if !ok {
		binaryPath, err := p.fetcher.Fetch(eksctl.Name, version)
		if err != nil {
			return nil, err
		}

		envs, err := p.credentials.AwsEnv()
		if err != nil {
			return nil, err
		}

		e, err := eksctl.New(p.progress, binaryPath, envs)
		if err != nil {
			return nil, err
		}

		p.eksctl[version] = e
	}

	return p.eksctl[version], nil
}

// New returns a provider that knows how to fetch binaries and make
// them available for other commands
func New(progress io.Writer, credentials credentials.Provider, fetcher fetch.Provider) Provider {
	return &provider{
		progress:    progress,
		credentials: credentials,
		fetcher:     fetcher,
		eksctl:      map[string]*eksctl.Eksctl{},
	}
}
