package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetIdentityManager is the API route for the identity manager
const TargetIdentityManager = "identitymanagers/"

type identityManagerAPI struct {
	client *HTTPClient
}

func (a *identityManagerAPI) CreateIdentityPool(opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	into := &api.IdentityPool{}
	return into, a.client.DoPost(TargetIdentityManager, &opts, into)
}

// NewIdentityManagerAPI returns an initialised API client
func NewIdentityManagerAPI(client *HTTPClient) client.IdentityManagerAPI {
	return &identityManagerAPI{
		client: client,
	}
}
