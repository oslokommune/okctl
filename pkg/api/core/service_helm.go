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

func (s *helmService) CreateBlockstorageHelmChart(_ context.Context, opts api.CreateBlockstorageHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating input options")
	}

	h, err := s.run.CreateBlockstorageHelmChart(opts)
	if err != nil {
		return nil, errors.E(err, "creating blockstorage helm chart", errors.Internal)
	}

	return h, nil
}

func (s *helmService) CreateAutoscalerHelmChart(_ context.Context, opts api.CreateAutoscalerHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating input options")
	}

	h, err := s.run.CreateAutoscalerHelmChart(opts)
	if err != nil {
		return nil, errors.E(err, "creating autoscaler helm chart", errors.Internal)
	}

	return h, nil
}

func (s *helmService) CreateArgoCD(_ context.Context, opts api.CreateArgoCDOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating input options")
	}

	h, err := s.run.CreateArgoCD(opts)
	if err != nil {
		return nil, errors.E(err, "creating argocd helm chart")
	}

	return h, nil
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

	return h, nil
}

func (s *helmService) CreateAWSLoadBalancerControllerHelmChart(_ context.Context, opts api.CreateAWSLoadBalancerControllerHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating input options", errors.Invalid)
	}

	h, err := s.run.CreateAWSLoadBalancerControllerHelmChart(opts)
	if err != nil {
		return nil, errors.E(err, "creating aws load balancer controller helm chart", errors.Internal)
	}

	return h, nil
}

func (s *helmService) CreateExternalSecretsHelmChart(_ context.Context, opts api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("validating input options: %w", err)
	}

	h, err := s.run.CreateExternalSecretsHelmChart(opts)
	if err != nil {
		return nil, fmt.Errorf("creating external secrets helm chart: %w", err)
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
