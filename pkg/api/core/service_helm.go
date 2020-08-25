package core

import (
	"context"
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type helmService struct {
	run   api.HelmRun
	store api.HelmStore
}

func (s *helmService) CreateAlbIngressControllerHelmChart(_ context.Context, opts api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate input options")
	}

	h, err := s.run.CreateAlbIngressControllerHelmChart(opts)
	if err != nil {
		return nil, errors.E(err, "failed to create alb ingress controller helm chart")
	}

	err = s.store.SaveAlbIngressControllerHelmChar(h)
	if err != nil {
		return nil, errors.E(err, "failed to store alb ingress controller helm chart")
	}

	return h, nil
}

func (s *helmService) CreateExternalSecretsHelmChart(_ context.Context, opts api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate input options: %w", err)
	}

	h, err := s.run.CreateExternalSecretsHelmChart(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create external secrets helm chart: %w", err)
	}

	err = s.store.SaveExternalSecretsHelmChart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to store external secrets helm chart: %w", err)
	}

	return h, nil
}

// NewHelmService returns an initialised helm service
func NewHelmService(run api.HelmRun, store api.HelmStore) api.HelmService {
	return &helmService{
		run:   run,
		store: store,
	}
}
