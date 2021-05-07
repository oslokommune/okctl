package core

import (
	"context"
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type helmService struct {
	run api.HelmRun
}

func (s *helmService) CreateHelmRelease(_ context.Context, opts api.CreateHelmReleaseOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating input options", errors.Invalid)
	}

	h, err := s.run.CreateHelmRelease(opts)
	if err != nil {
		return nil, errors.E(err, fmt.Sprintf("creating helm release (%s): ", opts.ReleaseName), errors.Internal)
	}

	return h, nil
}

func (s *helmService) DeleteHelmRelease(_ context.Context, opts api.DeleteHelmReleaseOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating input options", errors.Invalid)
	}

	err = s.run.DeleteHelmRelease(opts)
	if err != nil {
		return errors.E(err, "removing helm release", errors.Internal)
	}

	return nil
}

// NewHelmService returns an initialised helm service
func NewHelmService(run api.HelmRun) api.HelmService {
	return &helmService{
		run: run,
	}
}
