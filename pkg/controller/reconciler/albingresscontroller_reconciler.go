package reconciler

import (
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// AlbIngressControllerResourceState contains runtime data necessary for Reconcile to do its job
type AlbIngressControllerResourceState struct {
	VpcID string
}

// albIngressReconciler contains service and metadata for the relevant resource
type albIngressReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	client         client.ALBIngressControllerService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *albIngressReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeALBIngress
}

// SetCommonMetadata stores common metadata for later use
func (z *albIngressReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *albIngressReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	state, ok := node.ResourceState.(AlbIngressControllerResourceState)
	if !ok {
		return result, errors.New("casting ALB Ingress Controller state")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateALBIngressController(z.commonMetadata.Ctx, client.CreateALBIngressControllerOpts{
			ID:    z.commonMetadata.ClusterID,
			VPCID: state.VpcID,
		})
		if err != nil {
			return result, fmt.Errorf("creating ALB Ingress controller: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteALBIngressController(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting ALB Ingress controller: %w", err)
		}
	}

	return result, nil
}

// NewALBIngressReconciler creates a new reconciler for the ALB Ingress controller resource
func NewALBIngressReconciler(client client.ALBIngressControllerService) Reconciler {
	return &albIngressReconciler{
		client: client,
	}
}
