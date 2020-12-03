package reconciler

import (
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ArgocdMetadata contains data known before anything has been done, which is needed in Reconcile()
type ArgocdMetadata struct {
	Organization string
}

// ArgocdResourceState contains runtime data needed in Reconcile()
type ArgocdResourceState struct {
	HostedZone *state.HostedZone
	Repository *client.GithubRepository

	UserPoolID string
	AuthDomain string
}

type argocdReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ArgoCDService
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
		return nil, errors.New("error casting argocd resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateArgoCD(z.commonMetadata.Ctx, client.CreateArgoCDOpts{
			ID:                 z.commonMetadata.ClusterId,
			Domain:             resourceState.HostedZone.Domain,
			FQDN:               resourceState.HostedZone.FQDN,
			HostedZoneID:       resourceState.HostedZone.ID,
			GithubOrganisation: resourceState.Repository.Organisation,
			UserPoolID:         resourceState.UserPoolID,
			AuthDomain:         resourceState.AuthDomain,
			Repository:         resourceState.Repository,
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
func NewArgocdReconciler(client client.ArgoCDService) *argocdReconciler {
	return &argocdReconciler{
		client: client,
	}
}
