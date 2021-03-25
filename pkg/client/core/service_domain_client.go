package core

import (
	"context"
	"errors"
	"io"

	"github.com/oslokommune/okctl/pkg/client/core/state/storm"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/client"
)

type domainService struct {
	spinner spinner.Spinner
	out     io.Writer
	ask     *ask.Ask
	api     client.DomainAPI
	store   client.DomainStore
	state   client.DomainState
	report  client.DomainReport
}

func (s *domainService) GetPrimaryHostedZone(_ context.Context) (*client.HostedZone, error) {
	return s.state.GetPrimaryHostedZone()
}

// DeletePrimaryHostedZone and all associated records
func (s *domainService) DeletePrimaryHostedZone(_ context.Context, opts client.DeletePrimaryHostedZoneOpts) error {
	err := s.spinner.Start("domain")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	hz, err := s.state.GetPrimaryHostedZone()
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil
		}

		return err
	}

	var reports []*store.Report

	opts.HostedZoneID = hz.HostedZoneID

	if hz.Managed {
		// HostedZone is managed by us, so delete it
		err = s.api.DeletePrimaryHostedZone(hz.Domain, opts)
		if err != nil {
			return err
		}

		report, err := s.store.RemoveHostedZone(hz.Domain)
		if err != nil {
			return err
		}

		reports = append(reports, report)
	}

	err = s.state.RemoveHostedZone(hz.Domain)
	if err != nil {
		return err
	}

	err = s.report.ReportDeletePrimaryHostedZone(reports)
	if err != nil {
		return err
	}

	return nil
}

func (s *domainService) CreatePrimaryHostedZone(_ context.Context, opts client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	err := s.spinner.Start("domain")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	primary, err := s.state.GetPrimaryHostedZone()
	if err == nil {
		return primary, nil
	}

	zone, err := s.api.CreatePrimaryHostedZone(opts)
	if err != nil {
		return nil, err
	}

	zone.IsDelegated = false

	r1, err := s.store.SaveHostedZone(zone)
	if err != nil {
		return nil, err
	}

	err = s.state.SaveHostedZone(zone)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreatePrimaryHostedZone(zone, []*store.Report{r1})
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *domainService) SetHostedZoneDelegation(_ context.Context, domain string, isDelegated bool) (err error) {
	hz, err := s.state.GetHostedZone(domain)
	if err != nil {
		return err
	}

	hz.IsDelegated = isDelegated

	err = s.state.UpdateHostedZone(hz)
	if err != nil {
		return err
	}

	err = s.report.ReportHostedZoneDelegation(hz)
	if err != nil {
		return err
	}

	return nil
}

// NewDomainService returns an initialised service
func NewDomainService(
	spinner spinner.Spinner,
	out io.Writer,
	ask *ask.Ask,
	api client.DomainAPI,
	store client.DomainStore,
	report client.DomainReport,
	state client.DomainState,
) client.DomainService {
	return &domainService{
		spinner: spinner,
		api:     api,
		out:     out,
		store:   store,
		report:  report,
		ask:     ask,
		state:   state,
	}
}
