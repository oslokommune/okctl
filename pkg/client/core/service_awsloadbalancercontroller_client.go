package core // nolint: dupl

import (
	"context"

	"github.com/oslokommune/okctl/pkg/helm/charts/awslbc"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type awsLoadBalancerControllerService struct {
	policy  client.ManagedPolicyService
	account client.ServiceAccountService
	helm    client.HelmService
}

func (s *awsLoadBalancerControllerService) DeleteAWSLoadBalancerController(ctx context.Context, id api.ID) error {
	config, err := clusterconfig.NewAWSLoadBalancerControllerServiceAccount(
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
		Name:   "aws-load-balancer-controller",
		Config: config,
	})
	if err != nil {
		return err
	}

	stackName := cfn.NewStackNamer().
		AWSLoadBalancerControllerPolicy(id.ClusterName)

	err = s.policy.DeletePolicy(ctx, client.DeletePolicyOpts{
		ID:        id,
		StackName: stackName,
	})
	if err != nil {
		return err
	}

	ch := awslbc.New(nil)

	err = s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          id,
		ReleaseName: ch.ReleaseName,
		Namespace:   ch.Namespace,
	})
	if err != nil {
		return err
	}

	return nil
}

//nolint:lll,funlen
func (s *awsLoadBalancerControllerService) CreateAWSLoadBalancerController(ctx context.Context, opts client.CreateAWSLoadBalancerControllerOpts) (*client.AWSLoadBalancerController, error) {
	b := cfn.New(
		components.NewAWSLoadBalancerControllerComposer(
			opts.ID.ClusterName,
		),
	)

	stackName := cfn.NewStackNamer().
		AWSLoadBalancerControllerPolicy(opts.ID.ClusterName)

	template, err := b.Build()
	if err != nil {
		return nil, err
	}

	policy, err := s.policy.CreatePolicy(ctx, client.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              stackName,
		PolicyOutputName:       "AWSLoadBalancerControllerPolicy",
		CloudFormationTemplate: template,
	})
	if err != nil {
		return nil, err
	}

	config, err := clusterconfig.NewAWSLoadBalancerControllerServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		policy.PolicyARN,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
	)
	if err != nil {
		return nil, err
	}

	account, err := s.account.CreateServiceAccount(ctx, client.CreateServiceAccountOpts{
		ID:        opts.ID,
		Name:      "aws-load-balancer-controller",
		PolicyArn: policy.PolicyARN,
		Config:    config,
	})
	if err != nil {
		return nil, err
	}

	ch := awslbc.New(awslbc.NewDefaultValues(opts.ID.ClusterName, opts.VPCID, opts.ID.Region))

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

	return &client.AWSLoadBalancerController{
		Policy:         policy,
		ServiceAccount: account,
		Chart:          chart,
	}, nil
}

// NewAWSLoadBalancerControllerService returns an initialised service
func NewAWSLoadBalancerControllerService(
	policy client.ManagedPolicyService,
	account client.ServiceAccountService,
	helm client.HelmService,
) client.AWSLoadBalancerControllerService {
	return &awsLoadBalancerControllerService{
		policy:  policy,
		account: account,
		helm:    helm,
	}
}
