package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client"
)

type manifestService struct {
	api   client.ManifestAPI
	store client.ManifestStore
}

func (s *manifestService) CreateExternalSecret(ctx context.Context, opts client.CreateExternalSecretOpts) (*client.ExternalSecret, error) {
	m, err := s.api.CreateExternalSecret(opts)
	if err != nil {
		return nil, err
	}

	manifest := &client.ExternalSecret{
		ID:        m.ID,
		Manifests: m.Manifests,
	}

	_, err = s.store.SaveExternalSecret(manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// NewManifestService returns an initialised service
func NewManifestService(api client.ManifestAPI, store client.ManifestStore) client.ManifestService {
	return &manifestService{
		api:   api,
		store: store,
	}
}
