package reconciliation

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// clusterReconciler contains service and metadata for the relevant resource
type clusterReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.ClusterService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *clusterReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeCluster
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *clusterReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *clusterReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		vpc, err := state.Vpc.GetVpc(
			cfn.NewStackNamer().Vpc(z.commonMetadata.Declaration.Metadata.Name),
		)
		if err != nil {
			return result, fmt.Errorf("getting vpc: %w", err)
		}

		_, err = z.client.CreateCluster(z.commonMetadata.Ctx, client.ClusterCreateOpts{
			ID:                z.commonMetadata.ClusterID,
			Cidr:              z.commonMetadata.Declaration.VPC.CIDR,
			Version:           constant.DefaultEKSKubernetesVersion,
			VpcID:             vpc.VpcID,
			VpcPrivateSubnets: vpc.PrivateSubnets,
			VpcPublicSubnets:  vpc.PublicSubnets,
		})
		if err != nil {
			return result, fmt.Errorf("creating cluster: %w", err)
		}
	case dependencytree.NodeStateAbsent:
		err = z.client.DeleteCluster(z.commonMetadata.Ctx, client.ClusterDeleteOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("deleting cluster: %w", err)
		}
	}

	return result, nil
}

// NewClusterReconciler creates a new reconciler for the cluster resource
func NewClusterReconciler(client client.ClusterService) reconciliation.Reconciler {
	return &clusterReconciler{
		client: client,
	}
}
