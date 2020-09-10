package core

import (
	"context"
	"io"

	"github.com/theckman/yacspin"

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type domainService struct {
	api       client.DomainAPI
	out       io.Writer
	store     client.DomainStore
	report    client.DomainReport
	repoState *state.Repository
	ask       *ask.Ask
	spinner   *yacspin.Spinner
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
		err = s.spinner.Pause()
		if err != nil {
			return nil, err
		}

		delegated, err := s.ask.ConfirmPostingNameServers(s.out, hz.HostedZone.Domain, hz.HostedZone.NameServers)
		if err != nil {
			return nil, err
		}

		err = s.spinner.Unpause()
		if err != nil {
			return nil, err
		}

		hz.IsDelegated = delegated
	}

	report, err := s.store.SaveHostedZone(hz)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreatePrimaryHostedZone(hz, report)
	if err != nil {
		return nil, err
	}

	return hz.HostedZone, nil
}

// NewDomainService returns an initialised service
func NewDomainService(
	out io.Writer,
	repoState *state.Repository,
	ask *ask.Ask,
	api client.DomainAPI,
	store client.DomainStore,
	report client.DomainReport,
	spinner *yacspin.Spinner,
) client.DomainService {
	return &domainService{
		out:       out,
		ask:       ask,
		repoState: repoState,
		api:       api,
		store:     store,
		report:    report,
		spinner:   spinner,
	}
}
