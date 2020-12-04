package reconciler

import (
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// VPCMetadata contains data extracted from the desired state
type VPCMetadata struct {
	Cidr             string
	HighAvailability bool
}

// vpcReconciler contains service and metadata for the relevant resource
type vpcReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.VPCService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *vpcReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *vpcReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	metadata, ok := node.Metadata.(VPCMetadata)
	if !ok {
		return nil, errors.New("unable to cast VPC metadata")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateVpc(z.commonMetadata.Ctx, api.CreateVpcOpts{
			ID:      z.commonMetadata.ClusterID,
			Cidr:    metadata.Cidr,
			Minimal: !metadata.HighAvailability,
		})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating vpc: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteVpc(z.commonMetadata.Ctx, api.DeleteVpcOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error deleting vpc: %w", err)
		}
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewVPCReconciler creates a new reconciler for the VPC resource
func NewVPCReconciler(client client.VPCService) Reconciler {
	return &vpcReconciler{
		client: client,
	}
}
