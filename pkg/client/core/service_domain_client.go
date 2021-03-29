package core

import (
	"context"
	"errors"

	"github.com/oslokommune/okctl/pkg/client/core/state/storm"

	"github.com/oslokommune/okctl/pkg/client"
)

type domainService struct {
	api   client.DomainAPI
	store client.DomainStore
	state client.DomainState
}

func (s *domainService) GetPrimaryHostedZone(_ context.Context) (*client.HostedZone, error) {
	return s.state.GetPrimaryHostedZone()
}

// DeletePrimaryHostedZone and all associated records
func (s *domainService) DeletePrimaryHostedZone(_ context.Context, opts client.DeletePrimaryHostedZoneOpts) error {
	hz, err := s.state.GetPrimaryHostedZone()
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil
		}

		return err
	}

	opts.HostedZoneID = hz.HostedZoneID

	if hz.Managed {
		// HostedZone is managed by us, so delete it
		err = s.api.DeletePrimaryHostedZone(hz.Domain, opts)
		if err != nil {
			return err
		}

		_, err := s.store.RemoveHostedZone(hz.Domain)
		if err != nil {
			return err
		}
	}

	err = s.state.RemoveHostedZone(hz.Domain)
	if err != nil {
		return err
	}

	return nil
}

func (s *domainService) CreatePrimaryHostedZone(_ context.Context, opts client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	// [Refactor] Reconciler is responsible for ordering operations
	//
	// We should be doing this check in the reconciler together with a
	// verification towards the AWS API. Keeping this here for the
	// time being, so we are compatible with expected behavior.
	{
		p, err := s.state.GetPrimaryHostedZone()
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			return nil, err
		}

		if err == nil {
			return p, nil
		}
	}

	zone, err := s.api.CreatePrimaryHostedZone(opts)
	if err != nil {
		return nil, err
	}

	zone.IsDelegated = false

	_, err = s.store.SaveHostedZone(zone)
	if err != nil {
		return nil, err
	}

	err = s.state.SaveHostedZone(zone)
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

	return nil
}

// NewDomainService returns an initialised service
func NewDomainService(
	api client.DomainAPI,
	store client.DomainStore,
	state client.DomainState,
) client.DomainService {
	return &domainService{
		api:   api,
		store: store,
		state: state,
	}
}
