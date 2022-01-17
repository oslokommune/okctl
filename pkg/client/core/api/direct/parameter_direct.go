package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type parameterAPIDirectClient struct {
	service api.ParameterService
}

func (p *parameterAPIDirectClient) DeleteSecret(opts api.DeleteSecretOpts) error {
	return p.service.DeleteSecret(context.Background(), opts)
}

func (p *parameterAPIDirectClient) CreateSecret(opts api.CreateSecretOpts) (*api.SecretParameter, error) {
	return p.service.CreateSecret(context.Background(), opts)
}

// NewParameterAPI returns an initialised API client that use core parameter service directly
func NewParameterAPI(service api.ParameterService) client.ParameterAPI {
	return &parameterAPIDirectClient{
		service: service,
	}
}
