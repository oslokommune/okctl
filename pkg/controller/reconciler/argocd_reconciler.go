package reconciler

import (
	"fmt"

	"github.com/miekg/dns"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ArgocdMetadata contains data known before anything has been done, which is needed in Reconcile()
type ArgocdMetadata struct {
	Organization      string
	PrimaryHostedZone string
}

// ArgocdResourceState contains runtime data needed in Reconcile()
type ArgocdResourceState struct {
	HostedZone *state.HostedZone
	Repository *client.GithubRepository

	UserPoolID string
	AuthDomain string
}

// argocdReconciler contains service and metadata for the relevant resource
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
	metadata, ok := node.Metadata.(ArgocdMetadata)
	if !ok {
		return nil, errors.New("casting ArgoCD metadata")
	}

	resourceState, ok := node.ResourceState.(ArgocdResourceState)
	if !ok {
		return nil, errors.New("casting ArgoCD resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateArgoCD(z.commonMetadata.Ctx, client.CreateArgoCDOpts{
			ID:                 z.commonMetadata.ClusterID,
			Domain:             metadata.PrimaryHostedZone,
			FQDN:               dns.Fqdn(metadata.PrimaryHostedZone),
			HostedZoneID:       resourceState.HostedZone.ID,
			GithubOrganisation: metadata.Organization,
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
func NewArgocdReconciler(client client.ArgoCDService) Reconciler {
	return &argocdReconciler{
		client: client,
	}
}
