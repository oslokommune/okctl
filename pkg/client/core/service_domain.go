package core

import (
	"context"
	"io"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type domainService struct {
	out    io.Writer
	ask    *ask.Ask
	api    client.DomainAPI
	store  client.DomainStore
	state  client.DomainState
	report client.DomainReport
}

func (s *domainService) CreatePrimaryHostedZone(_ context.Context, opts client.CreatePrimaryHostedZoneOpts) (*api.HostedZone, error) {
	zone, err := s.api.CreatePrimaryHostedZone(opts)
	if err != nil {
		return nil, err
	}

	hz := &client.HostedZone{
		Primary:    true,
		HostedZone: zone,
	}

	if !hz.IsDelegated {
		delegated, err := s.ask.ConfirmPostingNameServers(s.out, hz.HostedZone.Domain, hz.HostedZone.NameServers)
		if err != nil {
			return nil, err
		}

		hz.IsDelegated = delegated
	}

	r1, err := s.store.SaveHostedZone(hz)
	if err != nil {
		return nil, err
	}

	r2, err := s.state.SaveHostedZone(hz)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreatePrimaryHostedZone(hz, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return hz.HostedZone, nil
}

// NewDomainService returns an initialised service
func NewDomainService(
	out io.Writer,
	ask *ask.Ask,
	api client.DomainAPI,
	store client.DomainStore,
	report client.DomainReport,
	state client.DomainState,
) client.DomainService {
	return &domainService{
		api:    api,
		out:    out,
		store:  store,
		report: report,
		ask:    ask,
		state:  state,
	}
}
