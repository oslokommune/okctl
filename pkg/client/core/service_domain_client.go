package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"

	stormpkg "github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/client"
)

type domainService struct {
	service api.DomainService
	state   client.DomainState
}

func (s *domainService) GetPrimaryHostedZone(_ context.Context) (*client.HostedZone, error) {
	return s.state.GetPrimaryHostedZone()
}

// DeletePrimaryHostedZone and all associated records
func (s *domainService) DeletePrimaryHostedZone(context context.Context, opts client.DeletePrimaryHostedZoneOpts) error {
	hz, err := s.state.GetPrimaryHostedZone()
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	opts.HostedZoneID = hz.HostedZoneID

	if hz.Managed {
		// HostedZone is managed by us, so delete it
		err := s.service.DeleteHostedZone(context, api.DeleteHostedZoneOpts{
			ID:           opts.ID,
			HostedZoneID: opts.HostedZoneID,
			Domain:       hz.Domain,
		})
		if err != nil {
			return fmt.Errorf("deleting hosted zone: %w", err)
		}
	}

	err = s.state.RemoveHostedZone(hz.Domain)
	if err != nil {
		return err
	}

	return nil
}

func (s *domainService) CreatePrimaryHostedZone(context context.Context, opts client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	// [Refactor] Reconciler is responsible for ordering operations
	//
	// We should be doing this check in the reconciler together with a
	// verification towards the AWS API. Keeping this here for the
	// time being, so we are compatible with expected behavior.
	{
		p, err := s.state.GetPrimaryHostedZone()
		if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
			return nil, err
		}

		if err == nil {
			return p, nil
		}
	}

	hz, err := s.service.CreateHostedZone(context, api.CreateHostedZoneOpts{
		ID:     opts.ID,
		Domain: opts.Domain,
		FQDN:   opts.FQDN,
		NSTTL:  opts.NameServerTTL,
	})
	if err != nil {
		return nil, fmt.Errorf("creating hosted zone: %w", err)
	}

	zone := &client.HostedZone{
		ID:                     hz.ID,
		IsDelegated:            false,
		Primary:                true,
		Managed:                hz.Managed,
		FQDN:                   hz.FQDN,
		Domain:                 hz.Domain,
		HostedZoneID:           hz.HostedZoneID,
		NameServers:            hz.NameServers,
		StackName:              hz.StackName,
		CloudFormationTemplate: hz.CloudFormationTemplate,
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
func NewDomainService(service api.DomainService, state client.DomainState) client.DomainService {
	return &domainService{
		service: service,
		state:   state,
	}
}
