package reconciler

import (
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ClusterResourceState contains runtime data needed in Reconcile()
type ClusterResourceState struct {
	VPC api.Vpc
}

// ClusterReconciler contains service and metadata for the relevant resource
type ClusterReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ClusterService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *ClusterReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *ClusterReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	resourceState, ok := node.ResourceState.(ClusterResourceState)
	if !ok {
		return nil, errors.New("error casting cluster resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateCluster(z.commonMetadata.Ctx, api.ClusterCreateOpts{
			ID:                z.commonMetadata.ClusterID,
			Cidr:              resourceState.VPC.Cidr,
			VpcID:             resourceState.VPC.VpcID,
			VpcPrivateSubnets: resourceState.VPC.PrivateSubnets,
			VpcPublicSubnets:  resourceState.VPC.PublicSubnets,
		})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating cluster: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteCluster(z.commonMetadata.Ctx, api.ClusterDeleteOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error deleting cluster: %w", err)
		}
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewClusterReconciler creates a new reconciler for the cluster resource
func NewClusterReconciler(client client.ClusterService) *ClusterReconciler {
	return &ClusterReconciler{
		client: client,
	}
}
