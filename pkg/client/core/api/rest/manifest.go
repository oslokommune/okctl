package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetKubeExternalSecret is the API route for the manifest
const TargetKubeExternalSecret = "kube/externalsecrets/"

type manifestAPI struct {
	client *client.HTTPClient
}

func (a *manifestAPI) CreateExternalSecret(opts client.CreateExternalSecretOpts) (*api.Kube, error) {
	into := &api.Kube{}
	return into, a.client.DoPost(TargetKubeExternalSecret, opts, into)
}

// NewManifestAPI returns an initialised API client
func NewManifestAPI(client *client.HTTPClient) client.ManifestAPI {
	return &manifestAPI{
		client: client,
	}
}
