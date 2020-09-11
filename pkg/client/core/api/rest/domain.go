package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHostedZone matches the REST API route
const TargetHostedZone = "domains/hostedzones/"

type domainAPI struct {
	client *HTTPClient
}

func (a *domainAPI) CreatePrimaryHostedZone(opts client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	into := &api.HostedZone{}

	err := a.client.DoPost(TargetHostedZone, &api.CreateHostedZoneOpts{
		ID:     opts.ID,
		Domain: opts.Domain,
		FQDN:   opts.FQDN,
	}, into)
	if err != nil {
		return nil, err
	}

	return &client.HostedZone{
		IsDelegated: false,
		Primary:     true,
		HostedZone:  into,
	}, nil
}

// NewDomainAPI returns an initialised REST API client
func NewDomainAPI(client *HTTPClient) client.DomainAPI {
	return &domainAPI{
		client: client,
	}
}
