package reconsiler

import (
	"errors"
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// AlbIngressControllerResourceState contains runtime data necessary for Reconsile to do its job
type AlbIngressControllerResourceState struct {
	VpcID string
}

type albIngressReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata
	client client.ALBIngressControllerService
}

// SetCommonMetadata stores common metadata for later use
func (z *albIngressReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconsile knows how to do what is necessary to ensure the desired state is achieved
func (z *albIngressReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	state, ok := node.ResourceState.(AlbIngressControllerResourceState)
	if !ok {
	    return nil, errors.New("error casting ALB Ingress Controller state")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateALBIngressController(z.commonMetadata.Ctx, client.CreateALBIngressControllerOpts{
			ID:    z.commonMetadata.ClusterId,
			VPCID: state.VpcID,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating ALB Ingress controller: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteALBIngressController(z.commonMetadata.Ctx, z.commonMetadata.ClusterId)
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error deleting ALB Ingress controller: %w", err)
		}
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewALBIngressReconsiler creates a new reconsiler for the ALB Ingress controller resource
func NewALBIngressReconsiler(client client.ALBIngressControllerService) *albIngressReconsiler {
	return &albIngressReconsiler{
		client: client,
	}
}
