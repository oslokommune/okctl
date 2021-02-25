package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHelmKubePromStack matches the REST API route
const TargetHelmKubePromStack = "helm/kubepromstack/"

// TargetHelmLoki matches the REST API route
const TargetHelmLoki = "helm/loki/"

// TargetHelmReleases matches the REST API route
const TargetHelmReleases = "helm/releases/"

type monitoringAPI struct {
	client *HTTPClient
}

func (a *monitoringAPI) DeleteKubePromStack(opts api.DeleteHelmReleaseOpts) error {
	return a.client.DoDelete(TargetHelmReleases, &opts)
}

func (a *monitoringAPI) DeleteLoki(opts api.DeleteHelmReleaseOpts) error {
	return a.client.DoDelete(TargetHelmReleases, &opts)
}

func (a *monitoringAPI) CreateLoki(opts client.CreateLokiOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetHelmLoki, &opts, into)
}

func (a *monitoringAPI) CreateKubePromStack(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetHelmKubePromStack, &opts, into)
}

// NewMonitoringAPI returns an initialised service
func NewMonitoringAPI(client *HTTPClient) client.MonitoringAPI {
	return &monitoringAPI{
		client: client,
	}
}
