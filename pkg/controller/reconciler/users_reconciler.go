package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// UsersState contains runtime data needed in Reconcile()
type UsersState struct {
	UserPoolID string
	Users      []v1alpha1.ClusterUser
}

type usersReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.IdentityManagerService
}

// NodeType returns the resource node type
func (z *usersReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeUsers
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *usersReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *usersReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(UsersState)
	if !ok {
		return ReconcilationResult{}, fmt.Errorf("casting UsersState resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		for _, u := range z.commonMetadata.Declaration.Users {
			_, err = z.client.CreateIdentityPoolUser(z.commonMetadata.Ctx, client.CreateIdentityPoolUserOpts{
				ID:         z.commonMetadata.ClusterID,
				Email:      u.Email,
				UserPoolID: resourceState.UserPoolID,
			})
			if err != nil {
				return result, fmt.Errorf("creating user: %w", err)
			}
		}
	case resourcetree.ResourceNodeStateAbsent:
		return result, fmt.Errorf("not implemented")
	}

	return result, nil
}

// NewUsersReconciler creates a new reconciler for the Users resource
func NewUsersReconciler(client client.IdentityManagerService) Reconciler {
	return &usersReconciler{
		client: client,
	}
}
