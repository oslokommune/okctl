package reconciler

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/cfn"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/miekg/dns"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *argocdReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeArgoCD
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
// nolint: funlen
func (z *argocdReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		repo, err := z.githubClient.CreateGithubRepository(z.commonMetadata.Ctx, client.CreateGithubRepositoryOpts{
			ID:           z.commonMetadata.ClusterID,
			Host:         constant.DefaultGithubHost,
			Organization: z.commonMetadata.Declaration.Github.Organisation,
			Name:         z.commonMetadata.Declaration.Github.Repository,
		})
		if err != nil {
			return result, fmt.Errorf("fetching deploy key: %w", err)
		}

		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		im, err := state.IdentityManager.GetIdentityPool(
			cfn.NewStackNamer().IdentityPool(z.commonMetadata.Declaration.Metadata.Name),
		)
		if err != nil {
			return result, fmt.Errorf("getting identity pool: %w", err)
		}

		_, err = z.argocdClient.CreateArgoCD(z.commonMetadata.Ctx, client.CreateArgoCDOpts{
			ID:                 z.commonMetadata.ClusterID,
			Domain:             z.commonMetadata.Declaration.ClusterRootDomain,
			FQDN:               dns.Fqdn(z.commonMetadata.Declaration.ClusterRootDomain),
			HostedZoneID:       hz.HostedZoneID,
			GithubOrganisation: z.commonMetadata.Declaration.Github.Organisation,
			UserPoolID:         im.UserPoolID,
			AuthDomain:         im.AuthDomain,
			Repository:         repo,
		})

		if err != nil {
			// nolint: godox
			// TODO: Need to identify the correct error
			if strings.Contains(strings.ToLower(err.Error()), "timeout") {
				fmt.Println(z.commonMetadata.Out, fmt.Errorf("got ArgoCD timeout: %w", err).Error())

				return ReconcilationResult{
					Requeue:      true,
					RequeueAfter: 0,
				}, nil
			}

			return result, fmt.Errorf("creating argocd: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.argocdClient.DeleteArgoCD(z.commonMetadata.Ctx, client.DeleteArgoCDOpts{
			ID: z.commonMetadata.ClusterID,
		})
		if err != nil {
			return result, fmt.Errorf("deleting argocd: %w", err)
		}

		err = z.githubClient.DeleteGithubRepository(z.commonMetadata.Ctx, client.DeleteGithubRepositoryOpts{
			ID:           z.commonMetadata.ClusterID,
			Organisation: z.commonMetadata.Declaration.Github.Organisation,
			Name:         z.commonMetadata.Declaration.Github.Repository,
		})
		if err != nil {
			return result, fmt.Errorf("deleting github repository: %w", err)
		}
	}

	return result, nil
}

// NewArgocdReconciler creates a new reconciler for the ArgoCD resource
func NewArgocdReconciler(
	argocdClient client.ArgoCDService,
	githubClient client.GithubService,
) Reconciler {
	return &argocdReconciler{
		argocdClient: argocdClient,
		githubClient: githubClient,
	}
}
