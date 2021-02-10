package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetKubeExternalSecret is the API route for the manifest
const TargetKubeExternalSecret = "kube/externalsecrets/"

// TargetKubeNamespace is the api route for the namespace
const TargetKubeNamespace = "kube/namespaces/"

type manifestAPI struct {
	client *HTTPClient
}

func (a *manifestAPI) DeleteNamespace(opts api.DeleteNamespaceOpts) error {
	return a.client.DoDelete(TargetKubeNamespace, &opts)
}

func (a *manifestAPI) CreateExternalSecret(opts client.CreateExternalSecretOpts) (*api.ExternalSecretsKube, error) {
	into := &api.ExternalSecretsKube{}
	return into, a.client.DoPost(TargetKubeExternalSecret, &opts, into)
}

// NewManifestAPI returns an initialised API client
func NewManifestAPI(client *HTTPClient) client.ManifestAPI {
	return &manifestAPI{
		client: client,
	}
}
