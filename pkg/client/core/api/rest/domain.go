package rest

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/domain"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHostedZone matches the REST API route
const TargetHostedZone = "domains/hostedzones/"

type domainAPI struct {
	client *HTTPClient
}

func (a *domainAPI) CreatePrimaryHostedZone(opts client.CreatePrimaryHostedZoneOpts) (*api.HostedZone, error) {
	d, err := domain.NewDefaultWithSurvey(opts.ID.Repository, opts.ID.Environment)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain name: %w", err)
	}

	into := &api.HostedZone{}

	return into, a.client.DoPost(TargetHostedZone, &api.CreateHostedZoneOpts{
		ID:     opts.ID,
		Domain: d.Domain,
		FQDN:   d.FQDN,
	}, into)
}

// NewDomainAPI returns an initialised REST API client
func NewDomainAPI(client *HTTPClient) client.DomainAPI {
	return &domainAPI{
		client: client,
	}
}
