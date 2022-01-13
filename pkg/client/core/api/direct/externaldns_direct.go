package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type externalDNSDirectClient struct {
	service api.KubeService
}

func (e *externalDNSDirectClient) CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.ExternalDNSKube, error) {
	return e.service.CreateExternalDNSKubeDeployment(context.Background(), opts)
}

// NewExternalDNSAPI returns an initialised REST API client with core kubeservice
func NewExternalDNSAPI(service api.KubeService) client.ExternalDNSAPI {
	return &externalDNSDirectClient{
		service: service,
	}
}
