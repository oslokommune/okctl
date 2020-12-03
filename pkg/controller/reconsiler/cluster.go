package reconsiler

import (
	"errors"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ClusterResourceState contains runtime data needed in Reconsile()
type ClusterResourceState struct {
	VPC api.Vpc
}

type clusterReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ClusterService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *clusterReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconsile knows how to do what is necessary to ensure the desired state is achieved
func (z *clusterReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	resourceState, ok := node.ResourceState.(ClusterResourceState)
	if !ok {
		return nil, errors.New("error casting cluster resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateCluster(z.commonMetadata.Ctx, api.ClusterCreateOpts{
			ID:                z.commonMetadata.ClusterId,
			Cidr:              resourceState.VPC.Cidr,
			VpcID:             resourceState.VPC.VpcID,
			VpcPrivateSubnets: resourceState.VPC.PrivateSubnets,
			VpcPublicSubnets:  resourceState.VPC.PublicSubnets,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating cluster: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteCluster(z.commonMetadata.Ctx, api.ClusterDeleteOpts{ID: z.commonMetadata.ClusterId})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error deleting cluster: %w", err)
		}
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewClusterReconsiler creates a new reconsiler for the cluster resource
func NewClusterReconsiler(client client.ClusterService) *clusterReconsiler {
	return &clusterReconsiler{
		client: client,
	}
}

