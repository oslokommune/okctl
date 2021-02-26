package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetKubeExternalSecret is the API route for the manifest
const TargetKubeExternalSecret = "kube/externalsecrets/"

// TargetKubeNamespace is the api route for the namespace
const TargetKubeNamespace = "kube/namespaces/"

// TargetKubeStorageClasses is the API route for the storage classes
const TargetKubeStorageClasses = "kube/storageclasses/"

// TargetKubeNativeSecret is the API route for the native secrets
// nolint: gosec
const TargetKubeNativeSecret = "kube/nativesecrets/"

type manifestAPI struct {
	client *HTTPClient
}

func (a *manifestAPI) CreateNativeSecret(opts api.CreateNativeSecretOpts) (*api.NativeSecret, error) {
	into := &api.NativeSecret{}
	return into, a.client.DoPost(TargetKubeNativeSecret, &opts, into)
}

func (a *manifestAPI) DeleteNativeSecret(opts api.DeleteNativeSecretOpts) error {
	return a.client.DoDelete(TargetKubeNativeSecret, &opts)
}

func (a *manifestAPI) DeleteExternalSecret(opts api.DeleteExternalSecretsOpts) error {
	return a.client.DoDelete(TargetKubeExternalSecret, &opts)
}

func (a *manifestAPI) CreateStorageClass(opts api.CreateStorageClassOpts) (*api.StorageClassKube, error) {
	into := &api.StorageClassKube{}
	return into, a.client.DoPost(TargetKubeStorageClasses, &opts, into)
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
