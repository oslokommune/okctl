package reconciliation

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

type usersReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.IdentityManagerService
}

// NodeType returns the resource node type
func (z *usersReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeUsers
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *usersReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *usersReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		im, err := state.IdentityManager.GetIdentityPool(
			cfn.NewStackNamer().IdentityPool(z.commonMetadata.Declaration.Metadata.Name),
		)
		if err != nil {
			return result, fmt.Errorf("getting identity pool: %w", err)
		}

		for _, u := range z.commonMetadata.Declaration.Users {
			_, err = z.client.CreateIdentityPoolUser(z.commonMetadata.Ctx, client.CreateIdentityPoolUserOpts{
				ID:         z.commonMetadata.ClusterID,
				Email:      u.Email,
				UserPoolID: im.UserPoolID,
			})
			if err != nil {
				return result, fmt.Errorf("creating user: %w", err)
			}
		}
	case dependencytree.NodeStateAbsent:
		// We need to implement delete of individual users
		return result, nil
	}

	return result, nil
}

// NewUsersReconciler creates a new reconciler for the Users resource
func NewUsersReconciler(client client.IdentityManagerService) reconciliation.Reconciler {
	return &usersReconciler{
		client: client,
	}
}
