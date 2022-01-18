package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type serviceAccountAPIDirectClient struct {
	service api.ServiceAccountService
}

func (s *serviceAccountAPIDirectClient) CreateServiceAccount(opts api.CreateServiceAccountOpts) (*api.ServiceAccount, error) {
	return s.service.CreateServiceAccount(context.Background(), opts)
}

func (s *serviceAccountAPIDirectClient) DeleteServiceAccount(opts api.DeleteServiceAccountOpts) error {
	return s.service.DeleteServiceAccount(context.Background(), opts)
}

// NewServiceAccountAPI returns an initialised API client using core service
func NewServiceAccountAPI(service api.ServiceAccountService) client.ServiceAccountAPI {
	return &serviceAccountAPIDirectClient{
		service: service,
	}
}
