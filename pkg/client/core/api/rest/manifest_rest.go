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

// TargetKubeConfigMap is the API route for the native secrets
// nolint: gosec
const TargetKubeConfigMap = "kube/configmaps/"

// TargetKubeScaleDeployment is th API route for the scaling of a deployment
const TargetKubeScaleDeployment = "kube/scale/"

type manifestAPI struct {
	client *HTTPClient
}

func (a *manifestAPI) CreateNamespace(opts api.CreateNamespaceOpts) (*api.Namespace, error) {
	into := &api.Namespace{}
	return into, a.client.DoPost(TargetKubeNamespace, &opts, into)
}

func (a *manifestAPI) ScaleDeployment(opts api.ScaleDeploymentOpts) error {
	return a.client.DoPost(TargetKubeScaleDeployment, &opts, nil)
}

func (a *manifestAPI) CreateConfigMap(opts api.CreateConfigMapOpts) (*api.ConfigMap, error) {
	into := &api.ConfigMap{}
	return into, a.client.DoPost(TargetKubeConfigMap, &opts, into)
}

func (a *manifestAPI) DeleteConfigMap(opts api.DeleteConfigMapOpts) error {
	return a.client.DoDelete(TargetKubeConfigMap, &opts)
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
