package reconciler

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"

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
func (z *argocdReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(ArgocdResourceState)
	if !ok {
		return result, errors.New("casting ArgoCD resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		repository := client.NewGithubRepository(
			z.commonMetadata.ClusterID,
			constant.DefaultGithubHost,
			z.commonMetadata.Declaration.Github.Organisation,
			z.commonMetadata.Declaration.Github.Repository,
		)

		var key *client.GithubDeployKey

		key, err = z.githubClient.CreateDeployKey(z.commonMetadata.Ctx, repository)
		if err != nil {
			return result, fmt.Errorf("fetching deploy key: %w", err)
		}

		repository.DeployKey = key

		_, err = z.argocdClient.CreateArgoCD(z.commonMetadata.Ctx, client.CreateArgoCDOpts{
			ID:                 z.commonMetadata.ClusterID,
			Domain:             z.commonMetadata.Declaration.ClusterRootURL,
			FQDN:               dns.Fqdn(z.commonMetadata.Declaration.ClusterRootURL),
			HostedZoneID:       resourceState.HostedZone.ID,
			GithubOrganisation: z.commonMetadata.Declaration.Github.Organisation,
			UserPoolID:         resourceState.UserPoolID,
			AuthDomain:         resourceState.AuthDomain,
			Repository:         repository,
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
		return result, errors.New("deletion of the argocd resource is not implemented")
	}

	return result, nil
}

// NewArgocdReconciler creates a new reconciler for the ArgoCD resource
func NewArgocdReconciler(argocdClient client.ArgoCDService, githubClient client.GithubService) Reconciler {
	return &argocdReconciler{
		argocdClient: argocdClient,
		githubClient: githubClient,
	}
}
