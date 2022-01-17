package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type managedPolicyAPIDirectClient struct {
	service api.ManagedPolicyService
}

func (m *managedPolicyAPIDirectClient) CreatePolicy(opts api.CreatePolicyOpts) (*api.ManagedPolicy, error) {
	return m.service.CreatePolicy(context.Background(), opts)
}

func (m *managedPolicyAPIDirectClient) DeletePolicy(opts api.DeletePolicyOpts) error {
	return m.service.DeletePolicy(context.Background(), opts)
}

// NewManagedPolicyAPI returns an initialised API client that use core service directly
func NewManagedPolicyAPI(service api.ManagedPolicyService) client.ManagedPolicyAPI {
	return &managedPolicyAPIDirectClient{
		service: service,
	}
}
