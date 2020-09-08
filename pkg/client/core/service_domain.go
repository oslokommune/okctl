package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type domainService struct {
	api   client.DomainAPI
	store client.DomainStore
}

func (s *domainService) CreateHostedZone(_ context.Context, opts api.CreateHostedZoneOpts) (*api.HostedZone, error) {
	domain, err := s.api.CreateHostedZone(opts)
	if err != nil {
		return nil, err
	}

	_, err = s.store.SaveHostedZone(domain)
	if err != nil {
		return nil, err
	}

	return domain, nil
}

// NewDomainService returns an initialised service
func NewDomainService(api client.DomainAPI, store client.DomainStore) client.DomainService {
	return &domainService{
		api:   api,
		store: store,
	}
}
