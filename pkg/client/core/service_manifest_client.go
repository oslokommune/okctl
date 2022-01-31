package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type manifestService struct {
	service api.KubeService
	state   client.ManifestState
}

func (s *manifestService) ScaleDeployment(context context.Context, opts api.ScaleDeploymentOpts) error {
	return s.service.ScaleDeployment(context, opts)
}

func (s *manifestService) CreateConfigMap(context context.Context, opts client.CreateConfigMapOpts) (*client.KubernetesManifest, error) {
	cm, err := s.service.CreateConfigMap(context, api.CreateConfigMapOpts{
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

func (s *manifestService) DeleteConfigMap(context context.Context, opts client.DeleteConfigMapOpts) error {
	err := s.service.DeleteConfigMap(context, api.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
	})
	if err != nil && !errors.IsKind(err, errors.NotExist) {
		return err
	}

	err = s.state.RemoveKubernetesManifests(opts.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestService) DeleteExternalSecret(context context.Context, opts client.DeleteExternalSecretOpts) error {
	err := s.service.DeleteExternalSecrets(context, api.DeleteExternalSecretsOpts{
		ID:        opts.ID,
		Manifests: opts.Secrets,
	})
	if err != nil && !errors.IsKind(err, errors.NotExist) {
		return err
	}

	err = s.state.RemoveKubernetesManifests(opts.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestService) CreateStorageClass(context context.Context, opts api.CreateStorageClassOpts) (*client.KubernetesManifest, error) {
	sc, err := s.service.CreateStorageClass(context, opts)
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

func (s *manifestService) CreateNamespace(context context.Context, opts api.CreateNamespaceOpts) (*client.KubernetesManifest, error) {
	ns, err := s.service.CreateNamespace(context, opts)
	if err != nil {
		return nil, err
	}

	m := &client.KubernetesManifest{
		ID:        opts.ID,
		Name:      opts.Namespace,
		Namespace: opts.Namespace,
		Type:      client.ManifestTypeNamespace,
		Content:   ns.Manifest,
	}

	err = s.state.SaveKubernetesManifests(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *manifestService) DeleteNamespace(context context.Context, opts api.DeleteNamespaceOpts) error {
	err := s.service.DeleteNamespace(context, opts)
	if err != nil && !errors.IsKind(err, errors.NotExist) {
		return err
	}

	err = s.state.RemoveKubernetesManifests(opts.Namespace)
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestService) CreateExternalSecret(context context.Context, featureOpts client.CreateExternalSecretOpts) (*client.KubernetesManifest, error) {
	toolOpts := api.CreateExternalSecretsOpts{
		ID:       featureOpts.ID,
		Manifest: featureOpts.Manifest,
	}

	m, err := s.service.CreateExternalSecrets(context, toolOpts)
	if err != nil {
		return nil, err
	}

	manifest := &client.KubernetesManifest{
		ID:        featureOpts.ID,
		Name:      featureOpts.Name,
		Namespace: featureOpts.Namespace,
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
func NewManifestService(service api.KubeService, state client.ManifestState) client.ManifestService {
	return &manifestService{
		service: service,
		state:   state,
	}
}
