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

func (s *manifestService) CreateNativeSecret(_ context.Context, opts client.CreateNativeSecretOpts) (*client.NativeSecret, error) {
	err := s.spinner.Start("native-secret")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	secret, err := s.api.CreateNativeSecret(api.CreateNativeSecretOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Data:      opts.Data,
		Labels:    opts.Labels,
	})
	if err != nil {
		return nil, err
	}

	sec := &client.NativeSecret{
		ID:        secret.ID,
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Manifest:  secret.Manifest,
	}

	report, err := s.store.SaveNativeSecret(sec)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveNativeSecret(sec, report)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

func (s *manifestService) DeleteNativeSecret(_ context.Context, opts client.DeleteNativeSecretOpts) error {
	err := s.spinner.Start("native-secret")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteNativeSecret(api.DeleteNativeSecretOpts{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
	})

	report, err := s.store.RemoveNativeSecret(opts.Name, opts.Namespace)
	if err != nil {
		return err
	}

	return s.report.RemoveNativeSecret(report)
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

func (s *manifestService) DeleteNamespace(_ context.Context, opts api.DeleteNamespaceOpts) error {
	err := s.spinner.Start("namespace")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	return s.api.DeleteNamespace(opts)
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
