package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/oslokommune/okctl/pkg/helm/charts/autoscaler"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type autoscalerService struct {
	policy  client.ManagedPolicyService
	account client.ServiceAccountService
	helm    client.HelmService
}

func (s *autoscalerService) DeleteAutoscaler(ctx context.Context, id api.ID) error {
	config, err := clusterconfig.NewAutoscalerServiceAccount(
		id.ClusterName,
		id.Region,
		"n/a",
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
	)
	if err != nil {
		return err
	}

	err = s.account.DeleteServiceAccount(ctx, client.DeleteServiceAccountOpts{
		ID:     id,
		Name:   "autoscaler",
		Config: config,
	})
	if err != nil {
		return err
	}

	stackName := cfn.NewStackNamer().
		AutoscalerPolicy(id.ClusterName)

	err = s.policy.DeletePolicy(ctx, client.DeletePolicyOpts{
		ID:        id,
		StackName: stackName,
	})
	if err != nil {
		return err
	}

	err = s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          id,
		ReleaseName: autoscaler.ReleaseName,
		Namespace:   autoscaler.Namespace,
	})
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen
func (s *autoscalerService) CreateAutoscaler(ctx context.Context, opts client.CreateAutoscalerOpts) (*client.Autoscaler, error) {
	b := cfn.New(
		components.NewAutoscalerPolicyComposer(opts.ID.ClusterName),
	)

	stackName := cfn.NewStackNamer().
		AutoscalerPolicy(opts.ID.ClusterName)

	template, err := b.Build()
	if err != nil {
		return nil, err
	}

	policy, err := s.policy.CreatePolicy(ctx, client.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              stackName,
		PolicyOutputName:       "AutoscalerPolicy",
		CloudFormationTemplate: template,
	})
	if err != nil {
		return nil, err
	}

	config, err := clusterconfig.NewAutoscalerServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		policy.PolicyARN,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
	)
	if err != nil {
		return nil, err
	}

	sa, err := s.account.CreateServiceAccount(ctx, client.CreateServiceAccountOpts{
		ID:        opts.ID,
		Name:      "autoscaler",
		PolicyArn: policy.PolicyARN,
		Config:    config,
	})
	if err != nil {
		return nil, err
	}

	ch := autoscaler.New(
		autoscaler.NewDefaultValues(opts.ID.Region, opts.ID.ClusterName, "autoscaler"),
		constant.DefaultChartApplyTimeout,
	)

	values, err := ch.ValuesYAML()
	if err != nil {
		return nil, err
	}

	chart, err := s.helm.CreateHelmRelease(ctx, client.CreateHelmReleaseOpts{
		ID:             opts.ID,
		RepositoryName: ch.RepositoryName,
		RepositoryURL:  ch.RepositoryURL,
		ReleaseName:    ch.ReleaseName,
		Version:        ch.Version,
		Chart:          ch.Chart,
		Namespace:      ch.Namespace,
		Values:         values,
	})
	if err != nil {
		return nil, err
	}

	return &client.Autoscaler{
		Policy:         policy,
		ServiceAccount: sa,
		Chart:          chart,
	}, nil
}

// NewAutoscalerService returns an initialised service
func NewAutoscalerService(
	policy client.ManagedPolicyService,
	account client.ServiceAccountService,
	helm client.HelmService,
) client.AutoscalerService {
	return &autoscalerService{
		policy:  policy,
		account: account,
		helm:    helm,
	}
}
