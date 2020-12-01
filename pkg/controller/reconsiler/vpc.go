package reconsiler

import (
	"fmt"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// VPCMetadata contains data extracted from the desired state
type VPCMetadata struct {
	Cidr string
	HighAvailability bool
}

type vpcReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata
	
	client client.VPCService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *vpcReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconsile knows how to ensure the desired state is achieved
func (z *vpcReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	metadata, ok := node.Metadata.(VPCMetadata)
	if !ok {
	    return nil, errors.New("unable to cast VPC metadata")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateVpc(z.commonMetadata.Ctx, api.CreateVpcOpts{
			ID:      z.commonMetadata.ClusterId,
			Cidr:    metadata.Cidr,
			Minimal: !metadata.HighAvailability,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating vpc: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteVpc(z.commonMetadata.Ctx, api.DeleteVpcOpts{ ID: z.commonMetadata.ClusterId})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error deleting vpc: %w", err)
		}
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewVPCReconsiler creates a new reconsiler for the VPC resource
func NewVPCReconsiler(client client.VPCService) *vpcReconsiler {
	return &vpcReconsiler{
		client: client,
	}
}
