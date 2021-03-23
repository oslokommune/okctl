package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// PostgresState contains runtime data needed in Reconcile()
type PostgresState struct {
	VpcID             string
	DBSubnetGroupName string
	DBSubnetIDs       []string
	DBSubnetCIDRs     []string
}

type postgresReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

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

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *postgresReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(PostgresState)
	if !ok {
		return ReconcilationResult{}, fmt.Errorf("casting PostgresState resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		if z.commonMetadata.Declaration.Databases != nil {
			for _, db := range z.commonMetadata.Declaration.Databases.Postgres {
				_, err = z.client.CreatePostgresDatabase(z.commonMetadata.Ctx, client.CreatePostgresDatabaseOpts{
					ID:                z.commonMetadata.ClusterID,
					ApplicationName:   db.Name,
					UserName:          db.User,
					VpcID:             resourceState.VpcID,
					DBSubnetGroupName: resourceState.DBSubnetGroupName,
					DBSubnetIDs:       resourceState.DBSubnetIDs,
					DBSubnetCIDRs:     resourceState.DBSubnetCIDRs,
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
