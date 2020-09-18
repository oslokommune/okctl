package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/client"
)

type manifestService struct {
	spinner spinner.Spinner
	api     client.ManifestAPI
	store   client.ManifestStore
	report  client.ManifestReport
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
