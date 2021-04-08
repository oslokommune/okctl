package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/oslokommune/okctl/pkg/helm/charts/blockstorage"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"

	"github.com/oslokommune/okctl/pkg/kube/manifests/storageclass"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type blockstorageService struct {
	policy  client.ManagedPolicyService
	account client.ServiceAccountService
	helm    client.HelmService
	kube    client.ManifestService
}

func (s *blockstorageService) DeleteBlockstorage(ctx context.Context, id api.ID) error {
	chart := blockstorage.New(nil)

	err := s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          id,
		ReleaseName: chart.ReleaseName,
		Namespace:   chart.Namespace,
	})
	if err != nil {
		return err
	}

	config, err := clusterconfig.NewBlockstorageServiceAccount(
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
		Name:   "blockstorage",
		Config: config,
	})
	if err != nil {
		return err
	}

	stackName := cfn.NewStackNamer().
		BlockstoragePolicy(id.Repository, id.Environment)

	err = s.policy.DeletePolicy(ctx, client.DeletePolicyOpts{
		ID:        id,
		StackName: stackName,
	})
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen
func (s *blockstorageService) CreateBlockstorage(ctx context.Context, opts client.CreateBlockstorageOpts) (*client.Blockstorage, error) {
	b := cfn.New(
		components.NewBlockstoragePolicyComposer(opts.ID.Repository, opts.ID.Environment),
	)

	stackName := cfn.NewStackNamer().
		BlockstoragePolicy(opts.ID.Repository, opts.ID.Environment)

	template, err := b.Build()
	if err != nil {
		return nil, err
	}

	policy, err := s.policy.CreatePolicy(ctx, client.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              stackName,
		PolicyOutputName:       "BlockstoragePolicy",
		CloudFormationTemplate: template,
	})
	if err != nil {
		return nil, err
	}

	config, err := clusterconfig.NewBlockstorageServiceAccount(
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
		Name:      "blockstorage",
		PolicyArn: policy.PolicyARN,
		Config:    config,
	})
	if err != nil {
		return nil, err
	}

	ch := blockstorage.New(blockstorage.NewDefaultValues(opts.ID.Region, opts.ID.ClusterName, "blockstorage"))

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

	a := &client.Blockstorage{
		Policy:         policy,
		ServiceAccount: sa,
		Chart:          chart,
	}

	_, err = s.kube.CreateStorageClass(ctx, api.CreateStorageClassOpts{
		ID:          opts.ID,
		Name:        "ebs-sc",
		Parameters:  storageclass.NewEBSParameters(),
		Annotations: storageclass.DefaultStorageClassAnnotation(),
	})
	if err != nil {
		return nil, fmt.Errorf("creating default storage class: %w", err)
	}

	return a, nil
}

// NewBlockstorageService returns an initialised service
func NewBlockstorageService(
	policy client.ManagedPolicyService,
	account client.ServiceAccountService,
	helm client.HelmService,
	kube client.ManifestService,
) client.BlockstorageService {
	return &blockstorageService{
		policy:  policy,
		account: account,
		helm:    helm,
		kube:    kube,
	}
}
