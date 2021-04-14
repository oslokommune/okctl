package reconciler

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// autoscalerReconciler contains service and metadata for the relevant resource
type autoscalerReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	stateHandlers  *clientCore.StateHandlers

	client client.AutoscalerService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *autoscalerReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeAutoscaler
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *autoscalerReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// SetStateHandlers sets the state handlers
func (z *autoscalerReconciler) SetStateHandlers(handlers *clientCore.StateHandlers) {
	z.stateHandlers = handlers
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *autoscalerReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateAutoscaler(z.commonMetadata.Ctx, client.CreateAutoscalerOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("creating autoscaler: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteAutoscaler(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting autoscaler: %w", err)
		}
	}

	return result, nil
}

// NewAutoscalerReconciler creates a new reconciler for the autoscaler resource
func NewAutoscalerReconciler(client client.AutoscalerService) Reconciler {
	return &autoscalerReconciler{
		client: client,
	}
}
