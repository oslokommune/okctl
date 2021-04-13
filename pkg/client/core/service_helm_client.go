package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type helmService struct {
	api   client.HelmAPI
	state client.HelmState
}

func (h *helmService) CreateHelmRelease(_ context.Context, opts client.CreateHelmReleaseOpts) (*client.Helm, error) {
	r, err := h.api.CreateHelmRelease(api.CreateHelmReleaseOpts{
		ID:             opts.ID,
		RepositoryName: opts.RepositoryName,
		RepositoryURL:  opts.RepositoryURL,
		ReleaseName:    opts.ReleaseName,
		Version:        opts.Version,
		Chart:          opts.Chart,
		Namespace:      opts.Namespace,
		Values:         opts.Values,
	})
	if err != nil {
		return nil, err
	}

	release := &client.Helm{
		ID:      r.ID,
		Release: r.Release,
		Chart:   r.Chart,
	}

	err = h.state.SaveHelmRelease(release)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func (h *helmService) DeleteHelmRelease(_ context.Context, opts client.DeleteHelmReleaseOpts) error {
	err := h.api.DeleteHelmRelease(api.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: opts.ReleaseName,
		Namespace:   opts.Namespace,
	})
	if err != nil {
		return err
	}

	return h.state.RemoveHelmRelease(opts.ReleaseName)
}

// NewHelmService returns an initialised helm service
func NewHelmService(api client.HelmAPI, state client.HelmState) client.HelmService {
	return &helmService{
		api:   api,
		state: state,
	}
}
