package reconciliation

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/cfn"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// PostgresReconcilerState contains the required state for
// adding or removing a postgres database
type PostgresReconcilerState struct {
	DB v1alpha1.ClusterDatabasesPostgres
}

type postgresReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.ComponentService
}

// NodeType returns the resource node type
func (z *postgresReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypePostgresInstance
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *postgresReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *postgresReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	data, ok := node.Data.(*PostgresReconcilerState)
	if !ok {
		return result, fmt.Errorf("getting postgres data")
	}

	switch node.State {
	case dependencytree.NodeStatePresent:
		vpc, err := state.Vpc.GetVpc(
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

		_, err = z.client.CreatePostgresDatabase(z.commonMetadata.Ctx, client.CreatePostgresDatabaseOpts{
			ID:                z.commonMetadata.ClusterID,
			ApplicationName:   data.DB.Name,
			UserName:          data.DB.User,
			VpcID:             vpc.VpcID,
			DBSubnetGroupName: vpc.DatabaseSubnetsGroupName,
			DBSubnetIDs:       ids,
			DBSubnetCIDRs:     cidrs,
			Namespace:         data.DB.Namespace,
		})
		if err != nil {
			return result, fmt.Errorf("creating postgres database: %w", err)
		}
	case dependencytree.NodeStateAbsent:
		vpc, err := state.Vpc.GetVpc(
			cfn.NewStackNamer().Vpc(z.commonMetadata.Declaration.Metadata.Name),
		)
		if err != nil {
			return result, fmt.Errorf("getting vpc: %w", err)
		}

		err = z.client.DeletePostgresDatabase(z.commonMetadata.Ctx, client.DeletePostgresDatabaseOpts{
			ID:              z.commonMetadata.ClusterID,
			ApplicationName: data.DB.Name,
			Namespace:       data.DB.Namespace,
			VpcID:           vpc.VpcID,
		})
		if err != nil {
			return result, fmt.Errorf("deleting database: %s, got: %w", data.DB.Name, err)
		}

		return result, nil
	}

	return result, nil
}

// NewPostgresReconciler creates a new reconciler for the Postgres resource
func NewPostgresReconciler(client client.ComponentService) reconciliation.Reconciler {
	return &postgresReconciler{
		client: client,
	}
}
