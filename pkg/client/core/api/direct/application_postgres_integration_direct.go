package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type applicationPostgresIntegrationDirectClient struct {
	service api.KubeService
}

func (a *applicationPostgresIntegrationDirectClient) DisableEarlyTCPDemux(ctx context.Context, id api.ID) error {
	return a.service.DisableEarlyDEMUX(ctx, id)
}

// NewApplicationPostgresIntegrationAPI initializes an application postgres integration API that uses the server side service directly
func NewApplicationPostgresIntegrationAPI(service api.KubeService) client.ApplicationPostgresAPI {
	return &applicationPostgresIntegrationDirectClient{
		service: service,
	}
}
