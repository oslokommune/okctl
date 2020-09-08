package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetDomain matches the REST API route
const TargetDomain = "domains/"

type domainAPI struct {
	client *client.HTTPClient
}

func (a *domainAPI) CreateDomain(opts api.CreateDomainOpts) (*api.Domain, error) {
	into := &api.Domain{}
	return into, a.client.DoPost(TargetDomain, opts, into)
}

// NewDomainAPI returns an initialised REST API client
func NewDomainAPI(client *client.HTTPClient) client.DomainAPI {
	return &domainAPI{
		client: client,
	}
}
