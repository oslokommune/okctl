package core // nolint: dupl

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type awsLoadBalancerControllerService struct {
	spinner spinner.Spinner
	api     client.AWSLoadBalancerControllerAPI
	store   client.AWSLoadBalancerControllerStore
	report  client.AWSLoadBalancerControllerReport
}

func (s *awsLoadBalancerControllerService) DeleteAWSLoadBalancerController(_ context.Context, id api.ID) error {
	err := s.spinner.Start("aws-load-balancer-controller")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteAWSLoadBalancerControllerServiceAccount(id)
	if err != nil {
		return err
	}

	err = s.api.DeleteAWSLoadBalancerControllerPolicy(id)
	if err != nil {
		return err
	}

	report, err := s.store.RemoveAWSLoadBalancerController(id)
	if err != nil {
		return err
	}

	err = s.report.ReportDeleteAWSLoadBalancerController(report)
	if err != nil {
		return err
	}

	return nil
}

// nolint: lll
func (s *awsLoadBalancerControllerService) CreateAWSLoadBalancerController(_ context.Context, opts client.CreateAWSLoadBalancerControllerOpts) (*client.AWSLoadBalancerController, error) {
	err := s.spinner.Start("aws-load-balancer-controller")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	policy, err := s.api.CreateAWSLoadBalancerControllerPolicy(api.CreateAWSLoadBalancerControllerPolicyOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	account, err := s.api.CreateAWSLoadBalancerControllerServiceAccount(api.CreateAWSLoadBalancerControllerServiceAccountOpts{
		CreateServiceAccountBaseOpts: api.CreateServiceAccountBaseOpts{
			ID:        opts.ID,
			PolicyArn: policy.PolicyARN,
		},
	})
	if err != nil {
		return nil, err
	}

	chart, err := s.api.CreateAWSLoadBalancerControllerHelmChart(api.CreateAWSLoadBalancerControllerHelmChartOpts{
		ID:    opts.ID,
		VpcID: opts.VPCID,
	})
	if err != nil {
		return nil, err
	}

	AWSLoadBalancerController := &client.AWSLoadBalancerController{
		Policy:         policy,
		ServiceAccount: account,
		Chart:          chart,
	}

	report, err := s.store.SaveAWSLoadBalancerController(AWSLoadBalancerController)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateAWSLoadBalancerController(AWSLoadBalancerController, report)
	if err != nil {
		return nil, err
	}

	return AWSLoadBalancerController, nil
}

// NewAWSLoadBalancerControllerService returns an initialised service
func NewAWSLoadBalancerControllerService(
	spinner spinner.Spinner,
	api client.AWSLoadBalancerControllerAPI,
	store client.AWSLoadBalancerControllerStore,
	report client.AWSLoadBalancerControllerReport,
) client.AWSLoadBalancerControllerService {
	return &awsLoadBalancerControllerService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
