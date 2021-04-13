package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetKubeExternalDNS matches the REST API route
	TargetKubeExternalDNS = "kube/externaldns/"
)

type externalDNSAPI struct {
	client *HTTPClient
}

func (a *externalDNSAPI) CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.ExternalDNSKube, error) {
	into := &api.ExternalDNSKube{}
	return into, a.client.DoPost(TargetKubeExternalDNS, &opts, into)
}

// NewExternalDNSAPI returns an initialised REST API client
func NewExternalDNSAPI(client *HTTPClient) client.ExternalDNSAPI {
	return &externalDNSAPI{
		client: client,
	}
}
