package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHelmArgoCD matches the REST API route
const TargetHelmArgoCD = "helm/argocd/"

type argoCDAPI struct {
	client *HTTPClient
}

func (a *argoCDAPI) CreateArgoCD(opts api.CreateArgoCDOpts) (*client.ArgoCD, error) {
	chart := &api.Helm{}

	err := a.client.DoPost(TargetHelmArgoCD, &opts, chart)
	if err != nil {
		return nil, err
	}

	return &client.ArgoCD{
		Chart: chart,
	}, nil
}

// NewArgoCDAPI returns an initialised service
func NewArgoCDAPI(client *HTTPClient) client.ArgoCDAPI {
	return &argoCDAPI{
		client: client,
	}
}
