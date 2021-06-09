package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// vpcReconciler contains service and metadata for the relevant resource
type vpcReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.VPCService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *vpcReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeVPC
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *vpcReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *vpcReconciler) Reconcile(node *dependencytree.Node, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err = z.client.CreateVpc(z.commonMetadata.Ctx, client.CreateVpcOpts{
			ID:      z.commonMetadata.ClusterID,
			Cidr:    z.commonMetadata.Declaration.VPC.CIDR,
			Minimal: !z.commonMetadata.Declaration.VPC.HighAvailability,
		})
		if err != nil {
			return result, fmt.Errorf("creating vpc: %w", err)
		}
	case dependencytree.NodeStateAbsent:
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
