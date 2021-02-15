package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// autoscalerReconciler contains service and metadata for the relevant resource
type autoscalerReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.AutoscalerService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *autoscalerReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *autoscalerReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateAutoscaler(z.commonMetadata.Ctx, client.CreateAutoscalerOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("creating autoscaler: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteAutoscaler(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("deleting autoscaler: %w", err)
		}
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewAutoscalerReconciler creates a new reconciler for the autoscaler resource
func NewAutoscalerReconciler(client client.AutoscalerService) Reconciler {
	return &autoscalerReconciler{
		client: client,
	}
}
