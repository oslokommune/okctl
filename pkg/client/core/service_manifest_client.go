package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/client"
)

type manifestService struct {
	spinner spinner.Spinner
	api     client.ManifestAPI
	store   client.ManifestStore
	report  client.ManifestReport
}

func (s *manifestService) ScaleDeployment(_ context.Context, opts api.ScaleDeploymentOpts) error {
	err := s.spinner.Start("scale-deployment")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	return s.api.ScaleDeployment(opts)
}

func (s *manifestService) CreateConfigMap(_ context.Context, opts client.CreateConfigMapOpts) (*client.ConfigMap, error) {
	err := s.spinner.Start("native-secret")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	secret, err := s.api.CreateConfigMap(api.CreateConfigMapOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Data:      opts.Data,
		Labels:    opts.Labels,
	})
	if err != nil {
		return nil, err
	}

	sec := &client.ConfigMap{
		ID:        secret.ID,
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Manifest:  secret.Manifest,
	}

	report, err := s.store.SaveConfigMap(sec)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveConfigMap(sec, report)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

func (s *manifestService) DeleteConfigMap(_ context.Context, opts client.DeleteConfigMapOpts) error {
	err := s.spinner.Start("native-secret")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteConfigMap(api.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
	})

	report, err := s.store.RemoveConfigMap(opts.Name, opts.Namespace)
	if err != nil {
		return err
	}

	return s.report.RemoveConfigMap(report)
}

func (s *manifestService) DeleteExternalSecret(_ context.Context, opts client.DeleteExternalSecretOpts) error {
	err := s.spinner.Start("storage-class")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteExternalSecret(api.DeleteExternalSecretsOpts{
		ID:        opts.ID,
		Manifests: opts.Secrets,
	})
	if err != nil {
		return err
	}

	report, err := s.store.RemoveExternalSecret(opts.Secrets)
	if err != nil {
		return err
	}

	return s.report.RemoveExternalSecret(report)
}

func (s *manifestService) CreateStorageClass(_ context.Context, opts api.CreateStorageClassOpts) (*client.StorageClass, error) {
	err := s.spinner.Start("storage-class")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	sc, err := s.api.CreateStorageClass(opts)
	if err != nil {
		return nil, err
	}

	storage := &client.StorageClass{
		ID:       sc.ID,
		Name:     sc.Name,
		Manifest: sc.Manifest,
	}

	report, err := s.store.SaveStorageClass(storage)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveStorageClass(storage, report)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *manifestService) CreateNamespace(_ context.Context, opts api.CreateNamespaceOpts) (*client.Namespace, error) {
	err := s.spinner.Start("namespace")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	ns, err := s.api.CreateNamespace(opts)
	if err != nil {
		return nil, err
	}

	namespace := &client.Namespace{
		ID:        ns.ID,
		Namespace: ns.Namespace,
		Manifest:  ns.Manifest,
	}

	report, err := s.store.SaveNamespace(namespace)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveNamespace(namespace, report)
	if err != nil {
		return nil, err
	}

	return namespace, nil
}

func (s *manifestService) DeleteNamespace(_ context.Context, opts api.DeleteNamespaceOpts) error {
	err := s.spinner.Start("namespace")
	if err != nil {
		return err
	}

	defer func() {
		_ = s.spinner.Stop()
	}()

	err = s.api.DeleteNamespace(opts)
	if err != nil {
		return err
	}

	report, err := s.store.RemoveNamespace(opts.Namespace)
	if err != nil {
		return err
	}

	return s.report.RemoveNamespace(opts.Namespace, report)
}

func (s *manifestService) CreateExternalSecret(_ context.Context, opts client.CreateExternalSecretOpts) (*client.ExternalSecret, error) {
	err := s.spinner.Start("parameter")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	m, err := s.api.CreateExternalSecret(opts)
	if err != nil {
		return nil, err
	}

	manifest := &client.ExternalSecret{
		ID:        m.ID,
		Manifests: m.Manifests,
	}

	report, err := s.store.SaveExternalSecret(manifest)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveExternalSecret(manifest, report)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// NewManifestService returns an initialised service
func NewManifestService(spinner spinner.Spinner, api client.ManifestAPI, store client.ManifestStore, report client.ManifestReport) client.ManifestService {
	return &manifestService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
