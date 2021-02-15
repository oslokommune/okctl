package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type autoscalerService struct {
	spinner spinner.Spinner
	api     client.AutoscalerAPI
	store   client.AutoscalerStore
	report  client.AutoscalerReport
}

func (s *autoscalerService) DeleteAutoscaler(_ context.Context, id api.ID) error {
	err := s.spinner.Start("autoscaler")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteAutoscalerServiceAccount(id)
	if err != nil {
		return err
	}

	err = s.api.DeleteAutoscalerPolicy(id)
	if err != nil {
		return err
	}

	report, err := s.store.RemoveAutoscaler(id)
	if err != nil {
		return err
	}

	err = s.report.ReportDeleteAutoscaler(report)
	if err != nil {
		return err
	}

	return nil
}

func (s *autoscalerService) CreateAutoscaler(_ context.Context, opts client.CreateAutoscalerOpts) (*client.Autoscaler, error) {
	err := s.spinner.Start("autoscaler")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	policy, err := s.api.CreateAutoscalerPolicy(api.CreateAutoscalerPolicy{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	sa, err := s.api.CreateAutoscalerServiceAccount(api.CreateAutoscalerServiceAccountOpts{
		CreateServiceAccountOpts: api.CreateServiceAccountOpts{
			ID:        opts.ID,
			PolicyArn: policy.PolicyARN,
		},
	})
	if err != nil {
		return nil, err
	}

	chart, err := s.api.CreateAutoscalerHelmChart(api.CreateAutoscalerHelmChartOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	a := &client.Autoscaler{
		Policy:         policy,
		ServiceAccount: sa,
		Chart:          chart,
	}

	report, err := s.store.SaveAutoscaler(a)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateAutoscaler(a, report)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// NewAutoscalerService returns an initialised service
func NewAutoscalerService(
	spinner spinner.Spinner,
	api client.AutoscalerAPI,
	store client.AutoscalerStore,
	report client.AutoscalerReport,
) client.AutoscalerService {
	return &autoscalerService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
