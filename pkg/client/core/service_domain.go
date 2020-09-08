package core

import (
	"context"
	"io"

	"github.com/oslokommune/okctl/pkg/config/repository"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type domainService struct {
	api       client.DomainAPI
	out       io.Writer
	store     client.DomainStore
	repoState *repository.Data
	ask       *ask.Ask
}

func (s *domainService) CreatePrimaryHostedZone(_ context.Context, opts client.CreatePrimaryHostedZoneOpts) (*api.HostedZone, error) {
	hz, err := s.store.GetPrimaryHostedZone(opts.ID)
	if err != nil {
		return nil, err
	}

	if hz == nil {
		zone, err := s.api.CreatePrimaryHostedZone(opts)
		if err != nil {
			return nil, err
		}

		hz = &client.HostedZone{
			IsDelegated: false,
			Primary:     true,
			HostedZone:  zone,
		}
	}

	if !hz.IsDelegated {
		delegated, err := s.ask.ConfirmPostingNameServers(s.out, hz.HostedZone.Domain, hz.HostedZone.NameServers)
		if err != nil {
			return nil, err
		}

		hz.IsDelegated = delegated
	}

	_, err = s.store.SaveHostedZone(hz)
	if err != nil {
		return nil, err
	}

	return hz.HostedZone, nil
}

// NewDomainService returns an initialised service
func NewDomainService(out io.Writer, repoState *repository.Data, ask *ask.Ask, api client.DomainAPI, store client.DomainStore) client.DomainService {
	return &domainService{
		out:       out,
		ask:       ask,
		repoState: repoState,
		api:       api,
		store:     store,
	}
}
