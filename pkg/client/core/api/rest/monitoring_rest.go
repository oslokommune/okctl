package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHelmKubePromStack matches the REST API route
const TargetHelmKubePromStack = "helm/kubepromstack/"

// TargetHelmLoki matches the REST API route
const TargetHelmLoki = "helm/loki/"

// TargetHelmPromtail matches the REST API route
const TargetHelmPromtail = "helm/promtail/"

type monitoringAPI struct {
	client *HTTPClient
}

func (a *monitoringAPI) CreateTempo(opts api.CreateHelmReleaseOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetHelmReleases, &opts, into)
}

func (a *monitoringAPI) DeleteTempo(opts api.DeleteHelmReleaseOpts) error {
	return a.client.DoDelete(TargetHelmReleases, &opts)
}

func (a *monitoringAPI) DeletePromtail(opts api.DeleteHelmReleaseOpts) error {
	return a.client.DoDelete(TargetHelmReleases, &opts)
}

func (a *monitoringAPI) CreatePromtail(opts client.CreatePromtailOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetHelmPromtail, &opts, into)
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
