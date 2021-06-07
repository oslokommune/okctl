package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// vpcReconciler contains service and metadata for the relevant resource
type vpcReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.VPCService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *vpcReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeVPC
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *vpcReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *vpcReconciler) Reconcile(node *resourcetree.ResourceNode, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateVpc(z.commonMetadata.Ctx, client.CreateVpcOpts{
			ID:      z.commonMetadata.ClusterID,
			Cidr:    z.commonMetadata.Declaration.VPC.CIDR,
			Minimal: !z.commonMetadata.Declaration.VPC.HighAvailability,
		})
		if err != nil {
			return result, fmt.Errorf("creating vpc: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteVpc(z.commonMetadata.Ctx, client.DeleteVpcOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("deleting vpc: %w", err)
		}
	}

	return result, nil
}

// NewVPCReconciler creates a new reconciler for the VPC resource
func NewVPCReconciler(client client.VPCService) reconciliation.Reconciler {
	return &vpcReconciler{
		client: client,
	}
}
