package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type helmService struct {
	service api.HelmService
	state   client.HelmState
}

func (h *helmService) CreateHelmRelease(context context.Context, opts client.CreateHelmReleaseOpts) (*client.Helm, error) {
	r, err := h.service.CreateHelmRelease(context, api.CreateHelmReleaseOpts{
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

func (h *helmService) DeleteHelmRelease(context context.Context, opts client.DeleteHelmReleaseOpts) error {
	err := h.service.DeleteHelmRelease(context, api.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: opts.ReleaseName,
		Namespace:   opts.Namespace,
	})
	if err != nil {
		return err
	}

	return h.state.RemoveHelmRelease(opts.ReleaseName)
}

func (h *helmService) GetHelmRelease(context context.Context, opts client.GetHelmReleaseOpts) (*client.Helm, error) {
	release, err := h.service.GetHelmRelease(context, api.GetHelmReleaseOpts{
		ClusterID:   opts.ClusterID,
		ReleaseName: opts.ReleaseName,
		Namespace:   opts.Namespace,
	})
	if err != nil {
		return nil, err
	}

	return &client.Helm{
		ID:      release.ID,
		Release: release.Release,
		Chart:   release.Chart,
	}, nil
}

// NewHelmService returns an initialised helm service
func NewHelmService(service api.HelmService, state client.HelmState) client.HelmService {
	return &helmService{
		service: service,
		state:   state,
	}
}
