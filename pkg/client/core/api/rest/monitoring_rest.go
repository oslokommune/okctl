package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetKubePromStack matches the REST API route
const TargetKubePromStack = "helm/kubepromstack/"

type monitoringAPI struct {
	client *HTTPClient
}

func (a *monitoringAPI) CreateKubePromStack(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetKubePromStack, &opts, into)
}

// NewMonitoringAPI returns an initialised service
func NewMonitoringAPI(client *HTTPClient) client.MonitoringAPI {
	return &monitoringAPI{
		client: client,
	}
}
