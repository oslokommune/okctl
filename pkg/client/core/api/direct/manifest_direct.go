package direct

import (
	"context"

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

type manifestAPIDirectClient struct {
	service api.KubeService
}

func (a *manifestAPIDirectClient) CreateNamespace(opts api.CreateNamespaceOpts) (*api.Namespace, error) {
	return a.service.CreateNamespace(context.Background(), opts)
}

func (a *manifestAPIDirectClient) ScaleDeployment(opts api.ScaleDeploymentOpts) error {
	return a.service.ScaleDeployment(context.Background(), opts)
}

func (a *manifestAPIDirectClient) CreateConfigMap(opts api.CreateConfigMapOpts) (*api.ConfigMap, error) {
	return a.service.CreateConfigMap(context.Background(), opts)
}

func (a *manifestAPIDirectClient) DeleteConfigMap(opts api.DeleteConfigMapOpts) error {
	return a.service.DeleteConfigMap(context.Background(), opts)
}

func (a *manifestAPIDirectClient) DeleteExternalSecret(opts api.DeleteExternalSecretsOpts) error {
	return a.service.DeleteExternalSecrets(context.Background(), opts)
}

func (a *manifestAPIDirectClient) CreateStorageClass(opts api.CreateStorageClassOpts) (*api.StorageClassKube, error) {
	return a.service.CreateStorageClass(context.Background(), opts)
}

func (a *manifestAPIDirectClient) DeleteNamespace(opts api.DeleteNamespaceOpts) error {
	return a.service.DeleteNamespace(context.Background(), opts)
}

func (a *manifestAPIDirectClient) CreateExternalSecret(opts client.CreateExternalSecretOpts) (*api.ExternalSecretsKube, error) {
	serviceOpts := api.CreateExternalSecretsOpts{
		ID:       opts.ID,
		Manifest: opts.Manifest,
	}

	return a.service.CreateExternalSecrets(context.Background(), serviceOpts)
}

// NewManifestAPI returns an initialised API client that use core service directly
func NewManifestAPI(service api.KubeService) client.ManifestAPI {
	return &manifestAPIDirectClient{
		service: service,
	}
}
