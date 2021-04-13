package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type externalDNSService struct {
	api   client.ExternalDNSAPI
	state client.ExternalDNSState

	policy  client.ManagedPolicyService
	account client.ServiceAccountService
}

func (s *externalDNSService) DeleteExternalDNS(ctx context.Context, id api.ID) error {
	config, err := clusterconfig.NewExternalDNSServiceAccount(
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
		Name:   "external-dns",
		Config: config,
	})
	if err != nil {
		return err
	}

	stackName := cfn.NewStackNamer().
		ExternalDNSPolicy(id.ClusterName)

	err = s.policy.DeletePolicy(ctx, client.DeletePolicyOpts{
		ID:        id,
		StackName: stackName,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveExternalDNS()
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen
func (s *externalDNSService) CreateExternalDNS(ctx context.Context, opts client.CreateExternalDNSOpts) (*client.ExternalDNS, error) {
	b := cfn.New(
		components.NewExternalDNSPolicyComposer(opts.ID.ClusterName),
	)

	stackName := cfn.NewStackNamer().
		ExternalDNSPolicy(opts.ID.ClusterName)

	template, err := b.Build()
	if err != nil {
		return nil, err
	}

	policy, err := s.policy.CreatePolicy(ctx, client.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              stackName,
		PolicyOutputName:       "ExternalDNSPolicy",
		CloudFormationTemplate: template,
	})
	if err != nil {
		return nil, err
	}

	config, err := clusterconfig.NewExternalDNSServiceAccount(
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
		Name:      "external-dns",
		PolicyArn: policy.PolicyARN,
		Config:    config,
	})
	if err != nil {
		return nil, err
	}

	kube, err := s.api.CreateExternalDNSKubeDeployment(api.CreateExternalDNSKubeDeploymentOpts{
		ID:           opts.ID,
		HostedZoneID: opts.HostedZoneID,
		DomainFilter: opts.Domain,
	})
	if err != nil {
		return nil, err
	}

	externalDNS := &client.ExternalDNS{
		Policy:         policy,
		ServiceAccount: account,
		Kube: &client.ExternalDNSKube{
			ID:           kube.ID,
			HostedZoneID: kube.HostedZoneID,
			DomainFilter: kube.DomainFilter,
			Manifests:    kube.Manifests,
		},
	}

	err = s.state.SaveExternalDNS(externalDNS)
	if err != nil {
		return nil, err
	}

	return externalDNS, nil
}

// NewExternalDNSService returns an initialised service
func NewExternalDNSService(
	api client.ExternalDNSAPI,
	state client.ExternalDNSState,
	policy client.ManagedPolicyService,
	account client.ServiceAccountService,
) client.ExternalDNSService {
	return &externalDNSService{
		api:     api,
		state:   state,
		policy:  policy,
		account: account,
	}
}
