package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// autoscalerReconciler contains service and metadata for the relevant resource
type autoscalerReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.AutoscalerService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *autoscalerReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeAutoscaler
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *autoscalerReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *autoscalerReconciler) Reconcile(node *dependencytree.Node, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err = z.client.CreateAutoscaler(z.commonMetadata.Ctx, client.CreateAutoscalerOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("creating autoscaler: %w", err)
		}
	case dependencytree.NodeStateAbsent:
		err = z.client.DeleteAutoscaler(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting autoscaler: %w", err)
		}
	}

	return result, nil
}

// NewAutoscalerReconciler creates a new reconciler for the autoscaler resource
func NewAutoscalerReconciler(client client.AutoscalerService) reconciliation.Reconciler {
	return &autoscalerReconciler{
		client: client,
	}
}
