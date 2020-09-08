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

func (s *domainService) CreateDomain(_ context.Context, opts api.CreateDomainOpts) (*api.Domain, error) {
	domain, err := s.api.CreateDomain(opts)
	if err != nil {
		return nil, err
	}

	_, err = s.store.SaveDomain(domain)
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
