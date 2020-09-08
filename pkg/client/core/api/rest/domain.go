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

func (a *domainAPI) CreateHostedZone(opts api.CreateHostedZoneOpts) (*api.HostedZone, error) {
	into := &api.HostedZone{}
	return into, a.client.DoPost(TargetHostedZone, &opts, into)
}

// NewDomainAPI returns an initialised REST API client
func NewDomainAPI(client *HTTPClient) client.DomainAPI {
	return &domainAPI{
		client: client,
	}
}
