package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetKubePromStack matches the REST API route
const TargetKubePromStack = "helm/kubepromstack/"

type kubePromStackAPI struct {
	client *HTTPClient
}

func (a *kubePromStackAPI) CreateKubePromStack(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetKubePromStack, &opts, into)
}

// NewKubePromStackAPI returns an initialised service
func NewKubePromStackAPI(client *HTTPClient) client.KubePromStackAPI {
	return &kubePromStackAPI{
		client: client,
	}
}
