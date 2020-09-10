package rest

import (
	"github.com/oslokommune/okctl/pkg/ask"

	"github.com/oslokommune/okctl/pkg/domain"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHostedZone matches the REST API route
const TargetHostedZone = "domains/hostedzones/"

type domainAPI struct {
	client *HTTPClient
	ask    *ask.Ask
}

func (a *domainAPI) CreatePrimaryHostedZone(opts client.CreatePrimaryHostedZoneOpts) (*api.HostedZone, error) {
	d, err := a.ask.Domain(domain.Default(opts.ID.Repository, opts.ID.Environment))
	if err != nil {
		return nil, err
	}

	into := &api.HostedZone{}

	return into, a.client.DoPost(TargetHostedZone, &api.CreateHostedZoneOpts{
		ID:     opts.ID,
		Domain: d.Domain,
		FQDN:   d.FQDN,
	}, into)
}

// NewDomainAPI returns an initialised REST API client
func NewDomainAPI(ask *ask.Ask, client *HTTPClient) client.DomainAPI {
	return &domainAPI{
		client: client,
		ask:    ask,
	}
}
