package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHelmReleases matches the REST API route
const TargetHelmReleases = "helm/releases/"

type helmAPI struct {
	client *HTTPClient
}

func (h *helmAPI) CreateHelmRelease(opts api.CreateHelmReleaseOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, h.client.DoPost(TargetHelmReleases, &opts, into)
}

func (h *helmAPI) DeleteHelmRelease(opts api.DeleteHelmReleaseOpts) error {
	return h.client.DoDelete(TargetHelmReleases, &opts)
}

// NewHelmAPI returns an initialised API client
func NewHelmAPI(client *HTTPClient) client.HelmAPI {
	return &helmAPI{
		client: client,
	}
}
