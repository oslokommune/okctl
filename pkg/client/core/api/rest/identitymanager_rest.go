package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// nolint: golint
const (
	TargetIdentityPool       = "identitymanagers/pools/"
	TargetIdentityPoolClient = "identitymanagers/pools/clients/"
	TargetIdentityPoolUsers  = "identitymanagers/pools/users/"
)

type identityManagerAPI struct {
	client *HTTPClient
}

func (a *identityManagerAPI) DeleteIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) error {
	return a.client.DoDelete(TargetIdentityPoolClient, &opts)
}

func (a *identityManagerAPI) DeleteIdentityPool(opts api.DeleteIdentityPoolOpts) error {
	return a.client.DoDelete(TargetIdentityPool, &opts)
}

func (a *identityManagerAPI) CreateIdentityPoolUser(opts api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error) {
	into := &api.IdentityPoolUser{}
	return into, a.client.DoPost(TargetIdentityPoolUsers, &opts, into)
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
