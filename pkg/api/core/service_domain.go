package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type domainService struct {
	cloudProvider api.DomainCloudProvider
}

func (d *domainService) DeleteHostedZone(ctx context.Context, opts api.DeleteHostedZoneOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "failed to validate inputs")
	}

	err = d.cloudProvider.DeleteHostedZone(opts)
	if err != nil {
		return errors.E(err, "failed to delete hosted zone")
	}

	return nil
}

func (d *domainService) CreateHostedZone(_ context.Context, opts api.CreateHostedZoneOpts) (*api.HostedZone, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate hosted zone inputs")
	}

	domain, err := d.cloudProvider.CreateHostedZone(opts)
	if err != nil {
		return nil, errors.E(err, "failed to create hosted zone")
	}

	return domain, nil
}

// NewDomainService returns an initialised domain service
func NewDomainService(cloudProvider api.DomainCloudProvider) api.DomainService {
	return &domainService{
		cloudProvider: cloudProvider,
	}
}
