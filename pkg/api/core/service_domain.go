package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type domainService struct {
	cloudProvider api.DomainCloudProvider
	store         api.DomainStore
}

func (d *domainService) CreateHostedZone(_ context.Context, opts api.CreateHostedZoneOpts) (*api.HostedZone, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate domain inputs")
	}

	domain, err := d.cloudProvider.CreateHostedZone(opts)
	if err != nil {
		return nil, errors.E(err, "failed to create domain")
	}

	err = d.store.SaveHostedZone(domain)
	if err != nil {
		return nil, errors.E(err, "failed to store domain")
	}

	return domain, nil
}

// NewDomainService returns an initialised domain service
func NewDomainService(cloudProvider api.DomainCloudProvider, store api.DomainStore) api.DomainService {
	return &domainService{
		cloudProvider: cloudProvider,
		store:         store,
	}
}
