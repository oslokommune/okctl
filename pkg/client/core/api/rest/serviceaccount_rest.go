package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// ServiceAccountsTarget defines the REST API Route
const ServiceAccountsTarget = "serviceaccounts/"

type serviceAccountAPI struct {
	client *HTTPClient
}

func (m *serviceAccountAPI) CreateServiceAccount(opts api.CreateServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, m.client.DoPost(ServiceAccountsTarget, &opts, into)
}

func (m *serviceAccountAPI) DeleteServiceAccount(opts api.DeleteServiceAccountOpts) error {
	return m.client.DoDelete(ServiceAccountsTarget, &opts)
}

// NewServiceAccountAPI returns an initialised API client
func NewServiceAccountAPI(client *HTTPClient) client.ServiceAccountAPI {
	return &serviceAccountAPI{
		client: client,
	}
}
