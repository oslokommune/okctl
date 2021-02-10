package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type kubePrometheusStackService struct {
	spinner spinner.Spinner
	api     client.KubePrometheusStackAPI
	store   client.KubePrometheusStackStore
	report  client.KubePrometheusStackReport
}

//func (s *kubePrometheusStackService) DeleteKubePrometheusStack(_ context.Context, id api.ID) error {
//	err := s.spinner.Start("external-secrets")
//	if err != nil {
//		return err
//	}
//
//	defer func() {
//		err = s.spinner.Stop()
//	}()
//
//	err = s.api.DeleteKubePrometheusStackServiceAccount(id)
//	if err != nil {
//		return err
//	}
//
//	err = s.api.DeleteKubePrometheusStackPolicy(id)
//	if err != nil {
//		return err
//	}
//
//	report, err := s.store.RemoveKubePrometheusStack(id)
//	if err != nil {
//		return err
//	}
//
//	err = s.report.ReportDeleteKubePrometheusStack(report)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (s *kubePrometheusStackService) CreateKubePrometheusStack(_ context.Context, opts client.CreateKubePrometheusStackOpts) (*client.KubePrometheusStack, error) {
	err := s.spinner.Start("kubeprometheusstack")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	policy, err := s.api.CreateKubePrometheusStackPolicy(api.CreateKubePrometheusStackPolicyOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	sa, err := s.api.CreateKubePrometheusStackServiceAccount(api.CreateKubePrometheusStackServiceAccountOpts{
		CreateServiceAccountOpts: api.CreateServiceAccountOpts{
			ID:        opts.ID,
			PolicyArn: policy.PolicyARN,
		},
	})
	if err != nil {
		return nil, err
	}

	chart, err := s.api.CreateKubePrometheusStackHelmChart(api.CreateKubePrometheusStackHelmChartOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	externalSecrets := &client.KubePrometheusStack{
		Policy:         policy,
		ServiceAccount: sa,
		Chart:          chart,
	}

	report, err := s.store.SaveKubePrometheusStack(externalSecrets)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateKubePrometheusStack(externalSecrets, report)
	if err != nil {
		return nil, err
	}

	return externalSecrets, nil
}

// NewKubePrometheusStackService returns an initialised service
func NewKubePrometheusStackService(
	spinner spinner.Spinner,
	api client.KubePrometheusStackAPI,
	store client.KubePrometheusStackStore,
	report client.KubePrometheusStackReport,
) client.KubePrometheusStackService {
	return &kubePrometheusStackService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
