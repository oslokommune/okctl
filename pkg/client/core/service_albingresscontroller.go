package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type albIngressControllerService struct {
	api    client.ALBIngressControllerAPI
	store  client.ALBIngressControllerStore
	report client.ALBIngressControllerReport
}

func (s *albIngressControllerService) DeleteALBIngressController(_ context.Context, id api.ID) error {
	err := s.api.DeleteAlbIngressControllerPolicy(id)
	if err != nil {
		return err
	}

	err = s.api.DeleteAlbIngressControllerServiceAccount(id)
	if err != nil {
		return err
	}

	report, err := s.store.RemoveALBIngressController(id)
	if err != nil {
		return err
	}

	err = s.report.ReportDeleteALBIngressController(report)
	if err != nil {
		return err
	}

	return nil
}

func (s *albIngressControllerService) CreateALBIngressController(_ context.Context, opts client.CreateALBIngressControllerOpts) (*client.ALBIngressController, error) {
	policy, err := s.api.CreateAlbIngressControllerPolicy(api.CreateAlbIngressControllerPolicyOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	account, err := s.api.CreateAlbIngressControllerServiceAccount(api.CreateAlbIngressControllerServiceAccountOpts{
		CreateServiceAccountOpts: api.CreateServiceAccountOpts{
			ID:        opts.ID,
			PolicyArn: policy.PolicyARN,
		},
	})
	if err != nil {
		return nil, err
	}

	chart, err := s.api.CreateAlbIngressControllerHelmChart(api.CreateAlbIngressControllerHelmChartOpts{
		ID:    opts.ID,
		VpcID: opts.VPCID,
	})
	if err != nil {
		return nil, err
	}

	albIngressController := &client.ALBIngressController{
		Policy:         policy,
		ServiceAccount: account,
		Chart:          chart,
	}

	report, err := s.store.SaveALBIngressController(albIngressController)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateALBIngressController(albIngressController, report)
	if err != nil {
		return nil, err
	}

	return albIngressController, nil
}

// NewALBIngressControllerService returns an initialised service
func NewALBIngressControllerService(
	api client.ALBIngressControllerAPI,
	store client.ALBIngressControllerStore,
	report client.ALBIngressControllerReport,
) client.ALBIngressControllerService {
	return &albIngressControllerService{
		api:    api,
		store:  store,
		report: report,
	}
}
