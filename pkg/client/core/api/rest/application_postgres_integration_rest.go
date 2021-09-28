package rest

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetEarlyDemux matches the REST API route
	TargetEarlyDemux = "kube/earlydemux"
)

type applicationPostgresIntegrationAPI struct {
	client *HTTPClient
}

func (c *applicationPostgresIntegrationAPI) DisableEarlyTCPDemux(_ context.Context, clusterID api.ID) error {
	return c.client.DoDelete(TargetEarlyDemux, &clusterID)
}

// NewApplicationPostgresIntegrationAPI returns an initialised REST API invoker
func NewApplicationPostgresIntegrationAPI(client *HTTPClient) client.ApplicationPostgresAPI {
	return &applicationPostgresIntegrationAPI{
		client: client,
	}
}
