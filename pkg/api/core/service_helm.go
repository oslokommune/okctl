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

	err = s.store.SaveArgoCD(h)
	if err != nil {
		return nil, errors.E(err, "storing argocd helm chart")
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

	err = s.store.SaveAlbIngressControllerHelmChart(h)
	if err != nil {
		return nil, errors.E(err, "failed to store alb ingress controller helm chart")
	}

	return h, nil
}

func (s *helmService) CreateAWSLoadBalancerControllerHelmChart(_ context.Context, opts api.CreateAWSLoadBalancerControllerHelmChartOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate input options")
	}

	h, err := s.run.CreateAWSLoadBalancerControllerHelmChart(opts)
	if err != nil {
		return nil, errors.E(err, "creating aws load balancer controller helm chart")
	}

	err = s.store.SaveAWSLoadBalancerControllerHelmChart(h)
	if err != nil {
		return nil, errors.E(err, "storing aws load balancer controller helm chart")
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

func (s *helmService) CreateKubePrometheusStack(_ context.Context, opts api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf("validating input options: %w", err)
	}

	h, err := s.run.CreateKubePrometheusStack(opts)
	if err != nil {
		return nil, fmt.Errorf("creating kube prometheus stack helm chart: %w", err)
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
