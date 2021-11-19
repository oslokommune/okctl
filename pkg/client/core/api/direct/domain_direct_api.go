package direct

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type directClient struct {
	service api.DomainService
}

// CreatePrimaryHostedZone creates the hosted zone used for services like ArgoCD, Grafana, Cognito
func (d *directClient) CreatePrimaryHostedZone(opts client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	hz, err := d.service.CreateHostedZone(context.Background(), api.CreateHostedZoneOpts{
		ID:     opts.ID,
		Domain: opts.Domain,
		FQDN:   opts.FQDN,
		NSTTL:  opts.NameServerTTL,
	})
	if err != nil {
		return nil, fmt.Errorf("creating hosted zone: %w", err)
	}

	return &client.HostedZone{
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
	}, nil
}

// DeletePrimaryHostedZone deletes the hosted zone used for services like ArgoCD, Grafana, Cognito
func (d *directClient) DeletePrimaryHostedZone(domain string, opts client.DeletePrimaryHostedZoneOpts) error {
	err := d.service.DeleteHostedZone(context.Background(), api.DeleteHostedZoneOpts{
		ID:           opts.ID,
		HostedZoneID: opts.HostedZoneID,
		Domain:       domain,
	})
	if err != nil {
		return fmt.Errorf("deleting hosted zone: %w", err)
	}

	return nil
}

// NewDomainAPI initializes a domain API that uses the server side service directly
func NewDomainAPI(service api.DomainService) client.DomainAPI {
	return &directClient{
		service: service,
	}
}
