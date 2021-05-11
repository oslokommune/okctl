package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type externalSecretsService struct {
	policy  client.ManagedPolicyService
	account client.ServiceAccountService
	helm    client.HelmService
}

func (s *externalSecretsService) DeleteExternalSecrets(ctx context.Context, id api.ID) error {
	config, err := clusterconfig.NewExternalSecretsServiceAccount(
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
		Name:   "external-secrets",
		Config: config,
	})
	if err != nil {
		return err
	}

	stackName := cfn.NewStackNamer().
		ExternalSecretsPolicy(id.ClusterName)

	err = s.policy.DeletePolicy(ctx, client.DeletePolicyOpts{
		ID:        id,
		StackName: stackName,
	})
	if err != nil {
		return err
	}

	err = s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          id,
		ReleaseName: externalsecrets.ReleaseName,
		Namespace:   externalsecrets.Namespace,
	})
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen
func (s *externalSecretsService) CreateExternalSecrets(ctx context.Context, opts client.CreateExternalSecretsOpts) (*client.ExternalSecrets, error) {
	b := cfn.New(
		components.NewExternalSecretsPolicyComposer(
			opts.ID.ClusterName,
		),
	)

	stackName := cfn.NewStackNamer().
		ExternalSecretsPolicy(opts.ID.ClusterName)

	template, err := b.Build()
	if err != nil {
		return nil, err
	}

	policy, err := s.policy.CreatePolicy(ctx, client.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              stackName,
		PolicyOutputName:       "ExternalSecretsPolicy",
		CloudFormationTemplate: template,
	})
	if err != nil {
		return nil, err
	}

	config, err := clusterconfig.NewExternalSecretsServiceAccount(
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
		Name:      "external-secrets",
		PolicyArn: policy.PolicyARN,
		Config:    config,
	})
	if err != nil {
		return nil, err
	}

	ch := externalsecrets.New(externalsecrets.DefaultExternalSecretsValues(), constant.DefaultChartApplyTimeout)

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

	return &client.ExternalSecrets{
		Policy:         policy,
		ServiceAccount: sa,
		Chart:          chart,
	}, nil
}

// NewExternalSecretsService returns an initialised service
func NewExternalSecretsService(
	policy client.ManagedPolicyService,
	account client.ServiceAccountService,
	helm client.HelmService,
) client.ExternalSecretsService {
	return &externalSecretsService{
		policy:  policy,
		account: account,
		helm:    helm,
	}
}
