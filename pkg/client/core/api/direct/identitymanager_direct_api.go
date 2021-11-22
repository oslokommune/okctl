package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type identityManagerDirectClient struct {
	service api.IdentityManagerService
}

func (i identityManagerDirectClient) CreateIdentityPool(opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	return i.service.CreateIdentityPool(context.Background(), opts)
}

func (i identityManagerDirectClient) DeleteIdentityPool(opts api.DeleteIdentityPoolOpts) error {
	return i.service.DeleteIdentityPool(context.Background(), opts)
}

func (i identityManagerDirectClient) CreateIdentityPoolClient(opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	return i.service.CreateIdentityPoolClient(context.Background(), opts)
}

func (i identityManagerDirectClient) DeleteIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) error {
	return i.service.DeleteIdentityPoolClient(context.Background(), opts)
}

func (i identityManagerDirectClient) CreateIdentityPoolUser(opts api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error) {
	return i.service.CreateIdentityPoolUser(context.Background(), opts)
}

func (i identityManagerDirectClient) DeleteIdentityPoolUser(opts api.DeleteIdentityPoolUserOpts) error {
	return i.service.DeleteIdentityPoolUser(context.Background(), opts)
}

// NewIdentityManagerAPI returns an initialised API that user server side service directly
func NewIdentityManagerAPI(service api.IdentityManagerService) client.IdentityManagerAPI {
	return &identityManagerDirectClient{
		service: service,
	}
}
