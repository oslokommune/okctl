package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type manifestService struct {
	api   client.ManifestAPI
	state client.ManifestState
}

func (s *manifestService) ScaleDeployment(_ context.Context, opts api.ScaleDeploymentOpts) error {
	return s.api.ScaleDeployment(opts)
}

func (s *manifestService) CreateConfigMap(_ context.Context, opts client.CreateConfigMapOpts) (*client.KubernetesManifest, error) {
	cm, err := s.api.CreateConfigMap(api.CreateConfigMapOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Data:      opts.Data,
		Labels:    opts.Labels,
	})
	if err != nil {
		return nil, err
	}

	m := &client.KubernetesManifest{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Type:      client.ManifestTypeConfigMap,
		Content:   cm.Manifest,
	}

	err = s.state.SaveKubernetesManifests(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *manifestService) DeleteConfigMap(_ context.Context, opts client.DeleteConfigMapOpts) error {
	err := s.api.DeleteConfigMap(api.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveKubernetesManifests(opts.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestService) DeleteExternalSecret(_ context.Context, opts client.DeleteExternalSecretOpts) error {
	err := s.api.DeleteExternalSecret(api.DeleteExternalSecretsOpts{
		ID:        opts.ID,
		Manifests: opts.Secrets,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveKubernetesManifests(opts.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestService) CreateStorageClass(_ context.Context, opts api.CreateStorageClassOpts) (*client.KubernetesManifest, error) {
	sc, err := s.api.CreateStorageClass(opts)
	if err != nil {
		return nil, err
	}

	m := &client.KubernetesManifest{
		ID:      opts.ID,
		Name:    opts.Name,
		Type:    client.ManifestTypeStorageClass,
		Content: sc.Manifest,
	}

	err = s.state.SaveKubernetesManifests(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *manifestService) CreateNamespace(_ context.Context, opts api.CreateNamespaceOpts) (*client.KubernetesManifest, error) {
	ns, err := s.api.CreateNamespace(opts)
	if err != nil {
		return nil, err
	}

	m := &client.KubernetesManifest{
		ID:      opts.ID,
		Name:    opts.Namespace,
		Type:    client.ManifestTypeNamespace,
		Content: ns.Manifest,
	}

	err = s.state.SaveKubernetesManifests(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *manifestService) DeleteNamespace(_ context.Context, opts api.DeleteNamespaceOpts) error {
	err := s.api.DeleteNamespace(opts)
	if err != nil {
		return err
	}

	err = s.state.RemoveKubernetesManifests(opts.Namespace)
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestService) CreateExternalSecret(_ context.Context, opts client.CreateExternalSecretOpts) (*client.KubernetesManifest, error) {
	m, err := s.api.CreateExternalSecret(opts)
	if err != nil {
		return nil, err
	}

	manifest := &client.KubernetesManifest{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Type:      client.ManifestTypeExternalSecret,
		Content:   m.Content,
	}

	err = s.state.SaveKubernetesManifests(manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// NewManifestService returns an initialised service
func NewManifestService(
	api client.ManifestAPI,
	state client.ManifestState,
) client.ManifestService {
	return &manifestService{
		api:   api,
		state: state,
	}
}
