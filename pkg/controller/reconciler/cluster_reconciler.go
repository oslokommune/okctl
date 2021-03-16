package reconciler

import (
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ClusterResourceState contains runtime data needed in Reconcile()
type ClusterResourceState struct {
	VPC api.Vpc
}

// clusterReconciler contains service and metadata for the relevant resource
type clusterReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ClusterService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *clusterReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeCluster
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *clusterReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *clusterReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(ClusterResourceState)
	if !ok {
		return result, errors.New("casting cluster resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateCluster(z.commonMetadata.Ctx, api.ClusterCreateOpts{
			ID:                z.commonMetadata.ClusterID,
			Cidr:              z.commonMetadata.Declaration.VPC.CIDR,
			Version:           constant.DefaultEKSKubernetesVersion,
			VpcID:             resourceState.VPC.VpcID,
			VpcPrivateSubnets: resourceState.VPC.PrivateSubnets,
			VpcPublicSubnets:  resourceState.VPC.PublicSubnets,
		})
		if err != nil {
			return result, fmt.Errorf("creating cluster: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteCluster(z.commonMetadata.Ctx, api.ClusterDeleteOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("deleting cluster: %w", err)
		}
	}

	return result, nil
}

// NewClusterReconciler creates a new reconciler for the cluster resource
func NewClusterReconciler(client client.ClusterService) Reconciler {
	return &clusterReconciler{
		client: client,
	}
}
