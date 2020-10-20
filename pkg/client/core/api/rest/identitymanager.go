package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// nolint: golint
const (
	TargetIdentityPool       = "identitymanagers/pools/"
	TargetIdentityPoolClient = "identitymanagers/pools/clients/"
)

type identityManagerAPI struct {
	client *HTTPClient
}

func (a *identityManagerAPI) CreateIdentityPoolClient(opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	into := &api.IdentityPoolClient{}
	return into, a.client.DoPost(TargetIdentityPoolClient, &opts, into)
}

func (a *identityManagerAPI) CreateIdentityPool(opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	into := &api.IdentityPool{}
	return into, a.client.DoPost(TargetIdentityPool, &opts, into)
}

// NewIdentityManagerAPI returns an initialised API client
func NewIdentityManagerAPI(client *HTTPClient) client.IdentityManagerAPI {
	return &identityManagerAPI{
		client: client,
	}
}
