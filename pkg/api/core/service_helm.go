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
		return nil, errors.E(err, "creating helm release", errors.Internal)
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

func (s *helmService) CreatePromtailHelmChart(_ context.Context, opts api.CreatePromtailHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating input options", errors.Invalid)
	}

	h, err := s.run.CreatePromtailHelmChart(opts)
	if err != nil {
		return nil, errors.E(err, "creating promtail helm chart", errors.Internal)
	}

	return h, nil
}

func (s *helmService) CreateLokiHelmChart(_ context.Context, opts api.CreateLokiHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating input options", errors.Invalid)
	}

	h, err := s.run.CreateLokiHelmChart(opts)
	if err != nil {
		return nil, errors.E(err, "creating loki helm chart", errors.Internal)
	}

	return h, nil
}

func (s *helmService) CreateKubePrometheusStack(_ context.Context, opts api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("validating input options: %w", err)
	}

	h, err := s.run.CreateKubePromStack(opts)
	if err != nil {
		return nil, fmt.Errorf("creating kube prometheus stack helm chart: %w", err)
	}

	return h, nil
}

// NewHelmService returns an initialised helm service
func NewHelmService(run api.HelmRun) api.HelmService {
	return &helmService{
		run: run,
	}
}
