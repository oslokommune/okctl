package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type postgresReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	stateHandlers  *clientCore.StateHandlers

	client client.ComponentService
}

// NodeType returns the resource node type
func (z *postgresReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypePostgres
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *postgresReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// SetStateHandlers sets the state handlers
func (z *postgresReconciler) SetStateHandlers(handlers *clientCore.StateHandlers) {
	z.stateHandlers = handlers
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *postgresReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		if z.commonMetadata.Declaration.Databases != nil {
			vpc, err := z.stateHandlers.Vpc.GetVpc(
				cfn.NewStackNamer().Vpc(z.commonMetadata.Declaration.Metadata.Name),
			)
			if err != nil {
				return result, fmt.Errorf("getting vpc: %w", err)
			}

			cidrs := make([]string, len(vpc.DatabaseSubnets))
			ids := make([]string, len(vpc.DatabaseSubnets))

			for i, s := range vpc.DatabaseSubnets {
				cidrs[i] = s.Cidr
				ids[i] = s.ID
			}

			for _, db := range z.commonMetadata.Declaration.Databases.Postgres {
				_, err = z.client.CreatePostgresDatabase(z.commonMetadata.Ctx, client.CreatePostgresDatabaseOpts{
					ID:                z.commonMetadata.ClusterID,
					ApplicationName:   db.Name,
					UserName:          db.User,
					VpcID:             vpc.VpcID,
					DBSubnetGroupName: vpc.DatabaseSubnetsGroupName,
					DBSubnetIDs:       ids,
					DBSubnetCIDRs:     cidrs,
					Namespace:         db.Namespace,
				})
				if err != nil {
					return result, fmt.Errorf("creating postgres database: %w", err)
				}
			}
		}
	case resourcetree.ResourceNodeStateAbsent:
		return result, fmt.Errorf("not implemented")
	}

	return result, nil
}

// NewPostgresReconciler creates a new reconciler for the Postgres resource
func NewPostgresReconciler(client client.ComponentService) Reconciler {
	return &postgresReconciler{
		client: client,
	}
}
