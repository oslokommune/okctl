package reconsiler

import (
	"fmt"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ArgocdMetadata contains data known before anything has been done, which is needed in Reconsile()
type ArgocdMetadata struct {
	Organization string
}

// ArgocdResourceState contains runtime data needed in Reconsile()
type ArgocdResourceState struct {
	HostedZone *state.HostedZone
	Repository *client.GithubRepository
	
	UserPoolID string
	AuthDomain string
}

type argocdReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ArgoCDService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *argocdReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

/*
Reconsile knows how to do what is necessary to ensure the desired state is achieved
Dependent on:
- Github repo setup
- Cognito user pool
- Primary hosted Zone
 */
func (z *argocdReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	resourceState, ok := node.ResourceState.(ArgocdResourceState)
	if !ok {
		return nil, errors.New("error casting argocd resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateArgoCD(z.commonMetadata.Ctx, client.CreateArgoCDOpts{
			ID:                 z.commonMetadata.Id,
			Domain:             resourceState.HostedZone.Domain,
			FQDN:               resourceState.HostedZone.FQDN,
			HostedZoneID:       resourceState.HostedZone.ID,
			GithubOrganisation: resourceState.Repository.Organisation,
			UserPoolID:         resourceState.UserPoolID,
			AuthDomain:         resourceState.AuthDomain,
			Repository:         resourceState.Repository,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating argocd: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, errors.New("deletion of the argocd resource is not implemented")
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewArgocdReconsiler creates a new reconsiler for the ArgoCD resource
func NewArgocdReconsiler(client client.ArgoCDService) *argocdReconsiler {
	return &argocdReconsiler{
		client: client,
	}
}
