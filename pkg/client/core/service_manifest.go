package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client"
)

type manifestService struct {
	api    client.ManifestAPI
	store  client.ManifestStore
	report client.ManifestReport
}

func (s *manifestService) CreateExternalSecret(_ context.Context, opts client.CreateExternalSecretOpts) (*client.ExternalSecret, error) {
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
func NewManifestService(api client.ManifestAPI, store client.ManifestStore, report client.ManifestReport) client.ManifestService {
	return &manifestService{
		api:    api,
		store:  store,
		report: report,
	}
}
