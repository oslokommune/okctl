package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/miekg/dns"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ArgocdResourceState contains runtime data needed in Reconcile()
type ArgocdResourceState struct {
	HostedZone *state.HostedZone

	UserPoolID string
	AuthDomain string
}

// argocdReconciler contains service and metadata for the relevant resource
type argocdReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	argocdClient client.ArgoCDService
	githubClient client.GithubService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *argocdReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

/*
Reconcile knows how to do what is necessary to ensure the desired state is achieved
Dependent on:
- Github repo setup
- Cognito user pool
- Primary hosted Zone
*/
func (z *argocdReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	resourceState, ok := node.ResourceState.(ArgocdResourceState)
	if !ok {
		return nil, errors.New("casting ArgoCD resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		repository := client.NewGithubRepository(
			z.commonMetadata.ClusterID,
			config.DefaultGithubHost,
			z.commonMetadata.Declaration.Github.Organisation,
			z.commonMetadata.Declaration.Github.Repository,
		)

		key, err := z.githubClient.CreateDeployKey(z.commonMetadata.Ctx, repository)
		if err != nil {
			return nil, fmt.Errorf("fetching deploy key: %w", err)
		}

		repository.DeployKey = key

		_, err = z.argocdClient.CreateArgoCD(z.commonMetadata.Ctx, client.CreateArgoCDOpts{
			ID:                 z.commonMetadata.ClusterID,
			Domain:             z.commonMetadata.Declaration.PrimaryDNSZone.ParentDomain,
			FQDN:               dns.Fqdn(z.commonMetadata.Declaration.PrimaryDNSZone.ParentDomain),
			HostedZoneID:       resourceState.HostedZone.ID,
			GithubOrganisation: z.commonMetadata.Declaration.Github.Organisation,
			UserPoolID:         resourceState.UserPoolID,
			AuthDomain:         resourceState.AuthDomain,
			Repository:         repository,
		})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating argocd: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, errors.New("deletion of the argocd resource is not implemented")
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewArgocdReconciler creates a new reconciler for the ArgoCD resource
func NewArgocdReconciler(argocdClient client.ArgoCDService, githubClient client.GithubService) Reconciler {
	return &argocdReconciler{
		argocdClient: argocdClient,
		githubClient: githubClient,
	}
}
